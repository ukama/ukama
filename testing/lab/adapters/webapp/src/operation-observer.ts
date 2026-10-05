/* SPDX-License-Identifier: MPL-2.0
 * Read-only DOM observation and explicitly classified, scoped read faults.
 * No app state, backend mutation, clock manipulation or expected API values.
 */
import type { Page, Route, Request } from 'playwright';
import { Budget, WorkerError, object, keys, str, type ObjectValue } from './contract.js';

type Fault = { mode: string; path: string; epoch: number; applied: number; handler: (route: Route) => Promise<void> };
const faults = new WeakMap<Page, Fault>();
const endpoints = new WeakMap<Page, Map<string, Set<string>>>();
function statusRequest(request: Request): {operation:string;id:string} | undefined {
  try {
    if (request.method() !== 'POST' || !new URL(request.url()).pathname.endsWith('/graphql')) return;
    const body=object(request.postDataJSON()), operation=body.operationName;
    if (operation !== 'GetSiteOperationStatus' && operation !== 'GetNodeOperationStatus') return;
    if (typeof body.query !== 'string' || !new RegExp(`^\\s*query\\s+${operation}\\s*\\(`).test(body.query) || /\bmutation\b/.test(body.query)) return;
    return {operation,id:str(object(body.variables)[operation==='GetSiteOperationStatus'?'siteId':'nodeId'],'status entity')};
  } catch { return; }
}
export function trackStatusRequests(page: Page): void {
  if (endpoints.has(page)) return;
  const known=new Map<string,Set<string>>(); endpoints.set(page,known);
  page.on('framenavigated',frame=>{
    const fault=faults.get(page);
    if (fault && frame === page.mainFrame()) {fault.epoch++;fault.applied=0;}
  });
  page.on('request',request=>{
    const read=statusRequest(request); if(!read)return;
    const key=JSON.stringify(read), urls=known.get(key)??new Set<string>();
    urls.add(request.url());known.set(key,urls);
  });
}
export async function statusFault(page: Page, site: boolean, inputs: ObjectValue, budget: Budget): Promise<void> {
  const mode = str(inputs.value, 'value');
  if (!['read_error', 'idle', 'none'].includes(mode)) throw new WorkerError('INVALID_INPUT', 'Unknown status fault');
  const old = faults.get(page);
  if (old) { faults.delete(page); await page.unroute('**/graphql', old.handler); }
  if (mode === 'none') return;
  const entity = object(inputs.entity), id = str(entity.id, 'entity.id');
  const operation = site ? 'GetSiteOperationStatus' : 'GetNodeOperationStatus';
  const key=JSON.stringify({operation,id});
  const urls=await budget.poll(async()=>endpoints.get(page)?.get(key),value=>!!value?.size,'No operation-status read observed for the opened entity');
  if (urls!.size !== 1) throw new WorkerError('AMBIGUOUS_ENDPOINT','Entity status was read from multiple endpoints');
  const endpoint=Array.from(urls!)[0]!;
  const fault: Fault = {mode, path:new URL(page.url()).pathname, epoch:0, applied:0, handler:async route => {
    const request = route.request();
    const read=statusRequest(request);
    const epoch=fault.epoch;
    // The console sends these two named, single read operations. Batch,
    // foreign-origin, different-entity and mutation requests pass unchanged.
    if (faults.get(page) !== fault || new URL(page.url()).pathname !== fault.path || request.url() !== endpoint || read?.operation !== operation || read.id !== id) {
      await route.fallback(); return;
    }
    const value = site ? {siteId:id,busy:false,degraded:false,nodes:[],actions:{restartSite:{available:true},rf:{available:true},service:{available:true}}} : {nodeId:id,busy:false,operation:null};
    await route.fulfill({status:200,contentType:'application/json',headers:{'access-control-allow-origin':new URL(page.url()).origin,'access-control-allow-credentials':'true'},body:JSON.stringify(mode === 'read_error' ?
      {errors:[{message:'Lab controlled operation-status read failure'}]} : {data:{[site ? 'getSiteOperationStatus' : 'getNodeOperationStatus']:value}})});
    if (faults.get(page) === fault && fault.epoch === epoch) fault.applied++;
  }};
  faults.set(page, fault);
  await page.route('**/graphql', fault.handler);
}
export function faultState(page: Page): string {
  const fault = faults.get(page);
  return fault ? new URL(page.url()).pathname !== fault.path ? `${fault.mode}:wrong_scope` : fault.applied ? `${fault.mode}:applied` : `${fault.mode}:armed` : 'none';
}

// State lives in the document being observed. Reload/navigation loses it and
// fails an outstanding check instead of silently arming a new observation.
export async function watchOperation(page: Page, kind: 'restart' | 'timeout', inputs: ObjectValue): Promise<void> {
  const nodes = kind === 'restart' ? inputs.nodes : [];
  if (!Array.isArray(nodes) || (kind === 'restart' && nodes.length !== 3))
    throw new WorkerError('INVALID_INPUT', 'Restart observation requires exactly one owned trio');
  const trio = nodes.map(raw => {
    const n = object(raw); keys(n, ['id','name','type']);
    return {id:str(n.id,'node id'),name:str(n.name,'node name',true),type:str(n.type,'node type')};
  });
  if (kind === 'restart' && (new Set(trio.map(n=>n.id)).size !== 3 || ['tower','amplifier','controller'].some(type=>trio.filter(n=>n.type===type).length!==1)))
    throw new WorkerError('INVALID_INPUT', 'Restart observation requires unique tower, amplifier and controller');
  const result = await page.evaluate(({kind,trio}) => {
    type State = {kind:string;path:string;violation:string;busy:boolean;started:number;done:boolean;offline:string[];samples:number;timer:number;observer:MutationObserver;read:()=>void};
    const host = window as typeof window & {__ulabOperation?:State};
    if (host.__ulabOperation) return 'Observation already armed';
    const visible = (el:Element) => (el as HTMLElement).getClientRects().length > 0 && getComputedStyle(el).visibility !== 'hidden';
    const read = () => {
      const s = host.__ulabOperation!;
      if (s.done || s.violation) return;
      s.samples++;
      if (location.pathname !== s.path || document.visibilityState !== 'visible') {s.violation='Observation lost its visible detail view';return;}
      const controls = Array.from(document.querySelectorAll('main.main .pagehead button')).filter(visible).filter(e=>
        kind === 'restart' ? /^Site actions(?: • busy)?$/.test(e.textContent?.trim()??'') : /^(Restart node|Restarting…(?: \([^()]*\))?)$/.test(e.textContent?.trim()??''));
      if (controls.length !== 1) {s.violation='Operation control missing or ambiguous';return;}
      const control = controls[0] as HTMLButtonElement;
      const busy = kind === 'restart' ? control.textContent!.includes('busy') : control.textContent!.includes('Restarting');
      if (s.samples === 1 && busy) {s.violation='Observation must start before the operation';return;}
      if (busy && !s.busy) {s.busy=true;s.started=performance.now();}
      if (kind === 'restart') {
        let recovered = true;
        for (const n of trio) {
          const cards = Array.from(document.querySelectorAll('main.main .app-card')).filter(visible).filter(card=>
            Array.from(card.querySelectorAll(':scope > div:first-child > span')).some(e=>[n.id,n.name].filter(Boolean).includes(e.textContent?.trim()??'')));
          if (cards.length !== 1) {s.violation=`Missing or ambiguous ${n.type} identity`;return;}
          const status = cards[0]!.querySelector(':scope > div:first-child > span:last-child')?.textContent?.trim();
          if (n.type === 'controller' && status !== 'is online and well') {s.violation='Controller was not continuously rendered online';return;}
          if (!s.busy && status !== 'is online and well') {s.violation='Trio was not online before restart';return;}
          if (s.busy && status === 'is offline' && !s.offline.includes(n.type)) s.offline.push(n.type);
          if (status !== 'is online and well') recovered=false;
        }
        if (s.busy && !busy && recovered && s.offline.includes('tower') && s.offline.includes('amplifier')) s.done=true;
      } else {
        const disabled = control.disabled || control.getAttribute('aria-disabled') === 'true';
        if (!s.busy && disabled) {s.violation='Restart unavailable before observation';return;}
        if (s.busy && busy && !disabled) {s.violation='Busy restart was available';return;}
        if (s.busy && !busy && !disabled) {
          if (performance.now()-s.started < 9750) s.violation='Optimistic lock released before its ten-second timeout';
          else s.done=true;
        }
      }
    };
    const observer = new MutationObserver(read);
    const state:State={kind,path:location.pathname,violation:'',busy:false,started:0,done:false,offline:[],samples:0,timer:0,observer,read};
    host.__ulabOperation=state;
    read(); observer.observe(document.documentElement,{subtree:true,childList:true,characterData:true,attributes:true});
    state.timer=window.setInterval(read,50);
    // These are listeners, not synthetic focus/visibility changes.
    document.addEventListener('visibilitychange',read);
    window.addEventListener('pagehide',()=>{state.violation='Observation interrupted by navigation';});
    return state.violation;
  }, {kind,trio});
  if (result) throw new WorkerError('OBSERVATION_PRECONDITION', result);
}
export async function observation(page: Page, kind: 'restart' | 'timeout'): Promise<string> {
  return page.evaluate(kind => {
    const s = (window as typeof window & {__ulabOperation?:{kind:string;violation:string;done:boolean;busy:boolean;read:()=>void}}).__ulabOperation;
    if (!s || s.kind !== kind) return 'not armed';
    s.read();
    return s.violation ? `failed: ${s.violation}` : s.done ? 'complete' : s.busy ? 'observing' : 'armed';
  }, kind);
}

/* SPDX-License-Identifier: MPL-2.0
 * DOM-only observations; read faults are controlled evidence, never live credit.
 */
import type {Page, Route, ElementHandle} from 'playwright';
import {Budget, WorkerError} from './contract.js';
export const QUALITY_ACTIONS=['member_fault','member_release','member_retry','focus_cycle'];
export const QUALITY_LABELS=['Members state','Dialog focus','Dialog keyboard cycle','Dialog fits viewport','Action reachable'];
const endpoints=new WeakMap<Page,string>();
const observed=new WeakSet<Page>();
type Fault={url:string; endpoint:string; mode:string; count:number; pending:Array<{release:()=>void;done:Promise<void>}>; handler:(route:Route)=>Promise<void>};
const faults=new WeakMap<Page,Fault>();
const cycles=new WeakMap<Page,{url:string;dialog:ElementHandle}>();
function teamRead(route:Route):boolean {
 const r=route.request();try{const b=r.postDataJSON();return r.method()==='POST'&&b?.operationName==='TeamList'&&/^\s*query\s+TeamList(?:\s|\{)/.test(b.query)&&b.variables&&Object.keys(b.variables).length===0}catch{return false}
}
export function trackQualityReads(page:Page) {
 if(observed.has(page))return;observed.add(page);
 page.on('request',r=>{try{const b=r.postDataJSON();if(r.method()!=='POST'||b?.operationName!=='TeamList'||!/^\s*query\s+TeamList(?:\s|\{)/.test(b.query)||!b.variables||Object.keys(b.variables).length)return;
  endpoints.set(page,endpoints.has(page)&&endpoints.get(page)!==r.url()?'':r.url());
 }catch{/* non-JSON requests are not candidates */}});
 page.on('framenavigated',frame=>{if(frame===page.mainFrame()){const old=cycles.get(page);if(old)void old.dialog.dispose();cycles.delete(page)}});
}
export class UIQuality {
 constructor(private page:Page){}
 private main(){return this.page.locator('main.main')}
 private dialog(){return this.page.getByRole('dialog').filter({visible:true})}
 async run(view:string,action:string,value:string,budget:Budget) {
  if(action==='focus_cycle') {
   const dialog=this.dialog();if(await dialog.count()!==1)throw new WorkerError('DIALOG_REQUIRED','Keyboard cycle requires exactly one visible dialog');
   cycles.delete(this.page);
   const count=await dialog.locator('button,a[href],input,select,textarea,[tabindex]').evaluateAll(es=>es.filter(e=>e instanceof HTMLElement&&e.tabIndex>=0&&!e.matches(':disabled')&&e.getClientRects().length).length);
   if(!count||count>100)throw new WorkerError('FOCUS_PRECONDITION','Dialog requires 1..100 keyboard controls');
   const inside=()=>dialog.evaluate(e=>e.contains(document.activeElement));
   if(!await inside())throw new WorkerError('FOCUS_ESCAPED','Dialog did not receive focus');
   for(const key of ['Tab','Shift+Tab'])for(let n=0;n<count+2;n++){
    budget.remaining();await this.page.keyboard.press(key);
    if(!await inside())throw new WorkerError('FOCUS_ESCAPED','Keyboard traversal escaped the dialog');
   }
   cycles.set(this.page,{url:this.page.url(),dialog:(await dialog.elementHandle())!});return;
  }
  if(view!=='business_members')throw new WorkerError('WRONG_VIEW','Members fault/retry requires the members view');
  if(action==='member_fault') {
   if(!['loading','error'].includes(value)||faults.has(this.page))throw new WorkerError('INVALID_INPUT','Select loading or error once, then release the fault');
   const endpoint=endpoints.get(this.page);if(!endpoint)throw new WorkerError('MISSING_OBSERVER','Requires one observed TeamList endpoint');
   const f:Fault={url:this.page.url(),endpoint,mode:value,count:0,pending:[],handler:async()=>{}};
   f.handler=async route=>{
    if(this.page.url()!==f.url||route.request().url()!==f.endpoint||!teamRead(route)){await route.fallback();return}
    f.count++;
    if(f.mode==='loading'){let release!:()=>void;const gate=new Promise<void>(resolve=>{release=resolve});const done=gate.then(async()=>{if(!this.page.isClosed())await route.continue()});f.pending.push({release,done});await done;return}
    await route.fulfill({status:200,contentType:'application/json',body:JSON.stringify({data:{membersView:{orgName:'',team:{rows:[],error:{code:'LAB_CONTROLLED',message:'Lab controlled member read failure',service:'lab'}}}}})});
   };
   this.page.once('close',()=>f.pending.forEach(p=>p.release()));
   faults.set(this.page,f);await this.page.route('**/*',f.handler);
   await this.page.reload({waitUntil:'domcontentloaded',timeout:budget.remaining()});
   await budget.poll(async()=>f.count,n=>n>0,'Member read fault was not consumed');return;
  }
  if(action==='member_release') {
   const f=faults.get(this.page);if(!f||!f.count||this.page.url()!==f.url)throw new WorkerError('FAULT_PRECONDITION','No consumed member fault on this page');
   for(const p of f.pending)p.release();await Promise.all(f.pending.map(p=>p.done));
   await this.page.unroute('**/*',f.handler);faults.delete(this.page);return;
  }
  if(action==='member_retry') {
   if(faults.has(this.page))throw new WorkerError('FAULT_PRECONDITION','Release the controlled fault before retry');
   if(await this.observe('Members state','')!=='error')throw new WorkerError('RETRY_PRECONDITION','Retry requires the visible member error');
   await this.main().getByRole('button',{name:'Try again',exact:true}).click({timeout:budget.remaining()});return;
  }
  throw new WorkerError('INVALID_INPUT','Unknown UI quality action');
 }
 async observe(label:string,subject:string):Promise<string|null> {
  const dialog=this.dialog();
  if(label==='Members state') {
   const loading=await this.main().locator('.MuiSkeleton-root:visible').count()>0;
   const error=await this.main().getByText("Couldn't load members",{exact:true}).isVisible();
   const empty=await this.main().getByText('No members yet',{exact:true}).isVisible();
   const data=await this.main().locator('tbody > tr:visible').count()>0;
   if([loading,error,empty,data].filter(Boolean).length!==1)return 'ambiguous';
   return loading?'loading':error?'error':empty?'empty':'data';
  }
  if(label==='Action reachable') {
   const scope=await dialog.count()?dialog:this.main();
   const button=scope.getByRole('button',{name:subject,exact:true}).filter({visible:true});
   if(await button.count()>1)throw new WorkerError('AMBIGUOUS_LOCATOR','Action matches multiple buttons');
   if(await button.count()!==1)return null;
   try{await button.click({trial:true,timeout:500});return 'true'}catch{return 'false'}
  }
  if(await dialog.count()!==1)return null;
  if(label==='Dialog focus')return String(await dialog.evaluate(e=>e.contains(document.activeElement)));
  if(label==='Dialog fits viewport')return String(await dialog.evaluate(e=>{const r=e.getBoundingClientRect();return r.width>0&&r.left>=0&&r.right<=innerWidth&&r.top>=0&&r.bottom<=innerHeight&&e.scrollWidth<=e.clientWidth+1}));
  if(label==='Dialog keyboard cycle') {
   const proof=cycles.get(this.page);if(!proof||proof.url!==this.page.url())return null;
   return String(await dialog.evaluate((e,prior)=>e===prior&&e.contains(document.activeElement),proof.dialog));
  }
  throw new WorkerError('INVALID_INPUT','Unknown UI quality label');
 }
}

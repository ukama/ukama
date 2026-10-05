/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 * Read-only DOM assertions and closed navigation commands for the original UI.
 */
import type { Locator, Page, Route } from 'playwright';
import { ConsoleApp, getView, assertSession } from './console-app.js';
import { Budget, WorkerError, canonical, keys, normalize, object, str, type ObjectValue } from './contract.js';
type Entity = {ref:string;id:string;name?:string;network_ref?:string;site_ref?:string};
type Scope = {network:Entity;sites:Entity[];nodes:Entity[];site?:Entity;node?:Entity;view:string;path:string};
type State = {fault?:string;epoch?:number;applied:number;handler?:(route:Route)=>Promise<void>;selection?:string;pin?:number};
const state = new WeakMap<Page,State>();
function ctx(i:ObjectValue):Scope {
  const context=object(i.context,'context'); keys(context,['network','sites','nodes']);
  const entity=(raw:unknown):Entity=>{const e=object(raw,'entity');keys(e,['ref','id','name','network_ref','site_ref']);return {ref:str(e.ref,'ref'),id:str(e.id,'id'),...(e.name===undefined?{}:{name:str(e.name,'name')}),...(e.network_ref===undefined?{}:{network_ref:str(e.network_ref,'network_ref')}),...(e.site_ref===undefined?{}:{site_ref:str(e.site_ref,'site_ref')})};};
  if (!Array.isArray(context.sites)||!Array.isArray(context.nodes))throw new WorkerError('INVALID_INPUT','Inventory context requires sites and nodes');
  const network=entity(context.network),sites=context.sites.map(entity),nodes=context.nodes.map(entity),view=str(i.view,'view');
  if(!network.name||!['network_home','network_sites','network_nodes','network_site_detail','network_node_detail'].includes(view))throw new WorkerError('INVALID_INPUT','Unsupported inventory view or network');
  const site=i.site_ref===undefined?undefined:sites.find(s=>s.ref===str(i.site_ref,'site_ref'));
  const node=i.node_ref===undefined?undefined:nodes.find(n=>n.ref===str(i.node_ref,'node_ref'));
  if((i.site_ref!==undefined&&!site)||(i.node_ref!==undefined&&!node)||(site&&node)||(site&&site.network_ref!==network.ref)||(node&&node.network_ref!==network.ref))throw new WorkerError('WRONG_ENTITY','Inventory reference is outside selected world network');
  const detail=getView(view).detail,selected=detail==='site'?site:node;
  if(detail&&(!selected||!/^[A-Za-z0-9_-]+$/.test(selected.id)))throw new WorkerError('INVALID_INPUT','Detail needs a bound entity ID');
  return {network,sites,nodes,site,node,view,path:getView(view).path+(detail?'/'+encodeURIComponent(selected!.id):'')};
}
async function text(l:Locator):Promise<string|null>{const v=l.filter({visible:true}),n=await v.count();if(n>1)throw new WorkerError('AMBIGUOUS_LOCATOR','Inventory field is ambiguous');return n===1?normalize(await v.innerText({timeout:500})):null;}
export class Inventory {
  private s:State;
  constructor(private page:Page,private origin:string,private app:ConsoleApp){this.s=state.get(page)??{applied:0};state.set(page,this.s);}
  private main(){return this.page.locator('main.main');}
  private async ready(c:Scope){const u=new URL(this.page.url());return u.origin===this.origin&&u.pathname===c.path&&await this.main().isVisible()&&await text(this.page.locator('header .netswitch .nm'))===c.network.name;}
  private async require(c:Scope,b:Budget){await b.poll(()=>this.ready(c),Boolean,'Inventory route/network differs from the expected world scope');}
  private siteCard(c:Scope){if(!c.site)throw new WorkerError('INVALID_INPUT','Site field requires a site');return this.main().locator('.ecard[role="button"]').filter({has:this.page.getByText(c.site.name!,{exact:true})});}
  private info(){return this.main().locator('.card').filter({has:this.page.locator('.sec-title').getByText('Site information',{exact:true})});}
  private location(){return this.info().getByText('Location',{exact:true}).locator('..');}
  private panel(){return this.main().getByRole('button',{name:'Open site',exact:true}).locator('..');}
  private async arm(c:Scope){
    const foreign=[...c.sites.filter(s=>s.network_ref!==c.network.ref).map(s=>s.name!),...c.nodes.filter(n=>n.network_ref!==c.network.ref).map(n=>n.id)];
    if(!foreign.length)throw new WorkerError('INVALID_INPUT','Scope watch needs foreign world identities');
    await this.page.evaluate(({name,foreign})=>{
      const w=window as unknown as {__ulabScope?:{observer:MutationObserver;frame:number;violations:string[];target:string;samples:number}};
      if(w.__ulabScope){w.__ulabScope.observer.disconnect();cancelAnimationFrame(w.__ulabScope.frame);}
      const record={observer:null as unknown as MutationObserver,frame:0,violations:[] as string[],target:name,samples:0};w.__ulabScope=record;
      const shown=(e:Element)=>{const r=e.getBoundingClientRect(),s=getComputedStyle(e);return r.width>0&&r.height>0&&s.visibility!=='hidden'&&s.display!=='none';};
      const sample=()=>{if(document.querySelector('header .netswitch .nm')?.textContent?.trim()!==name)return;record.samples++;for(const el of Array.from(document.querySelectorAll('main.main .ecard, main.main .kv-row, main.main .detail-subrow, main.main .pagehead, main.main .leaflet-popup-content'))){
        if(!shown(el))continue;const t=(el as HTMLElement).innerText;for(const f of foreign)if(t.includes(f)&&!record.violations.includes(f))record.violations.push(f);
      }};
      record.observer=new MutationObserver(sample);record.observer.observe(document.documentElement,{subtree:true,childList:true,characterData:true,attributes:true});
      const tick=()=>{sample();record.frame=requestAnimationFrame(tick);};tick();
    },{name:c.network.name!,foreign});
  }
  private async mask(value:string){
    if(!['kpis','registry','reset'].includes(value))throw new WorkerError('INVALID_INPUT','Unknown controlled home fault');
    const epoch=(this.s.epoch??0)+1;this.s.epoch=epoch;
    if(this.s.handler)await this.page.unroute('**/graphql',this.s.handler);
    this.s.handler=undefined;this.s.fault=value;this.s.applied=0;
    if(value==='reset')return;
    this.s.handler=async route=>{
      let body;try{body=route.request().postDataJSON();}catch{await route.continue();return;}
      const query=String(body?.query??''),op=body?.operationName;
      if(!/^\s*(?:#[^\n]*\n\s*)*(?:query|fragment)\b/.test(query)||!((value==='kpis'&&op==='GetKpiValues')||(value==='registry'&&op==='SitesList'))){await route.continue();return;}
      const response=await route.fetch();const json=await response.json();
      if(this.s.epoch!==epoch){await route.fulfill({response,json});return;}
      if(value==='kpis'&&json.data?.getKpiValues){json.data.getKpiValues.values=[];this.s.applied++;}
      else if(value==='registry'&&json.data?.sitesView?.sites){json.data.sitesView.sites={sites:[],error:{code:'INVENTORY_TEST_UNAVAILABLE',message:'Controlled registry unavailable'}};this.s.applied++;}
      await route.fulfill({response,json});
    };
    await this.page.route('**/graphql',this.s.handler);
  }
  async run(i:ObjectValue,b:Budget){
    keys(i,['view','context','site_ref','node_ref','action','value']);const c=ctx(i),action=str(i.action,'action'),value=str(i.value??'','value',true);
    if(!['switch','switch_watch','search','direct','back','map_select','map_open','map_clear','stale_selection','mask_home'].includes(action))throw new WorkerError('INVALID_INPUT','Unknown inventory action');
    if(action==='switch'||action==='switch_watch'){
      await assertSession(this.page,this.origin,b);if(action==='switch_watch')await this.arm(c);
      await this.app.selectNetwork(c.network.name!,b);return;
    }
    if(action==='direct'){this.app.forgetView();await this.page.goto(this.origin+c.path,{waitUntil:'domcontentloaded',timeout:b.remaining()});return;}
    if(action==='back'){this.app.forgetView();await this.page.goBack({waitUntil:'domcontentloaded',timeout:b.remaining()});return;}
    await this.require(c,b);
    if(action==='search'){
      if(c.view!=='network_sites'||(c.site&&value))throw new WorkerError('INVALID_INPUT','Search is a literal or owned site name on Sites');
      await this.main().getByPlaceholder('Search sites',{exact:true}).fill(c.site?.name??value,{timeout:b.remaining()});return;
    }
    if(action==='stale_selection'){
      // Controlled stale preference inside this scenario's isolated browser;
      // no server resource is removed and no other preference is rewritten.
      await this.page.evaluate(()=>{const key='uk-ui-prefs',raw=localStorage.getItem(key);if(!raw)throw Error('Missing persisted UI preference');const p=JSON.parse(raw);if(!p.state||typeof p.state.networkId!=='string')throw Error('Unknown UI preference schema');p.state.networkId='ulab-removed-network-selection';localStorage.setItem(key,JSON.stringify(p));});
      this.app.forgetView();await this.page.reload({waitUntil:'domcontentloaded',timeout:b.remaining()});return;
    }
    if(action==='mask_home'){if(c.view!=='network_home')throw new WorkerError('INVALID_INPUT','Home faults require Home');await this.mask(value);return;}
    if(c.view!=='network_home')throw new WorkerError('INVALID_INPUT','Map actions require Home');
    if(action==='map_clear'){await this.panel().getByRole('button',{name:'Clear',exact:true}).click({timeout:b.remaining()});this.s.selection=undefined;this.s.pin=undefined;return;}
    if(!c.site)throw new WorkerError('INVALID_INPUT','Map selection requires site');
    if(action==='map_open'){
      if(await text(this.panel().locator(':scope > div > div:first-child > span'))!==c.site.name)throw new WorkerError('WRONG_ENTITY','Selected map site differs from the requested site');
      this.app.forgetView();await this.panel().getByRole('button',{name:'Open site',exact:true}).click({timeout:b.remaining()});return;
    }
    const pins=this.main().locator('.leaflet-container .uk-map-pin');
    await pins.first().waitFor({timeout:b.remaining()});const count=await pins.count();
    for(let n=0;n<count;n++){
      await pins.nth(n).click({timeout:b.remaining()});
      const popup=this.main().locator('.leaflet-popup-content');await popup.waitFor({timeout:b.remaining()});
      if(await text(popup.locator(':scope > div > div:first-child'))===c.site.name){this.s.selection=c.site.ref;this.s.pin=n;return;}
    }
    throw new WorkerError('MAP_SITE_MISSING','No visible map marker opens the requested site',{site:c.site.name,count});
  }
  async check(i:ObjectValue,b:Budget){
    keys(i,['view','context','site_ref','node_ref','label','expected','requirement']);const c=ctx(i),label=str(i.label,'label'),expected=i.expected;
    if(!/^WEB-[A-Z0-9-]+$/.test(str(i.requirement,'requirement'))||!(typeof expected==='string'||(Array.isArray(expected)&&expected.every(v=>typeof v==='string'))))throw new WorkerError('INVALID_INPUT','Invalid inventory expectation');
    const actual=await b.poll(async():Promise<unknown>=>{
      if(label==='Selection valid'){
        const url=new URL(this.page.url());if(url.origin!==this.origin||url.pathname!==c.path)return null;
        const switcher=this.page.locator('header button.netswitch'),name=await text(switcher.locator('.nm'));if(!name)return null;
        await switcher.click({timeout:b.remaining()});const count=await this.page.getByRole('menuitem').filter({has:this.page.getByText(name,{exact:true})}).count();
        await this.page.keyboard.press('Escape');return String(count===1);
      }
      if(!await this.ready(c))return null;
      if(label==='Path')return new URL(this.page.url()).pathname;
      if(label==='Selected network')return text(this.page.locator('header .netswitch .nm'));
      if(label==='Scope leaks'){
        const watch=await this.page.evaluate(()=>{const w=window as unknown as {__ulabScope?:{target:string;violations:string[];samples:number}};return w.__ulabScope?{target:w.__ulabScope.target,violations:w.__ulabScope.violations,samples:w.__ulabScope.samples}:null;});
        if(!watch||watch.target!==c.network.name||watch.samples<1)throw new WorkerError('MISSING_OBSERVER','No scope observer for this switch');
        if(watch.violations.length)throw new WorkerError('TRANSIENT_SCOPE_LEAK','Foreign world data was visible under the selected network',watch);
        return 'none';
      }
      if(label==='KPI fault applied')return this.s.applied>0?this.s.fault:null;
      if(await this.main().locator('.MuiSkeleton-root:visible').count())return null;
      if(label==='Site names'){
        const cards=this.main().locator('.ecard[role="button"]:visible');
        if(!await cards.count())return await this.main().getByText(/^(No sites yet|No matching sites)$/).isVisible()?[]:null;
        const names=await cards.locator(':scope > div:first-child > div:first-child > div:last-child > div:first-child').allInnerTexts();
        if(names.length!==await cards.count())throw new WorkerError('MALFORMED_INVENTORY_CARD','Site cards do not each expose exactly one name');
        return names.map(normalize).sort();
      }
      if(label==='Node IDs'){
        const cards=this.main().locator('.ecard[role="button"]:visible');if(!await cards.count())return await this.main().getByText('No nodes yet',{exact:true}).isVisible()?[]:null;
        const ids=await cards.locator('.tnum').allInnerTexts();
        if(ids.length!==await cards.count()||ids.some(t=>normalize(t).split(' · ').length!==2))throw new WorkerError('MALFORMED_INVENTORY_CARD','Node cards do not each expose exactly one model and serial');
        return ids.map(t=>normalize(t).split(' · ').at(-1)!).sort();
      }
      if(label==='Empty search')return String(await this.main().getByText('No matching sites',{exact:true}).isVisible()&&await this.main().locator('.ecard[role="button"]:visible').count()===0);
      if(label==='Header count'){
        const count=await text(this.main().locator('.pagehead .pagetitle .cnt'));
        return count??(await this.main().getByText(c.view==='network_sites'?'No sites yet':'No nodes yet',{exact:true}).isVisible()?'0':null);
      }
      if(label==='Site title')return text(this.main().locator('.pagehead .pagetitle'));
      if(label==='Site location')return c.view==='network_sites'?text(this.siteCard(c).locator(':scope > div:first-child > div:first-child > div:last-child > div:nth-child(2)')):text(this.location().locator(':scope > div:nth-child(2)'));
      if(label==='Site coordinates')return text(this.location().locator(':scope > .tnum'));
      if(label==='Site node count')return c.view==='network_sites'?text(this.siteCard(c).locator(':scope > div:last-child > span')):String(await this.info().getByText('Nodes',{exact:true}).locator('..').locator(':scope > div > span.tnum').count());
      if(label==='Site node IDs')return (await this.info().getByText('Nodes',{exact:true}).locator('..').locator(':scope > div > span.tnum').allInnerTexts()).map(normalize).sort();
      if(label==='Site status')return text((c.view==='network_sites'?this.siteCard(c):this.main().locator('.detail-subrow')).locator('.MuiChip-label'));
      if(label==='Node serial')return text(this.main().locator('.kv-row').filter({has:this.page.getByText('Serial #',{exact:true})}).locator(':scope > span > .tnum'));
      if(label==='Map count')return await this.main().locator('.leaflet-container').isVisible()?String(await this.main().locator('.leaflet-container .uk-map-pin').count()):null;
      if(label==='Map clear')return String(await this.main().getByRole('button',{name:'Open site',exact:true}).count()===0);
      if(label==='Map selection')return text(this.panel().locator(':scope > div > div:first-child > span'));
      if(label==='Map location')return text(this.panel().locator(':scope > div > div:nth-child(2)'));
      if(label==='Map color'){
        if(this.s.selection!==c.site?.ref||this.s.pin===undefined)throw new WorkerError('WRONG_ENTITY','Map color requires prior selection of the same site');
        const svg=this.main().locator('.leaflet-container .uk-map-pin').nth(this.s.pin).locator('svg');
        if(!await svg.isVisible())return null;
        return svg.getAttribute('fill');
      }
      throw new WorkerError('INVALID_INPUT','Unsupported inventory label');
    },v=>v!==null&&canonical(v)===canonical(expected),'Visible inventory differs from the scenario world or expectation');
    return {expected,actual};
  }
}

/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 * Browser-only team/preferences/support checks. Probes never forward writes.
 */
import type {Page, Request, Route} from 'playwright';
import {Budget, WorkerError, normalize} from './contract.js';
export const TEAM_ACTIONS=['settings_tab','invite_probe','map_preference','support_lookup','support_node','support_copy','support_restart_probe','support_failure','support_clear_fault'];
export const TEAM_LABELS=['Member row','Invitation probe','Map preference','Support identity','Support field','Support summary','Support state','Support restart policy','Support fault'];
const roles:Record<string,string>={Owner:'ROLE_OWNER',Administrator:'ROLE_ADMIN','Network owner':'ROLE_NETWORK_OWNER',Vendor:'ROLE_VENDOR'};
const reads=new WeakMap<Page,Map<string,{url:string,network:string}>>();
const invitation=new WeakMap<Page,{url:string,email:string,role:string}>();
const preferences=new WeakMap<Page,string>();
const summaries=new WeakMap<Page,{url:string,identity:string,text:string}>();
const restart=new WeakMap<Page,{url:string,identity:string,policy:string}>();
const faults=new WeakMap<Page,{path:string,count:number,handler:(route:Route)=>Promise<void>}>();
const readOnlyPages=new WeakSet<Page>();
function envelope(request:Request):{name:string,variables:Record<string,unknown>,query:string}|null {
 try {const b=request.postDataJSON();if(!b||Array.isArray(b)||typeof b.query!=='string'||typeof b.operationName!=='string'||!b.variables||Array.isArray(b.variables)||typeof b.variables!=='object')return null;return {name:b.operationName,variables:b.variables,query:b.query.replace(/^\s*(?:#[^\n]*\n\s*)*/,'').trim()};}catch{return null;}
}
const namedQuery=(r:Request)=>{const b=envelope(r);return b && new RegExp('^query\\s+'+b.name.replace(/[^A-Za-z0-9_]/g,'!')+'(?:\\s|\\(|\\{)').test(b.query)?b:null;};
export function trackTeamReads(page:Page) {
 if(reads.has(page))return;const observed=new Map<string,{url:string,network:string}>();reads.set(page,observed);
 page.on('framenavigated',frame=>{
  if(frame!==page.mainFrame())return;
  invitation.delete(page);summaries.delete(page);restart.delete(page);
  const fault=faults.get(page);if(fault)fault.count=0;
 });
 page.on('request',request=>{
  const b=namedQuery(request);if(request.method()!=='POST'||!b||!['TeamList','NetworkCustomers','SitesList','NodesList'].includes(b.name))return;
  if(b.name==='TeamList'?Object.keys(b.variables).length!==0:Object.keys(b.variables).join()!=='networkId'||typeof b.variables.networkId!=='string')return;
  const key=b.name+':'+String(b.variables.networkId??'');const previous=observed.get(key);
  observed.set(key,{url:previous&&previous.url!==request.url()?'':request.url(),network:String(b.variables.networkId??'')});
 });
}
async function passRead(route:Route) {
 const r=route.request();
 // GraphQL POST queries are read-only; unknown POST/PUT/PATCH/DELETE and
 // GraphQL-over-GET URLs are blocked during mutation-policy probes.
 if(r.method()==='OPTIONS'||((r.method()==='GET'||r.method()==='HEAD')&&!/[?&](query|operationName)=/i.test(r.url()))||namedQuery(r)) {await route.fallback();return true;}
 return false;
}
async function readOnly(page:Page) {
 if(readOnlyPages.has(page))return;
 // Retain the write barrier through the rest of the page lifetime. Delayed
 // requests must not escape after a probe returns or its temporary route ends.
 await page.route('**/*',async route=>{if(!await passRead(route))await route.abort('blockedbyclient')});
 readOnlyPages.add(page);
}
async function singleText(locator:ReturnType<Page['locator']>) {
 const v=locator.filter({visible:true}),n=await v.count();if(n>1)throw new WorkerError('AMBIGUOUS_LOCATOR','Team/support observation is ambiguous');return n?normalize(await v.innerText({timeout:500})):null;
}
export class TeamSupport {
 constructor(private page:Page){}
 private main(){return this.page.locator('main.main')}
 private dialog(){return this.page.getByRole('dialog').filter({visible:true})}
 private field(label:string){return this.dialog().locator('.ff').filter({has:this.page.locator('.ff-label').filter({hasText:new RegExp('^'+label+'(?:\\s*\\*)?$')})})}
 private detail(){return this.main().locator('.tile-grid').first().locator(':scope > .card.card-pad').filter({has:this.page.locator('.tile-grid')})}
 private identity(){return singleText(this.detail().locator(':scope > div:first-child > span:first-child'))}
 private async pref(){const choice=this.main().getByRole('radiogroup',{name:'Map view',exact:true}).locator('[role="radio"][aria-checked="true"]');return singleText(choice.locator(':scope > div:first-child'))}
 async run(view:string,action:string,label:string,value:string,budget:Budget,networkId?:string) {
  const click=(loc:ReturnType<Page['locator']>)=>loc.click({timeout:budget.remaining()});
  if(action==='settings_tab') {
   if(!view.endsWith('_settings')||!['Preferences','My account','Organization'].includes(value))throw new WorkerError('INVALID_INPUT','Unsupported settings tab');
   await click(this.main().getByRole('tab',{name:value,exact:true}));return;
  }
  if(action==='invite_probe') {
   if(view!=='business_members'||!roles[label]||!/^[-a-z0-9._+]+@[-a-z0-9.]+\.test$/.test(value))throw new WorkerError('INVALID_INPUT','Invitation probe needs a listed role and a lowercase reserved .test address');
   const endpoint=reads.get(this.page)?.get('TeamList:')?.url;if(!endpoint)throw new WorkerError('MISSING_OBSERVER','Invitation probe requires one observed TeamList endpoint');
   invitation.delete(this.page);let matched=0,blocked=0;const path=this.page.url();
   const errors=this.page.getByText('Lab controlled invitation rejection',{exact:true}).filter({visible:true});
   const previous=await errors.elementHandles();
   const handler=async(route:Route)=>{
    if(await passRead(route))return;
    const r=route.request(),b=envelope(r),data=b?.variables.data as Record<string,unknown>|undefined;
    const correct=this.page.url()===path&&r.url()===endpoint&&r.method()==='POST'&&b?.name==='CreateInvitation'&&/^mutation\s+CreateInvitation(?:\s|\()/.test(b.query)&&Object.keys(b.variables).join()==='data'&&data&&Object.keys(data).sort().join()==='email,name,role'&&data.email===value&&data.name===value.split('@')[0]&&data.role===roles[label];
    if(correct&&++matched===1)await route.fulfill({status:200,contentType:'application/json',body:JSON.stringify({errors:[{message:'Lab controlled invitation rejection'}],data:null})});
    else{blocked++;await route.abort('blockedbyclient');}
   };
   await readOnly(this.page);
   await this.page.route('**/*',handler);
   try {
    await click(this.main().locator('.pagehead').getByRole('button',{name:'Invite member',exact:true}));
    await this.field('Email').locator('input').fill(value,{timeout:budget.remaining()});
    await this.field('Role').locator('select').selectOption({label},{timeout:budget.remaining()});
    await click(this.dialog().getByRole('button',{name:'Invite member',exact:true}));
    await budget.poll(async()=>{
     for(const current of await errors.elementHandles()) {
      let old=false;for(const node of previous)if(await current.evaluate((e,before)=>e===before,node))old=true;
      await current.dispose();if(!old)return true;
     }
     return false;
    },Boolean,'Invitation did not show a fresh rejection for this attempt');
    if(matched!==1||blocked||!await this.dialog().isVisible())throw new WorkerError('UNSAFE_INVITATION','Invitation probe did not prove one exact rejected attempt');
    invitation.set(this.page,{url:path,email:value,role:label});
   } finally {await this.page.unroute('**/*',handler);for(const node of previous)await node.dispose();}
   return;
  }
  if(action==='map_preference') {
   if(!view.endsWith('_settings'))throw new WorkerError('WRONG_VIEW','Map preference requires settings');
   await click(this.main().getByRole('tab',{name:'Preferences',exact:true}));
   const current=await this.pref();if(!current)throw new WorkerError('UNSUPPORTED_LOCATOR','Map preference has no unique selected radio');
   if(!preferences.has(this.page))preferences.set(this.page,current);
   const target=value==='Restore'?preferences.get(this.page)!:value;
   if(!['Street','Satellite','Terrain'].includes(target))throw new WorkerError('INVALID_INPUT','Unsupported map preference');
   const radio=this.main().getByRole('radiogroup',{name:'Map view',exact:true}).getByRole('radio').filter({has:this.page.getByText(target,{exact:true})});
   await click(radio);await budget.poll(()=>this.pref(),v=>v===target,'Map preference selection did not update');return;
  }
  if(!['business_support','network_support'].includes(view))throw new WorkerError('WRONG_VIEW','Support operation requires a support view');
  if(action==='support_node') {
   if(view!=='network_support'||!/^[A-Za-z0-9_-]+$/.test(value))throw new WorkerError('INVALID_INPUT','Node lookup requires a bound node identifier');
   summaries.delete(this.page);restart.delete(this.page);
   await this.main().getByPlaceholder('Search site or node',{exact:true}).fill('',{timeout:budget.remaining()});
   const list=this.main().locator('.card').filter({has:this.page.locator('.sec-title').filter({hasText:/^Sites & nodes$/})}).getByRole('button');
   const count=await list.count();if(count>200)throw new WorkerError('UNSUPPORTED_LOCATOR','Support list exceeds the bounded scan');
   const matches:number[]=[];
   for(let n=0;n<count;n++) {await click(list.nth(n));if(await this.observe(view,'Support field','Node ID')===value)matches.push(n);}
   if(matches.length!==1)throw new WorkerError('WRONG_ENTITY','Support list does not identify exactly one owned node');
   await click(list.nth(matches[0]!));
   await budget.poll(()=>this.observe(view,'Support field','Node ID'),id=>id===value,'Selected support node changed');return;
  }
  if(action==='support_lookup') {
   summaries.delete(this.page);restart.delete(this.page);
   await this.main().getByPlaceholder(view==='network_support'?'Search site or node':'Search customer by name or phone',{exact:true}).fill(value,{timeout:budget.remaining()});
   const search=this.main().getByRole('button',{name:'Search',exact:true});if(await search.isEnabled())await click(search);
   const cards=this.main().locator('.card').filter({has:this.page.locator('.sec-title').filter({hasText:/^(Customers|Sites & nodes)$/})});
   const results=cards.getByRole('button').filter({has:this.page.getByText(value,{exact:true})});
   if(await results.count()>1)throw new WorkerError('AMBIGUOUS_LOCATOR','Support lookup matches more than one exact identity');
   if(await results.count()===1)await click(results);return;
  }
  if(action==='support_copy') {
   const identity=await this.identity();if(!identity)throw new WorkerError('UNSUPPORTED_LOCATOR','No unique support detail to copy');
   await this.page.context().grantPermissions(['clipboard-read','clipboard-write'],{origin:new URL(this.page.url()).origin});
   // Clear only the isolated browser context clipboard, so a broken copy
   // cannot reuse a previous customer's summary.
   await this.page.evaluate(()=>navigator.clipboard.writeText(''));
   await click(this.detail().getByRole('button',{name:'Copy summary',exact:true}));
   const copied=await budget.poll(()=>this.page.evaluate(()=>navigator.clipboard.readText()),v=>v.length>0,'Support did not copy a summary');
   if(await this.identity()!==identity)throw new WorkerError('WRONG_ENTITY','Support selection changed during copy');
   summaries.set(this.page,{url:this.page.url(),identity,text:copied});return;
  }
  if(action==='support_restart_probe') {
   if(view!=='network_support'||!['Restart node','Restart site'].includes(value))throw new WorkerError('INVALID_INPUT','Unsupported support restart probe');
   const identity=await this.identity();if(!identity)throw new WorkerError('WRONG_ENTITY','No selected support identity');
   restart.delete(this.page);const control=this.detail().getByRole('button',{name:value,exact:true});
   if(!await control.isVisible())throw new WorkerError('UNSUPPORTED_LOCATOR','Restart control is absent');
   if(await control.isDisabled()){restart.set(this.page,{url:this.page.url(),identity,policy:'locked'});return;}
   let writes=0;const handler=async(route:Route)=>{if(await passRead(route))return;writes++;await route.abort('blockedbyclient');};
   await readOnly(this.page);
   await this.page.route('**/*',handler);
   try {
    await click(control);
    const policy=await budget.poll(async()=>{
     if(writes)return 'unguarded';
     const dialog=this.dialog();if(await dialog.count()!==1)return null;
     const title=await singleText(dialog.getByRole('heading'));
     return title===value&&await dialog.getByRole('button',{name:'Cancel',exact:true}).isVisible()?'confirmation':null;
    },v=>v!==null,'Support restart did not show confirmation or a blocked write');
    if(policy==='confirmation')await click(this.dialog().getByRole('button',{name:'Cancel',exact:true}));
    if(await this.identity()!==identity)throw new WorkerError('WRONG_ENTITY','Restart probe selection changed');
    restart.set(this.page,{url:this.page.url(),identity,policy:writes?'unguarded':policy!});
   }finally{await this.page.unroute('**/*',handler);}
   return;
  }
  if(action==='support_clear_fault') {const f=faults.get(this.page);if(f)await this.page.unroute('**/*',f.handler);faults.delete(this.page);await this.page.reload({waitUntil:'domcontentloaded',timeout:budget.remaining()});return;}
  if(action==='support_failure') {
   if(faults.has(this.page))throw new WorkerError('INVALID_INPUT','Support fault already armed');
   if(!networkId)throw new WorkerError('INVALID_INPUT','Support fault requires the owned network ID');
   const names=view==='business_support'?['NetworkCustomers']:['SitesList','NodesList'];
   const sources=names.map(n=>({name:n,...reads.get(this.page)?.get(n+':'+networkId)}));
   if(sources.some(s=>!s.url))throw new WorkerError('MISSING_OBSERVER','Support fault requires its observed scoped query endpoints');
   const path=this.page.url();const record={path,count:0,handler:async(route:Route)=>{
    const r=route.request(),b=namedQuery(r);
    if(this.page.url()===path&&r.method()==='POST'&&b&&Object.keys(b.variables).join()==='networkId'&&b.variables.networkId===networkId&&sources.some(s=>s.name===b.name&&s.url===r.url())) {record.count++;await route.fulfill({status:200,contentType:'application/json',body:JSON.stringify({errors:[{message:'Lab controlled support read failure'}],data:null})});}
    else await route.fallback();
   }};
   faults.set(this.page,record);await this.page.route('**/*',record.handler);await this.page.reload({waitUntil:'domcontentloaded',timeout:budget.remaining()});return;
  }
  throw new WorkerError('INVALID_INPUT','Unknown team/support action');
 }
 async observe(view:string,label:string,subject:string):Promise<string|null> {
  if(label==='Member row') {
   if(view!=='business_members')throw new WorkerError('WRONG_VIEW','Member row requires Members');
   if(await this.main().locator('.MuiSkeleton-root:visible').count())return null;
   const row=this.main().getByRole('table').locator('tbody > tr').filter({has:this.page.getByText(subject,{exact:true})});
   if(await row.count()>1)throw new WorkerError('AMBIGUOUS_LOCATOR','Duplicate member email');
   if(!await row.isVisible())return null;const cells=row.locator(':scope > td');
   if(await cells.count()!==5)return null;
   const name=await singleText(cells.nth(0).locator(':scope > div > div > div:first-child'));
   const email=await singleText(cells.nth(0).locator(':scope > div > div > div:last-child'));
   const role=await singleText(cells.nth(1).locator(':scope > div:first-child'));
   const status=await singleText(cells.nth(3));return name&&email&&role&&status?[name,email,role,status].join('|'):null;
  }
  if(label==='Invitation probe') {const p=invitation.get(this.page);return p&&p.url===this.page.url()&&p.email===subject?'rejected':'not applied';}
  if(label==='Map preference')return this.pref();
  if(!['business_support','network_support'].includes(view))throw new WorkerError('WRONG_VIEW','Support assertion requires support');
  if(label==='Support identity')return this.identity();
  if(label==='Support field')return singleText(this.detail().locator(':scope > .tile-grid > div').filter({has:this.page.getByText(subject,{exact:true})}).locator(':scope > div.tnum'));
  if(label==='Support summary') {
   const summary=summaries.get(this.page);if(!summary||summary.url!==this.page.url()||summary.identity!==await this.identity())throw new WorkerError('MISSING_OBSERVER','Summary requires copying this selected entity');
   const lines=summary.text.split(/\r?\n/).filter(l=>l.startsWith(subject+': '));return lines.length===1?normalize(lines[0]!.slice(subject.length+2)):null;
  }
  if(label==='Support restart policy') {const p=restart.get(this.page);return p&&p.url===this.page.url()&&p.identity===await this.identity()?p.policy:null;}
  if(label==='Support fault'){const f=faults.get(this.page);return f&&f.path===this.page.url()&&f.count>0?'applied':'not applied';}
  if(label==='Support state') {
   if(await this.main().getByText('Loading…',{exact:true}).isVisible())return null;
   const error=await this.main().getByText("Couldn't load support",{exact:true}).isVisible();
   const empty=await this.main().getByText('No match',{exact:true}).isVisible();
   const identity=await this.identity();if(Number(error)+Number(empty)+Number(!!identity)!==1)return null;
   return error?'error':empty?'no match':'match';
  }
  throw new WorkerError('INVALID_INPUT','Unknown team/support observation');
 }
}

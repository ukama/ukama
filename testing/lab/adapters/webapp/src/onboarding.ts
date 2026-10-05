/* SPDX-License-Identifier: MPL-2.0
 * Staged browser onboarding. Responses establish cleanup ownership only;
 * assertions read rendered UI. Fault cases are controlled evidence, never live.
 */
import type { Locator, Page, Route } from 'playwright';
import { ConsoleApp } from './console-app.js';
import { Creation } from './provisioning.js';
import { Budget, WorkerError, keys, normalize, object, str, type ObjectValue } from './contract.js';
const headings = ['Let\'s set up your network','Name your network','Select a network','Install your site','Bring your site online','Name your site','Configure site settings','Upload SIMs',"You're all set!"];
const fields = ['Network name','Site name'];
const components = ['Switch','Backhaul','Power'];
const steps = ['Checking your nodes','Creating your site','Confirming your site'];
const readiness = ['Turn on your tower unit','Turn on your amplifier unit','Turn on your controller unit','Confirming your site location'];
const buttons = ['Get started','Next','Name site','Finish setup','Go to Console','Skip for now','Check now'];
const faults = ['transport','rejection','missing_amplifier','offline_tower','unready_controller','missing_location'];
const escape = (s:string)=>s.replace(/[.*+?^${}()|[\]\\]/g,'\\$&');
async function textOne(locator:Locator):Promise<string|null> {
 const visible=locator.filter({visible:true}), n=await visible.count();
 if(n>1)throw new WorkerError('AMBIGUOUS_LOCATOR','Onboarding UI field is ambiguous');
 return n===1 ? normalize(await visible.innerText({timeout:500})) : null;
}
export class Onboarding {
 private world:ObjectValue={};
 private jobs=new Map<string,{creation:Creation;command:number;done:boolean;intercepted:boolean}>();
 private counts={network:0,site:0};
 private fault=''; private consumed=false; private controlled=false;
 private faultError:unknown;
 constructor(private page:Page,private origin:string,private app:ConsoleApp,private directory:string) {
  page.on('request',r=>{try{const q=r.postDataJSON()?.query??'';if(/\baddNetwork\s*\(/.test(q))this.counts.network++;if(/\baddSite\s*\(/.test(q))this.counts.site++;}catch{}});
 }
 private cfg(){return this.page.locator('main.cfg-root .cfg-body');}
 private net(){return object(this.world.network,'planned network');}
 private site(){return object(this.world.site,'planned site');}
 private field(label:string){if(!fields.includes(label))throw new WorkerError('INVALID_INPUT','Unknown onboarding field');return this.cfg().locator('label.ff').filter({has:this.page.locator('.ff-label').filter({hasText:new RegExp(`^${escape(label)}\\s*\\*?$`)})});}
 private async scope(){if(new URL(this.page.url()).origin!==this.origin)throw new WorkerError('AUTH_REQUIRED','Onboarding left the configured console');}
 private async title(value:string,b:Budget){await this.cfg().getByRole('heading',{name:value,exact:true}).waitFor({timeout:b.remaining()});}
 private async go(path:string,b:Budget){const response=await this.page.goto(this.origin+path,{waitUntil:'domcontentloaded',timeout:b.remaining()});if(response&&response.status()>=400)throw new WorkerError('APP_UNAVAILABLE','Onboarding navigation failed',response.status());await this.scope();}
 private async identity(){await this.scope();const p=new URL(this.page.url()).searchParams;
  if(p.get('networkid')!==this.net().id)throw new WorkerError('WRONG_NETWORK','Wizard lost the planned network');
  if(new URL(this.page.url()).pathname==='/configure/site/settings'&&(p.get('nid')!==this.site().tower_id || await textOne(this.cfg().locator('.cfg-readonly'))!==this.site().tower_id || p.get('sitename')!==this.site().name))throw new WorkerError('WRONG_TOWER','Wizard lost the planned site/tower; no hidden repair is allowed');
 }
 // Invalid-name probes must never leak a resource if the product validator
 // regresses. Any attempted create fails the run and is blocked before forwarding.
 private validationGuard = async (route:Route) => {
  let query='';try{query=route.request().postDataJSON()?.query??''}catch{}
  if(/\b(addNetwork|addSite)\s*\(/.test(query)){
   this.faultError=new WorkerError('VALIDATION_MUTATION','Invalid-name validation attempted a backend create');
   await route.abort('blockedbyclient');
  }else await route.fallback();
 };
 private route = async (route:Route) => {
  try {
   let body:any;try{body=route.request().postDataJSON()}catch{return route.continue()}
   const query=body?.query??'';
   if(this.fault==='rejection'&&/\baddSite\s*\(/.test(query)) {
    const job=this.jobs.get('site');if(!job||body.variables?.data?.name!==this.site().name||body.variables?.data?.network_id!==this.net().id)throw new WorkerError('FAULT_SCOPE','Rejection fault must match its planned site');
    this.fault='';this.consumed=true;job.intercepted=true;job.creation.markIntercepted();
    await route.fulfill({status:200,contentType:'application/json',body:JSON.stringify({errors:[{message:'Lab controlled site rejection'}],data:{addSite:null}})});return;
   }
   if(this.fault==='transport'&&/\baddSite\s*\(/.test(query)) {
    if(body.variables?.data?.name!==this.site().name||body.variables?.data?.network_id!==this.net().id)throw new WorkerError('FAULT_SCOPE','Transport fault must match its planned site');
    this.fault='';this.consumed=true;
    // Forward this one browser request once, then hide its response. The app
    // must reconcile through its own polling; the adapter never retries it.
    const response=await route.fetch({maxRetries:0,maxRedirects:0});await response.dispose();await route.abort('timedout');return;
   }
   if(['missing_amplifier','offline_tower','unready_controller','missing_location'].includes(this.fault)&&/\bgetNodes\s*\(/.test(query)) {
    const response=await route.fetch({maxRetries:0,maxRedirects:0}), payload=await response.json();
    if(!Array.isArray(payload?.data?.getNodes?.nodes))throw new WorkerError('FAULT_SCOPE','Node query has no compatible visible-data response');
    const site=this.site();let found=false;
    payload.data.getNodes.nodes=payload.data.getNodes.nodes.filter((n:any)=>{if(this.fault==='missing_amplifier'&&n.id===site.amplifier_id){found=true;return false}return true}).map((n:any)=>{
     if(this.fault==='offline_tower'&&n.id===site.tower_id){n.status={...n.status,connectivity:'Offline'};found=true}
     if(this.fault==='unready_controller'&&n.id===site.controller_id){n.status={...n.status,state:'Initializing'};found=true}
     if(this.fault==='missing_location'&&n.id===site.tower_id){n.latitude='0';n.longitude='0';found=true}return n;
    });
    if(!found)throw new WorkerError('FAULT_SCOPE','Planned runtime node was not present for the controlled fault');
    this.consumed=true;await route.fulfill({response,json:payload});await response.dispose();return;
   }
   await route.continue();
  } catch(error){this.faultError=error;await route.abort().catch(()=>{});}
 };
 private async monitor(){
  await this.page.evaluate((titles)=>{
   const w=window as any;w.__ulabOnboardObserver?.disconnect();w.__ulabOnboard={order:[],done:[],invalid:false};
   const scan=()=>{
    const state=w.__ulabOnboard;
    const list=Array.from(document.querySelectorAll('main.cfg-root .cfg-step-item')).filter(e=>(e as HTMLElement).offsetParent!==null);
    for(const e of list){const title=e.querySelector('.cfg-step-title')?.textContent?.trim(),status=e.getAttribute('data-state');if(!titles.includes(title??''))continue;
     if(status==='active'&&!state.order.includes(title)){if(titles[state.order.length]!==title)state.invalid=true;state.order.push(title)}
     if(status==='done'&&!state.done.includes(title)){if(!state.order.includes(title))state.invalid=true;state.done.push(title)}
    }
   };w.__ulabOnboardObserver=new MutationObserver(scan);w.__ulabOnboardObserver.observe(document.body,{childList:true,subtree:true,attributes:true,characterData:true});scan();
  },steps);
 }
 async run(i:ObjectValue,command:number,b:Budget):Promise<ObjectValue> {
  keys(i,['view','action','label','value','context','creation']);if(i.view!=='configure')throw new WorkerError('INVALID_INPUT','Onboarding requires configure view');
  this.world=object(i.context,'context');keys(this.world,['network','site']);
  const action=str(i.action,'action'),label=str(i.label??'','label',true),value=str(i.value??'','value',true);
  if(this.faultError)throw this.faultError;
  const click=(name:string)=>this.cfg().getByRole('button',{name,exact:true}).click({timeout:b.remaining()});
  let pending=false,creation_command_id:number|undefined,bindings:ObjectValue[]=[];
  if(action==='open') {
   const routes:Record<string,string>={overview:'/configure',network:'/configure/network',add_network:'/configure/network?flow=add-network',select_network:'/configure/select-network?flow=install-site',install:'/configure/install?flow=install-site&networkid='+encodeURIComponent(String(this.world.network?this.net().id:'')),sims:'/configure/sims',complete:'/configure/complete'};
   if(!routes[value])throw new WorkerError('INVALID_INPUT','Unknown configure entry');
   if(['network','add_network'].includes(value)&&this.world.network&&!this.net().id){
    await this.go('/configure/select-network',b);
    await b.poll(()=>textOne(this.cfg().locator('h1')),t=>t==='Select a network'||t==='Name your network','Network preflight did not settle');
    if(await this.cfg().getByRole('radio').filter({has:this.page.getByText(String(this.net().name),{exact:true})}).count())throw new WorkerError('NAME_EXISTS','Planned network already exists');
   }
   await this.go(routes[value]!,b);await b.poll(()=>textOne(this.cfg().locator('h1')),t=>t!==null&&headings.includes(t),'Configure page did not load');
  } else if(action==='arm_fault') {
   if(!faults.includes(value)||this.fault||[...this.jobs.values()].some(j=>!j.done))throw new WorkerError('INVALID_INPUT','Arm one controlled fault before submission');
   this.site();this.fault=value;this.consumed=false;this.controlled=true;await this.page.unroute('**/*',this.route);await this.page.route('**/*',this.route);
  } else if(action==='clear_fault') {this.fault='';await this.page.unroute('**/*',this.route);}
  else {
   await this.scope();
   if(action==='fill'){await this.field(label).locator('input').fill(value,{timeout:b.remaining()});}
   else if(action==='validate_name'){
    const input=this.field(label).locator('input'), name=await input.inputValue();
    if(/^[a-z0-9-]{3,40}$/.test(name))throw new WorkerError('INVALID_INPUT','Validation probe accepts invalid names only; valid submission requires a creation intent');
    await this.page.unroute('**/*',this.validationGuard);await this.page.route('**/*',this.validationGuard);
    await click(label==='Network name'?'Create network':'Name site');
   } else if(action==='click'){
    if(value==='Installed'){await this.cfg().getByRole('checkbox',{name:"I've installed and powered on all my units",exact:true}).check({timeout:b.remaining()});}
    else {if(!buttons.includes(value))throw new WorkerError('INVALID_INPUT','Unsupported onboarding button');await click(value);}
   } else if(action==='choose_network') {
    await this.title('Select a network',b);const radio=this.cfg().getByRole('radio').filter({has:this.page.getByText(String(this.net().name),{exact:true})});await radio.click({timeout:b.remaining()});
    await b.poll(()=>radio.getAttribute('aria-checked'),v=>v==='true','Network radio did not select');await click('Continue');await this.title('Install your site',b);await this.identity();
   } else if(action==='select_component') {
    if(!components.includes(label))throw new WorkerError('INVALID_INPUT','Unknown component');await this.title('Configure site settings',b);await this.identity();await this.cfg().locator(`select[name="${label.toLowerCase()}Id"]`).selectOption({label:value},{timeout:b.remaining()});
   } else if(action==='submit_network'||action==='submit_site') {
    const kind=action==='submit_network'?'network':'site';if(this.jobs.has(kind))throw new WorkerError('DUPLICATE_MUTATION','Creation was already submitted');
    const entity=kind==='network'?this.net():this.site();const intent=object(i.creation,'creation');keys(intent,['kind','ref','name']);
    if(intent.kind!==kind||intent.ref!==entity.ref||intent.name!==entity.name||entity.id)throw new WorkerError('OWNERSHIP_CONFLICT','Creation intent differs from unbound world entity');
    const inputs=kind==='network'?{ref:entity.ref,name:entity.name}:{ref:entity.ref,name:entity.name,network_name:this.net().name,network_id:this.net().id,tower_id:entity.tower_id};
    // Write a prepared receipt before UI guards, so a pre-submit failure can
    // release the owned runtime without treating an untouched form as uncertain.
    const creation=new Creation(this.page,inputs,kind,this.directory,command);
    this.jobs.set(kind,{creation,command,done:false,intercepted:false});
    await this.title(kind==='network'?'Name your network':'Configure site settings',b);
    if(kind==='network' && await this.field('Network name').locator('input').inputValue()!==entity.name)throw new WorkerError('WRONG_NAME','Network field differs from planned name');
    if(kind==='site'){await this.identity();await this.monitor();}
    await this.page.unroute('**/*',this.validationGuard);
    creation.begin();creation.markSubmitted();pending=true;await click(kind==='network'?'Create network':'Create site');
   } else if(action==='retry_site') {
    const job=this.jobs.get('site');if(!job?.intercepted||job.done||!this.consumed||this.fault)throw new WorkerError('UNCERTAIN_CREATION','Retry requires a known rejection intercepted before any backend submission');
    await this.identity();if(await textOne(this.cfg().locator('.cfg-error'))!=='Lab controlled site rejection')throw new WorkerError('WRONG_STATE','Controlled backend rejection is not visible');
    job.intercepted=false;await this.monitor();job.creation.markSubmitted();await click('Try again');
   } else if(action==='finish_network'||action==='finish_site') {
    const kind=action==='finish_network'?'network':'site',job=this.jobs.get(kind);if(!job||job.done)throw new WorkerError('WRONG_STATE','No pending creation');
    await this.title(kind==='network'?'Install your site':'Upload SIMs',b);bindings=await job.creation.finish(b);creation_command_id=job.command;job.done=true;
    if(kind==='network'&&new URL(this.page.url()).searchParams.get('networkid')!==bindings[0]!.id)throw new WorkerError('WRONG_NETWORK','Created network differs from the wizard selection');
   } else if(action==='reload'){await this.page.reload({waitUntil:'domcontentloaded',timeout:b.remaining()});}
   else if(action==='site_detail') {
    const site=this.site();if(!site.id)throw new WorkerError('WRONG_STATE','Site has not completed');
    await this.app.open({view:'network_sites',network_name:this.net().name},b);const card=this.page.locator('main.main .ecard[role="button"]').filter({has:this.page.getByText(String(site.name),{exact:true})});
    await card.click({timeout:b.remaining()});await b.poll(()=>Promise.resolve(new URL(this.page.url()).pathname),v=>v==='/network/sites/'+site.id,'Site detail identity differs');
    await this.page.locator('main.main').getByRole('heading',{name:String(site.name),exact:true}).waitFor({timeout:b.remaining()});
   } else throw new WorkerError('INVALID_INPUT','Unknown onboarding action');
  }
  return {actual:{executed:true,...(pending?{pending:true}:{}),...(creation_command_id?{creation_command_id}:{}),controlled:this.controlled},bindings};
 }
 async check(i:ObjectValue,b:Budget):Promise<ObjectValue>{
  keys(i,['view','label','subject','expected','requirement','context']);if(i.view!=='configure')throw new WorkerError('INVALID_INPUT','Onboarding assertion requires configure view');this.world=object(i.context,'context');
  const label=str(i.label,'label'),subject=str(i.subject??'','subject',true),expected=str(i.expected,'expected',true);
  const actual=await b.poll(async()=>{
   if(this.faultError)throw this.faultError;await this.scope();
   if(label==='Mutation count'){if(!['network','site'].includes(subject))throw new WorkerError('INVALID_INPUT','Unknown mutation');return String(this.counts[subject as 'network'|'site']);}
   if(label==='Fault consumed')return String(this.consumed);
   if(label==='Progress order'||label==='Progress complete'){
    const progress=await this.page.evaluate(()=>(window as any).__ulabOnboard??null);
    if(!progress)return null;
    if(label==='Progress order')return progress.invalid?'invalid':progress.order.join(' > ');
    return String(!progress.invalid&&progress.done.join('|')===steps.join('|'));
   }
   const cfg=await this.cfg().isVisible(),dashboard=await this.page.locator('header.topbar').isVisible()&&await this.page.locator('main.main').isVisible();
   if(!cfg&&!dashboard)return null;
   if(label==='Heading')return textOne(cfg?this.cfg().locator('h1'):this.page.locator('main.main h1'));
   if(label==='Path')return new URL(this.page.url()).pathname;
   if(label==='Step')return cfg?textOne(this.page.locator('.cfg-step-label')):null;
   if(label==='Network ID'){if(!cfg)return null;return new URL(this.page.url()).searchParams.get('networkid');}
   if(label==='Selected network')return dashboard?textOne(this.page.locator('header.topbar .netswitch .nm')):null;
   if(label==='Field error')return cfg?textOne(this.field(subject).locator('.ff-err')):null;
   if(label==='Field value')return cfg&&await this.field(subject).locator('input').isVisible()?this.field(subject).locator('input').inputValue():null;
   if(label==='Button state'){
    if(!cfg)return null;if(![...buttons,'Create network','Create site','Try again'].includes(subject))throw new WorkerError('INVALID_INPUT','Unknown button');const button=this.cfg().getByRole('button',{name:subject,exact:true});
    if(await this.cfg().locator('.MuiSkeleton-root:visible').count())return null;return !(await button.isVisible())?'absent':await button.isEnabled()?'enabled':'disabled';
   }
   if(label==='Readiness'){if(!readiness.includes(subject))throw new WorkerError('INVALID_INPUT','Unknown readiness step');const row=this.cfg().locator('.cfg-step-item').filter({has:this.page.locator('.cfg-step-title').filter({hasText:new RegExp('^'+escape(subject)+'$')})});return await row.isVisible()?row.getAttribute('data-state'):null;}
   if(label==='Readiness count')return textOne(this.cfg().locator('.cfg-checklist-count'));
   if(label==='Tower')return new URL(this.page.url()).pathname==='/configure/site/settings'?textOne(this.cfg().locator('.cfg-readonly')):null;
   if(label==='Component'){if(!components.includes(subject))throw new WorkerError('INVALID_INPUT','Unknown component');const select=this.cfg().locator(`select[name="${subject.toLowerCase()}Id"]`);return await select.isVisible()?normalize(await select.locator('option:checked').textContent()??''):null;}
   if(label==='Error')return textOne(this.cfg().locator('.cfg-error'));
   if(label==='SIM guidance')return new URL(this.page.url()).pathname==='/configure/sims'?textOne(this.cfg().locator('.cfg-copy')):null;
   if(label==='SIM upload present')return new URL(this.page.url()).pathname==='/configure/sims'?String(await this.cfg().locator('input[type=file]').count()>0):null;
   if(['Site name','Site nodes','Coordinates','Saved component'].includes(label)){
    if(!dashboard||new URL(this.page.url()).pathname!=='/network/sites/'+this.site().id)return null;
    if(label==='Site name')return textOne(this.page.locator('main.main h1'));
    const main=this.page.locator('main.main');
    if(label==='Site nodes'){
     const ids=[this.site().tower_id,this.site().amplifier_id,this.site().controller_id];
     for(const id of ids)if(await main.locator('.site-top .tnum').filter({hasText:new RegExp('^'+escape(String(id))+'$')}).filter({visible:true}).count()!==1)return null;
     return ids.join(' | ');
    }
    if(label==='Coordinates')return textOne(main.locator('.site-top .tnum').filter({hasText:/^-?\d+(?:\.\d+)?,\s*-?\d+(?:\.\d+)?$/}));
    if(!['Switch','Backhaul','Charge controller'].includes(subject))throw new WorkerError('INVALID_INPUT','Unknown saved component');
    return textOne(main.locator('.comp-tile').filter({has:this.page.locator('.comp-tile-label').filter({hasText:new RegExp('^'+escape(subject)+'$')})}).locator('.comp-tile-sub'));
   }
   throw new WorkerError('UNSUPPORTED_LOCATOR','Unknown onboarding assertion');
  },v=>v!==null&&v===expected,'Onboarding UI differs from the declared expectation');
  return {expected,actual};
 }
}

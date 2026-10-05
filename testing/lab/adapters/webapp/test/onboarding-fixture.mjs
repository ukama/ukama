/* SPDX-License-Identifier: MPL-2.0
 * Controlled contract fixture with source-shaped DOM. This is not React/Apollo
 * console verification; modes deliberately model both expected and broken UI.
 */
import { provisioningFixture } from './provisioning-fixture.mjs';
function client(){
 let db,selected=sessionStorage.getItem('network')||'existing',siteName='',tower='',settings={};
 const main=document.querySelector('main'), body=()=>main.querySelector('.cfg-body'),params=()=>new URLSearchParams(location.search);
 const wait=ms=>new Promise(r=>setTimeout(r,ms));
 const btn=(name,fn,parent=body()||main)=>{const b=document.createElement('button');b.textContent=name;b.onclick=fn;parent.append(b);return b};
 const move=async path=>{history.pushState({},'',path);await render()};
 const gql=async(op,data={})=>{const r=await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({query:`mutation ${op}($data:Input!){${op}(data:$data){id name}}`,variables:{data}})});const json=await r.json();if(json.errors)throw Error(json.errors[0].message);return json.data[op]};
 const field=label=>{const el=document.createElement('label');el.className='ff';el.innerHTML=`<div class="ff-label">${label} *</div><input><div class="ff-err"></div>`;body().append(el);return el.querySelector('input')};
 const valid=input=>{const v=input.value,e=v.length<3?'At least 3 characters':v.length>40?'At most 40 characters':!/^[a-z0-9-]+$/.test(v)?'Lowercase letters, numbers, and "-" only (no spaces).':'';input.parentElement.querySelector('.ff-err').textContent=e;return !e};
 const title=s=>body().innerHTML=`<h1 class="cfg-title">${s}</h1>`;
 const route=path=>path+'?flow=install-site&networkid='+selected;
 const remember=()=>sessionStorage.setItem('network',selected);
 const progress=['Checking your nodes','Creating your site','Confirming your site'];
 const states=(active,done=false)=>{let ul=body().querySelector('.cfg-steps-list');if(!ul){ul=document.createElement('ul');ul.className='cfg-steps-list';body().append(ul)}ul.innerHTML=progress.map((t,i)=>`<li class="cfg-step-item" data-state="${done||i<active?'done':i===active?'active':'pending'}"><span class="cfg-step-title">${t}</span></li>`).join('')};
 async function submit(){
  const submit=body().querySelector('button');submit.disabled=true;body().querySelector('.cfg-error')?.remove();states(0);await wait(100);states(1);await wait(100);
  try{
   const data={name:siteName,network_id:selected,access_id:tower,...settings};
   try{await gql('addSite',data);if(db.mode==='duplicate-site')await gql('addSite',data)}catch(e){if(e.message==='Lab controlled site rejection')throw e;if(db.mode==='duplicate-site')await gql('addSite',data);const sites=await gql('getSites',{networkId:selected});if(!sites.sites.some(s=>s.name===siteName))throw e}
   if(db.mode!=='missing-confirm')states(2);await wait(100);states(3,true);await wait(100);await move('/configure/sims');
  }catch(e){const err=document.createElement('div');err.className='cfg-error';err.textContent=e.message;body().append(err);submit.textContent='Try again';submit.disabled=false}
 }
 async function readiness(){
  const result=await gql('getNodes',{}),nodes=result.nodes;
  const get=t=>nodes.find(n=>n.type===t),ready=n=>n?.status?.connectivity==='Online'&&n.status.state==='Ready';
  const t=get('Tnode'),a=get('Anode'),c=get('Cnode');const flags=[ready(t),ready(a),ready(c),!!t&&Number(t.latitude)!==0&&Number(t.longitude)!==0];
  tower=db.mode==='wrong-tower'?'other-tower':t?.id;const ok=flags.every(Boolean);
  title(ok||db.mode==='premature'?'Name your site':'Bring your site online');
  if(ok||db.mode==='premature'){
   const input=field('Site name');btn('Name site',()=>{if(valid(input)){siteName=input.value;void move(route('/configure/site/settings')+'&nid='+tower+'&sitename='+siteName)}});
  }else{
   const names=['Turn on your tower unit','Turn on your amplifier unit','Turn on your controller unit','Confirming your site location'];let active=false;
   body().insertAdjacentHTML('beforeend',`<div class="cfg-checklist-count">${flags.filter(Boolean).length}/4 done</div><ul class="cfg-steps-list">${names.map((n,i)=>{const state=flags[i]?'done':!active?'active':'pending';if(!flags[i])active=true;return `<li class="cfg-step-item" data-state="${state}"><span class="cfg-step-title">${n}</span></li>`}).join('')}</ul>`);btn('Check now',readiness);
  }
 }
 async function render(){
  db=await(await fetch('/state')).json();const path=location.pathname,cfg=path.startsWith('/configure');
  document.querySelector('header').hidden=cfg;document.querySelector('aside').hidden=cfg;main.className=cfg?'cfg-root':'main';main.innerHTML=cfg?'<div class="cfg-step-label"></div><div class="cfg-body"></div>':'';
  document.querySelector('.nm').textContent=db.networks.find(n=>n.id===selected)?.name||'';
  document.querySelector('aside').innerHTML='<a href="/network">Home</a><a href="/network/sites">Sites</a>';
  document.querySelectorAll('aside a').forEach(a=>a.onclick=e=>{e.preventDefault();void move(a.getAttribute('href'))});
  if(cfg){let step=path.includes('/network')?1:path.includes('/install')?2:path.endsWith('/site')?3:path.endsWith('/settings')?4:5;document.querySelector('.cfg-step-label').textContent=`STEP ${step}/5`}
  if(path==='/configure'){title("Let's set up your network");btn('Get started',()=>move('/configure/network'))}
  else if(path==='/configure/network'){
   title('Name your network');const input=field('Network name');btn('Create network',async()=>{if(!valid(input)&&db.mode!=='broken-validation')return;const n=await gql('addNetwork',{name:input.value});selected=n.id;remember();await move(route('/configure/install'))});
  }else if(path==='/configure/select-network'){
   title('Select a network');for(const n of db.networks){const b=btn(n.name,()=>{selected=n.id;remember();for(const r of body().querySelectorAll('[role=radio]'))r.setAttribute('aria-checked',String(r===b))});b.role='radio';b.setAttribute('aria-checked','false')}
   btn('Continue',()=>move(route('/configure/install')));
  }else if(path==='/configure/install'){
   selected=params().get('networkid');title('Install your site');body().insertAdjacentHTML('beforeend',`<label><input type="checkbox">I've installed and powered on all my units</label>`);const b=btn('Next',()=>move(route('/configure/site')));b.disabled=true;body().querySelector('input').onchange=e=>b.disabled=!e.target.checked;
  }else if(path==='/configure/site'){selected=params().get('networkid');await readiness()}
  else if(path==='/configure/site/settings'){
   selected=params().get('networkid');tower=params().get('nid');siteName=params().get('sitename');title('Configure site settings');body().insertAdjacentHTML('beforeend',`<div class="cfg-readonly">${tower}</div>`);
   for(const k of ['switch','backhaul','power']){const s=document.createElement('select');s.name=k+'Id';s.innerHTML=`<option value="${k}-a">Default ${k}</option><option value="${k}-b">Test ${k}</option>`;settings[k+'_id']=k+'-a';s.onchange=()=>settings[k+'_id']=s.value;body().append(s)}btn('Create site',submit);
  }else if(path==='/configure/sims'){title('Upload SIMs');body().insertAdjacentHTML('beforeend','<p class="cfg-copy">Upload your SIMs later from Console Manage → SIM pool.</p>');btn('Finish setup',()=>move('/configure/complete'))}
  else if(path==='/configure/complete'){title("You're all set!");btn('Go to Console',()=>move('/network'))}
  else if(path==='/network/sites'){
   for(const s of db.sites.filter(s=>s.network_id===selected)){const card=document.createElement('div');card.className='ecard';card.role='button';card.textContent=s.name;card.onclick=()=>move('/network/sites/'+s.id);main.append(card)}
  }else if(path.startsWith('/network/sites/')){
   const s=db.sites.find(s=>path.endsWith('/'+s.id));main.innerHTML=`<h1>${s.name}</h1><div class="site-top">${db.nodes.map(n=>`<span class="tnum">${n.id}</span>`).join('')}<span class="tnum">${db.mode==='wrong-coordinates'?'0, 0':'-1.67, 29.22'}</span></div>`;
   for(const[k,label]of[['switch','Switch'],['backhaul','Backhaul'],['power','Charge controller']]){main.insertAdjacentHTML('beforeend',`<div class="comp-tile"><span class="comp-tile-label">${label}</span><span class="comp-tile-sub">${db.mode==='wrong-component'?'Default':db.saved[k+'_id'].endsWith('-b')?'Test':'Default'} ${k}</span></div>`)}
  }else{main.innerHTML='<h1>Network</h1>'}
 }
 document.querySelector('.netswitch').onclick=()=>{const menu=document.querySelector('#menu');menu.innerHTML='';for(const n of db.networks){const b=btn(n.name,()=>{selected=n.id;remember();menu.innerHTML='';void render()},menu);b.role='menuitem'}};
 document.querySelectorAll('.viewseg button').forEach(b=>b.onclick=()=>move('/'+b.textContent.toLowerCase()));
 if(location.pathname==='/')history.replaceState({},'', '/network');void render();
}
export async function onboardingFixture(mode=''){
 const app=await provisioningFixture(mode,{
  get(req,res){if(req.method!=='GET'||req.url==='/state')return false;res.setHeader('content-type','text/html');res.end(`<!doctype html><style>button,a{margin:8px}.tnum,.comp-tile,.cfg-step-item{display:block;padding:4px}h1{font-size:24px}[hidden]{display:none!important}</style><header class="topbar"><button class="netswitch"><span class="nm"></span></button><div class="viewseg"><button>Network</button><button>Business</button></div></header><aside class="sidebar"></aside><main></main><div id="menu"></div><script>(${client.toString()})()</script>`);return true},
  graphql(body,db){const q=body.query||'';if(q.includes('getNodes('))return{data:{getNodes:{nodes:db.nodes.map(n=>({...n,type:{'Tower node':'Tnode','Amplifier node':'Anode','Controller node':'Cnode'}[n.type],status:{connectivity:'Online',state:'Ready'},latitude:'-1.67',longitude:'29.22'}))}}};if(q.includes('getSites('))return{data:{getSites:{sites:db.sites}}};if(q.includes('addSite('))db.saved=body.variables.data;return null}
 });
 app.prepare=async()=>({ULAB_ONBOARD_SWITCH:'Test switch',ULAB_ONBOARD_BACKHAUL:'Test backhaul',ULAB_ONBOARD_POWER:'Test power',ULAB_ONBOARD_COORDINATES:'-1.67, 29.22'});return app;
}

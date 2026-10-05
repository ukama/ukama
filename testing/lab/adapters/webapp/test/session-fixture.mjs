/* SPDX-License-Identifier: MPL-2.0
 * Two-origin controlled auth/DOM fixture. Synthetic cookies here are test-only;
 * production worker scenarios always use manually captured real account state.
 * This validates lab behavior, not the uploaded React/Apollo/auth implementation.
 */
import http from 'node:http';
import { once } from 'node:events';
import { writeFile } from 'node:fs/promises';
import { join } from 'node:path';
const roles={OWNER:'Owner',ADMIN:'Admin',NETWORK_OWNER:'Network owner',VENDOR:'Vendor',MEMBER:'Member',WELCOME:'Member',LOGOUT:'Owner',EXPIRED:'Owner',NO_ORG:'Member',NO_ROLE:'Member',POLICY:'Member'};
const esc=s=>String(s).replaceAll('&','&amp;').replaceAll('<','&lt;').replaceAll('"','&quot;');
const cookie=(name,value)=>({name,value,domain:'127.0.0.1',path:'/',expires:-1,httpOnly:true,secure:false,sameSite:'Lax'});
export async function sessionFixture(mode='') {
 const org='Lab Test Organization',country='United States',currency='USD ($)',accounts={};
 for(const [id,role] of Object.entries(roles)){const person=['LOGOUT','EXPIRED'].includes(id)?'OWNER':id;accounts[id]={id,name:person.toLowerCase()+' test',email:person.toLowerCase()+'@example.test',role,org,country,currency,welcome:id==='WELCOME',valid:id!=='EXPIRED'};}
 const db={refresh:0,mints:0,logout:0,welcome:0,documents:0,operations:[]};
 const token=a=>Buffer.from([a.id,a.org,a.id,a.name,a.email,a.role,'true',String(a.welcome),'US','USD',Math.floor(Date.now()/1000)+3600].join(';')).toString('base64')+'.fixture-signature';
 const parsed=req=>Object.fromEntries((req.headers.cookie||'').split(';').filter(x=>x.includes('=')).map(x=>{const n=x.indexOf('=');return[x.slice(0,n).trim(),x.slice(n+1)]}));
 const redir=(res,to,cookies)=>{res.writeHead(302,{Location:to,...(cookies?{'Set-Cookie':cookies}:{})});res.end()};
 const clear='token=; Path=/; Max-Age=0; HttpOnly; SameSite=Lax';
 let origin,authOrigin;
 const html=(res,body,script='',status=200)=>{res.writeHead(status,{'Content-Type':'text/html; charset=utf-8','Cache-Control':'no-store'});res.end(`<!doctype html><html><head><title>Controlled lab auth fixture</title></head><body>${body}<script>${script}</script></body></html>`)};
 const serverAuth=http.createServer((req,res)=>{
  if(req.url==='/user/logout'){db.logout++;const id=parsed(req).ukama_session;if(accounts[id]&&mode!=='logout-leak')accounts[id].valid=false;redir(res,authOrigin+'/auth/login',mode==='logout-leak'?undefined:['ukama_session=; Path=/; Max-Age=0; HttpOnly; SameSite=Lax',clear]);return}
  html(res,'<h1>Sign in</h1><p>Controlled authentication service</p>','',mode==='auth-error'?404:200);
 });serverAuth.listen(0,'127.0.0.1');await once(serverAuth,'listening');authOrigin=`http://127.0.0.1:${serverAuth.address().port}`;
 const server=http.createServer((req,res)=>{
  const path=new URL(req.url,'http://fixture').pathname,cookies=parsed(req),account=accounts[cookies.ukama_session];db.documents++;
  if(path==='/api/auth/logout'){res.setHeader('Set-Cookie',mode==='logout-leak'?[]:[clear]);res.end('{}');return}
  if(path==='/api/auth/refresh'){db.refresh++;redir(res,mode==='redirect-loop'?origin+'/api/auth/refresh':origin+'/',[clear]);return}
  if(path==='/api/auth/welcome'){db.welcome++;if(mode==='welcome-error'){res.writeHead(502);res.end('{}');return}if(account&&mode!=='welcome-not-persisted')account.welcome=false;res.setHeader('Set-Cookie',clear);res.end('{}');return}
  if(path==='/graphql'){const valid=account?.valid&&cookies.token?.endsWith('.fixture-signature');res.setHeader('Content-Type','application/json');res.end(JSON.stringify(valid?{data:{}}:{errors:[{extensions:{code:'UNAUTHENTICATED'}}]}));return}
  if(path==='/unauthorized'){
   html(res,`<main><div>Your account isn't set up for this console</div><a href="${authOrigin}/user/logout">Log out</a><a href="mailto:support@ukama.com?subject=Access">Contact us</a></main>`+(mode==='unauthorized-leak'?'<header class="topbar">owner test</header><main class="main">OWNER</main><aside class="sidebar">Lab Test Organization</aside>':''));return;
  }
  if(!cookies.ukama_session){redir(res,authOrigin+'/auth/login',[clear]);return}
  let claims;try{const fields=Buffer.from((cookies.token||'').split('.')[0],'base64').toString().split(';');if(fields.length===11&&Number(fields[10])>Date.now()/1000)claims={...accounts[fields[0]],welcome:fields[7]==='true'};}catch{}
  if(!claims){
   if(!account?.valid||['NO_ORG','NO_ROLE'].includes(account.id)){redir(res,origin+'/unauthorized',[clear]);return}
   if(mode==='recovery-broken'){redir(res,authOrigin+'/auth/login',[clear]);return}
   db.mints++;claims={...account};res.setHeader('Set-Cookie',`token=${token(account)}; Path=/; HttpOnly; SameSite=Lax`);
  }
  if(claims.welcome&&path!=='/welcome'){redir(res,origin+'/welcome');return}
  if(!claims.welcome&&path==='/welcome'){redir(res,origin+'/');return}
  if(path==='/'){redir(res,origin+'/business');return}
  if(path==='/welcome'){
   html(res,`<main class="welcome-root"><h1 class="welcome-title">Welcome to Ukama!</h1>${[['Network operating country',country],['Organization name',org],['Role',claims.name+' | '+claims.role]].map(([k,v])=>`<div><div class="welcome-field-label">${k}</div><div class="welcome-field-value">${esc(v)}</div></div>`).join('')}<button>Continue</button><p class="welcome-error" hidden></p></main>`, `document.querySelector('button').onclick=async()=>{const r=await fetch('/api/auth/welcome',{method:'POST'});if(r.ok)location.assign('/');else{const e=document.querySelector('.welcome-error');e.hidden=false;e.textContent="Couldn't save your confirmation. Please try again."}}`);return;
  }
  if(path==='/business/manage/billing'){html(res,'<h1>404</h1><h2>This page could not be found.</h2>','',mode==='wrong-404'?200:404);return}
  const profile={...claims};if(mode==='wrong-identity')profile.name='another account';if(mode==='wrong-org')profile.org='another organization';if(mode==='wrong-role')profile.role='Owner';
  const routeControls={'/business/manage/members':'Invite member','/business/manage/data-plans':'Create plan','/business/manage/sim-pool':'Upload SIMs','/customer/customers':'Add customer'};
  const action=routeControls[path],button=action&&!(account.id==='POLICY'&&mode!=='policy-leak')?`<button>${action}</button>`:'';
  const body=`<header class="topbar"><div class="viewseg"><button>Business</button><button>Network</button><button>Customer</button></div><span title="Organization: ${esc(profile.org)}">${esc(profile.org)}</span><button class="avatar" aria-haspopup="menu">Account</button></header><aside class="sidebar"><a href="/business/settings">Settings</a>${mode==='billing-link'?'<a href="/business/manage/billing">Billing</a>':''}</aside><main class="main"><div class="pagehead"><h1>${path.endsWith('/settings')?'Settings':'Members'}</h1>${button}</div>${path.endsWith('/settings')?'<div role="tablist"><button role="tab">My account</button><button role="tab">Organization</button><button role="tab">Preferences</button></div><div id="fields"></div>':''}</main>`;
  const script=`const profile=${JSON.stringify(profile)},auth=${JSON.stringify(authOrigin)};
const fields=document.getElementById('fields');const render=(tab)=>{if(!fields)return;const pairs=tab==='Organization'?[['Organization name',profile.org],['Country',profile.country],['Currency',profile.currency]]:tab==='Preferences'?[]:[['Full name',profile.name],['Email',profile.email],['Role',profile.role],['Email verified','Yes']];fields.replaceChildren(...pairs.map(([label,value])=>{const d=document.createElement('div');d.className='card card-pad';const l=document.createElement('label');l.className='flabel';l.textContent=label;const v=document.createElement('div');v.className='ff-readonly';v.textContent=value;d.append(l,v);return d}))};render('My account');document.querySelectorAll('[role=tab]').forEach(b=>b.onclick=()=>render(b.textContent));
document.querySelector('.avatar').onclick=()=>{document.querySelector('[role=menu]')?.remove();const menu=document.createElement('div');menu.setAttribute('role','menu');for(const value of [profile.name,profile.email+' · '+profile.role]){const p=document.createElement('p');p.className='MuiTypography-root';p.textContent=value;menu.append(p)}const out=document.createElement('button');out.setAttribute('role','menuitem');out.textContent='Log out';out.onclick=async()=>{await fetch('/api/auth/logout',{method:'POST'});localStorage.removeItem('uk-ui-prefs');sessionStorage.clear();location.assign(auth+'/user/logout')};menu.append(out);document.body.append(menu)};
fetch('/graphql',{method:'POST'}).then(r=>r.json()).then(x=>{if(x.errors?.some(e=>e.extensions.code==='UNAUTHENTICATED')){${mode==='stale-access'?'':'location.assign("/api/auth/refresh")'}}});`;
  html(res,body,script);
 });server.listen(0,'127.0.0.1');await once(server,'listening');origin=`http://127.0.0.1:${server.address().port}`;
 const states=Object.fromEntries(Object.entries(accounts).map(([id,a])=>[id,{cookies:[cookie('ukama_session',id),...(['NO_ORG','NO_ROLE'].includes(id)?[]:[cookie('token',token(a))])],origins:[]}]));
 return {origin,authOrigin,state:states.OWNER,states,db,accounts,prepare:async dir=>{
  const env={ULAB_WEBAPP_AUTH_ORIGIN:authOrigin,ULAB_AUTH_ORG:org,ULAB_AUTH_COUNTRY:country,ULAB_AUTH_CURRENCY:currency};
  for(const [id,a] of Object.entries(accounts)){const path=join(dir,id.toLowerCase()+'.json');await writeFile(path,JSON.stringify(states[id]),{mode:0o600});env['ULAB_AUTH_'+id+'_STATE']=path;env['ULAB_AUTH_'+id+'_NAME']=a.name;env['ULAB_AUTH_'+id+'_EMAIL']=a.email;}
  for(const key of ['MEMBERS','PLANS','SIMS','CUSTOMERS']){env['ULAB_POLICY_'+key+'_SURFACE']='dashboard';env['ULAB_POLICY_'+key+'_CONTROL']='absent';}return env;
 },close:async()=>{server.closeAllConnections();serverAuth.closeAllConnections();await Promise.all([new Promise(r=>server.close(r)),new Promise(r=>serverAuth.close(r))])}};
}

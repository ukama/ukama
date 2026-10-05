/* SPDX-License-Identifier: MPL-2.0 — controlled DOM, not a live console. */
import {commerceFixture} from './commerce-fixture.mjs';
function install(){
 const prior=window.renderOperations;
 window.renderOperations=async args=>{
  await prior(args);const {db,main,menu,button,render,selected}=args,path=location.pathname;
  const move=p=>{history.pushState({},'',p);void render()};const lens=path.split('/')[1];
  for(const [name,p] of [['Settings',`/${lens}/settings`],['Support',`/${lens}/support`]]){const a=document.createElement('a');a.href=p;a.textContent=name;a.onclick=e=>{e.preventDefault();move(p)};document.querySelector('aside').append(a)}
  const query=async(name,variables={})=>fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:name,query:`query ${name}${name==='TeamList'?'': '($networkId: String!)'} { fixture }`,variables})}).then(r=>r.json());
  const field=(d,label,type='input',options='')=>{const f=document.createElement('label');f.className='ff';f.innerHTML=`<div class="ff-label">${label}</div><${type}${type==='input'?' type="email"':''}>${options}</${type}><span class="ff-err"></span>`;d.append(f);return f.querySelector(type)};
  if(path==='/business/manage/members'){
   main.innerHTML='<div class="MuiSkeleton-root">Loading</div>';const members=await query('TeamList');
   if(members.data?.membersView?.team?.error){main.innerHTML="<div>Couldn't load members</div>";button('Try again',()=>{if(db.mode!=='broken-retry')void render()},main);return;}
   main.innerHTML='<div class="pagehead"></div><table><thead><tr><th>Member</th><th>Role</th><th>Member since</th><th>Status</th><th></th></tr></thead><tbody></tbody></table>';
   let rows=db.members;if(db.mode==='duplicate-member')rows=[...rows,rows[0]];
   for(const m of rows)main.querySelector('tbody').innerHTML+=`<tr><td><div><span class="av-sm">L</span><div><div>${m.name}</div><div class="muted">${m.email}</div></div></div></td><td><div>${db.mode==='wrong-role'?'Vendor':m.role}</div><div>Role description</div></td><td>Oct 1, 2026</td><td><span>${m.status}</span></td><td><button>More actions</button></td></tr>`;
   if(db.mode==='page-error')queueMicrotask(()=>{throw Error('controlled browser error')});
   button('Invite member',()=>{
    const opener=document.activeElement,d=document.createElement('div');d.role='dialog';d.innerHTML='<h2>Invite member</h2>';menu.append(d);d.tabIndex=-1;d.style.cssText='position:fixed;inset:10% auto auto 3%;width:90vw;max-width:520px;background:white;padding:8px;box-sizing:border-box';if(db.mode==='wide-dialog')d.style.width='800px';if(db.mode==='wide-dialog')d.style.maxWidth='none';d.focus();
    d.addEventListener('keydown',e=>{if(e.key==='Escape'){d.remove();if(db.mode!=='lost-focus')opener.focus()}if(e.key==='Tab'&&db.mode!=='untrapped'){const list=[...d.querySelectorAll('input,select,button')].filter(e=>!e.disabled);let i=list.indexOf(document.activeElement);if(e.shiftKey&&i<=0){e.preventDefault();list.at(-1).focus()}else if(!e.shiftKey&&(i<0||i===list.length-1)){e.preventDefault();list[0].focus()}}});
    const email=field(d,'Email'),role=field(d,'Role','select','<option value="">Select a role</option>'+['Owner','Administrator','Network owner','Vendor'].map(r=>`<option>${r}</option>`).join(''));
    const submit=button('Invite member',async()=>{const map={Owner:'ROLE_OWNER',Administrator:'ROLE_ADMIN','Network owner':'ROLE_NETWORK_OWNER',Vendor:'ROLE_VENDOR'};
     const data={email:email.value.toLowerCase(),name:email.value.split('@')[0],role:db.mode==='wrong-invite-role'?'ROLE_OWNER':map[role.value]};
     const response=await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'CreateInvitation',query:'mutation CreateInvitation($data: CreateInvitationInputDto!) { createInvitation(data: $data) { id email } }',variables:{data}})}).then(r=>r.json());
     if(response.errors){const error=document.createElement('div');error.textContent=db.mode==='hidden-invite-error'?'':response.errors[0].message;menu.append(error)}
    },d);submit.disabled=true;
    const validate=()=>{const valid=/^[^@ ]+@[^@ ]+\.[^@ ]+$/.test(email.value);email.parentElement.querySelector('.ff-err').textContent=valid?'':'Enter a valid email';submit.disabled=db.mode==='invalid-invite-enabled'?false:!(valid&&role.value)};email.oninput=validate;role.onchange=validate;
    button('Cancel',()=>d.remove(),d);
   },main.querySelector('.pagehead'));
  }
  if(['business','network','customer'].includes(lens)&&path.endsWith('/settings')) {
   main.innerHTML='<div role="tablist"></div><div id="settings"></div>';
   for(const name of ['My account','Organization','Preferences'])button(name,()=>{
    const panel=main.querySelector('#settings');panel.innerHTML='';if(name!=='Preferences')return;
    const group=document.createElement('div');group.role='radiogroup';group.setAttribute('aria-label','Map view');panel.append(group);
    let choice=db.mode==='lost-preference'?'Terrain':localStorage.getItem('map-view')||'Terrain';
    for(const name of ['Street','Satellite','Terrain']){const radio=document.createElement('div');radio.role='radio';radio.tabIndex=0;radio.setAttribute('aria-checked',String(choice===name));radio.innerHTML=`<div>${name}</div><div>Map preview</div>`;radio.onclick=()=>{group.querySelectorAll('[role="radio"]').forEach(r=>r.setAttribute('aria-checked',String(r===radio)));localStorage.setItem('map-view',name)};group.append(radio)}
   },main.querySelector('[role="tablist"]')).role='tab';
  }
  if(path.endsWith('/support')) {
   const network=lens==='network',names=network?['SitesList','NodesList']:['NetworkCustomers'];
   main.innerHTML='<div>Loading…</div>';const results=await Promise.all(names.map(name=>query(name,{networkId:selected})));
   if(results.some(r=>r.errors)&&db.mode!=='hidden-support-error'){main.innerHTML='<div>Couldn\'t load support</div>';return;}
   const source=results.some(r=>r.errors)?[]:network?[...db.sites.filter(s=>s.network_id===selected).map(s=>({kind:'site',id:s.id,name:s.name,fields:[['Site ID',s.id],['Customers','0'],['Nodes','3']]})),...db.nodes.filter(n=>db.sites.some(s=>s.id===n.site&&s.network_id===selected)).map(n=>({kind:'node',id:n.id,name:db.mode==='named-node'?'Custom '+n.id:n.id,fields:[['Node ID',n.id],['Site',db.sites.find(s=>s.id===n.site)?.name]]}))]:db.subscribers.filter(s=>s.network_id===selected).map(s=>{const sim=db.sims.find(x=>x.subscriber_id===s.uuid),plan=db.plans.find(p=>p.uuid===sim?.package_id);return {kind:'customer',id:s.uuid,name:s.name,sim,email:s.email,plan:plan?.name||'—',fields:[['Package',plan?.name||'—'],['Last seen','Never']]}});
   main.innerHTML=`<input placeholder="${network?'Search site or node':'Search customer by name or phone'}"><div class="tile-grid"><div class="card card-pad"><div class="sec-title">${network?'Sites & nodes':'Customers'}</div><div id="results"></div></div><div id="detail"></div></div>`;
   const input=main.querySelector('input'),search=button('Search',()=>filter(),main);input.oninput=filter;
   function filter(){const rows=source.filter(r=>r.name.toLowerCase().includes(input.value.toLowerCase()));const list=main.querySelector('#results');list.innerHTML='';search.disabled=rows.length===0;
    for(const r of rows)button(r.name,()=>detail(r),list);
    if(rows.length)detail(rows[0]);else{list.textContent=`No matches for “${input.value}”.`;main.querySelector('#detail').className='';main.querySelector('#detail').innerHTML='<div>No match</div>'}
   }
   function detail(r){const d=main.querySelector('#detail');d.className='card card-pad';d.innerHTML=`<div><span>${db.mode==='wrong-support-identity'?'Foreign customer':r.name}</span><span>Active</span></div><div class="tile-grid">${r.fields.map(([k,v])=>`<div><div>${k}</div><div class="tnum">${db.mode==='wrong-support-id'&&k.endsWith('ID')?'foreign-id':v}</div></div>`).join('')}</div>`;
    button('Copy summary',async()=>{if(db.mode==='stale-summary')return;await navigator.clipboard.writeText([`Org: Fixture Organization`,`Network ID: ${selected}`,`Timestamp: ${new Date().toISOString()}`,`${r.kind==='customer'?'Customer':r.kind==='site'?'Site':'Node'}: ${r.name}`,`ICCID: ${db.mode==='wrong-iccid'?'foreign':r.sim?.iccid||'—'}`,`Plan: ${r.plan||'—'}`].join('\n'))},d);
    if(network){const restart=button('Restart '+r.kind,async()=>{
     if(db.mode==='unguarded-restart'){await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'RestartNode',query:'mutation RestartNode($data: Input!) { restartNode(data: $data) { success } }',variables:{data:{nodeId:r.id}}})});return;}
     const modal=document.createElement('div');modal.role='dialog';modal.innerHTML=`<h2>Restart ${r.kind}</h2>`;if(db.mode==='late-restart')setTimeout(()=>fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'RestartNode',query:'mutation RestartNode($data: Input!) { restartNode(data: $data) { success } }',variables:{data:{nodeId:r.id}}})}).catch(()=>{}),200);button('Cancel',()=>modal.remove(),modal);button('Restart',()=>{},modal);menu.append(modal);
    },d);restart.disabled=db.mode==='locked-restart';}
   }
   filter();
  }
 };
}
export async function teamSupportFixture(mode='') {
 const app=await commerceFixture(mode,{script:`(${install.toString()})();`,init(db){db.members=[{name:'Lab owner',email:'owner@example.test',role:'Owner',status:'Active'},{name:'Lab operator',email:'operator@example.test',role:'Admin',status:'Active'},{name:'Lab pending',email:'pending@example.test',role:'Vendor',status:'Invited'}];db.interceptedWrites=[]},graphql(body,db){
  if(['TeamList','SitesList','NodesList','NetworkCustomers'].includes(body.operationName))return {data:{fixture:true}};
  if(['CreateInvitation','RestartNode','RestartSite'].includes(body.operationName)){db.interceptedWrites.push(body);return {data:{unexpected:true}};}
 }});
 app.prepare=async()=>({ULAB_TEAM_PENDING_NAME:'Lab pending',ULAB_TEAM_PENDING_EMAIL:'pending@example.test',ULAB_TEAM_PENDING_ROLE:'Vendor',ULAB_TEAM_MEMBER_EMAIL:'owner@example.test',ULAB_TEAM_MEMBER_ROW:'Lab owner|owner@example.test|Owner|Active',ULAB_TEAM_SECOND_EMAIL:'operator@example.test',ULAB_TEAM_SECOND_ROW:'Lab operator|operator@example.test|Admin|Active'});return app;
}

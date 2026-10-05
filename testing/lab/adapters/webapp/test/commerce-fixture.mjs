/* SPDX-License-Identifier: MPL-2.0
 * Source-shaped controlled DOM fixture. Never counts as target-app coverage.
 */
import { provisioningFixture } from './provisioning-fixture.mjs';
function client() {
  const terms = minutes => ({1440:'1 day',10080:'1 week',43200:'1 month'})[minutes] || `${minutes} minutes`;
  const unit = minutes => ({1440:'day',10080:'week',43200:'month'})[minutes] || 'period';
  const bytes = n => n >= 1073741824 ? `${Math.round(n / 1073741824 * 100)/100} GB` : n >= 1048576 ? `${n / 1048576} MB` : `${n} B`;
  window.renderOperations = ({db,main,menu,button,graphql,render}) => {
    const path = location.pathname, selected = sessionStorage.getItem('network') || 'existing';
    const move = p => { history.pushState({},'',p); void render(); };
    document.querySelectorAll('.viewseg button').forEach(b=>b.onclick=()=>move(b.textContent==='Customer'?'/customer/customers':'/'+b.textContent.toLowerCase()));
    const lens=path.split('/')[1];
    const nav = document.querySelector('aside');
    const home=[...nav.querySelectorAll('a')].find(a=>a.textContent==='Home');if(home){home.href='/'+lens;home.onclick=e=>{e.preventDefault();move('/'+lens)}}
    for (const [n,p] of [['Data plans',`/${lens}/manage/data-plans`],['Customers',`/${lens}/customers`],['SIM pool',`/${lens}/manage/sim-pool`],['Revenue',`/${lens}/revenue`],['Packages',`/${lens}/packages`]]) {
      const a=document.createElement('a');a.href=p;a.textContent=n;a.onclick=e=>{e.preventDefault();move(p)};nav.append(a);
    }
    const modal = title => { const d=document.createElement('div');d.role='dialog';d.innerHTML=`<h2>${title}</h2>`;menu.append(d);return d; };
    const field = (d,label,html) => {const f=document.createElement('label');f.className='ff';f.innerHTML=`<div class="ff-label">${label}</div>${html}`;d.append(f);return f.querySelector('input,select')};
    const addCancel=d=>button('Cancel',()=>d.remove(),d);
    const plans = () => db.plans.filter(p=>!p.networkId||p.networkId===selected);
    const planOptions=()=>'<option value="">Select a plan</option>'+plans().map(p=>`<option value="${p.uuid}">${p.name} · $${p.amount}/${unit(p.duration)}</option>`).join('');
    const refresh=async()=>{Object.assign(db,await (await fetch('/state')).json());await render();};
    function drawer(sub) {
      document.querySelectorAll('.MuiDrawer-paper').forEach(e=>e.remove());
      const d=document.createElement('div');d.className='MuiDrawer-paper';document.body.append(d);d.innerHTML=`<div>${sub.name}</div>`;
      const sim=db.sims.find(s=>s.subscriber_id===sub.uuid), entries=db.entitlements.filter(p=>p.sim===sim?.id);
      const active=entries.find(e=>e.active),plan=db.plans.find(p=>p.uuid===active?.packageId);
      const cap=plan?plan.dataVolume*(plan.dataUnit==='GB'?1073741824:1048576):0;
      d.innerHTML+=`<div class="card card-pad"><div><span>${plan?.name||'No plan'}</span></div><div class="tnum">${bytes(db.usage)}${cap?' of '+bytes(cap)+' used this cycle':' used'}</div></div>`;
      for(const [label,value] of [['ICCID',sim?.iccid||'—'],['SIM status',sim?sim.status[0].toUpperCase()+sim.status.slice(1):'Inactive'],['Phone',''],['Total usage',bytes(db.usage)]]) d.innerHTML+=`<div><span>${label}</span><span class="tnum">${value}</span></div>`;
      d.innerHTML+='<div>Packages</div>';button('Close',()=>d.remove(),d);
      for(const e of entries) {
        const p=db.plans.find(p=>p.uuid===e.packageId),card=document.createElement('div');card.className='card card-pad';card.innerHTML=`<div><div>${p.name}</div><div class="tnum">Oct 1, 2026 – ${new Date(Date.UTC(2026,9,1)+p.duration*60000).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric',timeZone:'UTC'})}</div></div><div><span class="MuiChip-label">${e.active?'Current':'Upcoming'}</span></div>`;d.append(card);
        button('Package options',()=>{const item=button('View receipt',()=>{item.remove();receipt(sim,p)},menu);item.role='menuitem'},card).setAttribute('aria-label','Package options');
      }
      const after=async()=>{await refresh();drawer(db.subscribers.find(s=>s.uuid===sub.uuid))};
      if(!sim) button('Allocate SIM',()=>{
        const dialog=modal('Allocate a SIM'), p=field(dialog,'Data plan',`<select>${planOptions()}</select>`), i=field(dialog,'SIM',`<select>${db.pool.map(s=>`<option>${s.iccid}</option>`).join('')}</select>`);addCancel(dialog);
        button('Allocate SIM',async()=>{await graphql('allocateSim',{subscriber_id:sub.uuid,network_id:selected,package_id:p.value,iccid:i.value});if(db.mode==='late-sim')return;dialog.remove();await after()},dialog);
      },d);
      else {
        button('Top up data',()=>{
          const dialog=modal('Top up data'),p=field(dialog,'Data plan',`<select>${planOptions()}</select>`);addCancel(dialog);
          button('Top up',async()=>{const plan=plans().find(x=>x.uuid===p.value);await graphql('addPayment',{itemId:p.value,sim:sim.id,payerEmail:sub.email,amount:String(plan.amount),currency:'USD'});if(db.mode==='late-payment')return;dialog.remove();await after()},dialog);
        },d);
        const name=sim.status==='active'?'Deactivate SIM':'Activate SIM';button(name,()=>{const dialog=modal(name);addCancel(dialog);button(name,async()=>{await graphql('toggleSimServiceStatus',{sim_id:sim.id,status:sim.status==='active'?'service_off':'service_on'});dialog.remove();await after()},dialog)},d);
      }
    }
    function receipt(sim,plan) {
      const d=modal('Payment receipt'),p=db.payments.find(p=>p.sim===sim.id&&p.itemId===plan.uuid);
      if(!p) d.innerHTML+='<div>No receipt found for this package. It may have been allocated without a recorded payment.</div>';
      else { d.innerHTML+=`<span>Completed</span><div><div>Method</div><div>Cash</div></div><div>${plan.name}</div><div><span>Total paid</span><span>$${Number(p.amount).toFixed(2)}</span></div><div><div>Payment ID</div><div>${db.mode==='wrong-receipt'?'foreign-payment':p.id}</div></div>`; }
      button('Close',()=>d.remove(),d);
    }
    if (['/business', '/business/revenue', '/business/packages'].includes(path)) {
      const purchases = db.payments.filter(p => db.sims.some(s => s.id === p.sim && s.network_id === selected));
      const revenue = purchases.reduce((sum, p) => sum + Number(p.amount), 0);
      const shownRevenue = db.mode === 'wrong-revenue' ? revenue + 1 : revenue;
      const customers = db.subscribers.filter(s => s.network_id === selected).length;
      const paid = new Set(purchases.map(p => db.sims.find(s => s.id === p.sim)?.subscriber_id)).size;
      const money = n => '$' + n.toLocaleString('en-US', {maximumFractionDigits: 2});
      const tiles = path === '/business' ? [['Revenue', money(shownRevenue)], ['Customers', customers]] :
        path.endsWith('/revenue') ? [['Revenue', money(shownRevenue)], ['Purchases', purchases.length], ['Avg purchase', money(purchases.length ? revenue / purchases.length : 0)], ['Paid customers', paid]] :
        [['Package revenue', money(shownRevenue)], ['Packages sold', purchases.length]];
      main.innerHTML = tiles.map(([label, value]) => `<div class="MuiCard-root"><div><span>${label}</span></div><div><span class="tnum">${value}</span></div></div>`).join('');
      if (path.endsWith('/packages')) {
        main.innerHTML += '<table><thead><tr><th>Package</th><th>Price</th><th>Sold</th><th>Revenue</th><th>Share</th><th>Status</th></tr></thead><tbody></tbody></table>';
        for (const p of plans()) {
          const sold = purchases.filter(s => s.itemId === p.uuid), total = sold.reduce((sum, s) => sum + Number(s.amount), 0);
          const row = document.createElement('tr');
          row.innerHTML = `<td>${p.name}</td><td>${money(p.amount)}</td><td>${sold.length || '—'}</td><td>${money(db.mode === 'wrong-plan-revenue' && sold.length ? total + 1 : total)}</td><td>${revenue ? Math.round(total / revenue * 100) : 0}%</td><td>Active</td>`;
          main.querySelector('tbody').append(row);
        }
      }
    } else if(path.endsWith('/manage/data-plans')) {
      main.innerHTML='<div class="pagehead"></div>';button('Create plan',()=>{
        const d=modal('Create data plan'),n=field(d,'Data plan name','<input>'),a=field(d,'Price','<input type="number">'),v=field(d,'Data volume','<input type="number">'),u=field(d,'Unit','<select><option>GB</option><option>MB</option></select>'),days=field(d,'Validity','<select><option value="1">Daily (1 day)</option><option value="7">Weekly (7 days)</option><option value="30">Monthly (30 days)</option></select>');
        const org=field(d,'Available to all networks','<input type="checkbox" aria-label="Available to all networks">'),net=field(d,'Network',`<select>${db.networks.map(x=>`<option value="${x.id}">${x.name}</option>`).join('')}</select>`);d.innerHTML+='<div>✓ Name is available</div>';
        // innerHTML above reparses nodes; read current form values when submitting.
        addCancel(d);button('Create plan',async()=>{
          const f=[...d.querySelectorAll('input,select')];await graphql('addPackage',{name:f[0].value,amount:+f[1].value,dataVolume:+f[2].value,dataUnit:f[3].value,duration:+f[4].value*1440,currency:'USD',country:'USA',networkId:f[5].checked?'':f[6].value});
          if(db.mode==='late-plan')return;d.remove();await refresh();
        },d);
      },main.firstChild);
      for(const p of db.plans) {
        const card=document.createElement('div');card.className='card';card.innerHTML=`<div></div><div class="card-pad"><div><div>${p.name}</div></div><div><span class="tnum">$${p.amount}</span><span> / ${unit(p.duration)}</span></div><div>${p.dataVolume} ${p.dataUnit} data · ${terms(p.duration)} validity</div><div><span title="scope">${db.networks.find(n=>n.id===p.networkId)?.name||'All networks'}</span></div></div>`;main.append(card);
        button('Plan actions',()=>{const item=button('Edit plan',()=>{item.remove();const d=modal('Edit data plan');field(d,'Data plan name',`<input value="${p.name}">`);for(const [l,v] of [['Price',p.amount],['Data volume',p.dataVolume],['Unit',p.dataUnit],['Validity',terms(p.duration)]])field(d,l,`<div class="ff-readonly">${v}</div>`);addCancel(d)},menu);item.role='menuitem'},card.querySelector('.card-pad > div')).setAttribute('aria-label','Plan actions');
      }
    } else if(path.endsWith('/customers')) {
      main.innerHTML='<div class="pagehead"></div><table><tbody></tbody></table>';button('Add customer',()=>{
        const d=modal('Add customer'),first=field(d,'First name','<input>'),last=field(d,'Last name','<input>'),email=field(d,'Email','<input>');field(d,'Data plan',`<select>${planOptions()}</select>`);addCancel(d);
        button('Add customer',async()=>{await graphql('addSubscriber',{name:first.value+' '+last.value,email:email.value,network_id:selected});d.remove();await refresh()},d);
      },main.firstChild);
      for(const sub of db.subscribers.filter(s=>s.network_id===selected)) {const row=document.createElement('tr');row.innerHTML=`<td><b>${sub.name}</b></td>`;row.onclick=()=>drawer(sub);main.querySelector('tbody').append(row)}
    } else if(path.endsWith('/manage/sim-pool')) {
      main.innerHTML='<div class="pagehead"></div><table><tbody></tbody></table>';button('Upload SIMs',()=>{const d=modal('Upload SIMs'),file=document.createElement('input');file.type='file';d.append(file);button('Upload',async()=>{await graphql('uploadSims',{csv:await file.files[0].text()});d.remove();await refresh()},d)},main.firstChild);
      for(const sim of db.pool){const row=document.createElement('tr');row.innerHTML=`<td>${sim.iccid}</td><td><span class="MuiChip-label">${sim.assigned?'Assigned':'Available'}</span></td>`;main.querySelector('tbody').append(row)}
    }
  };
}
export const commerceFixture = mode => provisioningFixture(mode, {
  script:`(${client.toString()})();`,
  init(db){ Object.assign(db,{factory:[],plans:[],pool:[],subscribers:[],sims:[],payments:[],entitlements:[],usage:0}); },
  get(req,res,db) {
    if (req.method !== 'GET' || !req.url.startsWith('/v1/')) return false;
    if (req.url.startsWith('/v1/sims/csv')) {res.setHeader('content-type','text/csv');res.end('iccid,imsi\n'+db.factory.map(s=>s.iccid+','+s.imsi).join('\n')+'\n');}
    else if (req.url.startsWith('/v1/sims?')) {res.setHeader('content-type','application/json');res.end(JSON.stringify({sims:db.factory}));}
    else if (req.url.startsWith('/v1/asr/')) res.end('{}');
    else {res.statusCode=404;res.end('{}');}
    return true;
  },
  request(url,body,db) {
    if(url==='/v1/sims'){db.factory.push(body);return true;}
    if(url==='/runtime' && body.action==='traffic.sh')db.usage=Number(body.args[1])*1024*1024;
    return false;
  },
  graphql(body,db) {
    const op=/\b(addPackage|deletePackage|addSubscriber|deleteSubscriber|allocateSim|addPayment|toggleSimServiceStatus|uploadSims|getPackagesForSim|unsetPackageInUseForSim|removePackageForSim|deleteSim)\s*\(/.exec(body.query)?.[1];if(!op)return;
    const data=body.variables?.data||{},id=/\w+(?:Id|_id):\s*"([^"]+)"/.exec(body.query)?.[1];let value;
    db.operations.push({op,data,id});
    if(op==='addPackage'){value={...data,uuid:'plan-'+(db.plans.length+1)};if(mode==='wrong-minutes')value.duration/=1440;db.plans.push(value);if(mode==='opaque-plan')value={name:data.name};}
    if(op==='deletePackage'){db.plans=db.plans.filter(x=>x.uuid!==id);value={uuid:id}}
    if(op==='addSubscriber'){value={...data,uuid:'sub-'+(db.subscribers.length+1)};db.subscribers.push(value)}
    if(op==='deleteSubscriber'){db.subscribers=db.subscribers.filter(x=>x.uuid!==id);value={success:true}}
    if(op==='allocateSim'){value={...data,id:'sim-'+(db.sims.length+1),status:'active'};db.sims.push(value);db.pool.find(s=>s.iccid===data.iccid).assigned=true;db.entitlements.push({sim:value.id,packageId:data.package_id,active:true})}
    if(op==='addPayment'){value={...data,id:'payment-'+(db.payments.length+1)};db.payments.push(value);if(mode!=='no-entitlement')db.entitlements.push({sim:data.sim,packageId:data.itemId,active:false});if(mode==='opaque-payment')value={}}
    if(op==='toggleSimServiceStatus'){const sim=db.sims.find(s=>s.id===(data.sim_id||id));if(sim)sim.status=data.status==='service_on'?'active':'inactive';value={success:true}}
    if(op==='uploadSims'){for(const line of data.csv.trim().split('\n').slice(1)){const iccid=line.split(',')[0];db.pool.push({iccid})}value={iccid:db.pool.map(s=>s.iccid)}}
    if(op==='getPackagesForSim')value={packages:[]};
    if(op==='deleteSim'){db.sims=db.sims.filter(s=>s.id!==id);value={simId:id}}
    return {data:{[op]:value}};
  }
});

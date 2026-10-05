/* SPDX-License-Identifier: MPL-2.0
 * Source-shaped controlled DOM fixture. Never counts as target-app coverage.
 */
import { provisioningFixture } from './provisioning-fixture.mjs';
function receiptPdfFixture(lines) {
  const escape=s=>Buffer.from(s,'latin1').toString('latin1').replace(/[\\()]/g,'\\$&');
  const content='BT /F1 11 Tf 40 550 Td 24 TL '+lines.map((s,i)=>(i?'T* ':'')+'('+escape(s)+') Tj').join('\n')+' ET';
  const objects=['<< /Type /Catalog /Pages 2 0 R >>','<< /Type /Pages /Kids [3 0 R] /Count 1 >>','<< /Type /Page /Parent 2 0 R /MediaBox [0 0 420 595] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>','<< /Type /Font /Subtype /Type1 /BaseFont /Courier /Encoding /WinAnsiEncoding >>',`<< /Length ${Buffer.byteLength(content,'latin1')} >>\nstream\n${content}\nendstream`];
  let body='%PDF-1.4\n',offsets=[0];for(const [i,o] of objects.entries()){offsets.push(Buffer.byteLength(body,'latin1'));body+=`${i+1} 0 obj\n${o}\nendobj\n`;}
  const xref=Buffer.byteLength(body,'latin1');body+='xref\n0 6\n0000000000 65535 f \n'+offsets.slice(1).map(o=>String(o).padStart(10,'0')+' 00000 n ').join('\n')+`\ntrailer << /Size 6 /Root 1 0 R >>\nstartxref\n${xref}\n%%EOF\n`;return Buffer.from(body,'latin1');
}
function client() {
  const terms = minutes => ({1440:'1 day',10080:'1 week',43200:'1 month'})[minutes] || `${minutes} minutes`;
  const unit = minutes => ({1440:'day',10080:'week',43200:'month'})[minutes] || 'period';
  const bytes = n => n >= 1073741824 ? `${Math.round(n / 1073741824 * 100)/100} GB` : n >= 1048576 ? `${n / 1048576} MB` : `${n} B`;
  window.renderOperations = async ({db,main,menu,button,graphql,render}) => {
    const path = location.pathname, selected = sessionStorage.getItem('network') || 'existing';
    const move = p => { history.pushState({},'',p); void render(); };
    document.querySelectorAll('.viewseg button').forEach(b=>b.onclick=()=>move(b.textContent==='Customer'?'/customer/customers':'/'+b.textContent.toLowerCase()));
    const lens=path.split('/')[1];
    const nav = document.querySelector('aside');
    const home=[...nav.querySelectorAll('a')].find(a=>a.textContent==='Home');if(home){home.href='/'+lens;home.onclick=e=>{e.preventDefault();move('/'+lens)}}
    for (const [n,p] of [['Members',`/${lens}/manage/members`],['Data plans',`/${lens}/manage/data-plans`],['Customers',`/${lens}/customers`],['SIM pool',`/${lens}/manage/sim-pool`],['Revenue',`/${lens}/revenue`],['Packages',`/${lens}/packages`]]) {
      const a=document.createElement('a');a.href=p;a.textContent=n;a.onclick=e=>{e.preventDefault();move(p)};nav.append(a);
    }
    document.querySelector('[aria-label="Open navigation"]')?.remove();
    const mobile=button('Open navigation',()=>{const d=document.createElement('div');d.className='mobile-nav';d.innerHTML=nav.innerHTML;document.body.append(d);d.querySelectorAll('a').forEach(a=>a.onclick=e=>{e.preventDefault();d.remove();move(a.getAttribute('href'))})},document.querySelector('header'));mobile.setAttribute('aria-label','Open navigation');
    document.onkeydown=e=>{if(e.ctrlKey&&e.key==='k'){e.preventDefault();const d=modal('Jump to'),input=document.createElement('input');input.placeholder='Jump to a page…';d.append(input);input.focus();button('Data plans',()=>{d.remove();move('/business/manage/data-plans')},d)}};
    const modal = title => { const opener=document.activeElement, d=document.createElement('div');d.role='dialog';d.innerHTML=`<h2>${title}</h2>`;menu.append(d);d.tabIndex=-1;d.focus(); const escape=e=>{if(e.key==='Escape'){d.remove();if(db.mode!=='lost-focus')opener?.focus();document.removeEventListener('keydown',escape)}};document.addEventListener('keydown',escape);return d; };
    const field = (d,label,html) => {const f=document.createElement('label');f.className='ff';f.innerHTML=`<div class="ff-label">${label}</div>${html}`;d.append(f);return f.querySelector('input,select')};
    const addCancel=d=>button('Cancel',()=>d.remove(),d);
    const plans = () => db.plans.filter(p=>db.mode==='scope-leak'||!p.networkId||p.networkId===selected);
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
      if(!sim) button('Allocate a SIM',()=>{
        const dialog=modal('Allocate a SIM'), p=field(dialog,'Data plan',`<select>${planOptions()}</select>`), i=field(dialog,'SIM',`<select><option value="">Auto-assign from pool</option>${db.pool.filter(s=>db.mode==='allocated-option-leak'||!s.assigned&&!s.failed).map(s=>`<option value="${s.iccid}">${s.iccid}</option>`).join('')}</select>`);addCancel(dialog);
        button('Allocate SIM',async()=>{await graphql('allocateSim',{subscriber_id:sub.uuid,network_id:selected,package_id:p.value,iccid:i.value});if(db.mode==='late-sim')return;dialog.remove();await after()},dialog);
      },d);
      else {
        button('Top up',()=>{
          const dialog=modal('Top up data'),p=field(dialog,'Data plan',`<select>${planOptions()}</select>`);addCancel(dialog);
          const submit=button('Top up',async()=>{
            if(db.mode!=='double-payment')submit.disabled=true;
            const plan=plans().find(x=>x.uuid===p.value);
            const response=await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'addPayment',query:'mutation addPayment($data: AddPaymentInputDto!) { addPayment(data: $data) { id } }',variables:{data:{itemId:p.value,sim:sim.id,payerEmail:sub.email,amount:String(plan.amount),currency:'USD'}}})}).then(r=>r.json());
            if(response.errors){submit.disabled=false;const error=document.createElement('div');error.textContent=db.mode==='hidden-payment-error'?'':response.errors[0].message;menu.append(error);if(db.mode==='error-grants-entitlement')await graphql('fixtureEntitlement',{sim:sim.id,itemId:p.value});return;}
            if(db.mode==='late-payment')return;dialog.remove();await after();
          },dialog);submit.disabled=true;p.onchange=()=>{submit.disabled=!p.value};
        },d);
        const name=sim.status==='active'?'Deactivate SIM':'Activate SIM';button(name,()=>{const dialog=modal(name);addCancel(dialog);button(name,async()=>{await graphql('toggleSimServiceStatus',{sim_id:sim.id,status:sim.status==='active'?'service_off':'service_on'});dialog.remove();await after()},dialog)},d);
      }
    }
    function receipt(sim,plan) {
      const d=modal('Payment receipt'),p=db.payments.find(p=>p.sim===sim.id&&p.itemId===plan.uuid);
      if(!p) d.innerHTML+='<div>No receipt found for this package. It may have been allocated without a recorded payment.</div>';
      else {
        const date=new Intl.DateTimeFormat('en-GB',{day:'2-digit',month:'short',year:'numeric',hour:'2-digit',minute:'2-digit',timeZone:'UTC',hour12:false}).format(new Date(p.paidAt))+' UTC';
        const id=db.mode==='wrong-receipt'?'foreign-payment':p.id;
        d.innerHTML+=`<div><span>Fixture Organization</span><span>Completed</span></div><div><div>Receipt no</div><div>${id.slice(0,8)}</div></div><div><div>Paid on</div><div>${date}</div></div><div><div>Method</div><div>Cash</div></div><div><div>Billed to</div><div>Walk-in customer</div><div>SIM ${sim.id.slice(0,8)} · package top-up</div></div><div><span>Description</span><span>Amount</span></div><div><div><div>${plan.name}</div><div>Data package · qty 1</div></div><div>$${Number(p.amount).toFixed(2)}</div></div><div><span>Total paid</span><span>$${Number(p.amount).toFixed(2)}</span></div><div><div>Payment ID</div><div>${id}</div></div><div>Auto-generated · not a tax invoice</div>`;
        button('Download',()=>{const a=document.createElement('a');a.href='/receipt/'+encodeURIComponent(p.id);a.download=`receipt-${p.id.slice(0,8)}.pdf`;a.click()},d);
      }
      button('Close',()=>d.remove(),d);
    }
    if(path==='/network'){
      const active=db.sims.filter(s=>s.network_id===selected&&s.status==='active').length;
      const wrong=db.mode==='foreign-network-usage';
      const fields=[['Active customers',active],['Data volume',bytes(active||wrong?db.usage:0)]];
      main.innerHTML=fields.map(([label,value])=>`<div class="MuiCard-root"><div><span>${label}</span></div><div>${value}</div></div>`).join('');
    } else if (['/business', '/business/revenue', '/business/packages'].includes(path)) {
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
      if(path.endsWith('/revenue')){const head=document.createElement('div');head.className='pagehead';main.prepend(head);let range='Last 30 days';const chip=button(range,()=>{for(const v of ['Last 24h','Last 7 days','Last 30 days']){const item=button(v,()=>{range=v;chip.textContent=v;menu.innerHTML=''},menu);item.role='menuitem'}},head)}
      if (path.endsWith('/packages')) {
        main.innerHTML += '<table><thead><tr><th>Package</th><th>Price</th><th>Sold</th><th>Revenue</th><th>Share</th><th>Status</th></tr></thead><tbody></tbody></table>';
        for (const p of plans()) {
          const sold = purchases.filter(s => s.itemId === p.uuid), total = sold.reduce((sum, s) => sum + Number(s.amount), 0);
          const row = document.createElement('tr');
          row.innerHTML = `<td>${p.name}</td><td>${money(p.amount)}</td><td>${sold.length || '—'}</td><td>${money(db.mode === 'wrong-plan-revenue' && sold.length ? total + 1 : total)}</td><td>${revenue ? Math.round(total / revenue * 100) : 0}%</td><td>Active</td>`;
          main.querySelector('tbody').append(row);
        }
      }
    } else if(path.endsWith('/manage/members')) {
      main.innerHTML='<h1>Members</h1>';
    } else if(path.endsWith('/manage/data-plans')) {
      main.innerHTML='<div class="pagehead"></div>';button('Create plan',()=>{
        const d=modal('Create data plan'),n=field(d,'Data plan name','<input>'),a=field(d,'Price','<input type="number">'),v=field(d,'Data volume','<input type="number">'),u=field(d,'Unit','<select><option>GB</option><option>MB</option></select>'),days=field(d,'Validity','<select><option value="1">Daily (1 day)</option><option value="7">Weekly (7 days)</option><option value="30">Monthly (30 days)</option></select>');
        const org=field(d,'Available to all networks','<input type="checkbox" aria-label="Available to all networks">'),net=field(d,'Network',`<select>${db.networks.map(x=>`<option value="${x.id}">${x.name}</option>`).join('')}</select>`);
        const state=document.createElement('div');d.append(state);
        addCancel(d);const submit=button('Create plan',async()=>{
          await graphql('addPackage',{name:n.value,amount:+a.value,dataVolume:+v.value,dataUnit:u.value,duration:+days.value*1440,currency:'USD',country:'USA',networkId:org.checked?'':net.value});
          if(db.mode==='late-plan')return;d.remove();await refresh();
        },d);
        let nameState='idle',checked='',revision=0;
        const readName=async()=>{const rev=++revision;checked=n.value.trim();if(!checked){nameState='idle';validate();return;}nameState='checking';validate();const response=await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'isPackageNameAvailable',query:'query isPackageNameAvailable($name: String!) { isPackageNameAvailable(name: $name) { name isAvailable } }',variables:{name:checked}})}).then(r=>r.json()).catch(()=>({errors:[{}]}));if(rev!==revision)return;nameState=response.errors?'idle':response.data.isPackageNameAvailable.isAvailable?'available':'taken';validate()};
        const validate=()=>{for(const f of [n,a,v]){f.parentNode.querySelector('.ff-err')?.remove();if(!f.value.trim()||f!==n&&+f.value<=0){const err=document.createElement('div');err.className='ff-err';err.textContent=f===n?'Plan name is required':'Must be > 0';f.parentNode.append(err)}}state.textContent=nameState==='taken'?'That plan name is already taken':nameState==='checking'?'Checking availability…':nameState==='available'?'✓ Name is available':'';submit.disabled=db.mode!=='validation-broken'&&(!n.value.trim()||+a.value<=0||+v.value<=0||nameState==='taken'||db.mode!=='name-fail-open'&&nameState!=='available')};
        d.addEventListener('input',()=>{validate();if(n.value.trim()!==checked)void readName()});d.addEventListener('change',validate);validate();

      },main.firstChild);
      for(const p of db.plans) {
        const card=document.createElement('div');card.className='card';card.innerHTML=`<div></div><div class="card-pad"><div><div>${p.name}</div></div><div><span class="tnum">$${p.amount}</span><span> / ${unit(p.duration)}</span></div><div>${p.dataVolume} ${p.dataUnit} data · ${terms(p.duration)} validity</div><div><span title="scope">${db.networks.find(n=>n.id===p.networkId)?.name||'All networks'}</span></div></div>`;main.append(card);
        button('Plan actions',()=>{const item=button('Edit plan',()=>{item.remove();const d=modal('Edit data plan');const name=field(d,'Data plan name',`<input value="${p.name}">`);const available=document.createElement('div');d.append(available);name.oninput=()=>{available.textContent='✓ Name is available'};button('Save changes',async()=>{await graphql('updatePackage',{packageId:p.uuid,name:name.value});d.remove();await refresh()},d);for(const [l,v] of [['Price','$'+p.amount],['Data volume',p.dataVolume],['Unit',p.dataUnit],['Validity',terms(p.duration)]])field(d,l,`<div class="ff-readonly">${v}</div>`);addCancel(d)},menu);item.role='menuitem'},card.querySelector('.card-pad > div')).setAttribute('aria-label','Plan actions');
      }
    } else if(path.endsWith('/customers')) {
      main.innerHTML='<div class="pagehead"></div><input placeholder="Search name or phone"><table><thead><tr><th>Customer</th><th><button>Active plan</button></th><th>Data usage</th><th><button>SIM</button></th><th>Last seen</th></tr></thead><tbody></tbody></table><div class="empty"></div>';
      button('Add customer',()=>{
        if(!db.pool.filter(s=>!s.assigned&&!s.failed).length||!plans().length){const toast=document.createElement('div');toast.textContent=!db.pool.filter(s=>!s.assigned&&!s.failed).length?'No SIMs available — please upload SIMs to your SIM pool first.':'No data plans yet — please create a data plan before adding customers.';menu.append(toast);return;}
        const d=modal('Add customer'),first=field(d,'First name','<input>'),last=field(d,'Last name','<input>'),email=field(d,'Email','<input>');field(d,'Data plan',`<select>${planOptions()}</select>`);addCancel(d);
        const submit=button('Add customer',async()=>{await graphql('addSubscriber',{name:first.value+' '+last.value,email:email.value,network_id:selected});d.remove();await refresh()},d);
        const validate=()=>{email.parentNode.querySelector('.ff-err')?.remove();const invalid=email.value&&!/^[^@]+@[^@]+\.[^@]+$/.test(email.value);if(invalid){const err=document.createElement('div');err.className='ff-err';err.textContent='Enter a valid email';email.parentNode.append(err)}submit.disabled=!first.value.trim()||!!invalid};d.addEventListener('input',validate);validate();
      },main.firstChild);
      let search='',descending=false,simFilter='All',planFilter='All';
      const rows=()=>{let source=db.subscribers.filter(s=>db.mode==='customer-leak'||s.network_id===selected);let shown=source.filter(s=>s.name.toLowerCase().includes(search.toLowerCase())).filter(s=>{const sim=db.sims.find(x=>x.subscriber_id===s.uuid);return(simFilter==='All'||sim?.status===simFilter.toLowerCase())&&(planFilter!=='No plan'||!sim)}).sort((a,b)=>a.name.localeCompare(b.name)*(descending?-1:1));main.querySelector('tbody').innerHTML='';main.querySelector('table').hidden=!shown.length;main.querySelector('.empty').textContent=shown.length?'':source.length?'No customers match':'No customers yet';for(const sub of shown){const row=document.createElement('tr');row.role='button';row.tabIndex=0;row.innerHTML=`<td><div><span>LU</span><div><div>${sub.name}</div><div></div></div></div></td><td>${db.mode==='wrong-customer-plan'?'Foreign plan':db.plans.find(p=>p.uuid===db.entitlements.find(e=>e.active&&e.sim===db.sims.find(s=>s.subscriber_id===sub.uuid)?.id)?.packageId)?.name||'No plan'}</td><td>—</td><td>Inactive</td><td>—</td>`;row.onclick=()=>drawer(sub);main.querySelector('tbody').append(row)}};
      main.querySelector('input').oninput=e=>{search=e.target.value;rows()};main.querySelector('th').onclick=()=>{descending=!descending;rows()};
      for(const [column,opts] of [['SIM',['All','active','inactive','suspended']],['Active plan',['All','No plan']]]){const opener=[...main.querySelectorAll('th button')].find(b=>b.textContent===column);opener.onclick=()=>{for(const val of opts){const item=button(val,()=>{if(column==='SIM')simFilter=val;else planFilter=val;menu.innerHTML='';rows()},menu);item.role='menuitem'}}}rows();
    } else if(path.endsWith('/manage/sim-pool')) {
      main.innerHTML='<div class="MuiSkeleton-root">Loading</div>';
      const response=await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'SimPoolOverview',query:'query SimPoolOverview($simType: String!, $limit: Int!) { simPoolView(simType: $simType) { sims(limit: $limit) { sims { iccid } } } }',variables:{simType:'test',limit:100}})}).then(r=>r.json());
      main.innerHTML='<div class="pagehead"></div><div class="kpis"></div><div class="filters"></div><table><tbody></tbody></table><div class="pool-empty"></div><div class="tbl-foot"><span class="tnum"></span></div>';
      button('Upload SIMs',()=>{const d=modal('Upload SIMs'),file=document.createElement('input');file.type='file';d.append(file);button('Upload',async()=>{await graphql('uploadSims',{csv:await file.files[0].text()});d.remove();await refresh()},d)},main.firstChild);
      const status=s=>s.failed?'Faulty':s.assigned?'Assigned':'Available';
      const failed=!!response.data?.simPoolView?.sims?.error;
      for(const name of ['Assigned','Available','Faulty']) {const n=db.pool.filter(s=>status(s)===name).length;main.querySelector('.kpis').innerHTML+=`<div class="MuiCard-root"><div><span>${name}</span></div><div>${failed?'—':db.mode==='wrong-pool-count'?n+1:n}</div></div>`;}
      if(failed && db.mode!=='pool-fail-empty'){main.querySelector('.pool-empty').textContent="Couldn't load SIMs";return;}
      const sims=failed?[]:db.pool;
      const show=filter=>{main.querySelector('tbody').innerHTML='';const rows=sims.filter(s=>filter==='All statuses'||status(s)===filter);for(const sim of rows){const row=document.createElement('tr');row.innerHTML=`<td>${sim.iccid}</td><td><span class="MuiChip-label">${status(sim)}</span></td>`;main.querySelector('tbody').append(row)}main.querySelector('.pool-empty').textContent=rows.length?'':sims.length?'No SIMs match':'No SIMs';main.querySelector('.tbl-foot .tnum').textContent=`Showing ${rows.length} of ${sims.length}`};
      for(const filter of ['All statuses','Available','Assigned','Faulty'])button(filter,()=>show(filter),main.querySelector('.filters'));show('All statuses');
    }
  };
}
export const commerceFixture = (mode, extension = {}) => provisioningFixture(mode, {
  script:`(${client.toString()})();${extension.script || ""}`,
  init(db){ Object.assign(db,{factory:[],plans:[],pool:[],subscribers:[],sims:[],payments:[],entitlements:[],usage:0}); extension.init?.(db); },
  get(req,res,db) {
    if(req.method==='GET' && req.url.startsWith('/receipt/')){
      const id=decodeURIComponent(req.url.slice('/receipt/'.length)),p=db.payments.find(p=>p.id===id),plan=db.plans.find(x=>x.uuid===p.itemId);
      const date=new Intl.DateTimeFormat('en-GB',{day:'2-digit',month:'short',year:'numeric',hour:'2-digit',minute:'2-digit',timeZone:'UTC',hour12:false}).format(new Date(p.paidAt))+' UTC';
      const money='$'+Number(p.amount).toFixed(2);
      const lines=['Payment receipt','Fixture Organization','Completed','RECEIPT NO    PAID ON    METHOD',p.id.slice(0,8)+'    '+date+'    Cash','BILLED TO','Walk-in customer','SIM '+p.sim.slice(0,8)+' · package top-up','DESCRIPTION    AMOUNT',plan.name, money,'Data package · qty 1','Total paid',mode==='wrong-pdf-total'?'$999.00':money,'PAYMENT ID',mode==='wrong-pdf-id'?'foreign-payment':p.id,'Auto-generated · not a tax invoice'];
      const pdf=receiptPdfFixture(lines);
      res.setHeader('content-type',mode==='html-download'?'text/html':'application/pdf');res.setHeader('content-disposition',`attachment; filename="receipt-${p.id.slice(0,8)}.pdf"`);res.end(mode==='html-download'?'<html>not a receipt</html>':pdf);return true;
    }
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
    const extra=extension.graphql?.(body,db);if(extra)return extra;
    if(body.operationName==='isPackageNameAvailable')return {data:{isPackageNameAvailable:{name:body.variables.name,isAvailable:!db.plans.some(p=>p.name===body.variables.name)}}};
    if(body.operationName==='SimPoolOverview')return {data:{simPoolView:{sims:{error:null}}}};
    const op=/\b(fixtureEntitlement|updatePackage|addPackage|deletePackage|addSubscriber|deleteSubscriber|allocateSim|addPayment|toggleSimServiceStatus|uploadSims|getPackagesForSim|unsetPackageInUseForSim|removePackageForSim|deleteSim)\s*\(/.exec(body.query)?.[1];if(!op)return;
    const data=body.variables?.data||{},id=/\w+(?:Id|_id):\s*"([^"]+)"/.exec(body.query)?.[1];let value;
    db.operations.push({op,data,id});
    if(op==='fixtureEntitlement'){db.entitlements.push({sim:data.sim,packageId:data.itemId,active:false});value={ok:true};}
    if(op==='updatePackage'){value=db.plans.find(p=>p.uuid===data.packageId);if(mode!=='rename-lost')value.name=data.name;if(mode==='rename-terms-drift')value.amount++;}
    if(op==='addPackage'){value={...data,uuid:'plan-'+(db.plans.length+1)};if(mode==='wrong-minutes')value.duration/=1440;db.plans.push(value);if(mode==='opaque-plan')value={name:data.name};}
    if(op==='deletePackage'){db.plans=db.plans.filter(x=>x.uuid!==id);value={uuid:id}}
    if(op==='addSubscriber'){value={...data,uuid:'sub-'+(db.subscribers.length+1)};db.subscribers.push(value)}
    if(op==='deleteSubscriber'){db.subscribers=db.subscribers.filter(x=>x.uuid!==id);value={success:true}}
    if(op==='allocateSim'){data.iccid ||= db.pool.find(s=>!s.assigned&&!s.failed)?.iccid;value={...data,id:'sim-'+(db.sims.length+1),status:'active'};db.sims.push(value);db.pool.find(s=>s.iccid===data.iccid).assigned=true;db.entitlements.push({sim:value.id,packageId:data.package_id,active:true})}
    if(op==='addPayment'){value={...data,paidAt:mode==='old-receipt-date'?0:Date.now(),id:'payment-'+(db.payments.length+1)};db.payments.push(value);if(mode!=='no-entitlement')db.entitlements.push({sim:data.sim,packageId:data.itemId,active:false});if(mode==='opaque-payment')value={}}
    if(op==='toggleSimServiceStatus'){const sim=db.sims.find(s=>s.id===(data.sim_id||id));if(sim)sim.status=data.status==='service_on'?'active':'inactive';value={success:true}}
    if(op==='uploadSims'){for(const line of data.csv.trim().split('\n').slice(1)){const iccid=line.split(',')[0];db.pool.push({iccid})}value={iccid:db.pool.map(s=>s.iccid)}}
    if(op==='getPackagesForSim')value={packages:[]};
    if(op==='deleteSim'){db.sims=db.sims.filter(s=>s.id!==id);value={simId:id}}
    return {data:{[op]:value}};
  }
});

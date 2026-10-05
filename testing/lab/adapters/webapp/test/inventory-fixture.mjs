/* SPDX-License-Identifier: MPL-2.0
 * Intended DOM behavior plus adverse variants; never live console evidence.
 */
import { provisioningFixture } from './provisioning-fixture.mjs';
function client(){
  let range='Last 24h',q='',rerender,active,latest,selectedSite;
  const esc=s=>String(s).replaceAll('&','&amp;').replaceAll('<','&lt;').replaceAll('"','&quot;');
  let previousNodes;
  window.updateOperations=db=>{const fingerprint=JSON.stringify(db.nodes);if(previousNodes!==undefined&&previousNodes!==fingerprint&&location.pathname==='/network')void rerender?.();previousNodes=fingerprint;};
  window.inventorySwitch=()=>{
    const prefs=JSON.parse(localStorage.getItem('uk-ui-prefs')||'{"state":{}}');prefs.state.networkId=sessionStorage.getItem('network');localStorage.setItem('uk-ui-prefs',JSON.stringify(prefs));
    if(latest?.mode==='no-detail-redirect')return;
    if(/^\/network\/(sites|nodes)\//.test(location.pathname))history.replaceState({},'',location.pathname.split('/').slice(0,3).join('/'));
  };
  const prefs=JSON.parse(localStorage.getItem('uk-ui-prefs')||'{"state":{}}');
  if(prefs.state.networkId==='ulab-removed-network-selection')sessionStorage.setItem('network','existing');
  if(!prefs.state.networkId)localStorage.setItem('uk-ui-prefs',JSON.stringify({state:{networkId:'existing',accent:'blue'},version:0}));
  window.renderOperations=async({db,main,menu,button,render,move,selected})=>{
    latest=db;active=selected;rerender=render;
    const path=location.pathname;
    const sites=db.sites.filter(s=>s.network_id===selected),nodes=db.nodes.filter(n=>sites.some(s=>s.id===n.site));
    const bad=db.mode;
    const status=s=>{const ns=db.nodes.filter(n=>n.site===s.id);return ns.every(n=>n.offline)?'Offline':ns.some(n=>n.offline)?'Degraded':'Online'};
    const header=(name,count)=>`<div class="pagehead"><div class="pagetitle">${name}${count===undefined?'':`<span class="cnt">${count||''}</span>`}</div></div>`;
    const chip=s=>`<div class="MuiChip-root"><span class="MuiChip-label">${s}</span></div>`;
    if(path==='/network/sites'){
      main.innerHTML=header('Sites',sites.length)+'<input placeholder="Search sites"><div class="tile-grid"></div>';
      const input=main.querySelector('input');input.value=q;const grid=main.lastChild;
      const draw=()=>{grid.innerHTML='';let list=sites.filter(s=>s.name.toLowerCase().includes(q.toLowerCase()));
        if(bad==='extra-site')list=[...list,{id:'alien',name:'alien-site'}];
        if(bad==='duplicate-site'&&list.length&&db.sites.length===4&&selected===db.networks[1]?.id)list.push(list[0]);
        for(const s of list){const card=document.createElement('div');card.className='card ecard';card.role='button';card.innerHTML=`<div><div><div>icon</div><div><div>${esc(s.name)}</div><div>${bad==='wrong-location'?'Wrong place':'Calgary'}</div></div></div>${chip(status(s))}</div><hr><div><span>3 nodes</span></div>`;card.onclick=()=>move('/network/sites/'+s.id);grid.append(card);}
        if(!list.length)grid.innerHTML='<div>'+(!sites.length?'No sites yet':'No matching sites')+'</div>';
      };input.oninput=()=>{q=input.value;draw()};draw();
    }else if(path.startsWith('/network/sites/')){
      const s=db.sites.find(s=>path.endsWith('/'+s.id));if(!s)return;
      const ns=db.nodes.filter(n=>n.site===s.id);
      main.innerHTML=header(esc(s.name))+`<div class="detail-subrow">${chip(status(s))}</div><div class="card"><div class="sec-title">Site information</div><div class="sec-body"><div><div><div>Nodes</div>${ns.map(n=>`<div><span class="tnum">${bad==='wrong-site-node'?'wrong-node':n.id}</span><span> · ${n.type}</span></div>`).join('')}</div><div><div>Location</div><div>Calgary</div><div class="tnum">${bad==='wrong-coordinates'?'0, 0':'51.05, -114.07'}</div></div></div></div></div>`;
    }else if(path==='/network/nodes'&&bad==='malformed-node'){
      const e=document.createElement('div');e.className='ecard';e.role='button';e.textContent='Malformed node';main.append(e);
    }else if(path==='/network/nodes'&&bad==='foreign-node'){
      const foreign=db.nodes.find(n=>n.site&&!sites.some(s=>s.id===n.site));if(foreign){const e=document.createElement('div');e.className='ecard';e.role='button';e.innerHTML=`<div class="tnum">Tower node · ${foreign.id}</div>`;main.append(e);}
    }else if(path==='/network'){
      const req=async op=>{const r=await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:op,query:`query ${op} { placeholder }`,variables:{networkId:selected}})});return (await r.json()).data};
      const k=await req('GetKpiValues'),registry=await req('SitesList');
      const unavailable=!k.getKpiValues.values.length,missing=!!registry.sitesView.sites.error;
      const count=sites.filter(s=>!db.nodes.some(n=>n.site===s.id&&n.offline)).length;
      const vals={'Network uptime':unavailable?'—':({'Last 24h':'98.5%','Last 7 days':'97.2%','Last 30 days':'95.0%'}[range]),'Active customers':unavailable?'—':'0','Data volume':unavailable?'—':'0 B','Sites online':unavailable||missing?'—':`${count}/${sites.length}`};
      if(bad==='invented-zero'&&(unavailable||missing))vals['Sites online']='0/0';
      main.innerHTML=header('Home')+`<div class="pagehead"><button>${range}</button></div>`+Object.entries(vals).map(([label,value])=>`<div class="MuiCard-root"><div><span>${label}</span></div><div>${value}</div></div>`).join('')+'<div class="leaflet-container" style="height:140px"></div><div id="selection"></div>';
      main.querySelector('.pagehead button').onclick=()=>{menu.innerHTML='';for(const r of ['Last 24h','Last 7 days','Last 30 days']){const b=button(r,()=>{range=r;menu.innerHTML='';void render()},menu);b.role='menuitem'}};
      const map=main.querySelector('.leaflet-container');
      for(const s of sites){const pin=document.createElement('button');pin.className='uk-map-pin';pin.innerHTML=`<svg width="30" height="30" fill="${bad==='wrong-pin-color'?'var(--uk-error)':'var(--uk-success-bright)'}"><circle cx="15" cy="15" r="10"/></svg>`;
        pin.onclick=()=>{selectedSite=s;let popup=map.querySelector('.leaflet-popup-content');if(!popup){popup=document.createElement('div');popup.className='leaflet-popup-content';map.append(popup)}popup.innerHTML=`<div><div>${esc(s.name)}</div><div>Calgary</div></div>`;const panel=document.querySelector('#selection');panel.innerHTML=`<div><div><div><span>${esc(bad==='wrong-map-site'?'wrong-site':s.name)}</span>${chip('Online')}</div><div>Calgary</div></div></div>`;const box=panel.firstChild;button('Open site',()=>move('/network/sites/'+(bad==='wrong-map-link'?sites.at(-1).id:selectedSite.id)),box);button('Clear',()=>{selectedSite=undefined;panel.innerHTML=''},box)};map.append(pin);
      }
    }
    if(bad==='transient-leak'&&sites.length&&db.networks.filter(n=>n.id!=='existing').at(-1)?.id===selected){const foreign=db.sites.find(s=>s.network_id!==selected);if(foreign){const leak=document.createElement('div');leak.className='ecard';leak.textContent=foreign.name;main.append(leak);setTimeout(()=>leak.remove(),100);}}
  };
}
export async function inventoryFixture(mode=''){
  return provisioningFixture(mode,{script:`(${client.toString()})()`,graphql(body){
    if(body.operationName==='GetKpiValues')return {data:{getKpiValues:{values:[{kpi:'SITES_ONLINE',value:1}]}}};
    if(body.operationName==='SitesList')return {data:{sitesView:{sites:{sites:[],error:null}}}};
  }});
}

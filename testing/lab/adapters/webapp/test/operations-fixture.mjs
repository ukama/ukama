/* SPDX-License-Identifier: MPL-2.0
 * Controlled, source-shaped operation UI. Does not execute React or BFF code.
 * Product coverage still requires running the scenarios against the real stack.
 */
import { provisioningFixture } from './provisioning-fixture.mjs';
import { writeFile } from 'node:fs/promises';
import { join } from 'node:path';
function operationsClient() {
  let ctx, software = false, menuOpen = false, dialogOpen = false, confirm = '', path, component='Node';
  let statusUnknown=false, statusBusy=false, optimisticUntil=0, statusPending=false, statusFingerprint;
  const localPorts={};
  async function pollStatus() {
    if (!ctx || statusPending || !/\/network\/(sites|nodes)\/[^/]+$/.test(location.pathname)) return;
    statusPending=true;
    const before=location.pathname, site=before.startsWith('/network/sites/'), id=before.split('/').at(-1);
    const operationName=site?'GetSiteOperationStatus':'GetNodeOperationStatus',key=site?'siteId':'nodeId';
    try {
      const result=await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName,query:`query ${operationName}($${key}: String!) { ${site?'getSiteOperationStatus':'getNodeOperationStatus'}(data: { ${key}: $${key} }) { busy } }`,variables:{[key]:id}})}).then(r=>r.json());
      if(location.pathname!==before)return;
      statusUnknown=!!result.errors;
      statusBusy=!!result.data?.[site?'getSiteOperationStatus':'getNodeOperationStatus']?.busy;
      if(statusBusy)optimisticUntil=0;
    } finally {
      statusPending=false;
      const next=JSON.stringify([statusUnknown,statusBusy,Date.now()<optimisticUntil]);
      if(next!==statusFingerprint){statusFingerprint=next;refresh(ctx.db);}
    }
  }
  setInterval(()=>void pollStatus(),100);
  const send = async (action, data = {}) => {
    const id = location.pathname.split('/').at(-1);
    if(action==='port')localPorts[data.port]=data.value;
    if(action==='restart-node')optimisticUntil=Date.now()+(ctx.db.mode==='timeout-early'?500:ctx.db.mode==='timeout-stuck'?60000:10000);
    await fetch('/operation', { method:'POST', body:JSON.stringify({ action, id, ...data }) });
    void pollStatus();
  };
  const el = (tag, text, parent, attrs = {}) => {
    const node = document.createElement(tag); node.textContent=text;
    for(const [k,v] of Object.entries(attrs)) node.setAttribute(k,v);
    parent.append(node);return node;
  };
  const states={available:'Update available',running:'Updating…',failed:'Update failed',done:'Up to date'};
  function refresh(db) {
    if (!ctx) return;
    ctx.db=db;
    const id=location.pathname.split('/').at(-1),sitePage=location.pathname.startsWith('/network/sites/'),nodePage=location.pathname.startsWith('/network/nodes/');
    if (!sitePage&&!nodePage) return;
    const node = nodePage?db.nodes.find(n=>n.id===id):null;
    const site = sitePage?db.sites.find(s=>s.id===id):null;
    if ((!nodePage||!node)&&(!sitePage||!site)) return;
    const nodes = db.nodes.filter(n=>n.site===(site?.id??node?.site));
    const siteBusy=nodes.some(n=>n.busy), busy=nodePage?(statusBusy||Date.now()<optimisticUntil):siteBusy;
    const unknown=(db.mode==='status-error'||statusUnknown)&&db.mode!=='status-fail-open', blocked=unknown||busy||!!node?.offline;
    const reason=unknown?'Cannot verify operation status':node?.offline?`${node.id} is offline — it can't be restarted until it reconnects`:busy?'Restart in progress by fixture':'';
    const current=ctx.main.querySelector('[data-operations]');
    if(current) current.remove();
    const root=el('div','',ctx.main,{'data-operations':'',class:'pagehead'});
    if(nodePage) {
      ctx.main.querySelector('[title^="Connectivity:"]')?.setAttribute('title',`Connectivity: ${node.offline?'Offline':'Online'}`);
      el('span','Operational',el('div','',root,{class:'detail-subrow'}));
      const restart=el('button',busy?'Restarting… (1s)':'Restart node',root,{'aria-disabled':String(blocked),title:reason});
      restart.onclick=()=>{if(blocked)return;dialogOpen=true;confirm='';refresh(db);};
      const tab=el('button','Software',root,{role:'tab'});tab.onclick=()=>{software=true;refresh(db);};
      if(software){
        const card=el('div','example',root,{class:'app-card',role:'button','aria-label':'View example resources'}),status=node.software??'available';
        el('span',node.version??'1.0.0',el('div','Version: ',card),{class:'tnum'});
        const statusRow=el('div','',card);
        el('span',states[status]+(status==='available'||status==='failed'?' → 2.0.0':''),statusRow);
        if(status==='available'||status==='failed'){
          const update=el('button',status==='failed'?'Retry update':'Update Now',statusRow);update.disabled=blocked;
          update.onclick=()=>send('update',{app:'example',tag:'2.0.0'});
        }
      }
    } else {
      ctx.main.querySelector('.detail-subrow')?.remove();
      const status=el('div','',ctx.main,{class:'detail-subrow'});el('span',db.mode==='conflated-state'&&site.service===false?'Offline':'Online',status,{class:'MuiChip-label'});
      const b=el('button',busy?'Site actions • busy':'Site actions',root,{});b.onclick=()=>{menuOpen=true;refresh(db);};
      for(const name of ['Node','Switch']) {const b=el('button','',root,{class:'comp-tile'});el('div',name,b,{class:'comp-tile-label'});b.onclick=()=>{component=name;refresh(db);};}
      if(component==='Node')for(const n of nodes){
        if(db.mode==='missing-controller'&&n.type==='Controller node')continue;
        const card=el('div','',root,{class:'app-card'}),head=el('div','',card);
        el('span','',head);el('span',n.id,head);el('span',n.offline?'is offline':'is online and well',head);
        el('button','View node',el('div','',card)).onclick=()=>ctx.move('/network/nodes/'+n.id);
      }
      if(component==='Switch'){
        const controller=nodes.find(n=>n.type==='Controller node'), reason=controller?.offline?"Controller node is offline — ports can't be switched":unknown?'Cannot verify operation status':siteBusy?'Controller update in progress':'';
        for(const [port,name] of [[1,'Tnode PoE'],[2,'Cnode PoE'],[3,'Anode PoE'],[9,'Uplink SFP']]){
          const row=el('div','',root),head=el('div','',row);
          el('div',`Port ${port}: ${name}`,head);const label=el('label','',head,{title:reason});
          const on=(db.mode==='port-not-persisted'?localPorts[port]:site.ports?.[port])!==false;el('span',on?'On':'Off',label);
          const input=el('input','',label,{type:'checkbox'});input.checked=on;input.disabled=!!reason;
          input.onchange=()=>send('port',{port,value:input.checked});
        }
      }
    }
    ctx.menu.querySelector('[data-operations]')?.remove();
    const portal=el('div','',ctx.menu,{'data-operations':''});
    if(menuOpen&&sitePage){
      const menu=el('div','',portal,{role:'menu'});
      for(const label of ['Restart site','Radio','Cellular']){
        const towerOffline=nodes.find(n=>n.type==='Tower node')?.offline;
        const disabled=unknown||siteBusy||(label!=='Restart site'&&towerOffline);
        const item=el('div','',menu,{role:'menuitem','aria-disabled':String(disabled)});
        el('span',label,item,{class:'MuiListItemText-primary'});
        const controllerBusy=nodes.some(n=>n.type==='Controller node'&&n.busy);
        const why=unknown?'Cannot verify operation status':label!=='Restart site'&&towerOffline?"Tower node is offline — it can't be reached":siteBusy?(controllerBusy?'Controller update in progress':'Site restart in progress'):'';
        if(why)el('span',why,item,{class:'MuiListItemText-secondary'});
        if(label!=='Restart site'){
          const on=label==='Radio'?(site.radio??true):(site.service??true),input=el('input','',item,{type:'checkbox'});input.checked=on;input.disabled=disabled;
          el('span',on?'On':'Off',item,{});
          item.onclick=()=>{if(!disabled)void send(label==='Radio'?'radio':'service',{value:!on});};
        }else item.onclick=()=>{if(!disabled){menuOpen=false;dialogOpen=true;confirm='';refresh(db);}};
      }
    }
    if(dialogOpen){
      const dialog=el('div','',portal,{role:'dialog'});
      el('span',sitePage?'Restart site':'Restart node',dialog);
      if(sitePage){const input=el('input','',dialog);input.value=confirm;input.oninput=()=>{confirm=input.value;submit.disabled=blocked||confirm!==site.name;};}
      const cancel=el('button','Cancel',dialog);cancel.onclick=()=>{dialogOpen=false;refresh(db);};
      const submit=el('button',busy?(sitePage?'Site is busy':'Node is busy'):(sitePage?'Restart':'Restart node'),dialog);submit.disabled=blocked||(sitePage&&confirm!==site.name);
      submit.onclick=async()=>{if(submit.disabled)return;await send(sitePage?'restart-site':'restart-node');dialogOpen=false;refresh(db);};
    }
  }
  window.renderOperations = c => {
    if(path!==location.pathname){software=false;menuOpen=false;dialogOpen=false;confirm='';path=location.pathname;statusBusy=false;optimisticUntil=0;}
    ctx=c;refresh(c.db);void pollStatus();
  };
  // Do not rebuild the active textbox while the user types. Other controls
  // track server state continuously, matching live operation polling.
  let fingerprint;
  window.updateOperations = db => {
    const next=JSON.stringify([db.nodes,db.sites,db.mode]);
    if(next!==fingerprint){fingerprint=next;refresh(db);}
  };
  document.addEventListener('keydown',e=>{if(e.key==='Escape'&&menuOpen){menuOpen=false;refresh(ctx.db);}});
}
export async function operationsFixture(mode='') {
  const timers=new Set();
  const later=(fn,ms)=>{const timer=setTimeout(()=>{timers.delete(timer);fn();},ms);timers.add(timer);};
  const extension={
    script:`(${operationsClient.toString()})()`,
    graphql(body,db){
      if(!['GetNodeOperationStatus','GetSiteOperationStatus'].includes(body.operationName))return;
      const site=body.operationName==='GetSiteOperationStatus',id=body.variables[site?'siteId':'nodeId'];
      const nodes=db.nodes.filter(n=>site?n.site===id:n.id===id);
      return {data:{[site?'getSiteOperationStatus':'getNodeOperationStatus']:{busy:nodes.some(n=>n.busy)}}};
    },
    request(url,body,db){
      if(url!=='/operation')return false;
      db.operations.push({ui:body.action,id:body.id,app:body.app,tag:body.tag,value:body.value,port:body.port});
      const node=db.nodes.find(n=>n.id===body.id),site=db.sites.find(s=>s.id===body.id);
      const nodes=node?[node]:db.nodes.filter(n=>n.site===site?.id);
      if(body.action==='update'){
        node.busy=true;node.software='running';node.attempt=(node.attempt??0)+1;
        const attempt=node.attempt;
        later(()=>{if(node.attempt!==attempt||node.software!=='running')return;node.busy=false;node.software=mode==='update-failed'&&attempt===1?'failed':'done';if(node.software==='done')node.version=mode==='wrong-terminal-version'?'wrong-version':'2.0.0';},6000);
      } else if(body.action==='port'){
        site.ports??={};site.ports[body.port]=body.value;
        if(mode==='port-cross-talk')site.ports[body.port===1?2:1]=body.value;
      } else if(body.action==='radio'||body.action==='service'){
        nodes.forEach(n=>n.busy=true);
        later(()=>{site[body.action]=body.value;nodes.forEach(n=>n.busy=false);},500);
      } else {
        nodes.forEach(n=>n.busy=true);
        if(body.action==='restart-site'){
          for(const n of nodes)if(n.type!=='Controller node' || mode==='controller-outage')n.offline=true;
          if(mode==='no-restart-cycle')nodes.forEach(n=>n.offline=false);
          site.service=mode==='service-resumed';
        }
        later(()=>nodes.forEach(n=>{n.busy=false;n.offline=false;}),1500);
      }
      return true;
    },
    connectivity(node){if(node.offline&&node.software==='running'){node.software='failed';node.busy=false;}},
  };
  const fixture=await provisioningFixture(mode,extension),close=fixture.close;
  const cookie={name:'ukama_session',value:'fixture-primary',domain:'127.0.0.1',path:'/',expires:-1,httpOnly:true,secure:false,sameSite:'Lax'};
  fixture.state.cookies=[cookie];
  fixture.prepare=async dir=>{
    const state=join(dir,'peer.json');
    await writeFile(state,JSON.stringify({cookies:[{...cookie,value:mode==='same-peer'?'fixture-primary':'fixture-peer'}],origins:[]}));
    return {ULAB_WEBAPP_PEER_AUTH_STATE:state};
  };
  fixture.close=async()=>{for(const t of timers)clearTimeout(t);await close();};
  return fixture;
}

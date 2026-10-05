/* SPDX-License-Identifier: MPL-2.0
 * Controlled, source-shaped operation UI. Does not execute React or BFF code.
 * Product coverage still requires running the scenarios against the real stack.
 */
import { provisioningFixture } from './provisioning-fixture.mjs';
function operationsClient() {
  let ctx, software = false, menuOpen = false, dialogOpen = false, confirm = '', path;
  const send = async (action, data = {}) => {
    const id = location.pathname.split('/').at(-1);
    await fetch('/operation', { method:'POST', body:JSON.stringify({ action, id, ...data }) });
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
    const siteBusy=nodes.some(n=>n.busy), busy=nodePage?!!node.busy:siteBusy;
    const unknown=db.mode==='status-error', blocked=unknown||busy||!!node?.offline;
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
    }
    ctx.menu.querySelector('[data-operations]')?.remove();
    const portal=el('div','',ctx.menu,{'data-operations':''});
    if(menuOpen&&sitePage){
      const menu=el('div','',portal,{role:'menu'});
      for(const label of ['Restart site','Radio','Cellular']){
        const disabled=unknown||siteBusy;
        const item=el('div','',menu,{role:'menuitem','aria-disabled':String(disabled)});
        el('span',label,item,{class:'MuiListItemText-primary'});
        const controllerBusy=nodes.some(n=>n.type==='Controller node'&&n.busy);
        const why=unknown?'Cannot verify operation status':siteBusy?(controllerBusy?'Controller update in progress':'Site restart in progress'):'';
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
    if(path!==location.pathname){software=false;menuOpen=false;dialogOpen=false;confirm='';path=location.pathname;}
    ctx=c;refresh(c.db);
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
    request(url,body,db){
      if(url!=='/operation')return false;
      db.operations.push({ui:body.action,id:body.id,app:body.app,tag:body.tag,value:body.value});
      const node=db.nodes.find(n=>n.id===body.id),site=db.sites.find(s=>s.id===body.id);
      const nodes=node?[node]:db.nodes.filter(n=>n.site===site?.id);
      if(body.action==='update'){
        node.busy=true;node.software='running';node.attempt=(node.attempt??0)+1;
        const attempt=node.attempt;
        later(()=>{if(node.attempt!==attempt||node.software!=='running')return;node.busy=false;node.software=mode==='update-failed'&&attempt===1?'failed':'done';if(node.software==='done')node.version=mode==='wrong-terminal-version'?'wrong-version':'2.0.0';},6000);
      } else if(body.action==='radio'||body.action==='service'){
        nodes.forEach(n=>n.busy=true);
        later(()=>{site[body.action]=body.value;nodes.forEach(n=>n.busy=false);},500);
      } else {
        nodes.forEach(n=>n.busy=true);
        later(()=>nodes.forEach(n=>n.busy=false),1500);
      }
      return true;
    },
    connectivity(node){if(node.offline&&node.software==='running'){node.software='failed';node.busy=false;}},
  };
  const fixture=await provisioningFixture(mode,extension),close=fixture.close;
  fixture.close=async()=>{for(const t of timers)clearTimeout(t);await close();};
  return fixture;
}

/* SPDX-License-Identifier: MPL-2.0
 * Source-derived DOM/runtime fixture. Not console-app product coverage.
 */
import http from 'node:http';
import { once } from 'node:events';

function client() {
  let db, selected = 'existing', tower, siteName;
  const main = document.querySelector('main'), menu = document.querySelector('#menu');
  const button = (label, action, parent = main) => { const b = document.createElement('button'); b.textContent = label; b.onclick = action; parent.append(b); return b; };
  const move = async path => { history.pushState({}, '', path); await render(); };
  const graphql = async (op, data) => {
    const r = await fetch('/graphql', { method: 'POST', headers: { 'content-type': 'application/json' }, body: JSON.stringify({ operationName: op, query: `mutation ${op}($data: Input!) { ${op}(data: $data) { id name } }`, variables: { data } }) });
    return (await r.json()).data[op];
  };
  const read = async () => { db = await (await fetch('/state')).json(); };
  const kv = (label, value) => `<div class="kv-row"><span>${label}</span><span><span class="tnum">${value}</span></span></div>`;
  const name = () => db.networks.find(n => n.id === selected)?.name ?? 'existing-network';
  const params = () => new URLSearchParams(location.search);
  document.querySelectorAll('.viewseg button').forEach(b => b.onclick = () => move('/' + b.textContent.toLowerCase()));
  document.querySelector('.netswitch').onclick = () => {
    menu.innerHTML = '';
    for (const n of db.networks) {
      const b = button(n.name, () => { selected = n.id; menu.innerHTML = ''; void render(); }, menu); b.role = 'menuitem';
    }
    const b = button('Add network', () => {
      menu.innerHTML = '<div role="dialog"><h2>Add network</h2><input placeholder="network-name"></div>';
      const dialog = menu.firstChild;
      button('Create network', async () => {
        const n = await graphql('addNetwork', { name: dialog.querySelector('input').value, isDefault: false });
        if (db.mode === 'late-network') return;
        selected = n.id; menu.innerHTML = ''; await render();
      }, dialog);
    }, menu); b.role = 'menuitem';
  };
  async function render() {
    await read(); const path = location.pathname;
    document.querySelector('header').hidden = path.startsWith('/configure');
    document.querySelector('aside').hidden = path.startsWith('/configure');
    document.querySelector('.netswitch').hidden = path.includes('/manage/');
    document.querySelector('.nm').textContent = name();
    document.querySelector('aside').innerHTML = [['Home', '/network'], ['Sites', '/network/sites'], ['Nodes', '/network/nodes'], ['Node pool', '/network/manage/node-pool']].map(([n,p]) => `<a href="${p}">${n}</a>`).join('');
    document.querySelectorAll('aside a').forEach(a => a.onclick = e => { e.preventDefault(); void move(a.getAttribute('href')); });
    main.innerHTML = '';
    if (path === '/network' || path === '/business') {
      main.innerHTML = '<div class="MuiCard-root"><div><span>Sites online</span></div><div class="value"></div></div>'; updateKpi();
    } else if (path === '/network/manage/node-pool') {
      main.innerHTML = '<table><tbody></tbody></table>';
      for (const node of db.nodes.filter(n => n.type === 'Tower node' && !n.site)) {
        const row = document.createElement('tr'); row.innerHTML = `<td>${node.id}</td><td></td>`;
        button('Configure', () => { tower = db.mode === 'wrong-tower' ? 'other-tnode-999' : node.id; void move('/configure/select-network?flow=install-site&nid=' + tower); }, row.lastChild);
        main.querySelector('tbody').append(row);
      }
    } else if (path === '/network/sites') {
      const sites = db.sites.filter(s => s.network_id === selected);
      main.innerHTML = `<div class="pagehead"><div class="pagetitle">Sites<span class="cnt">${sites.length || ''}</span></div></div>`;
      if (!sites.length) main.innerHTML += '<div>No sites yet</div>';
      for (const site of sites) { const card = document.createElement('div'); card.className = 'ecard'; card.role = 'button'; card.innerHTML = `<span>${site.name}</span>`; card.onclick = () => move('/network/sites/' + site.id); main.append(card); }
    } else if (path.startsWith('/network/sites/')) {
      const site = db.sites.find(s => path.endsWith('/' + s.id)); main.innerHTML = `<h1>${site.name}</h1>`;
    } else if (path === '/network/nodes') {
      const nodes = db.nodes.filter(n => db.sites.some(s => s.id === n.site && s.network_id === selected));
      main.innerHTML = `<div class="pagehead"><div class="pagetitle">Nodes<span class="cnt">${nodes.length}</span></div></div>`;
      for (const node of nodes) {
        const card = document.createElement('div'); card.className = 'ecard'; card.role = 'button';
        card.innerHTML = `<div><span title="Connectivity: Online" class="dot"></span><span>${node.type}</span><div class="tnum">${node.type} · ${node.id}</div></div><hr><div><span>${db.sites.find(s=>s.id===node.site).name}</span></div>`;
        card.onclick = () => move('/network/nodes/' + node.id); main.append(card);
      }
    } else if (path.startsWith('/network/nodes/')) {
      const n = db.nodes.find(n => path.endsWith('/' + n.id));
      main.innerHTML = kv('Serial #', n.id) + kv('Model type', n.type) + kv('Site', db.sites.find(s=>s.id===n.site).name) + '<span title="Connectivity: Online" class="dot"></span>';
    } else if (path === '/configure/select-network') {
      main.innerHTML = '<h1>Select a network</h1><div role="radiogroup" aria-label="Network"></div>';
      for (const n of db.networks) { const b = button(n.name, () => { selected = n.id; }, main.lastChild); b.role = 'radio'; }
      button('Continue', () => move('/configure/install?networkid=' + selected + '&nid=' + tower));
    } else if (path === '/configure/install') {
      main.innerHTML = '<h1>Install your site</h1><label><input type="checkbox">I\'ve installed and powered on all my units</label>';
      button('Next', () => move('/configure/site?networkid=' + selected + '&nid=' + tower));
    } else if (path === '/configure/site') {
      main.innerHTML = '<h1>Name your site</h1><input placeholder="site-name">';
      button('Name site', () => { siteName = main.querySelector('input').value; void move('/configure/site/settings?networkid=' + selected + '&nid=' + tower); });
    } else if (path === '/configure/site/settings') {
      main.innerHTML = `<h1>Configure site settings</h1><div class="cfg-readonly">${params().get('nid')}</div>` + ['switch', 'backhaul', 'power'].map(f => `<select name="${f}Id"><option value="${f}-id">Default ${f}</option></select>`).join('');
      button('Create site', async () => {
        await graphql('addSite', { name: siteName, network_id: selected, access_id: tower, switch_id: 'switch-id', power_id: 'power-id', backhaul_id: 'backhaul-id' });
        if (db.mode === 'late-site') return;
        await move('/configure/sims');
      });
    } else if (path === '/configure/sims') { main.innerHTML = '<h1>Upload SIMs</h1>'; button('Finish setup', () => move('/configure/complete')); }
    else if (path === '/configure/complete') button('Go to Console', () => move('/network'));
  }
  function updateKpi() {
    const value = main.querySelector('.value'); if (!value) return;
    const sites = db.sites.filter(s=>s.network_id===selected);
    value.textContent = `${sites.filter(s=>!db.nodes.some(n=>n.site===s.id && n.offline)).length}/${sites.length}`;
  }
  setInterval(async () => { await read(); updateKpi(); }, 100);
  if (location.pathname === '/') history.replaceState({}, '', '/network');
  void render();
}
export async function provisioningFixture(mode = '') {
  const db = { mode, networks: [{ id: 'existing', name: 'existing-network' }], sites: [], nodes: [], operations: [], documents: 0 };
  const server = http.createServer(async (req,res) => {
    if (req.url === '/state') { res.setHeader('content-type','application/json'); res.end(JSON.stringify(db)); return; }
    if (req.method === 'POST') {
      let text = ''; for await (const chunk of req) text += chunk;
      const body = JSON.parse(text); let result = {};
      if (req.url === '/runtime') {
        db.operations.push({ runtime: body.action, args: body.args });
        if (body.action === 'build-and-start-site.sh') {
          for (const [kind,type] of [['tnode','Tower node'],['anode','Amplifier node'],['cnode','Controller node']]) db.nodes.push({ id: `uk-${kind}-${body.args[4]}`, type, site: null });
        }
        if (body.action === 'disconnect-node.sh' || body.action === 'reconnect-node.sh') {
          const n = db.nodes.find(n=>n.id === body.id); if(n) n.offline = body.action === 'disconnect-node.sh';
        }
      } else {
        const op = /\b(addNetwork|addSite|deleteNetwork|deleteSite|deleteNode|releaseNodeFromSite)\s*\(/.exec(body.query)?.[1];
        const data = body.variables?.data;
        const id = /id:\s*"([^"]+)"/.exec(body.query)?.[1];
        db.operations.push({ op, data, id });
        if (op === 'addNetwork') { const n = { id: `network-${db.networks.length}`, name: data.name }; db.networks.push(n); result = mode==='opaque-network' ? {id:n.id} : n; }
        if (op === 'addSite') { const s = { id: `site-${db.sites.length+1}`, name:data.name, network_id:data.network_id }; db.sites.push(s); for (const n of db.nodes.filter(n=>n.id.replace(/-(a|c)node-/,'-tnode-')===data.access_id)) n.site = s.id; result=s; }
        if (op === 'releaseNodeFromSite') { result={success:true}; }
        if (op === 'deleteNode') { db.nodes=db.nodes.filter(n=>n.id!==id); result={id}; }
        if (op === 'deleteSite') { db.sites=db.sites.filter(n=>n.id!==id); result={success:true}; }
        if (op === 'deleteNetwork') { db.networks=db.networks.filter(n=>n.id!==id); result={success:mode!=='cleanup-false'}; }
        result = { data: { [op]: result } };
      }
      res.setHeader('content-type','application/json'); res.end(JSON.stringify(result)); return;
    }
    db.documents++;
    res.setHeader('content-type','text/html; charset=utf-8'); res.end(`<!doctype html><html><style>.ecard,.MuiCard-root{border:1px solid;padding:8px;margin:8px}.dot{display:inline-block;width:8px;height:8px;background:green}button,a{margin:5px}.tnum{display:inline-block}</style><header class="topbar"><button class="netswitch"><span class="nm">existing-network</span></button><div class="viewseg"><button>Business</button><button>Network</button><button>Customer</button></div></header><aside class="sidebar"></aside><main class="main"></main><div id="menu"></div><script>(${client.toString()})()</script></html>`);
  });
  server.listen(0,'127.0.0.1'); await once(server,'listening');
  return { db, origin:`http://127.0.0.1:${server.address().port}`, close:()=>new Promise(resolve=>server.close(resolve)), state:{cookies:[],origins:[]} };
}

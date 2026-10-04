/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
// A local DOM fixture derived from console source, not the console app or BFF.
// Its runs only qualify worker behavior. Never count them as product coverage.
import http from 'node:http';
import { once } from 'node:events';

export async function fixture() {
  const metrics = { documents: 0, reloads: 0 };
  const server = http.createServer((req, res) => {
    if (req.url === '/private-error?token=do-not-record') { res.writeHead(503); res.end('private'); return; }
    if (!req.headers.cookie?.includes('fixture_session=valid') && !req.url?.startsWith('/login')) {
      res.writeHead(302, { location: '/login?token=do-not-record' }); res.end(); return;
    }
    if (req.url?.startsWith('/login')) { res.end('<h1>Sign in</h1>'); return; }
    metrics.documents++;
    res.setHeader('Content-Type', 'text/html');
    res.end(`<!doctype html><html lang="en"><head><title>Worker DOM fixture</title>
      <style>.card,.MuiCard-root{border:1px solid #888;padding:12px;margin:8px} button,a{margin:6px} aside{display:block} main{display:block}.tnum{display:inline-block}</style></head>
      <body><header class="topbar"><button class="netswitch"><span class="nm">other-network</span></button>
      <span class="orgchip" hidden>Org-wide</span><div class="viewseg"><button>Business</button><button>Network</button><button>Customer</button></div></header>
      <aside class="sidebar"></aside><main class="main"></main><div id="menu"></div>
      <script>
        let selected=localStorage.getItem('fixture_network')||'other-network';
        document.querySelector('.nm').textContent=selected;
        const main=document.querySelector('main');
        const nav={business:[['Home','/business'],['Revenue','/business/revenue'],['Customers','/business/customers'],['Packages','/business/packages'],['Data plans','/business/manage/data-plans'],['Members','/business/manage/members'],['SIM pool','/business/manage/sim-pool'],['Support','/business/support'],['Settings','/business/settings']],
          network:[['Home','/network'],['Sites','/network/sites'],['Nodes','/network/nodes'],['Customers','/network/customers'],['Node pool','/network/manage/node-pool'],['SIM pool','/network/manage/sim-pool'],['Support','/network/support'],['Settings','/network/settings']],
          customer:[['Customers','/customer/customers'],['Data plans','/customer/data-plans'],['Settings','/customer/settings']]};
        function goto(path){history.pushState({},'',path);render();}
        document.querySelectorAll('.viewseg button').forEach(b=>b.onclick=()=>goto(b.textContent==='Customer'?'/customer/customers':'/'+b.textContent.toLowerCase()));
        document.querySelector('.netswitch').onclick=()=>{
          document.querySelector('#menu').innerHTML=['other-network','lab-network','empty-network','loading-network','wrong-network','duplicate-network','disabled-network','hidden-network','switch-away','delayed-network','table-network','bad-table-network'].map(n=>'<button role="menuitem"><span>'+n+'</span><span>Default network</span></button>').join('');
          document.querySelectorAll('[role=menuitem]').forEach(b=>b.onclick=()=>{selected=b.firstChild.textContent;localStorage.setItem('fixture_network',selected);document.querySelector('.nm').textContent=selected;document.querySelector('#menu').innerHTML='';if(/\\/network\\/(nodes|sites)\\/.+/.test(location.pathname))goto(location.pathname.split('/').slice(0,3).join('/'));else render();});
        };
        const kpi=value=>'<div class="MuiCard-root"><div><span>Sites online</span></div><div><span class="tnum">'+value+'</span></div></div>';
        function render(){
          const path=location.pathname;const lens=path.split('/')[1]||'business';
          const links=nav[lens]||nav.business;
          document.querySelector('aside').innerHTML=links.map(([n,p])=>'<a href="'+p+'">'+n+'</a>').join('');
          document.querySelectorAll('aside a').forEach(a=>a.onclick=e=>{e.preventDefault();goto(a.getAttribute('href'));});
          const manage=path.includes('/manage/');document.querySelector('.netswitch').hidden=manage;document.querySelector('.orgchip').hidden=!manage;
          main.innerHTML='<div class="page"><div class="pagehead"><div class="pagetitle">Home</div></div></div>';
          const p=main.firstChild;
          if(path==='/network'||path==='/business'){
            p.innerHTML+=kpi(selected==='wrong-network'?'1/2':selected==='delayed-network'?'…':'2/2');
            if(selected==='duplicate-network')p.innerHTML+=kpi('2/2');
            if(selected==='delayed-network')setTimeout(()=>{const v=main.querySelector('.tnum');if(v)v.textContent='2/2';},300);
            if(selected==='switch-away')setTimeout(()=>{document.querySelector('.nm').textContent='other-network';},100);
          }else if(path==='/network/nodes'||path==='/network/sites'){
            const nodes=path.endsWith('nodes');const title=nodes?'Nodes':'Sites';
            p.querySelector('.pagetitle').textContent=title;
            if(selected==='loading-network'){p.innerHTML+='<div class="MuiSkeleton-root">Loading</div>';return;}
            if(selected==='empty-network'){p.innerHTML+='<div>No '+title.toLowerCase()+' yet</div>';return;}
            p.querySelector('.pagetitle').innerHTML+='<span class="cnt tnum">'+(nodes?'3':'1')+'</span>';
            const ids=nodes?['tower-001','amp-001','ctrl-001']:['site-001'];
            for(const id of ids){const card=document.createElement('div');card.className='card ecard';card.setAttribute('role','button');card.innerHTML='<span>Fixture</span><div class="tnum">Unit · '+id+'</div>';card.onclick=()=>goto(path+'/'+id);p.appendChild(card);}
          }else if(path.startsWith('/network/nodes/')){
            p.innerHTML+='<div class="kv-row"><span>Serial #</span><span><span class="tnum">'+path.split('/').pop()+'</span></span></div>';
            p.innerHTML+='<div class="kv-row"><span>Empty field</span><span><span class="tnum"></span></span></div>';
            if(selected!=='hidden-network')p.innerHTML+='<button '+(selected==='disabled-network'?'disabled':'')+'>Restart node</button>';
          }else if(path==='/business/manage/members'){
            p.innerHTML+='<table><thead><tr><th>Member</th></tr></thead><tbody>'+(selected==='table-network'?'<tr><td>One</td></tr><tr><td>Two</td></tr>':'')+'</tbody></table>';
            if(selected==='bad-table-network')p.innerHTML='<div class="MuiSkeleton-root">Loading</div>';
          }
        }
        if(location.pathname==='/')history.replaceState({},'','/business');render();
        console.error('Sensitive token must not appear in diagnostics: do-not-record');
        fetch('/private-error?token=do-not-record');
      </script></body></html>`);
  });
  server.listen(0, '127.0.0.1'); await once(server, 'listening');
  return { origin: `http://127.0.0.1:${server.address().port}`, metrics,
    close: () => new Promise(resolve => server.close(resolve)),
    state: { cookies: [{ name: 'fixture_session', value: 'valid', domain: '127.0.0.1', path: '/', expires: -1, httpOnly: true, secure: false, sameSite: 'Lax' }], origins: [] } };
}

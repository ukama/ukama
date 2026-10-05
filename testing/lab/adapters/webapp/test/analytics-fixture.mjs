/* SPDX-License-Identifier: MPL-2.0 — source-shaped SVG/DOM, not the React app. */
import {provisioningFixture} from './provisioning-fixture.mjs';
export function chartMarkup(title,kind='metric',mode='') {
 const unit=kind==='revenue'?'$':title==='Battery charge'?'%':'°C';
 const state=title==='Memory'?'empty':title==='Backhaul downlink'?'error':'data';
 const text=state==='empty'?'No data':"Couldn't load metric";
 const curve=mode==='zero-gap'?'M72,150L152,90L232,190L312,90L392,50':'M72,150L152,90M312,90L392,50';
 return `<div class="card card-pad"><div class="sec-head"><div class="sec-title">${title}</div><div role="group" aria-label="Time range"><button aria-pressed="true">Day</button><button aria-pressed="false">Week</button><button aria-pressed="false">Month</button></div></div><div>${state!=='data'?`<div>${text}</div>`:`<div class="recharts-wrapper"><svg class="recharts-surface" width="440" height="220"><defs><clipPath id="clip-${title.replaceAll(' ','-')}"><rect x="72" y="10" width="320" height="180"/></clipPath></defs><g class="recharts-xAxis"><text class="recharts-cartesian-axis-tick-value" x="72" y="215">12:00 PM</text><text class="recharts-cartesian-axis-tick-value" x="232" y="215">12:30 PM</text><text class="recharts-cartesian-axis-tick-value" x="365" y="215">1:00 PM</text></g><g class="recharts-yAxis"><text class="recharts-cartesian-axis-tick-value" x="0" y="190">0${unit}</text><text class="recharts-cartesian-axis-tick-value" x="0" y="90">80${unit}</text></g><path class="recharts-line-curve" fill="none" stroke="green" d="${curve}"/></svg><div class="recharts-tooltip-wrapper" style="visibility:hidden"></div></div><div><span>Below 60 ${unit}: Normal</span><span>60–80 ${unit}: High</span><span>Above 80 ${unit}: Critical</span></div>`}</div></div>`;
}
function install() {
 window.renderOperations=async({db,main,menu,button,render,move})=>{
  const path=location.pathname;
  const link=(name,path)=>{const a=document.createElement('a');a.href=path;a.textContent=name;a.onclick=e=>{e.preventDefault();move(path)};document.querySelector('aside').append(a)};
  if(path.startsWith('/business')){document.querySelector('aside').innerHTML='';link('Home','/business');link('Revenue','/business/revenue');link('Packages','/business/packages')}
  if(path==='/business/revenue')main.innerHTML=chartMarkup('Revenue trend','revenue',db.mode);
  if(path.startsWith('/network/sites/')) {
   for(const name of ['Batteries','Charge controller','Backhaul']){const b=button(name,()=>show(name==='Batteries'?'Battery charge':name==='Backhaul'?'Backhaul downlink':'Temperature'));b.className='comp-tile';b.innerHTML=`<div class="comp-tile-label">${name}</div>`}
   const panel=document.createElement('div');panel.id='charts';main.append(panel);show('Battery charge');
  }
  if(path.startsWith('/network/nodes/')){
   const tabs=document.createElement('div');tabs.role='tablist';main.append(tabs);
   for(const name of ['Overview','Resources','Software'])button(name,()=> name==='Software'?apps():show(name==='Resources'?'CPU':'Temperature'),tabs).role='tab';
   const rail=document.createElement('div');rail.innerHTML='<div class="card card-selectable" role="button"><div class="sec-title">Node health</div></div><div class="kv-row"><span>Temperature</span><span><span class="tnum">'+(db.mode==='wrong-unit'?'55.50 F':'55.50 °C')+'</span></span></div><div class="kv-row"><span>Memory</span><span><span class="tnum">—</span></span></div>';main.append(rail);
   rail.querySelector('[role="button"]').onclick=()=>show('Temperature');
   const panel=document.createElement('div');panel.id='charts';main.append(panel);show('Temperature');
  }
  function show(title){document.querySelector('#charts').innerHTML=chartMarkup(title,'metric',db.mode)+chartMarkup('Memory','metric',db.mode);bind()}
  function apps(){document.querySelector('#charts').innerHTML='';for(const name of ['metrics','controller']){const b=button(`View ${name} resources`,()=>{const d=document.createElement('div');d.role='dialog';d.innerHTML=`<h2>${db.mode==='wrong-app'?'foreign':name}</h2>`;for(const [k,v] of [['CPU','12.5%'],['Memory (RSS)','64.0 MB'],['Disk read','2.0 KB'],['Disk write','—']])d.innerHTML+=`<div><span>${k}</span><span class="tnum">${v}</span></div>`;button('Close',()=>d.remove(),d);menu.append(d)},document.querySelector('#charts'));b.className='app-card'}}
  function bind(){document.querySelectorAll('.card').forEach(card=>{
   card.querySelectorAll('[aria-pressed]').forEach(b=>b.onclick=()=>{card.querySelectorAll('[aria-pressed]').forEach(o=>o.setAttribute('aria-pressed',String(o===b))) });
   const svg=card.querySelector('svg');if(!svg)return;const tooltip=card.querySelector('.recharts-tooltip-wrapper'),title=card.querySelector('.sec-title').textContent;
   svg.onmousemove=e=>{const x=e.clientX-svg.getBoundingClientRect().x;const missing=x>192&&x<272;tooltip.style.visibility=missing&&db.mode!=='zero-gap'?'hidden':'visible';const time=x<192?'12:15 PM':'1:00 PM';tooltip.textContent=title==='Revenue trend'?'Oct 5 Revenue: $25.5':`Mon, Oct 5, 2026, ${db.mode==='wrong-time'?'9:15 PM':time} ${title}: ${db.mode==='wrong-tooltip'?'0.00':'55.50'} ${title==='Battery charge'?'%':'°C'} (normal)`};svg.onmouseleave=()=>tooltip.style.visibility='hidden';
  })}
  bind();
 };
}
export const analyticsFixture=mode=>provisioningFixture(mode,{script:`${chartMarkup.toString()};(${install.toString()})()`});

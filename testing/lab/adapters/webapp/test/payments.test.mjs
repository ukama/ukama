/* SPDX-License-Identifier: MPL-2.0 */
import test from 'node:test';
import assert from 'node:assert/strict';
import {EventEmitter} from 'node:events';
import {Budget} from '../dist/contract.js';
import {trackCommerceReads} from '../dist/commerce-faults.js';
import {paymentWindow,paymentDate,rejectPayment,paymentRejection} from '../dist/payments.js';
class Page extends EventEmitter {
 url(){return 'http://console.test/customer/customers'}
 async route(_,h){this.handler=h} async unroute(){this.handler=undefined}
 getByText(){return {filter(){return {async waitFor(){}}}}}
 getByRole(){return {filter(){return {async isVisible(){return true}}}}}
}
function page(){const p=new Page();trackCommerceReads(p);p.emit('request',{method:()=> 'POST',url:()=> 'http://bff.test/graphql',postDataJSON:()=>({operationName:'isPackageNameAvailable',query:'query isPackageNameAvailable($name: String!) {}',variables:{name:'owned'}})});return p;}
const inputs={plan:{id:'plan-1',amount:25,currency:'USD'},customer:{sim_id:'sim-1',email:'owned@example.test'}};
const request=(data={},url='http://bff.test/graphql',query='mutation addPayment($data: Input!) { addPayment(data: $data) { id } }')=>({url:()=>url,method:()=> 'POST',postDataJSON:()=>({query,variables:{data:{itemId:'plan-1',sim:'sim-1',payerEmail:'owned@example.test',amount:'25',currency:'USD',...data}}})});
const route=req=>({result:'',request:()=>req,async fallback(){this.result='forwarded'},async fulfill(){this.result='rejected'},async abort(){this.result='blocked'}});
test('receipt dates use the local submission window and reject invalid calendar/time values',()=>{
 const p=new Page(),now=Date.UTC(2026,9,5,12,30);paymentWindow(p,'payment-1',now,now+1000);
 assert.equal(paymentDate(p,'payment-1','05 Oct 2026, 12:30 UTC'),'within submission window');
 assert.equal(paymentDate(p,'payment-1','05 Oct 2025, 12:30 UTC'),'outside submission window');
 for(const date of ['31 Feb 2026, 12:30 UTC','05 Oct 2026, 99:30 UTC','05 Oct 2026, 12:60 UTC','—'])assert.equal(paymentDate(p,'payment-1',date),'unavailable');
 assert.equal(paymentDate(p,'foreign','05 Oct 2026, 12:30 UTC'),'unavailable');
});
test('controlled rejection never forwards an unrecognized or wrongly scoped write',async()=>{
 for(const req of [request({sim:'other'}),request({itemId:'other'}),request({amount:'250'}),request({},'http://foreign.test/graphql'),request({},undefined,'mutation updatePackage($data: Input!) {}')]){
  const p=page(),r=route(req);await assert.rejects(rejectPayment(p,inputs,()=>p.handler(r),new Budget(120)));
  assert.equal(r.result,'blocked');assert.equal(paymentRejection(p,'sim-1','plan-1'),'not applied');assert.equal(p.handler,undefined);
 }
});
test('controlled rejection proves one exact attempt, passes reads and rejects repeat attempts',async()=>{
 const p=page(),read=route(request({},undefined,'query getPackages($networkId: String) {}')),r=route(request());
 await rejectPayment(p,inputs,async()=>{await p.handler(read);await p.handler(r)},new Budget(500));
 assert.equal(read.result,'forwarded');assert.equal(r.result,'rejected');assert.equal(paymentRejection(p,'sim-1','plan-1'),'applied');assert.equal(p.handler,undefined);
 const q=page();await assert.rejects(rejectPayment(q,inputs,async()=>{await q.handler(route(request()));await q.handler(route(request()))},new Budget(500)),/duplicate writes/);
});

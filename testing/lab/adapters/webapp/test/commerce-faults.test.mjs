/* SPDX-License-Identifier: MPL-2.0 */
import test from 'node:test';
import assert from 'node:assert/strict';
import {EventEmitter} from 'node:events';
import {Budget} from '../dist/contract.js';
import {trackCommerceReads,commerceFault,commerceFaultState} from '../dist/commerce-faults.js';
class Page extends EventEmitter {
 path='http://console.test/business/manage/data-plans';
 url(){return this.path} mainFrame(){return this}
 async route(_,handler){this.handler=handler} async unroute(){this.handler=undefined}
}
const request=(extra={})=>({method:()=> 'POST',url:()=> 'http://bff.test/graphql',postDataJSON:()=>({operationName:'isPackageNameAvailable',query:'query isPackageNameAvailable($name: String!) { isPackageNameAvailable(name: $name) { isAvailable } }',variables:{name:'owned'}}),...extra});
const route=req=>({result:'',request:()=>req,async fallback(){this.result='fallback'},async fulfill(){this.result='fulfilled'},async abort(){this.result='aborted'}});
test('commerce faults require a passively observed unique endpoint and preserve foreign requests',async()=>{
 const p=new Page();trackCommerceReads(p);p.emit('request',request());
 await commerceFault(p,'name_failure','owned',new Budget(500));
 for(const req of [request({url:()=> 'http://foreign.test/graphql'}),request({method:()=> 'GET'}),request({postDataJSON:()=>({operationName:'isPackageNameAvailable',query:'mutation isPackageNameAvailable($name: String!) {}',variables:{name:'owned'}})}),request({postDataJSON:()=>({operationName:'isPackageNameAvailable',query:'query isPackageNameAvailable($name: String!) {}',variables:{name:'other'}})})]){
  const r=route(req);await p.handler(r);assert.equal(r.result,'fallback');assert.equal(commerceFaultState(p),'name_failure:armed');
 }
 const r=route(request());await p.handler(r);assert.equal(r.result,'fulfilled');assert.equal(commerceFaultState(p),'name_failure:applied');
 p.emit('framenavigated',p);assert.equal(commerceFaultState(p),'name_failure:armed');
 p.path='http://console.test/network';const off=route(request());await p.handler(off);assert.equal(off.result,'fallback');
 await commerceFault(p,'clear_fault',undefined,new Budget(500));assert.equal(commerceFaultState(p),'none');
});
test('ambiguous endpoints fail closed and pending requests release on close',async()=>{
 const p=new Page();trackCommerceReads(p);p.emit('request',request());p.emit('request',request({url:()=> 'http://other.test/graphql'}));
 await assert.rejects(commerceFault(p,'name_failure','owned',new Budget(500)),/multiple endpoints/);
 const q=new Page();trackCommerceReads(q);q.emit('request',request());await commerceFault(q,'name_pending','owned',new Budget(500));
 const r=route(request()), task=q.handler(r);assert.equal(commerceFaultState(q),'name_pending:applied');q.emit('close');await task;assert.equal(r.result,'aborted');
});

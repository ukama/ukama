/* SPDX-License-Identifier: MPL-2.0
 * Receipt acceptance reads the rendered dialog and the browser's real download.
 */
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import {writeFile,chmod} from 'node:fs/promises';
import {join} from 'node:path';
import type {Page,Route,Locator} from 'playwright';
import {Budget,WorkerError,normalize,type ObjectValue} from './contract.js';
import {commerceEndpoint} from './commerce-faults.js';
const exec=promisify(execFile);
const windows=new WeakMap<Page,Map<string,{start:number;end:number}>>();
const pdfs=new WeakMap<Page,{id:string;plan:string;path:string;fields:Record<string,string>}>();
const files=new WeakMap<Page,string[]>();
const rejected=new WeakMap<Page,{sim:string;plan:string;count:number}>();
const escape=(s:string)=>s.replace(/[.*+?^${}()|[\]\\]/g,'\\$&');
export function paymentWindow(page:Page,id:string,start:number,end:number) {
  const map=windows.get(page)??new Map();map.set(id,{start,end});windows.set(page,map);
}
export function paymentDate(page:Page,id:string,date:string):string {
  const w=windows.get(page)?.get(id);
  const match=/^(\d{2}) (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (\d{4}), (\d{2}):(\d{2}) UTC$/.exec(date);
  if(!w||!match)return 'unavailable';
  if(+match[1]!<1||+match[1]!>31||+match[4]!>23||+match[5]!>59)return 'unavailable';
  const ms=Date.UTC(+match[3]!,['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'].indexOf(match[2]!),+match[1]!,+match[4]!,+match[5]!);
  if(new Date(ms).getUTCDate()!==+match[1]!)return 'unavailable';
  return ms>=w.start-120000&&ms<=w.end+120000?'within submission window':'outside submission window';
}
export function paymentFiles(page:Page) {const list=files.get(page)??[];files.delete(page);return list;}
async function visible(target:Locator) {
  const v=target.filter({visible:true});
  if(await v.count()!==1)throw new WorkerError('RECEIPT_FIELD','Receipt field is missing or ambiguous');
  const value=normalize(await v.innerText());
  if(!value||value==='—')throw new WorkerError('RECEIPT_FIELD','Receipt field is empty');
  return value;
}
export async function receiptFields(page:Page,plan:string) {
  const d=page.getByRole('dialog').filter({visible:true});
  await visible(d.getByText('Payment receipt',{exact:true}));
  const meta=(label:string)=>d.getByText(label,{exact:true}).locator('..').locator(':scope > div').nth(1);
  const billed=d.getByText('Billed to',{exact:true}).locator('..');
  const status=d.getByText('Completed',{exact:true});
  return {
    id:await visible(meta('Payment ID')),number:await visible(meta('Receipt no')),date:await visible(meta('Paid on')),method:await visible(meta('Method')),
    payer:await visible(billed.locator(':scope > div').nth(1)),sim:await visible(billed.locator(':scope > div').nth(2)),
    status:await visible(status),organization:await visible(status.locator('..').locator(':scope > span').first()),
    plan:await visible(d.getByText(plan,{exact:true})),total:await visible(d.getByText('Total paid',{exact:true}).locator('..').locator(':scope > span').nth(1)),
    amount:await visible(d.getByText(plan,{exact:true}).locator('..').locator('..').locator(':scope > div').nth(1)),
    quantity:await visible(d.getByText('Data package · qty 1',{exact:true})),footer:await visible(d.getByText('Auto-generated · not a tax invoice',{exact:true}))
  };
}
export async function downloadReceipt(page:Page,plan:string,payment:string,directory:string,command:number,budget:Budget) {
  pdfs.delete(page);
  const fields=await receiptFields(page,plan);
  if(fields.id!==payment||fields.number!==payment.slice(0,8))throw new WorkerError('WRONG_RECEIPT','Receipt does not match the owned payment');
  const d=page.getByRole('dialog').filter({visible:true});
  const pending=page.waitForEvent('download',{timeout:budget.remaining()});
  // Attach a rejection handler immediately if the click itself fails.
  void pending.catch(()=>{});
  await d.getByRole('button',{name:'Download',exact:true}).click({timeout:budget.remaining()});
  const download=await pending;
  if(download.suggestedFilename()!==`receipt-${payment.slice(0,8)}.pdf`)throw new WorkerError('WRONG_DOWNLOAD','Unexpected receipt filename');
  let timer:ReturnType<typeof setTimeout>|undefined;
  const stream=await Promise.race([download.createReadStream(),new Promise<never>((_,reject)=>{
    timer=setTimeout(()=>{void download.cancel();reject(new WorkerError('DOWNLOAD_TIMEOUT','Receipt download exceeded its deadline'));},budget.remaining());
  })]).finally(()=>{if(timer)clearTimeout(timer);});
  if(!stream)throw new WorkerError('DOWNLOAD_FAILED','Receipt stream unavailable');
  const chunks:Buffer[]=[];let size=0;
  for await(const chunk of stream){budget.remaining();size+=chunk.length;if(size>5*1024*1024){stream.destroy();throw new WorkerError('DOWNLOAD_TOO_LARGE','Receipt exceeds 5 MiB');}chunks.push(Buffer.from(chunk));}
  const bytes=Buffer.concat(chunks);if(!bytes.subarray(0,5).equals(Buffer.from('%PDF-')))throw new WorkerError('INVALID_PDF','Download is not a PDF');
  const path=join(directory,`receipt-${command}.pdf`),textPath=join(directory,`receipt-${command}.txt`),png=join(directory,`receipt-${command}`);
  files.set(page,[path]);await writeFile(path,bytes,{mode:0o600,flag:'wx'});
  let text:string;
  try {
    const result=await exec('pdftotext',['-layout','-enc','UTF-8',path,'-'],{timeout:budget.remaining(),maxBuffer:2*1024*1024});text=result.stdout;
    await exec('pdftoppm',['-f','1','-singlefile','-scale-to','1200','-png',path,png],{timeout:budget.remaining(),maxBuffer:1024*1024});
  } catch {throw new WorkerError('PDF_TOOL_FAILED','Receipt extraction/rendering failed; install Poppler pdftotext and pdftoppm');}
  await writeFile(textPath,text,{mode:0o600,flag:'wx'});await chmod(png+'.png',0o600);files.get(page)!.push(textPath,png+'.png');
  const normalized=normalize(text).normalize('NFKC');
  for(const [label,value] of Object.entries(fields)) {
    const pattern=new RegExp(`(?:^|\\s)${escape(value.normalize('NFKC'))}(?:\\s|$)`);
    if(!pattern.test(normalized))throw new WorkerError('PDF_MISMATCH',`Downloaded receipt differs from visible ${label}`);
  }
  const section=(start:string,end:string)=>{
    const a=normalized.toLowerCase().indexOf(start.toLowerCase()),b=normalized.toLowerCase().indexOf(end.toLowerCase(),a+start.length);
    if(a<0||b<0||normalized.toLowerCase().indexOf(start.toLowerCase(),a+start.length)>=0)throw new WorkerError('PDF_MISMATCH','Missing or ambiguous receipt section');
    return normalized.slice(a+start.length,b).trim();
  };
  if(section('Total paid','Payment ID')!==fields.total || section('Payment ID',fields.footer)!==fields.id ||
     section('Billed to','Description')!==`${fields.payer} ${fields.sim}` ||
     section('METHOD','Billed to')!==`${fields.number} ${fields.date} ${fields.method}`)
    throw new WorkerError('PDF_MISMATCH','Downloaded receipt values are in different fields');
  // The dialog must still describe the same receipt after the asynchronous download.
  if(JSON.stringify(await receiptFields(page,plan))!==JSON.stringify(fields))throw new WorkerError('STALE_RECEIPT','Receipt changed during download');
  pdfs.set(page,{id:payment,plan,path,fields});
}
export async function receiptPdf(page:Page,plan:string,id:string) {
  const saved=pdfs.get(page);
  if(!saved||saved.id!==id||saved.plan!==plan)return 'not verified';
  return JSON.stringify(await receiptFields(page,plan))===JSON.stringify(saved.fields)?'matched':'stale';
}
// A controlled rejection never forwards a write. All unexpected writes are
// blocked too, so a stale selector/schema cannot create an unjournalled payment.
export async function rejectPayment(page:Page,inputs:ObjectValue,submit:()=>Promise<void>,budget:Budget) {
  const c=inputs.customer as ObjectValue,p=inputs.plan as ObjectValue;
  const path=page.url(),endpoint=commerceEndpoint(page);let invalid=false;
  const state={sim:String(c.sim_id),plan:String(p.id),count:0};rejected.set(page,state);
  const handler=async(route:Route)=>{
    const request=route.request();
    if(['GET','HEAD','OPTIONS'].includes(request.method())){await route.fallback();return;}
    let body:any;try{body=request.postDataJSON()}catch{}
    if(request.method()==='POST'&&typeof body?.query==='string'&&/^\s*query\b/.test(body.query)&&! /\bmutation\b/.test(body.query)){await route.fallback();return;}
    const data=body?.variables?.data;
    if(page.url()!==path||request.method()!=='POST'||request.url()!==endpoint||
      !/\bmutation\s+addPayment\s*\(/.test(body?.query??'')||! /\baddPayment\s*\(/.test(body?.query??'')||
      data?.itemId!==p.id||data?.sim!==c.sim_id||data?.payerEmail!==c.email||data?.amount!==String(p.amount)||data?.currency!==p.currency) {
      invalid=true;await route.abort('blockedbyclient');return;
    }
    await route.fulfill({status:200,contentType:'application/json',headers:{'access-control-allow-origin':new URL(page.url()).origin,'access-control-allow-credentials':'true'},body:JSON.stringify({errors:[{message:'Lab controlled payment rejection'}]})});state.count++;
  };
  await page.route('**/*',handler);
  try {
    await submit();
    await budget.poll(async()=>state.count,v=>v>0,'No owned payment request was rejected');
    if(invalid||state.count!==1)throw new WorkerError('REJECTION_SCOPE','Unexpected or duplicate writes attempted during controlled rejection');
    await page.getByText('Lab controlled payment rejection',{exact:true}).filter({visible:true}).waitFor({timeout:budget.remaining()});
    if(!await page.getByRole('dialog').filter({visible:true}).isVisible())throw new WorkerError('PAYMENT_ERROR_HIDDEN','Rejected payment closed its form');
    if(invalid||state.count!==1)throw new WorkerError('REJECTION_SCOPE','Unexpected or duplicate writes attempted during controlled rejection');
  } finally {await page.unroute('**/*',handler);}
}
export function paymentRejection(page:Page,sim:string,plan:string){const s=rejected.get(page);return s?.sim===sim&&s.plan===plan&&s.count===1?'applied':'not applied';}

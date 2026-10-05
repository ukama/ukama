/* SPDX-License-Identifier: MPL-2.0
 * Real matrix CLI -> C -> Chromium -> coverage report, controlled fixtures only.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp, mkdir, writeFile, readFile, copyFile, symlink, rm, chmod } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { spawn } from 'node:child_process';
import { commerceFixture } from '../../adapters/webapp/test/commerce-fixture.mjs';
const root=resolve(import.meta.dirname,'../..');
const call=(args,env=process.env)=>new Promise((resolve,reject)=>{const p=spawn('python3',args,{cwd:root,env});let output='';p.stdout.on('data',s=>output+=s);p.stderr.on('data',s=>output+=s);p.on('error',reject);p.on('close',code=>resolve({code,output}));});
test('matrix executes a fixture scenario and reports complete assertions without live credit',async t=>{
 const app=await commerceFixture(),dir=await mkdtemp(join(tmpdir(),'ulab-matrix-'));
 t.after(async()=>{await app.close();await rm(dir,{recursive:true,force:true})});
 for(const folder of ['utils/webapp','docs/webapp','scenarios/webapp/p0','adapters','scripts'])await mkdir(join(dir,folder),{recursive:true});
 for(const f of ['utils/webapp/matrix.py','utils/webapp-worker.sh'])await copyFile(join(root,f),join(dir,f));
 await chmod(join(dir,'utils/webapp-worker.sh'),0o700);
 await symlink(join(root,'adapters/webapp'),join(dir,'adapters/webapp'));
 const auth=join(dir,'auth.json');await writeFile(auth,JSON.stringify(app.state));
 const path='scenarios/webapp/p0/empty.yaml';
 const source=(await readFile(join(root,'scenarios/webapp/p0/expanded/wb-034-empty-customers.yaml'),'utf8')).replace('http://localhost:3000',app.origin).replace('.auth/owner.json',auth).replace('headless: false','headless: true').replace('scenario_timeout_seconds: 3600','scenario_timeout_seconds: 40');
 await writeFile(join(dir,path),source);
 await writeFile(join(dir,'docs/webapp/coverage.json'),JSON.stringify({requirements:[{id:'WEB-UI-001',requirement:'Partial empty state coverage',priority:'p0',automation:'planned',scenarios:[path]}]}));
 const out=join(dir,'matrix'),binary=process.env.ULAB_TEST_BINARY||join(root,'bin/ukama-lab');
 const result=await call([join(dir,'utils/webapp/matrix.py'),'run','--root',dir,'--binary',binary,'--out',out,'--browser','chromium','--evidence','fixture','--app-build','controlled-app','--backend-build','controlled-backend','--','--repo','/test/ukama','--scripts',join(dir,'scripts'),'--bff',app.origin+'/graphql','--factory-url',app.origin,'--warehouse-url',app.origin,'--asr-url',app.origin]);
 assert.equal(result.code,0,result.output);
 const manifest=JSON.parse(await readFile(join(out,'matrix.json'),'utf8'));assert.equal(manifest.attempts[0].exit_code,0);
 const reportDir=join(dir,'coverage');
 const report=await call([join(dir,'utils/webapp/matrix.py'),'report','--root',dir,'--out',reportDir,'--manifest',join(out,'matrix.json'),'--browser','chromium','--app-build','controlled-app','--backend-build','controlled-backend','--gate']);
 assert.equal(report.code,1,report.output);
 const coverage=JSON.parse(await readFile(join(reportDir,'coverage.json'),'utf8'));
 assert.equal(coverage.requirements[0].evidence[0].state,'passed',JSON.stringify(coverage));assert.equal(coverage.requirements[0].evidence[0].evidence,'fixture');assert.equal(coverage.totals.verified,0);assert.equal(coverage.gate_passed,false);
});

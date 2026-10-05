/* SPDX-License-Identifier: MPL-2.0
 * C runner + local Chromium + controlled runtime/BFF/DOM fixture.
 * This establishes integration behavior, not verified console-app coverage.
 */
import { mkdtemp, writeFile, readFile, rm, mkdir, copyFile, chmod } from 'node:fs/promises';
import { join, resolve } from 'node:path';
import { tmpdir } from 'node:os';
import { spawn } from 'node:child_process';
import { provisioningFixture } from '../../adapters/webapp/test/provisioning-fixture.mjs';
const root=resolve(import.meta.dirname,'../..');
const binary=process.env.ULAB_TEST_BINARY || join(root,'bin/ukama-lab');
export async function run(t, example, mode='', buildFail=false, fixture=provisioningFixture, transform=text=>text) {
  const app=await fixture(mode), dir=await mkdtemp(join(tmpdir(),'ulab-provision-'));
  t.after(async()=>{ await app.close(); if (!process.env.ULAB_KEEP_FIXTURE) await rm(dir,{recursive:true,force:true}); });
  const state=join(dir,'auth.json'); await writeFile(state,JSON.stringify(app.state));
  const scripts=join(dir,'scripts'); await mkdir(scripts);
  for(const name of ['start-media.sh','wait-media-ready.sh','start-ue.sh','wait-ues-attached.sh','traffic.sh','cleanup-ue.sh','ensure-network.sh','build-and-start-site.sh','wait-nodes-ready.sh','disconnect-node.sh','reconnect-node.sh','stop-node.sh','stop-media.sh','cleanup-network.sh']) {
    await writeFile(join(scripts,name), `#!/bin/sh
exec "${process.execPath}" "${join(root,'tests/webapp/runtime-fixture.mjs')}" "$0" "$@"
`); await chmod(join(scripts,name),0o700);
  }
  let text=(await readFile(join(root,'scenarios/webapp/p0',example.includes('/')?example:'network/'+example),'utf8')).replace('http://localhost:3000',app.origin).replace('.auth/owner.json',state).replace('headless: false','headless: true');
  text=text.replace('scenario_timeout_seconds: 3600','scenario_timeout_seconds: 40').replaceAll('timeout_seconds: 900','timeout_seconds: 15').replaceAll('timeout_seconds: 600','timeout_seconds: 15').replaceAll('timeout_seconds: 180','timeout_seconds: 15');
  text=transform(text);
  if(mode==='late-site'||mode==='late-network'||mode==='opaque-network') text=text.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 10').replaceAll('timeout_seconds: 30','timeout_seconds: 5').replaceAll('timeout_seconds: 15','timeout_seconds: 5').replaceAll('timeout_seconds: 900','timeout_seconds: 5');
  const file=join(dir,'scenario.yaml'); await writeFile(file,text);
  const runDir=join(dir,'test-run');
  const child=spawn(binary,['run',file,'--repo','/test/ukama','--scripts',scripts,'--warehouse-url',app.origin,'--factory-url',app.origin,'--asr-url',app.origin,'--bff',app.origin+'/graphql','--out',dir,'--run-id','test-run','--webapp-worker',join(root,'utils/webapp-worker.sh')],{cwd:root,env:{...process.env,ULAB_WEBAPP_AUTH_ORIGIN:app.authOrigin||app.origin,ULAB_WEBAPP_EXPECTED_CYCLE_USAGE:'64 MB of 1 GB used this cycle',ULAB_WEBAPP_EXPECTED_TOTAL_USAGE:'64 MB',ULAB_SOFTWARE_CURRENT_VERSION:'1.0.0',ULAB_SOFTWARE_TARGET_VERSION:'2.0.0',ULAB_CONTROLLER_RESTART_REASON:'Controller update in progress',ULAB_CONTROLLER_RF_REASON:'Controller update in progress',ULAB_CONTROLLER_SERVICE_REASON:'Controller update in progress',ULAB_FIXTURE_ORIGIN:app.origin,ULAB_FIXTURE_BUILD_FAIL:buildFail?'1':'0'}});
  let output=''; child.stdout.on('data',s=>output+=s); child.stderr.on('data',s=>output+=s);
  const timer=setTimeout(()=>child.kill('SIGTERM'),45000);
  const code=await new Promise((resolve,reject)=>{child.on('error',reject);child.on('close',resolve);});clearTimeout(timer);
  const report=JSON.parse(await readFile(join(runDir,'report.json'),'utf8').catch(()=>{throw new Error(output)}));
  const journal=JSON.parse(await readFile(join(runDir,'webapp-resources.json'),'utf8'));
  const commands=(await readFile(join(runDir,'webapp-commands.jsonl'),'utf8')).trim().split('\n').map(JSON.parse);
  return {app,code,report,journal,commands,output};
}

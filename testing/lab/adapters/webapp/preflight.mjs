/* SPDX-License-Identifier: MPL-2.0 — blank-page runtime probe, no credentials. */
import {chromium,firefox,webkit} from 'playwright';
const engines={chromium,firefox,webkit},names=process.argv.slice(2),results=[];
if(!names.length||new Set(names).size!==names.length||names.some(n=>!Object.hasOwn(engines,n)))process.exit(2);
for(const name of names){
 const specific=process.env[`ULAB_WEBAPP_${name.toUpperCase()}_EXECUTABLE_PATH`];
 if(names.length>1&&process.env.ULAB_WEBAPP_EXECUTABLE_PATH&&!specific){results.push({browser:name,state:'ambiguous_override'});continue}
 let browser;
 try{
  browser=await engines[name].launch({headless:true,timeout:10000,executablePath:specific||process.env.ULAB_WEBAPP_EXECUTABLE_PATH});
  const page=await browser.newPage({serviceWorkers:'block'});await page.setContent('<title>Local browser preflight</title>');
  results.push({browser:name,state:'ready',version:browser.version()});
 }catch{results.push({browser:name,state:'unavailable'})}
 finally{if(browser)await browser.close()}
}
process.stdout.write(JSON.stringify({browsers:results})+'\n');
process.exitCode=results.some(r=>r.state!=='ready')?1:0;

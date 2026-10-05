#!/usr/bin/env node
/* SPDX-License-Identifier: MPL-2.0
 * Controlled runtime double used only by provisioning integration tests.
 */
import { basename, join } from 'node:path';
import { mkdir, writeFile, readFile } from 'node:fs/promises';
const action=basename(process.argv[2]), args=process.argv.slice(3);
let id;
if (action==='build-and-start-site.sh') {
  const dir=join(args[3],'runtime-sites'); await mkdir(dir,{recursive:true});
  await mkdir(join(args[3],'runtime-nodes'),{recursive:true});
  await writeFile(join(dir,`${args[1]}.env`),['TNODE','CNODE','ANODE'].map(k=>`${k}_ID=uk-${k.toLowerCase()}-${args[4]}\n${k}_CONTAINER=container-${k}-${args[4]}`).join('\n')+'\n');
} else if (['disconnect-node.sh','reconnect-node.sh'].includes(action)) {
  const state=await readFile(join(args[1],'runtime-nodes',`${args[0]}.env`),'utf8');
  id=/^FACTORY_NODE_ID=(.+)$/m.exec(state)?.[1];
}
await fetch(process.env.ULAB_FIXTURE_ORIGIN+'/runtime',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({action,args,id})});
if (action==='build-and-start-site.sh' && process.env.ULAB_FIXTURE_BUILD_FAIL==='1') process.exitCode=1;

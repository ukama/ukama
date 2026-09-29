#!/usr/bin/env node
/* SPDX-License-Identifier: MPL-2.0
 * Development-only exporter: npm install graphql; node export-console.cjs SRC OUT
 * Runtime workload execution has no Node/npm dependency. */
const fs = require('node:fs');
const path = require('node:path');
const crypto = require('node:crypto');
const { parse, print, visit, Kind, buildSchema, validate } = require('graphql');
const [source, output, schemaFile] = process.argv.slice(2);
if (!source || !output) throw new Error('usage: export-console.cjs console/src/client/graphql workload/catalog/console.json');
const definitions = new Map();
const files = fs.readdirSync(source).filter(x => x.endsWith('.graphql') && x !== 'metrics.graphql').sort();
const sourceHashes = {};
for (const file of files) {
  const text = fs.readFileSync(path.join(source,file),'utf8');
  sourceHashes[file] = crypto.createHash('sha256').update(text).digest('hex');
  for (const def of parse(text).definitions) {
    if (def.name) definitions.set(def.name.value,def);
  }
}
// Pass a current BFF schema explicitly. The supplied console's offline schema
// predates its own GetKpiTimeSeries document and is not authoritative.
const schema = schemaFile ? buildSchema(fs.readFileSync(schemaFile,'utf8')) : null;
function document(name) {
  const op = definitions.get(name);
  if (!op || op.kind !== Kind.OPERATION_DEFINITION || op.operation !== 'query') throw new Error(`missing read operation ${name}`);
  const needed = new Set();
  function walk(def) {
    visit(def,{FragmentSpread(n) {
      if (!needed.has(n.name.value)) {
        needed.add(n.name.value);
        const fragment=definitions.get(n.name.value);
        if (!fragment) throw new Error(`missing fragment ${n.name.value}`);
        walk(fragment);
      }
    }});
  }
  walk(op);
  let doc={kind:Kind.DOCUMENT,definitions:[op,...[...needed].map(x=>definitions.get(x))]};
  // Apollo's default addTypename behavior for nested selections.
  doc=visit(doc,{SelectionSet(n,key,parent) {
    if (parent.kind===Kind.OPERATION_DEFINITION || n.selections.some(s=>s.kind===Kind.FIELD && s.name.value.startsWith('__'))) return;
    return {...n,selections:[...n.selections,{kind:Kind.FIELD,name:{kind:Kind.NAME,value:'__typename'}}]};
  }});
  if (schema) {
    const errors=validate(schema,doc);
    if(errors.length) throw new Error(`${name}: ${errors.map(x=>x.message).join('; ')}`);
  }
  return print(doc);
}
const required=(path,type)=>({path,type});
const net={networkId:'$networkId'};
const kpi=(keys,more={})=>({data:{keys,span:'last_30d',networkId:'$networkId',...more}});
const ops=[];
function add(id,name,variables,checks,allowed_errors=[]) {
  const query=document(name);
  ops.push({id,name,query,sha256:crypto.createHash('sha256').update(query).digest('hex'),variables,required:checks,allowed_errors});
}
add('sites-list','SitesList',net,[required('data.sitesView.sites.sites','array'),required('data.sitesView.nodeCounts.counts','array'),required('data.sitesView.customers.count','number')],[{path:'data.sitesView.kpis.error',code:'NOT_IMPLEMENTED'}]);
add('network-kpis','GetKpiValues',kpi(['NETWORK_UPTIME','ACTIVE_CUSTOMERS','DATA_USAGE','SITES_ONLINE']),[required('data.getKpiValues.values','array')]);
add('business-kpis','GetKpiValues',kpi(['REVENUE','CUSTOMERS','DATA_SOLD','NETWORK_UPTIME','SITES_ONLINE']),[required('data.getKpiValues.values','array')]);
add('revenue-kpis','GetKpiValues',kpi(['REVENUE','PAID_CUSTOMERS']),[required('data.getKpiValues.values','array')]);
add('revenue-count','GetKpiValues',kpi(['REVENUE'],{op:'COUNT'}),[required('data.getKpiValues.values','array')]);
add('revenue-average','GetKpiValues',kpi(['REVENUE'],{op:'AVG'}),[required('data.getKpiValues.values','array')]);
add('package-performance','GetPerformanceReport',{data:{report:'package_performance',span:'last_30d',networkId:'$networkId'}},[required('data.getPerformanceReport.rows','array')]);
add('revenue-series','GetKpiTimeSeries',kpi(['REVENUE'],{span:'daily',op:'SUM',from:'$from30d',to:'$mountedAt'}),[required('data.getKpiTimeSeries.values','array')]);
add('data-usage-kpis','GetKpiValues',kpi(['DATA_USAGE']),[required('data.getKpiValues.values','array')]);
add('customers','NetworkCustomers',net,[required('data.subscribersView.subscribers.subscribers','array'),required('data.subscribersView.plans.plans','array')]);
add('sim-usage','GetSimsUsageByNetwork',net,[required('data.getSimsUsageByNetwork','array')]);
add('site-detail','NetworkSiteDetail',{siteId:'$siteId'},[required('data.siteView.site.site','object'),required('data.siteView.nodes.nodes','array')]);
add('nodes-list','NodesList',net,[required('data.nodesView.nodes.nodes','array')],[{path:'data.nodesView.health.error',code:'NOT_IMPLEMENTED'}]);
add('node-detail','NodeDetail',{nodeId:'$nodeId'},[required('data.nodeView.node.node','object')]);
add('site-operation','GetSiteOperationStatus',{siteId:'$siteId'},[required('data.getSiteOperationStatus.busy','boolean')]);
add('node-operation','GetNodeOperationStatus',{nodeId:'$nodeId'},[required('data.getNodeOperationStatus.busy','boolean')]);
add('node-uptime','MetricsLast',{data:{keys:['$uptimeKey'],nodeId:'$nodeId'}},[required('data.metricsLast.metrics','array')]);
add('site-uptime','GetKpiTimeSeries',{data:{keys:['SITE_UPTIME'],span:'daily',networkId:'$networkId',siteId:'$siteId',from:'$from30d',to:'$mountedAt'}},[required('data.getKpiTimeSeries.values','array')]);
add('networks','getNetworks',{},[required('data.getNetworks.networks','array')]);
add('packages','getPackages',net,[required('data.getPackages.packages','array')]);
add('sim-pool','GetSimsFromPool',{data:{type:'ukama_data',status:'UNASSIGNED'}},[required('data.getSimsFromPool.sims','array')]);
fs.mkdirSync(path.dirname(output),{recursive:true});
fs.writeFileSync(output,JSON.stringify({format:1,source:'attached ukama-console snapshot 2026-09-28',sourceHashes,operations:ops},null,2)+'\n');
console.log(`Exported ${ops.length} operations${schema?' (schema validated)':''}`);

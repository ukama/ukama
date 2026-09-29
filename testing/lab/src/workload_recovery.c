/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include "util.h"
#include <stdlib.h>
#include <string.h>

/* Rehydrate only IDs recorded for this run. Never discover resources by a
 * broad prefix and never replay a create mutation after an ambiguous result. */
int workload_cleanup_run(const runner_opts_t *opts) {
    wl_config_t c={0}; wl_environment_t e={0}; ulab_error_t err={0};
    _Atomic sig_atomic_t stop=0;
    char path[ULAB_MAX_PATH],run_id[ULAB_MAX_ID];
    json_t *latest=json_object(); FILE *f=NULL;
    int rc=ULAB_ERR,complete=0,found_run=0;
    const char *dir=opts->scenario_path,*base=strrchr(dir,'/'); base=base?base+1:dir;
    if(!wl_safe_path(dir) || !*base || !wl_safe_path(opts->repo) || !wl_safe_path(opts->script_dir)) { wl_error(&err,"invalid recovery paths (no trailing slash)");goto done; }
    if(wl_path(path,sizeof(path),dir,"workload.resolved.json",&err) || wl_config_load(path,opts,&c,&err)) goto done;
    c.keep_environment=0;
    if(wl_path(path,sizeof(path),dir,"resources.jsonl",&err) || !(f=fopen(path,"r"))) {wl_error(&err,"cannot read resources.jsonl");goto done;}
    char *line=NULL;size_t cap=0;ssize_t length;
    while((length=getline(&line,&cap,f))>=0) {
        json_error_t je;json_t *r=json_loadb(line,(size_t)length,JSON_REJECT_DUPLICATES,&je);
        const char *kind=json_string_value(json_object_get(r,"kind"));
        if(!r || !kind) {json_decref(r);free(line);wl_error(&err,"invalid resource journal; inspect incomplete last line before cleanup");goto done;}
        if(!strcmp(kind,"run")) {
            const char *id=json_string_value(json_object_get(r,"run_id"));
            if(found_run || !id || strcmp(id,base) || ulab_copy(run_id,sizeof(run_id),id)) {json_decref(r);free(line);wl_error(&err,"journal/run directory identity mismatch");goto done;}
            found_run=1;
        } else if(!strcmp(kind,"cleanup")) {
            const char *state=json_string_value(json_object_get(r,"state"));complete=state && !strcmp(state,"complete");
        } else {
            json_t *n=json_object_get(r,"index");
            if(!json_is_integer(n) || json_integer_value(n)<0) {json_decref(r);free(line);wl_error(&err,"invalid resource index");goto done;}
            char key[128];snprintf(key,sizeof(key),"%s:%lld",kind,(long long)json_integer_value(n));
            json_object_set(latest,key,r);
        }
        json_decref(r);
    }
    free(line);fclose(f);f=NULL;
    if(!found_run) {wl_error(&err,"journal has no run identity");goto done;}
    if(complete) {printf("Cleanup already complete: %s\n",dir);rc=ULAB_OK;goto done;}
    pthread_mutex_init(&e.lock,NULL);e.initialized=1;e.config=&c;e.opts=opts;e.stop=&stop;
    strcpy(e.run_dir,dir);
    if(world_generate(c.environment,run_id,&e.world,&err)) goto done;
    const char *key;json_t *r;
    json_object_foreach(latest,key,r) {
        (void)key;
        const char *kind=json_string_value(json_object_get(r,"kind"));
        const char *id=json_string_value(json_object_get(r,"id"));
        const char *ref=json_string_value(json_object_get(r,"ref"));
        const char *state=json_string_value(json_object_get(r,"state"));
        size_t n=(size_t)json_integer_value(json_object_get(r,"index"));
        char *target=NULL;const char *expected=NULL;
#define RESOURCE(k,array,count) if(!strcmp(kind,k) && n<e.world.count) {target=e.world.array[n].bff_id;expected=e.world.array[n].ref;}
        RESOURCE("network",networks,network_count)
        RESOURCE("package",packages,package_count)
        RESOURCE("subscriber",subscribers,subscriber_count)
        RESOURCE("site",sites,site_count)
        RESOURCE("node",nodes,node_count)
        RESOURCE("ue",ues,ue_count)
#undef RESOURCE
        if(!target || !id || !ref || !state || strcmp(ref,expected) || ulab_copy(target,ULAB_MAX_ID,id)) {wl_error(&err,"invalid resource ownership entry");goto done;}
        if((!strcmp(state,"creating") || !strcmp(state,"allocating")) && !*id) {
            wl_error(&err,"unresolved %s %s: reconcile the create outcome and journal ID before cleanup; mutations are never replayed",kind,ref);goto done;
        }
        if(!strcmp(kind,"ue")) {
            ue_t *u=&e.world.ues[n];
#define FIELD(k,dest) do {const char *v=json_string_value(json_object_get(r,k));if(v && ulab_copy(dest,sizeof(dest),v)) {wl_error(&err,"oversized journal field");goto done;}} while(0)
            FIELD("iccid",u->iccid);FIELD("imsi",u->imsi);FIELD("pool_sim_id",u->pool_sim_id);FIELD("sim_package_id",u->sim_package_id);
#undef FIELD
            u->started=!strcmp(state,"attached") || !strcmp(state,"runtime_starting");
        }
    }
    if(runtime_init(&e.runtime,"virtual",opts->script_dir,dir,opts->repo) || runtime_load_workload_sites(&e.runtime,&e.world,&err)) goto done;
    e.runtime.execute=wl_environment_script;e.runtime.execute_ctx=&e;e.network_started=1;
    if(bff_init(&e.bff,opts->bff_url,dir) || !e.bff.authenticated) {wl_error(&err,"BFF authentication failed");goto done;}
    e.http.config=&c;e.http.auth=&e.bff;e.http.stop=&stop;e.http.phase="cleanup";
    e.bff.transport=wl_sync_transport;e.bff.transport_ctx=&e.http;
    if(wl_path(path,sizeof(path),dir,"resources.jsonl",&err) || !(e.journal=fopen(path,"a"))) {wl_error(&err,"cannot append cleanup journal");goto done;}
    rc=wl_environment_cleanup(&e,&err);
    printf("%s cleanup %s\n",rc?"FAIL":"PASS",dir);
done:
    if(err.msg[0]) fprintf(stderr,"%s\n",err.msg);
    if(f) fclose(f);
    wl_environment_close(&e);wl_config_free(&c);json_decref(latest);
    return rc;
}

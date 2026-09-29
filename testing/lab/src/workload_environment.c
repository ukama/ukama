/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include "sim_factory.h"
#include "util.h"
#include "log.h"
#include "report.h"
#include <errno.h>
#include <fcntl.h>
#include <spawn.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <unistd.h>
#include <ctype.h>

extern char **environ;
static void pause_briefly(void) { struct timespec t={0,20000000}; nanosleep(&t,NULL); }
pid_t wl_process_start(const char *scripts,const char *name,const char *args,
                       const char *log_path,ulab_error_t *err) {
    char script[ULAB_MAX_PATH], cmd[ULAB_MAX_QUERY*2];
    pid_t pid;
    posix_spawn_file_actions_t actions;
    posix_spawnattr_t attr;
    if(!wl_safe_path(scripts) || !wl_safe_path(name) || !wl_safe_path(log_path)) {
        wl_error(err,"workload runtime paths must contain only letters, digits, / _ . : -"); return -1;
    }
    for(const char *p=args;*p;p++) if(!isalnum((unsigned char)*p) && !strchr(" /_-.:=",*p)) {
        wl_error(err,"unsupported character in runtime argument"); return -1;
    }
    if(wl_path(script,sizeof(script),scripts,name,err) ||
       snprintf(cmd,sizeof(cmd),"exec %s %s",script,args)>=(int)sizeof(cmd)) return -1;
    char *argv[]={"sh","-c",cmd,NULL};
    posix_spawn_file_actions_init(&actions);
    posix_spawn_file_actions_addopen(&actions,STDOUT_FILENO,log_path,O_WRONLY|O_CREAT|O_TRUNC,0600);
    posix_spawn_file_actions_adddup2(&actions,STDOUT_FILENO,STDERR_FILENO);
    posix_spawnattr_init(&attr);
    posix_spawnattr_setflags(&attr,POSIX_SPAWN_SETPGROUP);
    posix_spawnattr_setpgroup(&attr,0);
    int rc=posix_spawnp(&pid,"sh",&actions,&attr,argv,environ);
    posix_spawnattr_destroy(&attr); posix_spawn_file_actions_destroy(&actions);
    if(rc) { wl_error(err,"cannot start %s: %s",name,strerror(rc)); return -1; }
    return pid;
}
static int process_wait(pid_t pid,double deadline,_Atomic sig_atomic_t *stop,
                        const char *script,ulab_error_t *err) {
    int status=0;
    double start=wl_now(),next_progress=start+WL_PROGRESS_SECONDS;
    for(;;) {
        pid_t rc=waitpid(pid,&status,WNOHANG);
        if(rc==pid) return WIFEXITED(status) && WEXITSTATUS(status)==0?ULAB_OK:wl_error(err,"runtime task failed (status=%d)",status);
        if(rc<0 && errno!=EINTR) return wl_error(err,"waitpid failed: %s",strerror(errno));
        if((stop && *stop) || wl_now()>deadline) {
            kill(-pid,SIGTERM);
            double grace=wl_now()+2;
            while(wl_now()<grace) { if(waitpid(pid,&status,WNOHANG)==pid) break; pause_briefly(); }
            kill(-pid,SIGKILL);
            while(waitpid(pid,&status,0)<0 && errno==EINTR) {}
            return wl_error(err,"runtime task %s",stop && *stop?"cancelled":"timed out");
        }
        if(script && wl_now()>=next_progress) {
            ulab_progress("runtime %s elapsed=%.0fs deadline_in=%.0fs",script,wl_now()-start,deadline-wl_now());
            next_progress=wl_now()+WL_PROGRESS_SECONDS;
        }
        pause_briefly();
    }
}
int wl_process_wait(pid_t pid,double deadline,_Atomic sig_atomic_t *stop,ulab_error_t *err) {
    return process_wait(pid,deadline,stop,NULL,err);
}
int wl_environment_script(void *ctx,const char *script,const char *args,ulab_error_t *err) {
    wl_environment_t *e=ctx;
    char file[256],path[ULAB_MAX_PATH];
    snprintf(file,sizeof(file),"env-%06llu-%s.log",(unsigned long long)++e->script_sequence,script);
    if(wl_path(path,sizeof(path),e->run_dir,file,err)) return ULAB_ERR;
    ulab_log_debug("runtime %s log=%s",script,path);
    pid_t pid=wl_process_start(e->runtime.script_dir,script,args,path,err);
    if(pid<0) return ULAB_ERR;
    int rc=process_wait(pid,wl_now()+e->config->job_timeout,e->stop,script,err);
    if(rc) ulab_status("ERROR","runtime %s failed; log=%s",script,path);
    return rc;
}

int wl_journal(wl_environment_t *e,const char *kind,size_t index,const char *state,ulab_error_t *err) {
    json_t *v=json_pack("{s:s,s:I,s:s,s:f}","kind",kind,"index",(json_int_t)index,"state",state,"monotonic_s",wl_now());
    const char *ref="",*id="";
    if(!strcmp(kind,"network")) { ref=e->world.networks[index].ref; id=e->world.networks[index].bff_id; }
    if(!strcmp(kind,"package")) { ref=e->world.packages[index].ref; id=e->world.packages[index].bff_id; }
    if(!strcmp(kind,"subscriber")) { ref=e->world.subscribers[index].ref; id=e->world.subscribers[index].bff_id; }
    if(!strcmp(kind,"site")) {
        site_t *s=&e->world.sites[index]; ref=s->ref; id=s->bff_id;
        json_object_set_new(v,"tower_id",json_string(s->tnode_id));
        json_object_set_new(v,"controller_id",json_string(s->cnode_id));
        json_object_set_new(v,"amplifier_id",json_string(s->anode_id));
    }
    if(!strcmp(kind,"node")) { ref=e->world.nodes[index].ref; id=e->world.nodes[index].bff_id; }
    if(!strcmp(kind,"ue")) {
        ue_t *u=&e->world.ues[index]; ref=u->ref; id=u->bff_id;
        json_object_set_new(v,"iccid",json_string(u->iccid)); json_object_set_new(v,"imsi",json_string(u->imsi));
        json_object_set_new(v,"pool_sim_id",json_string(u->pool_sim_id));
        json_object_set_new(v,"sim_package_id",json_string(u->sim_package_id));
    }
    json_object_set_new(v,"ref",json_string(ref)); json_object_set_new(v,"id",json_string(id));
    pthread_mutex_lock(&e->lock);
    int rc=!v || json_dumpf(v,e->journal,JSON_COMPACT) || fputc('\n',e->journal)==EOF || fflush(e->journal) || fsync(fileno(e->journal));
    pthread_mutex_unlock(&e->lock);
    json_decref(v);
    return rc?wl_error(err,"resource journal write failed"):ULAB_OK;
}
static int ambiguous(wl_environment_t *e) {
    return !strcmp(e->http.last_outcome,"timeout") || !strcmp(e->http.last_outcome,"transport_error") ||
        !strcmp(e->http.last_outcome,"http_error") || !strcmp(e->http.last_outcome,"cancelled") ||
        !strcmp(e->http.last_outcome,"invalid_response") || !strcmp(e->http.last_outcome,"graphql_error") ||
        !strcmp(e->http.last_outcome,"ok"); /* e.g. committed create with missing returned ID */
}
static int query(wl_environment_t *e,const char *name,const char *q,json_t *variables,json_t **out,ulab_error_t *err) {
    char *vars=json_dumps(variables,JSON_COMPACT);
    if(!vars) return wl_error(err,"out of memory");
    int rc=wl_sync_transport(&e->http,name,q,vars,out,err); free(vars); return rc;
}
static int confirm_site(wl_environment_t *e,site_t *site,ulab_error_t *err) {
    ulab_status("BACKEND","site %s addSite outcome=%s; reconcile existing site (no mutation replay)",site->ref,e->http.last_outcome);
    double end=wl_now()+240;
    const char *q="query GetSites($data: SitesInputDto!) { getSites(data:$data) { sites { id name } } }";
    json_t *vars=json_pack("{s:{s:s}}","data","networkId",e->world.networks[0].bff_id);
    int result=ULAB_ERR;
    while(!*e->stop && wl_now()<end) {
        json_t *root=NULL; ulab_error_t tmp={0};
        if(!query(e,"GetSites-confirm",q,vars,&root,&tmp)) {
            json_t *arr=json_object_get(json_object_get(json_object_get(root,"data"),"getSites"),"sites");
            for(size_t i=0;i<json_array_size(arr);i++) {
                json_t *s=json_array_get(arr,i);
                const char *name=json_string_value(json_object_get(s,"name")),*id=json_string_value(json_object_get(s,"id"));
                if(name && id && !strcmp(name,site->name) && !ulab_copy(site->bff_id,sizeof(site->bff_id),id)) result=ULAB_OK;
            }
        }
        json_decref(root);
        if(!result) break;
        double next=wl_now()+5; while(!*e->stop && wl_now()<next) pause_briefly();
    }
    json_decref(vars);
    if(result) { e->uncertain=1; return wl_error(err,"addSite outcome remains unknown for %s; no mutation replay",site->ref); }
    return ULAB_OK;
}
static int add_site(wl_environment_t *e,size_t index,ulab_error_t *err) {
    world_t *w=&e->world; site_t *s=&w->sites[index];
    size_t indices[3]={index*3,index*3+1,index*3+2};
    selector_result_t nodes={.idx=indices,.count=3};
    ulab_status("SITE","factory/build/start %s",s->ref);
    if(wl_journal(e,"site",index,"runtime_starting",err) || runtime_start_selected_site(&e->runtime,w,index,err)) return ULAB_ERR;
    ulab_status("SITE","%s tnode=%s cnode=%s anode=%s",s->ref,s->tnode_id,s->cnode_id,s->anode_id);
    for(size_t i=0;i<3;i++) {
        if(!wl_safe_path(w->nodes[indices[i]].id)) return wl_error(err,"invalid provider node identity");
        if(wl_journal(e,"node",indices[i],"registered",err)) return ULAB_ERR;
    }
    if(wl_journal(e,"site",index,"nodes_started",err) ||
       runtime_wait_nodes_ready(&e->runtime,w,&nodes,err) ||
       bff_wait_site_anchor_online(&e->bff,s,err) || wl_journal(e,"site",index,"creating",err)) return ULAB_ERR;
    ulab_status("SITE","add site %s",s->ref);
    if(bff_add_site(&e->bff,s,&w->networks[0],err)) {
        if(!ambiguous(e) || confirm_site(e,s,err)) return ULAB_ERR;
        err->msg[0]='\0';
    }
    if(wl_journal(e,"site",index,"created",err)) return ULAB_ERR;
    ulab_status("SITE","%s created id=%s",s->ref,s->bff_id);
    ulab_status("BACKEND","site %s wait for 3 nodes Online/Operational",s->ref);
    double end=wl_now()+240;
    size_t previous_ready=SIZE_MAX;
    for(;;) {
        size_t ready=0;
        for(size_t i=0;i<3;i++) {
            bff_node_status_t status={0};
            if(bff_get_node_status(&e->bff,&w->nodes[indices[i]],&status,err)) return ULAB_ERR;
            if(!strcmp(status.connectivity,"Online") && !strcmp(status.state,"Operational")) ready++;
        }
        if(ready!=previous_ready) {
            ulab_status("BACKEND","site %s nodes Online/Operational=%zu/3",s->ref,ready);
            previous_ready=ready;
        }
        if(ready==3) break;
        if(*e->stop || wl_now()>end) return wl_error(err,"site %s did not become operational",s->ref);
        double next=wl_now()+3; while(!*e->stop && wl_now()<next) pause_briefly();
    }
    if(e->config->service_enabled) {
        ulab_status("SERVICE","on %s",s->ref);
        if(bff_toggle_site_service(&e->bff,s,1,err)) return ULAB_ERR;
        size_t selected=index; selector_result_t sites={.idx=&selected,.count=1};
        if(runtime_wait_service_state(&e->runtime,w,&sites,1,err)) return ULAB_ERR;
    }
    if(e->config->environment->world.ues_per_site) {
        if(runtime_ensure_media(&e->runtime,err)) return ULAB_ERR;
        e->media_started=1;
    }
    int rc=wl_journal(e,"site",index,"ready",err);
    if(!rc) ulab_status("SITE","%s ready",s->ref);
    return rc;
}

static int provision(wl_environment_t *e,ulab_error_t *err) {
    ulab_status("SIM","provision batch count=%zu",e->batch_count);
    world_t view=e->world;
    ue_t *copies=calloc(e->batch_count,sizeof(*copies));
    char folder[ULAB_MAX_PATH],suffix[128],csv[ULAB_MAX_PATH];
    json_t *pool=NULL,*vars=NULL,*wanted=NULL;
    int rc=ULAB_ERR;
    if(!copies) return wl_error(err,"out of memory provisioning");
    for(size_t i=0;i<e->batch_count;i++) copies[i]=e->world.ues[e->batch[i]];
    view.ues=copies; view.ue_count=e->batch_count;
    snprintf(view.run_id,sizeof(view.run_id),"%.400s-b%06zu",e->world.run_id,++e->batch_sequence);
    snprintf(suffix,sizeof(suffix),"fixtures-%06zu",e->batch_sequence);
    if(wl_path(folder,sizeof(folder),e->run_dir,suffix,err) || ulab_mkdir_p(folder)) goto done;
    runner_opts_t factory_opts=*e->opts; factory_opts.workload_cancel=e->stop;
    if(*e->stop || sim_factory_prepare_world(&factory_opts,&view,folder,csv,sizeof(csv),err)) goto done;
    /* Persist identities before the first BFF mutation. */
    for(size_t i=0;i<e->batch_count;i++) {
        e->world.ues[e->batch[i]]=copies[i];
        if(wl_journal(e,"ue",e->batch[i],"factory_prepared",err)) goto done;
    }
    if(*e->stop || bff_upload_sims_from_csv(&e->bff,csv,e->opts->sim_type,err)) goto done;
    vars=json_pack("{s:{s:s,s:s}}","data","type",e->opts->sim_type,"status","UNASSIGNED");
    if(query(e,"GetSimsFromPool-setup","query GetSimsFromPool($data:GetSimsInput!){getSimsFromPool(data:$data){sims{id iccid}}}",vars,&pool,err)) goto done;
    json_t *arr=json_object_get(json_object_get(json_object_get(pool,"data"),"getSimsFromPool"),"sims");
    if(!json_is_array(arr)) { wl_error(err,"missing SIM pool list"); goto done; }
    wanted=json_object();
    for(size_t i=0;i<e->batch_count;i++) json_object_set_new(wanted,copies[i].iccid,json_integer((json_int_t)i));
    for(size_t i=0;i<json_array_size(arr);i++) {
        json_t *item=json_array_get(arr,i);
        const char *iccid=json_string_value(json_object_get(item,"iccid")),*id=json_string_value(json_object_get(item,"id"));
        json_t *slot=iccid?json_object_get(wanted,iccid):NULL;
        if(slot && id) {
            size_t n=(size_t)json_integer_value(slot);
            ulab_copy(copies[n].pool_sim_id,sizeof(copies[n].pool_sim_id),id);
            e->world.ues[e->batch[n]]=copies[n];
        }
    }
    json_decref(pool); pool=NULL;
    double next=wl_now();
    for(size_t i=0;i<e->batch_count;i++) {
        size_t n=e->batch[i]; ue_t *u=&e->world.ues[n];
        subscriber_t *sub=&e->world.subscribers[n];
        package_t *p=world_package_by_ref(&e->world,u->package_ref);
        while(!*e->stop && wl_now()<next) pause_briefly();
        next=wl_now()+1/e->config->provision_rate;
        if(*e->stop) { wl_error(err,"provisioning cancelled"); goto done; }
        if(!u->pool_sim_id[0] || !p) { wl_error(err,"prepared SIM/package missing for %s",u->ref); goto done; }
        if(wl_journal(e,"ue",n,"pool_uploaded",err) || wl_journal(e,"subscriber",n,"creating",err)) goto done;
        if(bff_add_subscriber(&e->bff,sub,&e->world.networks[0],err)) { if(ambiguous(e)) e->uncertain=1; goto done; }
        if(wl_journal(e,"subscriber",n,"created",err) || wl_journal(e,"ue",n,"allocating",err)) goto done;
        if(bff_allocate_sim_from_pool(&e->bff,u,sub,&e->world.networks[0],p,e->opts->sim_type,err)) { if(ambiguous(e)) e->uncertain=1; goto done; }
        if(wl_journal(e,"ue",n,"allocated",err)) goto done;
        int active=0;
        if(bff_get_packages_for_sim(&e->bff,u,p->bff_id,&active,err) || !active) {
            wl_error(err,"allocated SIM %s has no active package",u->ref); goto done;
        }
        if(wl_journal(e,"ue",n,"ready",err)) goto done;
    }
    rc=ULAB_OK;
done:
    if(!rc) ulab_status("SIM","provision batch ready count=%zu",e->batch_count);
    json_decref(pool); json_decref(vars); json_decref(wanted); free(copies); return rc;
}
static void *worker(void *ctx) {
    wl_environment_t *e=ctx;
    int rc=e->task_kind==1?add_site(e,e->site_index,&e->failure):provision(e,&e->failure);
    pthread_mutex_lock(&e->lock); e->result=rc; e->done=1; pthread_mutex_unlock(&e->lock);
    return NULL;
}
int wl_environment_open(wl_environment_t *e,wl_config_t *c,const runner_opts_t *opts,
                        wl_metrics_t *metrics,const char *run_dir,
                        _Atomic sig_atomic_t *stop,ulab_error_t *err) {
    char path[ULAB_MAX_PATH];
    const char *run_id=strrchr(run_dir,'/'); run_id=run_id?run_id+1:run_dir;
    memset(e,0,sizeof(*e)); pthread_mutex_init(&e->lock,NULL); e->initialized=1;
    e->config=c; e->opts=opts; e->metrics=metrics; e->stop=stop;
    strcpy(e->run_dir,run_dir);
    if(world_generate(c->environment,run_id,&e->world,err)) return ULAB_ERR;
    ulab_status("SCENARIO","%s",c->name);
    if(c->environment->description[0]) ulab_status("PURPOSE","%s",c->environment->description);
    report_world(&e->world);
    ulab_status("WORKLOAD","WORLD is planned capacity; initial sites=%zu provisioned=%zu attached=%zu",c->initial.sites,c->initial.provisioned,c->initial.attached);
    double planned_seconds=0;
    for(size_t i=0;i<c->phase_count;i++) planned_seconds+=c->phases[i].ramp+c->phases[i].hold;
    ulab_status("WORKLOAD","phases=%zu ramp+hold=%.0fs (%.1f min); setup, extra readiness waits and cleanup are additional",c->phase_count,planned_seconds,planned_seconds/60);
    ulab_status("OUTPUT","%s",run_dir);
    e->site_ready=calloc(e->world.site_count,1); e->provisioned=calloc(e->world.ue_count?:1,1);
    e->batch=calloc(c->provision_batch,sizeof(size_t));
    if(!e->site_ready || !e->provisioned || !e->batch) return wl_error(err,"out of memory");
    if(wl_path(path,sizeof(path),run_dir,"resources.jsonl",err)) return ULAB_ERR;
    e->journal=fopen(path,"a"); if(!e->journal) return wl_error(err,"cannot open resource journal");
    json_t *meta=json_pack("{s:s,s:s,s:s}","kind","run","run_id",run_id,"scenario",c->name);
    json_dumpf(meta,e->journal,JSON_COMPACT); fputc('\n',e->journal); fflush(e->journal); json_decref(meta);
    if(bff_init(&e->bff,opts->bff_url,run_dir)) return wl_error(err,"BFF authentication failed");
    if(!e->bff.authenticated || !e->bff.token[0]) return wl_error(err,"workload requires BFF authentication");
    e->http.config=c; e->http.auth=&e->bff; e->http.metrics=metrics; e->http.stop=stop; e->http.phase="setup";
    e->bff.transport=wl_sync_transport; e->bff.transport_ctx=&e->http;
    if(runtime_init(&e->runtime,"virtual",opts->script_dir,run_dir,opts->repo)) return wl_error(err,"runtime init failed");
    e->runtime.execute=wl_environment_script; e->runtime.execute_ctx=e;
    if(wl_journal(e,"network",0,"creating",err)) return ULAB_ERR;
    ulab_status("NETWORK","add network %s",e->world.networks[0].ref);
    if(bff_add_network(&e->bff,&e->world.networks[0],err)) { if(ambiguous(e)) e->uncertain=1; return ULAB_ERR; }
    if(wl_journal(e,"network",0,"created",err)) return ULAB_ERR;
    ulab_status("NETWORK","%s created",e->world.networks[0].ref);
    for(size_t i=0;i<e->world.package_count;i++) {
        if(wl_journal(e,"package",i,"creating",err)) return ULAB_ERR;
        ulab_status("PACKAGE","add package %s",e->world.packages[i].ref);
        if(bff_add_package(&e->bff,&e->world.packages[i],&e->world.networks[0],err)) { if(ambiguous(e)) e->uncertain=1; return ULAB_ERR; }
        if(wl_journal(e,"package",i,"created",err)) return ULAB_ERR;
    }
    if(runtime_ensure_network(&e->runtime,err)) return ULAB_ERR;
    e->network_started=1;
    if(c->environment->world.ues_per_site) {
        ulab_status("UE","prepare runtime image");
        char args[ULAB_MAX_ARGS];
        snprintf(args,sizeof(args),"%s %s",opts->repo,run_dir);
        if(wl_environment_script(e,"workload-prepare-ue.sh",args,err)) return ULAB_ERR;
    }
    return ULAB_OK;
}
void wl_environment_counts(wl_environment_t *e,wl_targets_t *a) {
    a->sites=a->provisioned=0;
    for(size_t i=0;i<e->world.site_count;i++) a->sites+=e->site_ready[i]!=0;
    for(size_t i=0;i<e->world.ue_count;i++) a->provisioned+=e->provisioned[i]!=0;
}
int wl_environment_step(wl_environment_t *e,const wl_targets_t *target,const char *phase,ulab_error_t *err) {
    if(e->running) {
        pthread_mutex_lock(&e->lock); int done=e->done; pthread_mutex_unlock(&e->lock);
        if(!done) return ULAB_OK;
        pthread_join(e->thread,NULL); e->running=0;
        if(e->result) { *err=e->failure; return ULAB_ERR; }
        if(e->task_kind==1) e->site_ready[e->site_index]=1;
        else for(size_t i=0;i<e->batch_count;i++) e->provisioned[e->batch[i]]=1;
    }
    wl_targets_t actual={0}; wl_environment_counts(e,&actual);
    if(*e->stop) return wl_error(err,"workload cancelled");
    e->task_kind=0;
    if(actual.sites<target->sites && wl_now()-e->last_site>=e->config->site_interval) {
        e->site_index=actual.sites; e->task_kind=1; e->last_site=wl_now();
    } else if(actual.provisioned<target->provisioned) {
        size_t need=target->provisioned-actual.provisioned;
        size_t per=e->config->environment->world.ues_per_site;
        e->batch_count=0;
        for(size_t slot=0;slot<per && e->batch_count<need && e->batch_count<e->config->provision_batch;slot++)
            for(size_t site=0;site<actual.sites && e->batch_count<need && e->batch_count<e->config->provision_batch;site++) {
                size_t i=site*per+slot;
                if(!e->provisioned[i]) e->batch[e->batch_count++]=i;
            }
        for(size_t i=e->world.site_count*per;i<e->world.ue_count && e->batch_count<need && e->batch_count<e->config->provision_batch;i++)
            if(!e->provisioned[i]) e->batch[e->batch_count++]=i;
        if(e->batch_count) e->task_kind=2;
    }
    if(!e->task_kind) return ULAB_OK;
    e->done=0; e->result=0; memset(&e->failure,0,sizeof(e->failure));
    ulab_copy(e->phase,sizeof(e->phase),phase); e->http.phase=e->phase;
    if(pthread_create(&e->thread,NULL,worker,e)) return wl_error(err,"cannot start environment worker");
    e->running=1;
    return ULAB_OK;
}
int wl_environment_verify(wl_environment_t *e,ulab_error_t *err) {
    ulab_status("VERIFY","check workload backend inventory");
    const char *queries[]={
        "query GetSites($data:SitesInputDto!){getSites(data:$data){sites{id}}}",
        "query GetSubscribersByNetwork($networkId:String!){getSubscribersByNetwork(networkId:$networkId){subscribers{uuid}}}",
        "query GetSimsByNetwork($networkId:String!){getSimsByNetwork(networkId:$networkId){sims{id}}}"
    };
    const char *fields[]={"getSites","getSubscribersByNetwork","getSimsByNetwork"};
    const char *lists[]={"sites","subscribers","sims"};
    e->http.phase="verification";
    for(size_t k=0;k<3;k++) {
        json_t *vars=k?json_pack("{s:s}","networkId",e->world.networks[0].bff_id):json_pack("{s:{s:s}}","data","networkId",e->world.networks[0].bff_id);
        json_t *root=NULL,*ids=json_object();
        int rc=query(e,fields[k],queries[k],vars,&root,err);json_decref(vars);
        if(rc) {json_decref(ids);return ULAB_ERR;}
        json_t *arr=json_object_get(json_object_get(json_object_get(root,"data"),fields[k]),lists[k]);
        if(!json_is_array(arr)) rc=wl_error(err,"verification: missing %s list",lists[k]);
        for(size_t i=0;i<json_array_size(arr);i++) {
            const char *id=json_string_value(json_object_get(json_array_get(arr,i),k==1?"uuid":"id"));
            if(id) json_object_set_new(ids,id,json_true());
        }
        size_t count=k==0?e->world.site_count:k==1?e->world.subscriber_count:e->world.ue_count,expected=0;
        for(size_t i=0;i<count;i++) {
            const char *id=k==0?e->world.sites[i].bff_id:k==1?e->world.subscribers[i].bff_id:e->world.ues[i].bff_id;
            if(*id) { expected++;if(!json_object_get(ids,id)) rc=wl_error(err,"verification: owned %s ID missing: %s",lists[k],id); }
        }
        json_t *event=json_pack("{s:s,s:I,s:I,s:b}","resource",lists[k],"expected",(json_int_t)expected,"returned",(json_int_t)json_array_size(arr),"passed",rc==0);
        wl_metrics_event(e->metrics,"inventory-check","verification",event);json_decref(event);
        json_decref(root);json_decref(ids);
        if(rc) return rc;
        ulab_status("VERIFY","%s owned=%zu present",lists[k],expected);
    }
    return ULAB_OK;
}
int wl_environment_cleanup(wl_environment_t *e,ulab_error_t *err) {
    if(!e->initialized) return ULAB_OK;
    if(e->running) { pthread_join(e->thread,NULL); e->running=0; }
    if(e->config->keep_environment) return ULAB_OK;
    _Atomic sig_atomic_t cleanup_stop=0;
    _Atomic sig_atomic_t *old=e->stop;
    e->stop=&cleanup_stop; e->http.stop=&cleanup_stop; e->http.phase="cleanup";
    int failures=e->uncertain;
    ulab_error_t tmp={0};
    if(e->runtime.repo[0]) {
        ulab_status("CLEANUP","detach UE sessions and collect diagnostics");
        /* A provider can fail after writing hardware IDs but before the
         * environment worker publishes them. Recover those owned IDs first. */
        if(runtime_load_workload_sites(&e->runtime,&e->world,&tmp)) failures++;
        size_t *indices=calloc(e->world.ue_count?e->world.ue_count:1,sizeof(size_t));
        selector_result_t selected={.idx=indices,.count=0};
        if(!indices) failures++;
        else for(size_t i=0;i<e->world.ue_count;i++) if(e->world.ues[i].started) indices[selected.count++]=i;
        if(indices) for(size_t i=0;i<selected.count;i++) {
            char args[ULAB_MAX_ARGS];
            snprintf(args,sizeof(args),"detach %s %s %s",e->world.ues[indices[i]].id,e->run_dir,e->world.ues[indices[i]].imsi);
            if(wl_environment_script(e,"workload-ue.sh",args,&tmp)) failures++;
        }
        if(e->config->environment->world.ues_per_site && e->config->cdr_wait) {
            double until=wl_now()+e->config->cdr_wait;
            double next_progress=wl_now()+WL_PROGRESS_SECONDS;
            ulab_status("CDR","wait %u sec for PCRF CDR publisher",e->config->cdr_wait);
            while(wl_now()<until) {
                if(wl_now()>=next_progress) {
                    ulab_progress("cleanup CDR wait remaining=%.0fs",until-wl_now());
                    next_progress=wl_now()+WL_PROGRESS_SECONDS;
                }
                pause_briefly();
            }
        }
        if(runtime_collect_cdr_diagnostics(&e->runtime,&e->world,&tmp)) failures++;
        if(indices && runtime_cleanup_selected_ues(&e->runtime,&e->world,&selected,&tmp)) failures++;
        free(indices);
    }
    if(e->bff.authenticated) {
        ulab_status("CLEANUP","delete backend resources");
        /* Legacy cleanup does not remove subscribers. Keep the network until
         * these workload-owned resources have also been removed. */
        world_t view=e->world; view.network_count=0;
        if(bff_cleanup_world(&e->bff,&view,&tmp)) failures++;
        for(size_t i=0;i<e->world.subscriber_count;i++) if(e->world.subscribers[i].bff_id[0]) {
            json_t *v=json_pack("{s:s}","subscriberId",e->world.subscribers[i].bff_id),*root=NULL;
            if(query(e,"deleteSubscriber","mutation deleteSubscriber($subscriberId:String!){deleteSubscriber(subscriberId:$subscriberId){success}}",v,&root,&tmp)) failures++;
            json_decref(v);json_decref(root);
        }
        if(e->world.networks && e->world.networks[0].bff_id[0]) {
            world_t network_view={0}; network_view.networks=e->world.networks; network_view.network_count=1;
            if(bff_cleanup_world(&e->bff,&network_view,&tmp)) failures++;
        }
    }
    if(e->network_started) {
        ulab_status("CLEANUP","stop media/nodes/network");
        if(runtime_cleanup_infra(&e->runtime,&e->world,&tmp)) failures++;
    }
    e->stop=old; e->http.stop=old;
    if(failures) return wl_error(err,"cleanup incomplete (%d failures; unresolved mutation=%d): %.700s",failures,e->uncertain,tmp.msg);
    if(e->journal) {
        json_t *v=json_pack("{s:s,s:s}","kind","cleanup","state","complete");
        json_dumpf(v,e->journal,JSON_COMPACT); fputc('\n',e->journal); fflush(e->journal); json_decref(v);
    }
    ulab_status("CLEANUP","complete");
    return ULAB_OK;
}
void wl_environment_close(wl_environment_t *e) {
    if(!e->initialized) return;
    if(e->running) pthread_join(e->thread,NULL);
    if(e->http.easy) curl_easy_cleanup(e->http.easy);
    if(e->journal) fclose(e->journal);
    runtime_close(&e->runtime); bff_close(&e->bff);
    world_free(&e->world); free(e->site_ready); free(e->provisioned); free(e->batch);
    pthread_mutex_destroy(&e->lock); memset(e,0,sizeof(*e));
}

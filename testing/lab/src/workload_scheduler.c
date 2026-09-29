/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include "util.h"
#include "log.h"
#include <errno.h>
#include <math.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <time.h>
#include <unistd.h>

typedef struct {
    int enabled,visible;
    size_t page,site,node,visits;
    double mounted_at,next_navigation,next[WL_MAX_PAGE_OPS],expires[WL_MAX_OPS];
    uint32_t keys[WL_MAX_OPS];
    unsigned pending[WL_MAX_OPS];
    int last_ok[WL_MAX_OPS];
    unsigned char mount[WL_MAX_PAGE_OPS];
} session_t;
typedef struct {
    pid_t pid;
    size_t ue;
    int kind;
    double deadline,kill_at;
    wl_sample_t sample;
    char result[ULAB_MAX_PATH];
} ue_job_t;
typedef struct {
    wl_environment_t *e;
    wl_catalog_t *cat;
    wl_http_t *http;
    session_t *sessions;
    ue_job_t *jobs;
    unsigned char *busy;
    double *next_traffic,*next_probe,*next_attach;
    uint64_t sequence;
    size_t attached,transfers;
    int runtime_failed;
    double next_start,rate_tokens,last_rate_tick;
    uint32_t random;
    double progress_at;
    uint64_t progress_reads;
} scheduler_t;
static double random_unit(scheduler_t *s) {
    uint32_t x=s->random; x^=x<<13; x^=x>>17; x^=x<<5; s->random=x;
    return x/4294967296.0;
}
static void announce(scheduler_t *s,const char *kind,const char *phase,const char *message) {
    json_t *v=json_pack("{s:s}","message",message);
    wl_metrics_event(s->e->metrics,kind,phase,v); json_decref(v);
}
static void mount_session(scheduler_t *s,session_t *v,size_t id,size_t sites,double now) {
    wl_config_t *c=s->e->config;
    v->page=(id+v->visits++)%c->page_count;
    v->site=(id+v->visits-1)%sites;
    v->node=v->site*3+(id%3);
    v->visible=random_unit(s)*100<c->visible_percent;
    v->mounted_at=(double)time(NULL);
    v->next_navigation=c->page_count>1?now+c->dwell*(0.8+0.4*random_unit(s)):INFINITY;
    for(size_t j=0;j<WL_MAX_PAGE_OPS;j++) { v->next[j]=now; v->mount[j]=1; }
}
static int consoles(scheduler_t *s,size_t wanted,size_t sites,const char *phase,double now,ulab_error_t *err) {
    wl_config_t *c=s->e->config;
    if(!sites) wanted=0;
    for(size_t i=0;i<c->max_sessions;i++) {
        session_t *v=&s->sessions[i];
        if(i>=wanted) { v->enabled=0; continue; }
        char actor[64]; snprintf(actor,sizeof(actor),"console-%06zu",i+1);
        if(!v->enabled) {
            v->enabled=1; mount_session(s,v,i,sites,now);
            int op=wl_catalog_operation(s->cat,"networks");
            if(op>=0) {
                json_t *vars=json_object();
                int rc=wl_http_submit(s->http,&s->cat->operations[op],vars,phase,actor,now,&v->pending[op],&v->last_ok[op],err);
                json_decref(vars); if(rc<0) return ULAB_ERR;
            }
        } else if(now>=v->next_navigation) mount_session(s,v,i,sites,now);
        int page_id=wl_catalog_page(s->cat,c->pages[v->page]);
        wl_page_t *page=&s->cat->pages[page_id];
        for(size_t j=0;j<page->count;j++) {
            wl_page_operation_t *po=&page->operations[j];
            if(now<v->next[j]) continue;
            wl_operation_t *op=&s->cat->operations[po->operation];
            json_t *vars=wl_bind_variables(op->variables,&s->e->world,v->site,v->node,v->mounted_at);
            char *encoded=json_dumps(vars,JSON_COMPACT|JSON_SORT_KEYS);
            if(!vars || !encoded) { json_decref(vars); free(encoded); return wl_error(err,"session allocation failed"); }
            uint32_t key=ulab_hash32(encoded,1); free(encoded);
            double due=v->next[j];
            int cached=v->last_ok[po->operation] && key==v->keys[po->operation] && now<v->expires[po->operation];
            /* cache expiry only suppresses mount fetches; poll ticks always
             * issue a request, matching Apollo polling. */
            int initial=v->mount[j];
            int send=(!initial || po->network_on_mount || !cached) &&
                (initial || v->visible || po->background_poll);
            if(send) {
                int rc=wl_http_submit(s->http,op,vars,phase,actor,due,
                    &v->pending[po->operation],&v->last_ok[po->operation],err);
                if(rc<0) { json_decref(vars); return ULAB_ERR; }
                if(rc==0) { v->keys[po->operation]=key; v->expires[po->operation]=now+po->ttl; }
            }
            json_decref(vars);
            double poll=po->optional_poll && !c->optional_polling?0:po->poll;
            v->next[j]=poll>0?now+poll:INFINITY;
            v->mount[j]=0;
        }
    }
    return ULAB_OK;
}
static int rate_requests(scheduler_t *s,double rps,size_t sites,const char *phase,double now,ulab_error_t *err) {
    double dt=now-s->last_rate_tick; s->last_rate_tick=now;
    if(rps<=0 || !sites) { s->rate_tokens=0; return ULAB_OK; }
    s->rate_tokens+=rps*dt;
    int index=wl_catalog_operation(s->cat,s->e->config->rate_operation);
    wl_operation_t *op=&s->cat->operations[index];
    size_t n=(size_t)floor(s->rate_tokens);
    if(n>4096) {
        wl_sample_t sample={0}; strcpy(sample.phase,phase); strcpy(sample.operation,op->id);
        strcpy(sample.actor,"rate"); strcpy(sample.outcome,"missed");
        sample.scheduled=now-dt; sample.ended=now; sample.weight=n-4096;
        wl_metrics_sample(s->e->metrics,&sample); s->rate_tokens-=(double)(n-4096); n=4096;
    }
    for(size_t i=0;i<n;i++) {
        size_t site=(size_t)(random_unit(s)*sites);
        json_t *template=s->e->config->rate_variables?s->e->config->rate_variables:op->variables;
        json_t *vars=wl_bind_variables(template,&s->e->world,site,site*3,(double)time(NULL));
        double due=now-fmax(0,s->rate_tokens-1)/rps;
        int rc=wl_http_submit(s->http,op,vars,phase,"rate",due,NULL,NULL,err);
        json_decref(vars); s->rate_tokens-=1;
        if(rc<0) return ULAB_ERR;
    }
    return ULAB_OK;
}
static int start_job(scheduler_t *s,ue_job_t *job,size_t slot,size_t index,int kind,const char *phase,ulab_error_t *err) {
    wl_environment_t *e=s->e; ue_t *u=&e->world.ues[index];
    char args[ULAB_MAX_ARGS],file[256],log[ULAB_MAX_PATH];
    const char *script="workload-ue.sh";
    memset(job,0,sizeof(*job)); job->ue=index; job->kind=kind;
    snprintf(file,sizeof(file),"ue-%08llu-%s",(unsigned long long)++s->sequence,u->ref);
    char suffix[280]; snprintf(suffix,sizeof(suffix),"%s.log",file);
    if(wl_path(log,sizeof(log),e->run_dir,suffix,err)) return ULAB_ERR;
    int n;
    if(kind==1) {
        n=snprintf(args,sizeof(args),"attach %s %s %s %s %s %s %s %s",e->opts->repo,u->ref,u->id,u->imsi,u->iccid,u->ip,u->site_ref,e->run_dir);
        u->started=1;
        if(wl_journal(e,"ue",index,"runtime_starting",err)) return ULAB_ERR;
    } else if(kind==2) {
        script="workload-traffic.sh"; snprintf(suffix,sizeof(suffix),"%s.json",file);
        if(wl_path(job->result,sizeof(job->result),e->run_dir,suffix,err)) return ULAB_ERR;
        n=snprintf(args,sizeof(args),"%s %llu %s %zu %s",u->id,(unsigned long long)e->config->traffic_mb,e->run_dir,slot,job->result);
    } else if(kind==3) n=snprintf(args,sizeof(args),"detach %s %s %s",u->id,e->run_dir,u->imsi);
    else n=snprintf(args,sizeof(args),"probe %s %s",u->id,e->run_dir);
    if(n<0 || (size_t)n>=sizeof(args)) return wl_error(err,"UE arguments too long");
    job->sample.scheduled=job->sample.started=wl_now();
    ulab_copy(job->sample.phase,sizeof(job->sample.phase),phase);
    ulab_copy(job->sample.actor,sizeof(job->sample.actor),u->ref);
    strcpy(job->sample.operation,kind==1?"ue_attach":kind==2?"ue_traffic":kind==3?"ue_detach":"ue_probe");
    job->pid=wl_process_start(e->runtime.script_dir,script,args,log,err);
    if(job->pid<0) { job->pid=0; return ULAB_ERR; }
    job->deadline=wl_now()+e->config->job_timeout; s->busy[index]=1;
    if(kind==2) s->transfers++;
    return ULAB_OK;
}
static int reap_jobs(scheduler_t *s,int cancelling,ulab_error_t *err) {
    wl_environment_t *e=s->e;
    for(size_t i=0;i<e->config->ue_concurrency;i++) {
        ue_job_t *j=&s->jobs[i]; int status=0;
        if(!j->pid) continue;
        pid_t rc=waitpid(j->pid,&status,WNOHANG);
        if(rc<0 && errno==EINTR) continue;
        if(rc==0) {
            if(cancelling || wl_now()>j->deadline) {
                if(!j->kill_at) { kill(-j->pid,SIGTERM); j->kill_at=wl_now()+2; }
                else if(wl_now()>j->kill_at) kill(-j->pid,SIGKILL);
            }
            continue;
        }
        int ok=rc>0 && WIFEXITED(status) && WEXITSTATUS(status)==0 && !j->kill_at;
        ue_t *u=&e->world.ues[j->ue];
        j->sample.ended=wl_now();
        strcpy(j->sample.outcome,ok?"ok":cancelling?"cancelled":j->kill_at?"timeout":"runtime_error");
        if(j->kind==2) {
            s->transfers--;
            json_error_t je; json_t *r=json_load_file(j->result,0,&je);
            json_t *bytes=json_object_get(json_object_get(json_object_get(r,"end"),"sum_received"),"bytes");
            if(ok && (!json_is_number(bytes) || json_number_value(bytes)<=0)) {
                ok=0; strcpy(j->sample.outcome,"invalid_traffic_result");
            }
            if(ok) j->sample.bytes=(size_t)json_number_value(bytes);
            json_decref(r);
            s->next_traffic[j->ue]=wl_now()+e->config->pause_min+random_unit(s)*(e->config->pause_max-e->config->pause_min);
        } else if(j->kind==1) {
            if(ok && !u->attached) { u->attached=1; s->attached++; }
            s->next_attach[j->ue]=wl_now()+30;
            s->next_probe[j->ue]=wl_now()+e->config->observe_interval;
            if(ok && wl_journal(e,"ue",j->ue,"attached",err)) return ULAB_ERR;
        } else if(j->kind==3 && ok) {
            if(u->attached) s->attached--;
            u->attached=0; u->started=0; s->next_attach[j->ue]=0;
            if(wl_journal(e,"ue",j->ue,"detached",err)) return ULAB_ERR;
        } else if(j->kind==4) {
            if(!ok && u->attached) { u->attached=0; s->attached--; }
            s->next_probe[j->ue]=wl_now()+e->config->observe_interval;
        }
        if(!ok && !cancelling) s->runtime_failed=1;
        wl_metrics_sample(e->metrics,&j->sample);
        s->busy[j->ue]=0; j->pid=0;
    }
    return ULAB_OK;
}
static size_t inflight_jobs(scheduler_t *s,int kind) {
    size_t n=0;
    for(size_t i=0;i<s->e->config->ue_concurrency;i++) if(s->jobs[i].pid && (!kind || s->jobs[i].kind==kind)) n++;
    return n;
}
static void show_progress(scheduler_t *s,const char *phase,double elapsed,double left,
                          const wl_targets_t *actual,const wl_targets_t *target,double now) {
    wl_read_counts_t reads;
    wl_metrics_read_counts(s->e->metrics,&reads);
    double rps=now>s->progress_at?(reads.completed-s->progress_reads)/(now-s->progress_at):0;
    const char *task=s->e->running?(s->e->task_kind==1?s->e->world.sites[s->e->site_index].ref:"provision-sims"):"none";
    if(strstr(phase,"/reach"))
        ulab_progress("%s %.0fs sites=%zu/%zu sims=%zu/%zu ues=%zu/%zu task=%s left=%.0fs reads=%llu errors=%llu",
            phase,elapsed,actual->sites,target->sites,actual->provisioned,target->provisioned,
            actual->attached,target->attached,task,fmax(0,left),
            (unsigned long long)reads.completed,(unsigned long long)reads.failed);
    else
        ulab_progress("%s %.0f/%.0fs reads=%llu errors=%llu missed=%llu rps=%.1f inflight=%zu sessions=%zu sites=%zu ues=%zu ue_jobs=%zu",
            phase,elapsed,elapsed+fmax(0,left),(unsigned long long)reads.completed,
            (unsigned long long)reads.failed,(unsigned long long)reads.missed,
            rps,wl_http_pending(s->http),actual->sessions,actual->sites,actual->attached,inflight_jobs(s,0));
    s->progress_at=now;s->progress_reads=reads.completed;
}
static int ues(scheduler_t *s,const wl_targets_t *target,const char *phase,double now,ulab_error_t *err) {
    wl_environment_t *e=s->e;
    size_t live=e->world.site_count*e->config->environment->world.ues_per_site;
    for(size_t slot=0;slot<e->config->ue_concurrency;slot++) {
        if(s->jobs[slot].pid) continue;
        size_t candidate=SIZE_MAX; int kind=0;
        size_t attaching=inflight_jobs(s,1),detaching=inflight_jobs(s,3);
        if(s->attached+attaching<target->attached && now>=s->next_start) {
            for(size_t i=0;i<live;i++) if(e->provisioned[i] && !e->world.ues[i].attached && !s->busy[i] && now>=s->next_attach[i]) { candidate=i; kind=1; break; }
        } else if(s->attached>target->attached+detaching) {
            for(size_t i=live;i>0;i--) if(e->provisioned[i-1] && e->world.ues[i-1].attached && !s->busy[i-1]) { candidate=i-1;kind=3;break; }
        }
        if(!kind) {
            for(size_t i=0;i<live;i++) if(e->provisioned[i] && e->world.ues[i].attached && !s->busy[i] && now>=s->next_probe[i]) { candidate=i;kind=4;break; }
        }
        if(!kind) {
            size_t active=0;
            double oldest=INFINITY;
            for(size_t i=0;i<live && active<target->active;i++) if(e->provisioned[i] && e->world.ues[i].attached) {
                active++;
                /* Oldest due first prevents low-index UEs monopolizing the
                 * worker pool when transfers take longer than think time. */
                if(!s->busy[i] && now>=s->next_traffic[i] && s->next_traffic[i]<oldest) {
                    candidate=i;kind=2;oldest=s->next_traffic[i];
                }
            }
        }
        if(kind && start_job(s,&s->jobs[slot],slot,candidate,kind,phase,err)) return ULAB_ERR;
        if(kind==1) s->next_start=now+1/e->config->attach_rate;
    }
    return ULAB_OK;
}
static wl_targets_t ramp_target(const wl_targets_t *from,const wl_targets_t *to,double f) {
    wl_targets_t t=*to;
#define RAMP(field) t.field=(size_t)floor((double)from->field+((double)to->field-from->field)*f+1e-9)
    RAMP(sites);RAMP(provisioned);RAMP(attached);RAMP(active);RAMP(sessions);
#undef RAMP
    t.rps=from->rps+(to->rps-from->rps)*f;
    return t;
}
static int reached(const wl_targets_t *a,const wl_targets_t *b) {
    return a->sites==b->sites && a->provisioned==b->provisioned && a->attached==b->attached && a->sessions==b->sessions && a->active==b->active;
}
int wl_schedule(wl_environment_t *e,wl_catalog_t *cat,ulab_error_t *err) {
    scheduler_t s={0}; wl_config_t *c=e->config; int rc=ULAB_OK;
    size_t reporting_phase=0;int reporting_hold=0,reporting_active=0;
    size_t n=e->world.ue_count?:1;
    s.e=e;s.cat=cat;s.random=c->environment->seed?:1;
    s.sessions=calloc(c->max_sessions?:1,sizeof(session_t));
    s.jobs=calloc(c->ue_concurrency,sizeof(ue_job_t)); s.busy=calloc(n,1);
    s.next_traffic=calloc(n,sizeof(double));s.next_probe=calloc(n,sizeof(double));s.next_attach=calloc(n,sizeof(double));
    s.http=wl_http_open(c,&e->bff,e->metrics,e->stop,err);
    if(!s.sessions || !s.jobs || !s.busy || !s.next_traffic || !s.next_probe || !s.next_attach || !s.http) { rc=wl_error(err,"scheduler allocation failed"); goto done; }
    wl_targets_t previous={0},actual={0},desired={0};
    double last_tick=wl_now(); s.last_rate_tick=s.progress_at=last_tick;
    for(size_t phase=0;phase<=c->phase_count;phase++) {
        int initial=phase==0;
        wl_phase_t init={.ramp=0,.reach=c->initial_timeout,.hold=0,.target=c->initial};
        strcpy(init.name,"setup");
        wl_phase_t *p=initial?&init:&c->phases[phase-1];
        double start=wl_now(),hold_start=0,next_progress=start+WL_PROGRESS_SECONDS;
        char label[ULAB_MAX_REF];
        reporting_phase=phase;reporting_hold=0;reporting_active=1;
        wl_metrics_stage_begin(e->metrics,phase,0,start);
        snprintf(label,sizeof(label),"%.100s/reach",p->name);
        announce(&s,"phase",label,"reach");
        ulab_status("PHASE","%s ramp=%.0fs reach_timeout=%.0fs target sites=%zu provisioned=%zu attached=%zu active=%zu sessions=%zu rate=%g/s",
            label,p->ramp,p->reach,p->target.sites,p->target.provisioned,p->target.attached,
            p->target.active,p->target.sessions,p->target.rps);
        for(;;) {
            double now=wl_now();
            if(*e->stop) { rc=wl_error(err,"workload interrupted"); goto finish; }
            double f=p->ramp>0?fmin(1,(now-start)/p->ramp):1;
            desired=ramp_target(&previous,&p->target,f);
            if(wl_environment_step(e,&desired,label,err) || reap_jobs(&s,0,err)) { rc=ULAB_ERR; goto finish; }
            wl_environment_counts(e,&actual); actual.attached=s.attached;
            if(ues(&s,&desired,label,now,err) || consoles(&s,desired.sessions,actual.sites,label,now,err) ||
               rate_requests(&s,desired.rps,actual.sites,label,now,err) || wl_http_poll(s.http,5,err)) { rc=ULAB_ERR; goto finish; }
            actual.sessions=actual.sites?desired.sessions:0;
            actual.active=fmin(desired.active,s.attached); actual.rps=0;
            if(now-last_tick>=1) {
                if(wl_metrics_tick(e->metrics,label,&desired,&actual,wl_http_pending(s.http),s.transfers,now)) { rc=wl_error(err,"metrics output failed"); goto finish; }
                last_tick=now;
            }
            if(!hold_start && f>=1 && reached(&actual,&p->target) && !e->running) {
                wl_metrics_stage_end(e->metrics,phase,0,now,1,&actual,NULL);
                reporting_hold=1;wl_metrics_stage_begin(e->metrics,phase,1,now);
                hold_start=now; snprintf(label,sizeof(label),"%.100s/hold",p->name);
                announce(&s,"phase",label,"hold");
                ulab_status("PHASE","%s target reached; hold=%.0fs",label,p->hold);
                next_progress=now+WL_PROGRESS_SECONDS;
            }
            if(!hold_start && now-start>p->reach) { rc=wl_error(err,"phase %s did not reach target populations",p->name); goto finish; }
            if(hold_start && !reached(&actual,&p->target)) {
                /* A population loss during hold is a reliability failure,
                 * even if the actor later reconnects successfully. */
                s.runtime_failed=1;
                wl_metrics_stage_fail(e->metrics,phase,1,"target population was lost during hold");
            }
            int complete=hold_start && now-hold_start>=p->hold;
            if(now>=next_progress) {
                double elapsed=now-(hold_start?hold_start:start);
                show_progress(&s,label,elapsed,(hold_start?p->hold:p->reach)-elapsed,&actual,&p->target,now);
                next_progress=now+WL_PROGRESS_SECONDS;
            }
            if(complete) {
                wl_metrics_stage_end(e->metrics,phase,1,now,1,&actual,NULL);reporting_active=0;
                wl_read_counts_t reads;wl_metrics_phase_counts(e->metrics,p->name,&reads);
                ulab_status("PHASE","%s complete; phase reads=%llu errors=%llu missed=%llu (final checks after drain)",p->name,
                    (unsigned long long)reads.completed,(unsigned long long)reads.failed,(unsigned long long)reads.missed);
                break;
            }
        }
        previous=p->target;
    }
finish:
    if(reporting_active) wl_metrics_stage_end(e->metrics,reporting_phase,reporting_hold,wl_now(),0,&actual,*e->stop?NULL:err->msg);
    announce(&s,"shutdown","drain",rc?err->msg:"workload phases completed");
    ulab_status("DRAIN","%s; inflight=%zu ue_jobs=%zu",rc?err->msg:"workload phases completed",wl_http_pending(s.http),inflight_jobs(&s,0));
    if(rc) *e->stop=1;
    double until=wl_now()+fmax(c->request_timeout,c->job_timeout)+2,next_progress=wl_now()+WL_PROGRESS_SECONDS;
    while((wl_http_pending(s.http)||inflight_jobs(&s,0)) && wl_now()<until) {
        ulab_error_t tmp={0};
        if(wl_http_poll(s.http,20,&tmp) || reap_jobs(&s,*e->stop,&tmp)) { if(!rc) *err=tmp;rc=ULAB_ERR;*e->stop=1; }
        if(wl_now()>=next_progress) {
            ulab_progress("drain inflight=%zu ue_jobs=%zu deadline_in=%.0fs",wl_http_pending(s.http),inflight_jobs(&s,0),fmax(0,until-wl_now()));
            next_progress=wl_now()+WL_PROGRESS_SECONDS;
        }
    }
    if(inflight_jobs(&s,0)) {
        for(size_t i=0;i<c->ue_concurrency;i++) if(s.jobs[i].pid) {
            kill(-s.jobs[i].pid,SIGKILL); wl_process_wait(s.jobs[i].pid,wl_now()+1,e->stop,err);s.jobs[i].pid=0;
        }
        rc=ULAB_ERR;
    }
    if(s.runtime_failed && !rc) rc=wl_error(err,"one or more UE operations or hold population checks failed");
    wl_http_close(s.http);s.http=NULL;
    wl_environment_counts(e,&actual);actual.attached=s.attached;
    wl_metrics_tick(e->metrics,"drain",&desired,&actual,0,0,wl_now());
done:
    wl_http_close(s.http);
    free(s.sessions);free(s.jobs);free(s.busy);free(s.next_traffic);free(s.next_probe);free(s.next_attach);
    return rc;
}

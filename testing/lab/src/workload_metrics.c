/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include "util.h"
#include "log.h"
#include <math.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#define HIST_BINS 4096
#define HIST_LOG 0.009950330853168082
typedef struct {
    uint64_t bins[HIST_BINS], count;
    double sum, max;
} histogram_t;
typedef struct row {
    char phase[ULAB_MAX_REF], operation[ULAB_MAX_REF], category[24];
    uint64_t completed, failed, missed, throttled, timeouts, gaps;
    histogram_t latency, dispatch;
    wl_sample_t first_error;
    struct row *next;
} row_t;
typedef struct {
    int entered, complete;
    double started, ended;
    wl_targets_t actual;
    char failure[ULAB_MAX_ERR];
} stage_t;
struct wl_metrics {
    FILE *requests, *series, *events;
    pthread_mutex_t lock;
    row_t *totals, *window;
    double start, last_tick;
    time_t started_at;
    char dir[ULAB_MAX_PATH];
    int io_failed;
    stage_t stages[WL_MAX_PHASES+1][2];
};
static void hist_add(histogram_t *h,double milliseconds) {
    double us=fmax(0,milliseconds*1000);
    size_t bin=(size_t)(log1p(us)/HIST_LOG);
    if(bin>=HIST_BINS) bin=HIST_BINS-1;
    h->bins[bin]++; h->count++; h->sum+=milliseconds;
    if(milliseconds>h->max) h->max=milliseconds;
}
static double percentile(histogram_t *h,double q) {
    uint64_t n=0, target=(uint64_t)ceil(h->count*q);
    if(!target) return 0;
    for(size_t i=0;i<HIST_BINS;i++) {
        n+=h->bins[i];
        if(n>=target) return fmin(h->max,expm1((i+1)*HIST_LOG)/1000);
    }
    return h->max;
}
static const char *category(const char *actor) {
    if(!strncmp(actor,"console",7)) return "console";
    if(!strncmp(actor,"rate",4)) return "rate";
    if(!strncmp(actor,"ue",2)) return "ue";
    return "environment";
}
static row_t *row_get(row_t **head,const wl_sample_t *s) {
    row_t *r;
    const char *kind=category(s->actor);
    for(r=*head;r;r=r->next)
        if(!strcmp(r->phase,s->phase) && !strcmp(r->operation,s->operation) && !strcmp(r->category,kind)) return r;
    r=calloc(1,sizeof(*r));
    if(!r) return NULL;
    ulab_copy(r->phase,sizeof(r->phase),s->phase);
    ulab_copy(r->operation,sizeof(r->operation),s->operation);
    ulab_copy(r->category,sizeof(r->category),kind);
    r->next=*head; *head=r;
    return r;
}
static void row_add(row_t *r,const wl_sample_t *s) {
    if(!strcmp(s->outcome,"missed") || !strcmp(s->outcome,"overlap_skipped")) { r->missed+=s->weight?s->weight:1; return; }
    r->completed++;
    r->gaps+=s->known_gaps;
    if(strcmp(s->outcome,"ok")) {
        r->failed++;
        if(!r->first_error.outcome[0]) r->first_error=*s;
    }
    if(!strcmp(s->outcome,"throttled")) r->throttled++;
    if(!strcmp(s->outcome,"timeout")) r->timeouts++;
    hist_add(&r->latency,(s->ended-s->started)*1000);
    hist_add(&r->dispatch,fmax(0,s->started-s->scheduled)*1000);
}
static void rows_free(row_t *r) { while(r) { row_t *next=r->next; free(r); r=next; } }
static void hist_merge(histogram_t *to,const histogram_t *from) {
    for(size_t i=0;i<HIST_BINS;i++) to->bins[i]+=from->bins[i];
    to->count+=from->count;to->sum+=from->sum;to->max=fmax(to->max,from->max);
}
static void row_merge(row_t *to,const row_t *from) {
    to->completed+=from->completed;to->failed+=from->failed;to->missed+=from->missed;
    to->throttled+=from->throttled;to->timeouts+=from->timeouts;to->gaps+=from->gaps;
    hist_merge(&to->latency,&from->latency);hist_merge(&to->dispatch,&from->dispatch);
}
static int write_json(FILE *f,json_t *v) {
    int rc=json_dumpf(v,f,JSON_COMPACT|JSON_ENSURE_ASCII);
    if(fputc('\n',f)==EOF) rc=-1;
    return rc;
}
wl_metrics_t *wl_metrics_open(const char *dir,const wl_config_t *c,ulab_error_t *err) {
    wl_metrics_t *m=calloc(1,sizeof(*m)); char path[ULAB_MAX_PATH];
    if(!m) return NULL;
    pthread_mutex_init(&m->lock,NULL); strcpy(m->dir,dir);
    m->start=m->last_tick=wl_now(); m->started_at=time(NULL);
    const char *names[]={"requests.jsonl","timeseries.jsonl","events.jsonl"};
    FILE **files[]={&m->requests,&m->series,&m->events};
    for(size_t i=0;i<3;i++) {
        if(wl_path(path,sizeof(path),dir,names[i],err) || !(*files[i]=fopen(path,"w"))) {
            wl_error(err,"cannot open workload output %s",names[i]); wl_metrics_close(m); return NULL;
        }
        setvbuf(*files[i],NULL,_IOFBF,256*1024);
    }
    if(wl_path(path,sizeof(path),dir,"workload.resolved.json",err) || json_dump_file(c->source,path,JSON_INDENT(2))) {
        wl_error(err,"cannot save resolved workload"); wl_metrics_close(m); return NULL;
    }
    return m;
}
void wl_metrics_phase_counts(wl_metrics_t *m,const char *phase,wl_read_counts_t *counts) {
    memset(counts,0,sizeof(*counts));
    pthread_mutex_lock(&m->lock);
    for(row_t *r=m->totals;r;r=r->next) {
        if(strcmp(r->category,"console") && strcmp(r->category,"rate")) continue;
        if(phase && (strncmp(r->phase,phase,strlen(phase)) || r->phase[strlen(phase)]!='/')) continue;
        counts->completed+=r->completed;
        counts->failed+=r->failed;
        counts->missed+=r->missed;
    }
    pthread_mutex_unlock(&m->lock);
}
void wl_metrics_read_counts(wl_metrics_t *m,wl_read_counts_t *counts) {
    wl_metrics_phase_counts(m,NULL,counts);
}
void wl_metrics_stage_begin(wl_metrics_t *m,size_t index,int hold,double now) {
    if(index>WL_MAX_PHASES) return;
    pthread_mutex_lock(&m->lock);
    stage_t *s=&m->stages[index][!!hold];s->entered=1;s->started=now;
    pthread_mutex_unlock(&m->lock);
}
void wl_metrics_stage_end(wl_metrics_t *m,size_t index,int hold,double now,
                          int complete,const wl_targets_t *actual,const char *failure) {
    if(index>WL_MAX_PHASES) return;
    pthread_mutex_lock(&m->lock);
    stage_t *s=&m->stages[index][!!hold];s->ended=now;s->complete=complete;s->actual=*actual;
    if(failure && *failure && !s->failure[0]) ulab_copy(s->failure,sizeof(s->failure),failure);
    pthread_mutex_unlock(&m->lock);
}
void wl_metrics_stage_fail(wl_metrics_t *m,size_t index,int hold,const char *failure) {
    if(index>WL_MAX_PHASES) return;
    pthread_mutex_lock(&m->lock);
    stage_t *s=&m->stages[index][!!hold];
    if(!s->failure[0]) ulab_copy(s->failure,sizeof(s->failure),failure);
    pthread_mutex_unlock(&m->lock);
}
void wl_metrics_sample(wl_metrics_t *m,const wl_sample_t *s) {
    json_t *v;
    if(!m) return;
    pthread_mutex_lock(&m->lock);
    row_t *total=row_get(&m->totals,s), *window=row_get(&m->window,s);
    if(!total || !window) { m->io_failed=1; pthread_mutex_unlock(&m->lock); return; }
    row_add(total,s); row_add(window,s);
    v=json_pack("{s:s,s:s,s:s,s:s,s:s,s:f,s:f,s:f,s:f,s:f,s:i,s:I,s:i,s:s}",
        "phase",s->phase,"operation",s->operation,"actor",s->actor,"outcome",s->outcome,
        "request_id",s->request_id,"scheduled_s",s->scheduled-m->start,
        "started_s",s->started>0?s->started-m->start:0,"ended_s",s->ended-m->start,
        "elapsed_ms",s->started>0?(s->ended-s->started)*1000:0,
        "dispatch_delay_ms",s->started>0?fmax(0,s->started-s->scheduled)*1000:0,
        "http_status",(int)s->http_status,"response_bytes",(json_int_t)s->bytes,
        "known_gaps",(int)s->known_gaps,"detail",s->detail);
    if(v) {
        json_object_set_new(v,"dns_ms",json_real(s->dns_ms));
        json_object_set_new(v,"connect_ms",json_real(s->connect_ms));
        json_object_set_new(v,"tls_ms",json_real(s->tls_ms));
        json_object_set_new(v,"ttfb_ms",json_real(s->ttfb_ms));
        json_object_set_new(v,"weight",json_integer((json_int_t)(s->weight?s->weight:1)));
    }
    if(!v || write_json(m->requests,v)) m->io_failed=1;
    json_decref(v);
    pthread_mutex_unlock(&m->lock);
}
void wl_metrics_event(wl_metrics_t *m,const char *kind,const char *phase,json_t *data) {
    if(!m) return;
    pthread_mutex_lock(&m->lock);
    json_t *v=json_pack("{s:s,s:s,s:f,s:O}","kind",kind,"phase",phase,"time_s",wl_now()-m->start,"data",data);
    if(!v || write_json(m->events,v)) m->io_failed=1;
    json_decref(v); fflush(m->events);
    pthread_mutex_unlock(&m->lock);
}
static json_t *row_json(row_t *r) {
    return json_pack("{s:s,s:s,s:s,s:I,s:I,s:I,s:I,s:I,s:I,s:f,s:f,s:f,s:f,s:f,s:f}",
        "phase",r->phase,"operation",r->operation,"category",r->category,
        "completed",(json_int_t)r->completed,"failed",(json_int_t)r->failed,
        "missed",(json_int_t)r->missed,"throttled",(json_int_t)r->throttled,
        "timeouts",(json_int_t)r->timeouts,"known_gaps",(json_int_t)r->gaps,
        "p50_ms",percentile(&r->latency,.50),"p95_ms",percentile(&r->latency,.95),
        "p99_ms",percentile(&r->latency,.99),"max_ms",r->latency.max,
        "mean_ms",r->latency.count?r->latency.sum/r->latency.count:0,
        "dispatch_p95_ms",percentile(&r->dispatch,.95));
}
static json_t *target_json(const wl_targets_t *t) {
    return json_pack("{s:I,s:I,s:I,s:I,s:I,s:f}","sites",(json_int_t)t->sites,
        "provisioned_ues",(json_int_t)t->provisioned,"attached_ues",(json_int_t)t->attached,
        "active_ues",(json_int_t)t->active,"console_sessions",(json_int_t)t->sessions,"rps",t->rps);
}
static int read_row(const row_t *r) {
    return !strcmp(r->category,"console") || !strcmp(r->category,"rate");
}
static double error_percent(const row_t *r) { return r->completed?100.0*r->failed/r->completed:0; }
static double missed_percent(const row_t *r) { return r->completed+r->missed?100.0*r->missed/(r->completed+r->missed):0; }
static int threshold_failed(row_t *r,const wl_config_t *c) {
    return error_percent(r)>c->max_error_percent || missed_percent(r)>c->max_missed_percent ||
        (c->p95_ms>0 && percentile(&r->latency,.95)>c->p95_ms);
}
/* Build verdicts only after all producers have drained. A completion belongs
 * to its dispatch stage, even when it arrives in a later phase. */
static json_t *stage_result(wl_metrics_t *m,const wl_config_t *c,size_t index,int hold,
                            const wl_phase_t *p,const wl_targets_t *from) {
    stage_t *s=&m->stages[index][hold];row_t total={0};
    char label[ULAB_MAX_REF],reason[ULAB_MAX_ERR]={0};
    snprintf(label,sizeof(label),"%.100s/%s",p->name,hold?"hold":"reach");
    json_t *checks=json_array();size_t failed_checks=0;uint64_t ue_failed=0;
    double worst_error=0,worst_missed=0,worst_p95=0;
    for(row_t *r=m->totals;r;r=r->next) {
        if(strcmp(r->phase,label)) continue;
        if(!strcmp(r->category,"ue")) ue_failed+=r->failed;
        if(!read_row(r)) continue;
        row_merge(&total,r);
        int failed=threshold_failed(r,c);failed_checks+=failed;
        double errors=error_percent(r),misses=missed_percent(r),p95=percentile(&r->latency,.95);
        worst_error=fmax(worst_error,errors);worst_missed=fmax(worst_missed,misses);worst_p95=fmax(worst_p95,p95);
        json_t *v=row_json(r);
        json_object_set_new(v,"status",json_string(failed?"FAIL":"PASS"));
        json_object_set_new(v,"error_percent",json_real(errors));
        json_object_set_new(v,"missed_percent",json_real(misses));
        json_array_append_new(checks,v);
    }
    double duration=s->entered?fmax(0,s->ended-s->started):0;
    int expects_reads=p->target.rps>0 || p->target.sessions>0 ||
        (!hold && p->ramp>0 && (from->rps>0 || from->sessions>0));
    const char *status="PASS";
    if(!s->entered) {status="NOT_RUN";strcpy(reason,"stage was not entered");}
    else if(s->failure[0]) {status="FAIL";ulab_copy(reason,sizeof(reason),s->failure);}
    else if(failed_checks || ue_failed) {
        status="FAIL";snprintf(reason,sizeof(reason),"%zu operation threshold checks failed; %llu UE operations failed",failed_checks,(unsigned long long)ue_failed);
    } else if(!s->complete || m->io_failed) {
        status="INCOMPLETE";strcpy(reason,!s->complete?"stage did not complete":"metrics output incomplete");
    } else if(expects_reads && (hold?p->hold>0:p->ramp>0) && duration>0 && !json_array_size(checks)) {
        status="NO_DATA";strcpy(reason,"no console/rate samples observed; thresholds not evaluated");
    } else if(!json_array_size(checks)) strcpy(reason,"population/lifecycle checks completed; no read threshold checks");
    else strcpy(reason,"all observed operation threshold checks passed");
    json_t *v=row_json(&total),*actual=target_json(&s->actual);
    json_object_del(actual,"rps"); /* rps is a scheduling target, not a population observation. */
    json_object_set_new(v,"phase",json_string(label));json_object_del(v,"category");json_object_del(v,"operation");
    json_object_set_new(v,"status",json_string(status));json_object_set_new(v,"reason",json_string(reason));
    json_object_set_new(v,"entered",json_boolean(s->entered));json_object_set_new(v,"complete",json_boolean(s->complete));
    json_object_set_new(v,"duration_sec",json_real(duration));
    json_object_set_new(v,"expected_population",target_json(&p->target));
    json_object_set_new(v,"actual_population",actual);
    json_object_set_new(v,"successful",json_integer((json_int_t)(total.completed-total.failed)));
    json_object_set_new(v,"successful_rps",json_real(duration>0?(total.completed-total.failed)/duration:0));
    json_object_set_new(v,"error_percent",json_real(error_percent(&total)));
    json_object_set_new(v,"missed_percent",json_real(missed_percent(&total)));
    json_object_set_new(v,"worst_error_percent",json_real(worst_error));
    json_object_set_new(v,"worst_missed_percent",json_real(worst_missed));
    json_object_set_new(v,"worst_p95_ms",json_real(worst_p95));
    json_object_set_new(v,"failed_checks",json_integer((json_int_t)failed_checks));
    json_object_set_new(v,"ue_failed",json_integer((json_int_t)ue_failed));
    json_object_set_new(v,"checks",checks);
    return v;
}
static json_t *phase_result(wl_metrics_t *m,const wl_config_t *c,size_t index,
                            const wl_phase_t *p,const wl_targets_t *from) {
    json_t *reach=stage_result(m,c,index,0,p,from),*hold=stage_result(m,c,index,1,p,from);
    const char *a=json_string_value(json_object_get(reach,"status")),*b=json_string_value(json_object_get(hold,"status"));
    const char *status="PASS";
    if(!strcmp(a,"FAIL") || !strcmp(b,"FAIL")) status="FAIL";
    else if(!strcmp(a,"NOT_RUN") && !strcmp(b,"NOT_RUN")) status="NOT_RUN";
    else if(!strcmp(a,"INCOMPLETE") || !strcmp(b,"INCOMPLETE") || !strcmp(a,"NOT_RUN") || !strcmp(b,"NOT_RUN")) status="INCOMPLETE";
    else if(!strcmp(a,"NO_DATA") || !strcmp(b,"NO_DATA")) status="NO_DATA";
    return json_pack("{s:s,s:s,s:o,s:o,s:o,s:o,s:f,s:f}","name",p->name,"status",status,
        "from",target_json(from),"target",target_json(&p->target),"reach",reach,"hold",hold,
        "ramp_seconds",p->ramp,"hold_seconds",p->hold);
}
int wl_metrics_tick(wl_metrics_t *m,const char *phase,const wl_targets_t *target,
                     const wl_targets_t *actual,size_t requests,size_t transfers,double now) {
    pthread_mutex_lock(&m->lock);
    for(row_t *r=m->window;r;r=r->next) {
        json_t *v=row_json(r);
        json_object_set_new(v,"kind",json_string("requests"));
        json_object_set_new(v,"time_s",json_real(now-m->start));
        json_object_set_new(v,"window_s",json_real(now-m->last_tick));
        if(write_json(m->series,v)) m->io_failed=1;
        json_decref(v);
    }
    rows_free(m->window); m->window=NULL;
    json_t *v=json_pack("{s:s,s:s,s:f,s:o,s:o,s:I,s:I}","kind","population","phase",phase,
        "time_s",now-m->start,"target",target_json(target),"actual",target_json(actual),
        "inflight_requests",(json_int_t)requests,"transfers",(json_int_t)transfers);
    if(!v || write_json(m->series,v)) m->io_failed=1;
    json_decref(v); m->last_tick=now;
    if(fflush(m->series) || fflush(m->requests)) m->io_failed=1;
    int rc=m->io_failed?ULAB_ERR:ULAB_OK;
    pthread_mutex_unlock(&m->lock); return rc;
}
static int html_report(wl_metrics_t *m,json_t *summary,const char *assets,ulab_error_t *err) {
    char path[ULAB_MAX_PATH], template_path[ULAB_MAX_PATH];
    FILE *out=NULL,*in=NULL; int rc=ULAB_ERR;
    if(wl_path(path,sizeof(path),m->dir,"workload.html",err) ||
       wl_path(template_path,sizeof(template_path),assets,"report.html",err)) return ULAB_ERR;
    out=fopen(path,"w"); in=fopen(template_path,"r");
    if(!out || !in) goto done;
    char buf[8192]; size_t n;
    while((n=fread(buf,1,sizeof(buf),in))) if(fwrite(buf,1,n,out)!=n) goto done;
    fclose(in); in=NULL;
    fputs("\n<script type=\"application/json\" id=\"summary-data\">",out);
    char *s=json_dumps(summary,JSON_COMPACT|JSON_ENSURE_ASCII);
    if(!s) goto done;
    for(char *p=s;*p;p++) { if(*p=='<') fputs("\\u003c",out); else fputc(*p,out); }
    free(s); fputs("</script>\n<script type=\"application/json\" id=\"series-data\">[\n",out);
    if(wl_path(path,sizeof(path),m->dir,"timeseries.jsonl",err)) goto done;
    in=fopen(path,"r"); if(!in) goto done;
    char *line=NULL; size_t cap=0; ssize_t len; int first=1;
    while((len=getline(&line,&cap,in))>=0) {
        (void)len;
        if(!first) fputs(",\n",out);
        first=0;
        for(char *p=line;*p;p++) { if(*p=='<') fputs("\\u003c",out); else fputc(*p,out); }
    }
    free(line);
    fputs("]</script>\n<script>renderWorkload();</script></body></html>\n",out);
    rc=ferror(out)||ferror(in)?ULAB_ERR:ULAB_OK;
done:
    if(in) fclose(in);
    if(out && fclose(out)) rc=ULAB_ERR;
    if(rc) wl_error(err,"failed to write workload.html");
    return rc;
}
static void add_reason(json_t *reasons,const char *reason) {
    if(!reason || !*reason) return;
    for(size_t i=0;i<json_array_size(reasons);i++)
        if(!strcmp(json_string_value(json_array_get(reasons,i)),reason)) return;
    json_array_append_new(reasons,json_string(reason));
}
static void verdict(json_t *summary,int result) {
    json_t *reasons=json_object_get(summary,"failure_reasons");
    char text[ULAB_MAX_ERR]={0};size_t used=0;
    for(size_t i=0;i<json_array_size(reasons) && used<sizeof(text)-1;i++) {
        int n=snprintf(text+used,sizeof(text)-used,"%s%s",used?"; ":"",json_string_value(json_array_get(reasons,i)));
        if(n<0 || (size_t)n>=sizeof(text)-used) break;
        used+=(size_t)n;
    }
    json_object_set_new(summary,"reason",json_string(text));
    json_object_set_new(summary,"passed",json_boolean(result==ULAB_OK));
    json_object_set_new(summary,"final_rc",json_integer(result));
}
static int save_report(wl_metrics_t *m,json_t *summary,const char *assets,const char *file,ulab_error_t *err) {
    char path[ULAB_MAX_PATH];
    if(!strcmp(file,"workload.html")) return html_report(m,summary,assets,err);
    if(wl_path(path,sizeof(path),m->dir,file,err)) return ULAB_ERR;
    if(!strcmp(file,"report.json")) return json_dump_file(summary,path,JSON_INDENT(2))?ULAB_ERR:ULAB_OK;
    FILE *f=fopen(path,"w");if(!f) return ULAB_ERR;
    int rc=wl_report_write(f,summary,0);
    if(fclose(f)) rc=ULAB_ERR;
    return rc;
}
int wl_metrics_finish(wl_metrics_t *m,const wl_config_t *c,int result,
                       int cleanup_failed,const char *reason,const char *assets,ulab_error_t *err) {
    json_t *summary=json_object(),*rows=json_array(),*violations=json_array(),*examples=json_array();
    json_t *reasons=json_array(),*writes=json_array(),*operations=json_array();
    row_t *totals=NULL;histogram_t read_latency={0};
    uint64_t requests=0,failed=0,missed=0;
    int execution_result=result;
    for(row_t *r=m->totals;r;r=r->next) {
        json_t *v=row_json(r);
        json_array_append_new(rows,v);
        if(r->first_error.outcome[0]) {
            wl_sample_t *s=&r->first_error;
            json_array_append_new(examples,json_pack("{s:s,s:s,s:s,s:s,s:i,s:s,s:s}",
                "phase",r->phase,"operation",r->operation,"category",r->category,
                "outcome",s->outcome,"http_status",(int)s->http_status,"detail",s->detail,"request_id",s->request_id));
        }
        if(strcmp(r->category,"console") && strcmp(r->category,"rate")) continue;
        requests+=r->completed; failed+=r->failed; missed+=r->missed;
        hist_merge(&read_latency,&r->latency);
        wl_sample_t key={0};strcpy(key.phase,"all");strcpy(key.actor,r->category);strcpy(key.operation,r->operation);
        row_t *total=row_get(&totals,&key);
        if(total) row_merge(total,r);else m->io_failed=1;
        double errors=r->completed?100.0*r->failed/r->completed:0;
        double misses=r->completed+r->missed?100.0*r->missed/(r->completed+r->missed):0;
        if(threshold_failed(r,c)) {
            json_t *violation=json_pack("{s:s,s:s,s:s,s:f,s:f,s:f,s:I,s:I,s:I}","phase",r->phase,
                "operation",r->operation,"category",r->category,"error_percent",errors,
                "missed_percent",misses,"p95_ms",percentile(&r->latency,.95),"completed",(json_int_t)r->completed,
                "failed",(json_int_t)r->failed,"missed",(json_int_t)r->missed);
            json_array_append_new(violations,violation);
        }
    }
    for(row_t *r=totals;r;r=r->next) {
        json_t *v=row_json(r);
        json_object_set_new(v,"error_percent",json_real(r->completed?100.0*r->failed/r->completed:0));
        json_array_append_new(operations,v);
    }
    rows_free(totals);
    if(fflush(m->requests)) m->io_failed=1;
    if(fflush(m->series)) m->io_failed=1;
    if(fflush(m->events)) m->io_failed=1;
    if(execution_result || cleanup_failed) add_reason(reasons,reason);
    if(execution_result && (!reason || !*reason)) add_reason(reasons,"workload execution or verification failed");
    if(cleanup_failed) add_reason(reasons,"cleanup failed");
    if(m->io_failed) add_reason(reasons,"metrics output incomplete");
    if(json_array_size(violations)) {
        char message[128];snprintf(message,sizeof(message),"%zu operation/phase threshold checks failed",json_array_size(violations));
        add_reason(reasons,message);
    }
    if(json_array_size(violations) || cleanup_failed || m->io_failed) result=ULAB_ERR;
    json_object_set_new(summary,"engine",json_string("workload-1"));
    json_object_set_new(summary,"lab_version",json_string(ULAB_VERSION));
    json_object_set_new(summary,"thresholds",json_pack("{s:f,s:f,s:f}","max_error_percent",c->max_error_percent,"max_missed_percent",c->max_missed_percent,"p95_ms",c->p95_ms));
    json_object_set_new(summary,"format",json_integer(1));
    json_object_set_new(summary,"kind",json_string("workload"));
    json_object_set_new(summary,"scenario",json_string(c->name));
    json_object_set_new(summary,"suite",json_string("workload"));
    json_object_set_new(summary,"output_dir",json_string(m->dir));
    json_object_set_new(summary,"execution",json_string(execution_result?"failed":"ok"));
    json_object_set_new(summary,"started_at",json_integer(m->started_at));
    json_object_set_new(summary,"duration_sec",json_real(wl_now()-m->start));
    json_object_set_new(summary,"cleanup",json_string(c->keep_environment?"retained":cleanup_failed?"failed":"ok"));
    json_object_set_new(summary,"metrics_complete",json_boolean(!m->io_failed));
    json_object_set_new(summary,"latency_note",json_string("Client elapsed time for all completed attempts, including errors/timeouts. Timeouts are censored at the deadline. Histogram percentile upper bounds have <=1% relative bin width above 1 ms. This is not server processing time."));
    json_object_set_new(summary,"requests",json_integer((json_int_t)requests));
    json_object_set_new(summary,"failed_requests",json_integer((json_int_t)failed));
    json_object_set_new(summary,"missed_requests",json_integer((json_int_t)missed));
    json_object_set_new(summary,"error_percent",json_real(requests?100.0*failed/requests:0));
    json_object_set_new(summary,"missed_percent",json_real(requests+missed?100.0*missed/(requests+missed):0));
    json_object_set_new(summary,"read_latency",json_pack("{s:f,s:f,s:f,s:f}","p50_ms",percentile(&read_latency,.5),
        "p95_ms",percentile(&read_latency,.95),"p99_ms",percentile(&read_latency,.99),"max_ms",read_latency.max));
    json_object_set_new(summary,"threshold_violations",violations);
    json_object_set_new(summary,"operations",rows);
    json_object_set_new(summary,"operation_totals",operations);
    json_object_set_new(summary,"error_examples",examples);
    json_object_set_new(summary,"failure_reasons",reasons);
    json_object_set_new(summary,"report_write_errors",writes);
    wl_phase_t setup={.target=c->initial};strcpy(setup.name,"setup");wl_targets_t zero={0};
    json_object_set_new(summary,"setup_result",phase_result(m,c,0,&setup,&zero));
    json_t *phases=json_array(),*phase_counts=json_object();
    const char *statuses[]={"PASS","FAIL","INCOMPLETE","NOT_RUN","NO_DATA"};
    for(size_t i=0;i<5;i++) json_object_set_new(phase_counts,statuses[i],json_integer(0));
    for(size_t i=0;i<c->phase_count;i++) {
        json_t *p=phase_result(m,c,i+1,&c->phases[i],i?&c->phases[i-1].target:&c->initial);
        const char *status=json_string_value(json_object_get(p,"status"));
        json_int_t n=json_integer_value(json_object_get(phase_counts,status));
        json_object_set_new(phase_counts,status,json_integer(n+1));json_array_append_new(phases,p);
    }
    json_object_set_new(summary,"phase_results",phases);json_object_set_new(summary,"phase_counts",phase_counts);
    json_object_set_new(summary,"phase_note",json_string("Chronological results, finalized after drain and attributed to dispatch stage. Reach includes the ramp and readiness wait; hold is the steady target. Limits apply independently to each operation/category. Phase status does not inherit unrelated cleanup or report failures. NO_DATA and NOT_RUN are not PASS. Successful req/s uses dispatch-stage duration; p95 includes failed attempts."));
    verdict(summary,result);
    const char *files[]={"report.json","report.txt","workload.html"};
    /* If a report write fails, refresh surviving artifacts with the final FAIL
     * and its reason. Never turn an output failure into a silent PASS. */
    for(size_t attempt=0;attempt<2;attempt++) {
        size_t before=json_array_size(writes);
        for(size_t i=0;i<3;i++) {
            if(save_report(m,summary,assets,files[i],err)) {
                char message[128];snprintf(message,sizeof(message),"failed to write %s",files[i]);
                add_reason(writes,message);add_reason(reasons,message);result=ULAB_ERR;verdict(summary,result);
            }
        }
        if(json_array_size(writes)==before) break;
    }
    ulab_progress_clear();
    wl_report_write(stdout,summary,20);fflush(stdout);
    json_decref(summary);return result;
}
void wl_metrics_close(wl_metrics_t *m) {
    if(!m) return;
    if(m->requests) fclose(m->requests);
    if(m->series) fclose(m->series);
    if(m->events) fclose(m->events);
    rows_free(m->totals); rows_free(m->window);
    pthread_mutex_destroy(&m->lock); free(m);
}

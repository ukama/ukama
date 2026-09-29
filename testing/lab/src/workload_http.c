/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include "util.h"
#include "log.h"
#include <stdlib.h>
#include <unistd.h>
#include <time.h>
#include <string.h>
#include <math.h>

typedef struct {
    char *data;
    size_t length, capacity, limit;
    int exceeded;
} buffer_t;
typedef struct transfer {
    CURL *easy;
    struct curl_slist *headers;
    char *body;
    buffer_t response;
    wl_sample_t sample;
    const wl_operation_t *operation;
    unsigned *pending;
    int *last_ok;
    struct transfer *next;
} transfer_t;
struct wl_http {
    CURLM *multi;
    const wl_config_t *config;
    bff_client_t *auth;
    wl_metrics_t *metrics;
    _Atomic sig_atomic_t *stop;
    transfer_t *active;
    size_t count;
    uint64_t sequence;
};
static size_t receive(void *data,size_t size,size_t count,void *ctx) {
    buffer_t *b=ctx;
    if(count && size>SIZE_MAX/count) return 0;
    size_t n=size*count;
    if(n>b->limit-b->length) { b->exceeded=1; return 0; }
    if(b->length+n+1>b->capacity) {
        size_t cap=b->capacity?b->capacity:4096;
        while(cap<b->length+n+1) cap*=2;
        if(cap>b->limit+1) cap=b->limit+1;
        char *p=realloc(b->data,cap);
        if(!p) return 0;
        b->data=p; b->capacity=cap;
    }
    memcpy(b->data+b->length,data,n); b->length+=n; b->data[b->length]='\0';
    return n;
}
static int progress(void *ctx,curl_off_t a,curl_off_t b,curl_off_t c,curl_off_t d) {
    _Atomic sig_atomic_t *stop=ctx;
    (void)a;(void)b;(void)c;(void)d;
    return stop && *stop;
}
typedef struct {
    wl_sync_http_t *http;
    const wl_sample_t *sample;
} sync_progress_t;
static int sync_progress(void *ctx,curl_off_t a,curl_off_t b,curl_off_t c,curl_off_t d) {
    sync_progress_t *p=ctx;
    wl_sync_http_t *h=p->http;
    double now=wl_now();
    (void)a;(void)b;(void)c;(void)d;
    if(now>=h->next_progress) {
        ulab_progress("%s BFF %s elapsed=%.0fs timeout=%.0fs",p->sample->phase,p->sample->operation,
            now-p->sample->started,h->config->setup_timeout);
        h->next_progress=now+WL_PROGRESS_SECONDS;
    }
    return h->stop && *h->stop;
}
static char *body(const char *query,json_t *vars,const char *operation_name) {
    json_t *o=json_pack("{s:s,s:O}","query",query,"variables",vars);
    if(operation_name) json_object_set_new(o,"operationName",json_string(operation_name));
    char *s=json_dumps(o,JSON_COMPACT); json_decref(o); return s;
}
static int configure(CURL *easy,const char *url,const char *token,
                      const char *request_id,char *payload,buffer_t *response,
                      double timeout,_Atomic sig_atomic_t *stop,
                      struct curl_slist **headers,ulab_error_t *err) {
    char value[8192];
    *headers=curl_slist_append(*headers,"Content-Type: application/json");
    if(snprintf(value,sizeof(value),"X-Session-Token: %s",token)>=(int)sizeof(value)) return wl_error(err,"session token too long");
    *headers=curl_slist_append(*headers,value);
    snprintf(value,sizeof(value),"X-Request-ID: %s",request_id);
    *headers=curl_slist_append(*headers,value);
    if(!*headers) return wl_error(err,"header allocation failed");
#define SET(k,v) do { if(curl_easy_setopt(easy,k,v)!=CURLE_OK) return wl_error(err,"curl option %s failed",#k); } while(0)
    SET(CURLOPT_URL,url); SET(CURLOPT_HTTPHEADER,*headers);
    SET(CURLOPT_POSTFIELDS,payload); SET(CURLOPT_POSTFIELDSIZE,(long)strlen(payload));
    SET(CURLOPT_WRITEFUNCTION,receive); SET(CURLOPT_WRITEDATA,response);
    SET(CURLOPT_TIMEOUT_MS,(long)(timeout*1000));
    SET(CURLOPT_CONNECTTIMEOUT_MS,(long)(fmin(timeout,10)*1000));
    SET(CURLOPT_NOSIGNAL,1L); SET(CURLOPT_ACCEPT_ENCODING,"");
    SET(CURLOPT_FOLLOWLOCATION,0L);
    SET(CURLOPT_XFERINFOFUNCTION,progress); SET(CURLOPT_XFERINFODATA,(void *)stop);
    SET(CURLOPT_NOPROGRESS,0L);
#undef SET
    return ULAB_OK;
}
static void timing(CURL *easy,wl_sample_t *s) {
    double v=0;
    curl_easy_getinfo(easy,CURLINFO_RESPONSE_CODE,&s->http_status);
    curl_easy_getinfo(easy,CURLINFO_NAMELOOKUP_TIME,&v); s->dns_ms=v*1000;
    curl_easy_getinfo(easy,CURLINFO_CONNECT_TIME,&v); s->connect_ms=v*1000;
    curl_easy_getinfo(easy,CURLINFO_APPCONNECT_TIME,&v); s->tls_ms=v*1000;
    curl_easy_getinfo(easy,CURLINFO_STARTTRANSFER_TIME,&v); s->ttfb_ms=v*1000;
}
int wl_sync_transport(void *ctx,const char *op,const char *query,const char *vars,
                       json_t **out,ulab_error_t *err) {
    wl_sync_http_t *h=ctx;
    struct curl_slist *headers=NULL;
    buffer_t response={0}; wl_sample_t s={0};
    json_error_t je; json_t *variables=json_loads(vars?vars:"{}",0,&je);
    char *payload;
    *out=NULL;
    if(!variables) return wl_error(err,"invalid variables for %s",op);
    payload=body(query,variables,NULL); json_decref(variables);
    if(!payload) return wl_error(err,"request allocation failed");
    if(!h->easy) h->easy=curl_easy_init();
    if(!h->easy) { free(payload); return wl_error(err,"curl init failed"); }
    curl_easy_reset(h->easy);
    response.limit=h->config->max_response_bytes;
    ulab_copy(s.phase,sizeof(s.phase),h->phase?h->phase:"setup");
    ulab_copy(s.operation,sizeof(s.operation),op); strcpy(s.actor,"environment");
    snprintf(s.request_id,sizeof(s.request_id),"ulab-%ld-%ld-env-%llu",(long)time(NULL),(long)getpid(),(unsigned long long)++h->sequence);
    if(configure(h->easy,h->auth->url,h->auth->token,s.request_id,payload,&response,
                  h->config->setup_timeout,h->stop,&headers,err)) {
        curl_slist_free_all(headers); free(payload); return ULAB_ERR;
    }
    while(wl_now()<h->next_request && !(h->stop && *h->stop)) { struct timespec t={0,10000000}; nanosleep(&t,NULL); }
    h->next_request=wl_now()+1/h->config->setup_max_rps;
    s.scheduled=s.started=wl_now();
    sync_progress_t progress_context={h,&s};
    if(!h->next_progress) h->next_progress=s.started+WL_PROGRESS_SECONDS;
    curl_easy_setopt(h->easy,CURLOPT_XFERINFOFUNCTION,sync_progress);
    curl_easy_setopt(h->easy,CURLOPT_XFERINFODATA,&progress_context);
    CURLcode code=curl_easy_perform(h->easy); s.ended=wl_now();
    timing(h->easy,&s); s.bytes=response.length;
    const char *outcome=wl_classify(s.http_status,code,response.data,response.length,NULL,out,&s.known_gaps,err);
    if(response.exceeded) { outcome="response_limit"; wl_error(err,"response exceeds configured limit"); }
    ulab_copy(s.outcome,sizeof(s.outcome),outcome); ulab_copy(h->last_outcome,sizeof(h->last_outcome),outcome);
    if(strcmp(outcome,"ok")) ulab_copy(s.detail,sizeof(s.detail),err->msg);
    wl_metrics_sample(h->metrics,&s);
    free(response.data); free(payload); curl_slist_free_all(headers);
    if(strcmp(outcome,"ok")) { json_decref(*out); *out=NULL; return ULAB_ERR; }
    return ULAB_OK;
}
wl_http_t *wl_http_open(const wl_config_t *c,bff_client_t *auth,wl_metrics_t *metrics,
                        _Atomic sig_atomic_t *stop,ulab_error_t *err) {
    wl_http_t *h=calloc(1,sizeof(*h));
    if(!h) return NULL;
    h->multi=curl_multi_init(); h->config=c; h->auth=auth; h->metrics=metrics; h->stop=stop;
    if(!h->multi) { free(h); wl_error(err,"curl multi init failed"); return NULL; }
    curl_multi_setopt(h->multi,CURLMOPT_MAX_TOTAL_CONNECTIONS,(long)c->max_inflight);
    curl_multi_setopt(h->multi,CURLMOPT_MAXCONNECTS,(long)c->max_inflight);
    return h;
}
int wl_http_submit(wl_http_t *h,const wl_operation_t *op,json_t *variables,
                    const char *phase,const char *actor,double scheduled,
                    unsigned *pending,int *last_ok,ulab_error_t *err) {
    wl_sample_t missed={0};
    if(h->count>=h->config->max_inflight || (pending && *pending)) {
        ulab_copy(missed.phase,sizeof(missed.phase),phase); ulab_copy(missed.operation,sizeof(missed.operation),op->id);
        ulab_copy(missed.actor,sizeof(missed.actor),actor);
        strcpy(missed.outcome,pending && *pending?"overlap_skipped":"missed");
        missed.scheduled=scheduled; missed.ended=wl_now();
        wl_metrics_sample(h->metrics,&missed); return 1;
    }
    transfer_t *t=calloc(1,sizeof(*t));
    if(!t) { wl_error(err,"request allocation failed"); return -1; }
    t->operation=op; t->pending=pending; t->last_ok=last_ok; t->response.limit=h->config->max_response_bytes;
    t->easy=curl_easy_init(); t->body=body(op->query,variables,op->name);
    ulab_copy(t->sample.phase,sizeof(t->sample.phase),phase);
    ulab_copy(t->sample.operation,sizeof(t->sample.operation),op->id);
    ulab_copy(t->sample.actor,sizeof(t->sample.actor),actor);
    snprintf(t->sample.request_id,sizeof(t->sample.request_id),"ulab-%ld-%ld-read-%llu",(long)time(NULL),(long)getpid(),(unsigned long long)++h->sequence);
    if(!t->easy || !t->body || configure(t->easy,h->auth->url,h->auth->token,t->sample.request_id,t->body,
        &t->response,h->config->request_timeout,h->stop,&t->headers,err)) goto bad;
    t->sample.scheduled=scheduled; t->sample.started=wl_now();
    curl_easy_setopt(t->easy,CURLOPT_PRIVATE,t);
    if(curl_multi_add_handle(h->multi,t->easy)!=CURLM_OK) { wl_error(err,"cannot queue request"); goto bad; }
    t->next=h->active; h->active=t; h->count++; if(pending) (*pending)++;
    if(last_ok) *last_ok=0;
    return ULAB_OK;
bad:
    if(t->easy) curl_easy_cleanup(t->easy);
    curl_slist_free_all(t->headers); free(t->body); free(t); return -1;
}
static void complete(wl_http_t *h,transfer_t *t,CURLcode code) {
    json_t *root=NULL; ulab_error_t err={0};
    t->sample.ended=wl_now(); timing(t->easy,&t->sample); t->sample.bytes=t->response.length;
    const char *outcome=wl_classify(t->sample.http_status,code,t->response.data,t->response.length,t->operation,&root,&t->sample.known_gaps,&err);
    if(t->response.exceeded) { outcome="response_limit"; wl_error(&err,"response exceeds configured limit"); }
    ulab_copy(t->sample.outcome,sizeof(t->sample.outcome),outcome);
    ulab_copy(t->sample.detail,sizeof(t->sample.detail),err.msg);
    wl_metrics_sample(h->metrics,&t->sample); json_decref(root);
    if(t->last_ok) *t->last_ok=!strcmp(outcome,"ok");
    if(t->pending) (*t->pending)--;
    curl_multi_remove_handle(h->multi,t->easy); curl_easy_cleanup(t->easy);
    curl_slist_free_all(t->headers); free(t->body); free(t->response.data);
    transfer_t **p=&h->active; while(*p && *p!=t) p=&(*p)->next;
    if(*p) *p=t->next;
    h->count--; free(t);
}
int wl_http_poll(wl_http_t *h,int wait_ms,ulab_error_t *err) {
    int running=0, messages=0, fds=0;
    CURLMsg *msg;
    if(curl_multi_perform(h->multi,&running)!=CURLM_OK) return wl_error(err,"curl multi perform failed");
    while((msg=curl_multi_info_read(h->multi,&messages))) if(msg->msg==CURLMSG_DONE) {
        transfer_t *t=NULL; curl_easy_getinfo(msg->easy_handle,CURLINFO_PRIVATE,&t);
        if(t) complete(h,t,msg->data.result);
    }
    if(wait_ms>0 && curl_multi_poll(h->multi,NULL,0,wait_ms,&fds)!=CURLM_OK) return wl_error(err,"curl multi poll failed");
    return ULAB_OK;
}
size_t wl_http_pending(wl_http_t *h) { return h?h->count:0; }
void wl_http_close(wl_http_t *h) {
    if(!h) return;
    while(h->active) complete(h,h->active,CURLE_ABORTED_BY_CALLBACK);
    curl_multi_cleanup(h->multi); free(h);
}

/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include "util.h"
#include <stdlib.h>
#include <string.h>
#include <time.h>

int wl_catalog_operation(const wl_catalog_t *c, const char *id) {
    for(size_t i=0;i<c->operation_count;i++) if(!strcmp(c->operations[i].id,id)) return (int)i;
    return -1;
}
int wl_catalog_page(const wl_catalog_t *c, const char *name) {
    for(size_t i=0;i<c->page_count;i++) if(!strcmp(c->pages[i].name,name)) return (int)i;
    return -1;
}
int wl_catalog_load(const wl_config_t *cfg, wl_catalog_t *c, ulab_error_t *err) {
    char path[ULAB_MAX_PATH];
    json_error_t je;
    json_t *arr, *profiles;
    memset(c,0,sizeof(*c));
    if(wl_path(path,sizeof(path),cfg->assets,"catalog/console.json",err)) return ULAB_ERR;
    c->manifest=json_load_file(path,JSON_REJECT_DUPLICATES,&je);
    if(!c->manifest) return wl_error(err,"cannot load catalog %s: %s",path,je.text);
    arr=json_object_get(c->manifest,"operations");
    if(!json_is_array(arr) || !json_array_size(arr) || json_array_size(arr)>WL_MAX_OPS) return wl_error(err,"invalid operation catalog");
    for(size_t i=0;i<json_array_size(arr);i++) {
        json_t *v=json_array_get(arr,i);
        wl_operation_t *op=&c->operations[c->operation_count++];
        const char *id=json_string_value(json_object_get(v,"id"));
        const char *name=json_string_value(json_object_get(v,"name"));
        const char *query=json_string_value(json_object_get(v,"query"));
        if(!id || !name || !query || strlen(id)>=sizeof(op->id) || strlen(name)>=sizeof(op->name) || strncmp(query,"query ",6))
            return wl_error(err,"catalog requires named query operations");
        strcpy(op->id,id); strcpy(op->name,name); op->query=strdup(query);
        op->variables=json_incref(json_object_get(v,"variables"));
        op->required=json_incref(json_object_get(v,"required"));
        op->allowed_errors=json_incref(json_object_get(v,"allowed_errors"));
        if(!op->query || !json_is_object(op->variables) || !json_is_array(op->required)) return wl_error(err,"invalid operation %s",id);
        for(size_t j=0;j<i;j++) if(!strcmp(c->operations[j].id,id)) return wl_error(err,"duplicate operation id");
    }
    if(wl_path(path,sizeof(path),cfg->assets,"profiles/console.yaml",err)) return ULAB_ERR;
    profiles=wl_yaml_load(path,err);
    if(!profiles) return ULAB_ERR;
    arr=json_object_get(profiles,"pages");
    if(!json_is_array(arr) || json_array_size(arr)>WL_MAX_PAGES) { json_decref(profiles); return wl_error(err,"invalid console profiles"); }
    for(size_t i=0;i<json_array_size(arr);i++) {
        wl_page_t *p=&c->pages[c->page_count++];
        json_t *v=json_array_get(arr,i), *ops=json_object_get(v,"operations");
        const char *name=json_string_value(json_object_get(v,"name"));
        if(!name || strlen(name)>=sizeof(p->name) || !json_is_array(ops) || json_array_size(ops)>WL_MAX_PAGE_OPS) goto bad;
        strcpy(p->name,name);
        for(size_t j=0;j<json_array_size(ops);j++) {
            json_t *x=json_array_get(ops,j);
            const char *id=json_string_value(json_object_get(x,"operation"));
            wl_page_operation_t *po=&p->operations[p->count++];
            int index=id?wl_catalog_operation(c,id):-1;
            if(index<0) goto bad;
            po->operation=(size_t)index;
            po->poll=json_number_value(json_object_get(x,"poll_seconds"));
            po->ttl=json_number_value(json_object_get(x,"ttl_seconds"));
            po->optional_poll=json_is_true(json_object_get(x,"optional_poll"));
            po->background_poll=json_is_true(json_object_get(x,"background_poll"));
            po->network_on_mount=json_is_true(json_object_get(x,"network_on_mount"));
            if(po->poll<0 || po->ttl<0) goto bad;
        }
    }
    json_decref(profiles);
    for(size_t i=0;i<cfg->page_count;i++) if(wl_catalog_page(c,cfg->pages[i])<0) return wl_error(err,"unknown console page %s",cfg->pages[i]);
    if(cfg->rate_operation[0] && wl_catalog_operation(c,cfg->rate_operation)<0) return wl_error(err,"unknown rate operation %s",cfg->rate_operation);
    return ULAB_OK;
bad:
    json_decref(profiles);
    return wl_error(err,"invalid console profile entry");
}
void wl_catalog_free(wl_catalog_t *c) {
    for(size_t i=0;i<c->operation_count;i++) {
        free(c->operations[i].query); json_decref(c->operations[i].variables);
        json_decref(c->operations[i].required); json_decref(c->operations[i].allowed_errors);
    }
    json_decref(c->manifest); memset(c,0,sizeof(*c));
}
json_t *wl_bind_variables(json_t *v, const world_t *w, size_t site,
                          size_t node, double mounted_at) {
    if(json_is_object(v)) {
        const char *k; json_t *value, *out=json_object();
        json_object_foreach(v,k,value) json_object_set_new(out,k,wl_bind_variables(value,w,site,node,mounted_at));
        return out;
    }
    if(json_is_array(v)) {
        json_t *out=json_array();
        for(size_t i=0;i<json_array_size(v);i++) json_array_append_new(out,wl_bind_variables(json_array_get(v,i),w,site,node,mounted_at));
        return out;
    }
    if(json_is_string(v)) {
        const char *s=json_string_value(v);
        if(!strcmp(s,"$networkId")) return json_string(w->networks[0].bff_id);
        if(!strcmp(s,"$siteId")) return json_string(site<w->site_count?w->sites[site].bff_id:"");
        if(!strcmp(s,"$nodeId")) return json_string(node<w->node_count?w->nodes[node].bff_id:"");
        if(!strcmp(s,"$uptimeKey")) return json_string(node%3==1?"ctl_uptime":"com_uptime");
        if(!strcmp(s,"$mountedAt") || !strcmp(s,"$from30d")) {
            char out[64]; struct tm tm; time_t t=(time_t)mounted_at;
            if(!strcmp(s,"$from30d")) t-=30*86400;
            gmtime_r(&t,&tm); strftime(out,sizeof(out),"%Y-%m-%dT%H:%M:%S.000Z",&tm);
            return json_string(out);
        }
    }
    return json_incref(v);
}

static json_t *at(json_t *root,const char *path) {
    char copy[512], *save=NULL, *part;
    if(strlen(path)>=sizeof(copy)) return NULL;
    strcpy(copy,path);
    for(part=strtok_r(copy,".",&save);part;part=strtok_r(NULL,".",&save)) {
        root=json_object_get(root,part);
        if(!root) break;
    }
    return root;
}
static int sections(json_t *v,const char *path,const wl_operation_t *op,
                     unsigned *gaps,ulab_error_t *err,unsigned depth) {
    if(depth>64) return wl_error(err,"response nesting too deep");
    if(json_is_object(v)) {
        const char *key; json_t *value;
        json_object_foreach(v,key,value) {
            char next[1024];
            if(snprintf(next,sizeof(next),"%s%s%s",path,*path?".":"",key)>=(int)sizeof(next)) return wl_error(err,"response path too long");
            if(!strcmp(key,"error") && !json_is_null(value)) {
                const char *code=json_string_value(json_object_get(value,"code"));
                int allowed=0;
                for(size_t i=0;op && i<json_array_size(op->allowed_errors);i++) {
                    json_t *rule=json_array_get(op->allowed_errors,i);
                    const char *rp=json_string_value(json_object_get(rule,"path"));
                    const char *rc=json_string_value(json_object_get(rule,"code"));
                    if(rp && rc && code && !strcmp(next,rp) && !strcmp(code,rc)) allowed=1;
                }
                if(allowed) { (*gaps)++; continue; }
                return wl_error(err,"%s: %s",next,code?code:"section error");
            }
            if(!strcmp(key,"success") && json_is_false(value)) return wl_error(err,"%s=false",next);
            if(sections(value,next,op,gaps,err,depth+1)) return ULAB_ERR;
        }
    } else if(json_is_array(v)) {
        for(size_t i=0;i<json_array_size(v);i++) if(sections(json_array_get(v,i),path,op,gaps,err,depth+1)) return ULAB_ERR;
    }
    return ULAB_OK;
}
const char *wl_classify(long status,int curl_code,const char *body,size_t len,
                        const wl_operation_t *op,json_t **root,unsigned *gaps,
                        ulab_error_t *err) {
    json_error_t je;
    *root=NULL; *gaps=0;
    if(curl_code) {
        wl_error(err,"%s",curl_easy_strerror((CURLcode)curl_code));
        if(curl_code==CURLE_OPERATION_TIMEDOUT) return "timeout";
        if(curl_code==CURLE_ABORTED_BY_CALLBACK) return "cancelled";
        return "transport_error";
    }
    if(status==429) { wl_error(err,"HTTP 429"); return "throttled"; }
    if(status<200 || status>=300) { wl_error(err,"HTTP %ld",status); return "http_error"; }
    *root=json_loadb(body?body:"",len,JSON_REJECT_DUPLICATES,&je);
    if(!json_is_object(*root)) { wl_error(err,"invalid JSON response"); return "invalid_response"; }
    json_t *errors=json_object_get(*root,"errors");
    if(errors && (!json_is_array(errors) || json_array_size(errors))) {
        const char *msg=json_string_value(json_object_get(json_array_get(errors,0),"message"));
        wl_error(err,"%s",msg?msg:"GraphQL error");
        const char *code=json_string_value(json_object_get(json_object_get(json_array_get(errors,0),"extensions"),"code"));
        if(code && (!strcmp(code,"GRAPHQL_PARSE_FAILED") || !strcmp(code,"GRAPHQL_VALIDATION_FAILED") ||
           !strcmp(code,"BAD_USER_INPUT") || !strcmp(code,"UNAUTHENTICATED") || !strcmp(code,"FORBIDDEN"))) return "graphql_rejected";
        return "graphql_error";
    }
    json_t *data=json_object_get(*root,"data");
    if(!json_is_object(data)) { wl_error(err,"missing data object"); return "invalid_response"; }
    if(sections(data,"data",op,gaps,err,0)) return "section_error";
    for(size_t i=0;op && i<json_array_size(op->required);i++) {
        json_t *rule=json_array_get(op->required,i);
        const char *path=json_string_value(json_object_get(rule,"path"));
        const char *type=json_string_value(json_object_get(rule,"type"));
        json_t *v=path?at(*root,path):NULL;
        int valid=v && !json_is_null(v);
        if(type && valid) {
            if(!strcmp(type,"array")) valid=json_is_array(v);
            else if(!strcmp(type,"object")) valid=json_is_object(v);
            else if(!strcmp(type,"number")) valid=json_is_number(v);
            else if(!strcmp(type,"boolean")) valid=json_is_boolean(v);
            else if(!strcmp(type,"string")) valid=json_is_string(v);
            else valid=0;
        }
        if(!valid) { wl_error(err,"missing/invalid %s",path?path:"required field"); return "invalid_response"; }
    }
    return "ok";
}

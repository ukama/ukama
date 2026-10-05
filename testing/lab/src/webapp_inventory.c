/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include "webapp.h"
#include <string.h>
#include <stdlib.h>

/* Expectations are owned-world identities, never list/query response values. */
static json_t *inventory_context(world_t *w, const char *ref) {
    json_t *ctx = json_object(), *sites = json_array(), *nodes = json_array();
    size_t i;
    network_t *net = world_network_by_ref(w, ref);
    if (!net || !net->bff_id[0]) goto fail;
    json_object_set_new(ctx,"network",json_pack("{s:s,s:s,s:s}","ref",net->ref,"id",net->bff_id,"name",net->name));
    for (i=0;i<w->site_count;i++) {
        site_t *s=&w->sites[i];
        if (!s->bff_id[0]) goto fail;
        json_array_append_new(sites,json_pack("{s:s,s:s,s:s,s:s}","ref",s->ref,"id",s->bff_id,"name",s->name,"network_ref",s->network_ref));
    }
    for (i=0;i<w->node_count;i++) {
        node_t *n=&w->nodes[i];
        if (!n->bff_id[0]) goto fail;
        json_array_append_new(nodes,json_pack("{s:s,s:s,s:s,s:s}","ref",n->ref,"id",n->bff_id,"site_ref",n->site_ref,"network_ref",n->network_ref));
    }
    json_object_set_new(ctx,"sites",sites); json_object_set_new(ctx,"nodes",nodes);
    return ctx;
fail:
    json_decref(ctx); json_decref(sites); json_decref(nodes); return NULL;
}
static json_t *inputs(world_t *w, const char *view, const selector_t *net, const selector_t *sites, const selector_t *nodes) {
    json_t *ctx=inventory_context(w,net->value), *out;
    site_t *site=world_site_by_ref(w,sites->value);
    node_t *node=world_node_by_ref(w,nodes->value);
    if (!ctx) return NULL;
    out=json_pack("{s:s,s:o}","view",view,"context",ctx);
    if (site) json_object_set_new(out,"site_ref",json_string(site->ref));
    if (node) json_object_set_new(out,"node_ref",json_string(node->ref));
    return out;
}
int webapp_inventory_event(const event_spec_t *e, world_t *w, json_t **out, ulab_error_t *err) {
    *out=inputs(w,e->view,&e->networks,&e->sites,&e->nodes);
    if (!*out) return webapp_error(err,"inventory context requires provisioned world identities");
    json_object_set_new(*out,"action",json_string(e->target));
    json_object_set_new(*out,"value",json_string(e->status));
    return ULAB_OK;
}
static int compare_string(const void *a, const void *b) { return strcmp(json_string_value(*(json_t *const *)a),json_string_value(*(json_t *const *)b)); }
static void sort_values(json_t *array) {
    size_t n=json_array_size(array),i;
    json_t **values=calloc(n?n:1,sizeof(*values));
    if (!values) return;
    for(i=0;i<n;i++) values[i]=json_incref(json_array_get(array,i));
    qsort(values,n,sizeof(*values),compare_string); json_array_clear(array);
    for(i=0;i<n;i++) json_array_append_new(array,values[i]);
    free(values);
}
int webapp_inventory_check(const check_spec_t *c, world_t *w, json_t **out, ulab_error_t *err) {
    json_t *expected=NULL;
    network_t *net=world_network_by_ref(w,c->networks.value);
    site_t *site=world_site_by_ref(w,c->sites.value);
    node_t *node=world_node_by_ref(w,c->nodes.value);
    char path[ULAB_MAX_PATH];
    size_t i;
    *out=inputs(w,c->view,&c->networks,&c->sites,&c->nodes);
    if (!*out) return webapp_error(err,"inventory context requires provisioned world identities");
    if (!c->key[0]) expected=json_string(c->expected);
    else if (!strcmp(c->key,"network_name")) expected=json_string(net->name);
    else if (!strcmp(c->key,"path")) {
        const char *base=!strcmp(c->view,"network_home")?"/network":strstr(c->view,"site")?"/network/sites":"/network/nodes";
        snprintf(path,sizeof(path),"%s%s%s",base,strstr(c->view,"_detail")?"/":"",strstr(c->view,"_detail")?(site?site->bff_id:node?node->bff_id:""):"");
        expected=json_string(path);
    } else if (!strcmp(c->key,"site_name") && site) {
        expected=!strcmp(c->label,"Site names")?json_pack("[s]",site->name):json_string(site->name);
    } else if (!strcmp(c->key,"node_id") && node) expected=json_string(node->bff_id);
    else if (!strcmp(c->key,"site_names")) {
        expected=json_array();
        for(i=0;i<w->site_count;i++) if(!strcmp(w->sites[i].network_ref,net->ref)) json_array_append_new(expected,json_string(w->sites[i].name));
    } else if (!strcmp(c->key,"node_ids") || (!strcmp(c->key,"site_node_ids") && site)) {
        expected=json_array();
        for(i=0;i<w->node_count;i++) if(!strcmp(w->nodes[i].network_ref,net->ref) && (!site || !strcmp(w->nodes[i].site_ref,site->ref))) json_array_append_new(expected,json_string(w->nodes[i].bff_id));
    }
    if (!expected) { json_decref(*out); *out=NULL; return webapp_error(err,"unresolved inventory expectation"); }
    if (json_is_array(expected)) sort_values(expected);
    json_object_set_new(*out,"expected",expected);
    json_object_set_new(*out,"label",json_string(c->label));
    json_object_set_new(*out,"requirement",json_string(c->requirement));
    return ULAB_OK;
}

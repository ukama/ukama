/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include "util.h"
#include <yaml.h>
#include <ctype.h>
#include <errno.h>
#include <stdarg.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <math.h>

double wl_now(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (double)ts.tv_sec + ts.tv_nsec / 1e9;
}
int wl_error(ulab_error_t *err, const char *fmt, ...) {
    va_list ap;
    va_start(ap, fmt);
    vsnprintf(err->msg, sizeof(err->msg), fmt, ap);
    va_end(ap);
    return ULAB_ERR;
}
int wl_path(char *out, size_t n, const char *dir, const char *name,
            ulab_error_t *err) {
    if (snprintf(out, n, "%s/%s", dir, name) >= (int)n)
        return wl_error(err, "path too long");
    return ULAB_OK;
}
int wl_safe_path(const char *v) {
    if (v == NULL || *v == '\0') return 0;
    for (; *v; v++)
        if (!isalnum((unsigned char)*v) && strchr("/_-.:", *v) == NULL)
            return 0;
    return 1;
}

/* Convert a bounded YAML tree to JSON. Reject cycles/aliases and duplicate
 * keys instead of accepting ambiguous configuration. No secret expansion. */
static json_t *yaml_value(yaml_document_t *doc, int id, unsigned char *seen,
                          size_t nodes, unsigned depth, ulab_error_t *err) {
    yaml_node_t *n;
    json_t *out = NULL;
    if (depth > 32 || id <= 0 || (size_t)id > nodes || seen[id]) {
        wl_error(err, "YAML nesting/alias limit exceeded");
        return NULL;
    }
    seen[id] = 1;
    n = yaml_document_get_node(doc, id);
    if (n->type == YAML_SCALAR_NODE) {
        const char *v = (const char *)n->data.scalar.value;
        char *end;
        double number;
        if (n->data.scalar.length > 1024 * 1024) goto bad;
        if (n->data.scalar.style == YAML_PLAIN_SCALAR_STYLE) {
            if (strcmp(v, "true") == 0) return json_true();
            if (strcmp(v, "false") == 0) return json_false();
            if (strcmp(v, "null") == 0 || strcmp(v, "~") == 0) return json_null();
            errno = 0;
            number = strtod(v, &end);
            if (*v && *end == '\0' && errno == 0 && isfinite(number)) {
                if (floor(number) == number && fabs(number) < 9e15)
                    return json_integer((json_int_t)number);
                return json_real(number);
            }
        }
        return json_stringn(v, n->data.scalar.length);
    }
    if (n->type == YAML_SEQUENCE_NODE) {
        yaml_node_item_t *it;
        out = json_array();
        for (it = n->data.sequence.items.start; it < n->data.sequence.items.top; it++) {
            json_t *v = yaml_value(doc, *it, seen, nodes, depth + 1, err);
            if (v == NULL || json_array_append_new(out, v)) goto bad;
        }
    } else if (n->type == YAML_MAPPING_NODE) {
        yaml_node_pair_t *it;
        out = json_object();
        for (it = n->data.mapping.pairs.start; it < n->data.mapping.pairs.top; it++) {
            yaml_node_t *key = yaml_document_get_node(doc, it->key);
            json_t *v;
            const char *k;
            if (key->type != YAML_SCALAR_NODE) goto bad;
            k = (const char *)key->data.scalar.value;
            if (json_object_get(out, k)) {
                wl_error(err, "duplicate YAML key: %s", k);
                goto bad;
            }
            v = yaml_value(doc, it->value, seen, nodes, depth + 1, err);
            if (v == NULL || json_object_set_new(out, k, v)) goto bad;
        }
    } else goto bad;
    return out;
bad:
    json_decref(out);
    if (!err->msg[0]) wl_error(err, "invalid YAML document");
    return NULL;
}

json_t *wl_yaml_load(const char *path, ulab_error_t *err) {
    FILE *f = fopen(path, "rb");
    yaml_parser_t p;
    yaml_document_t d, trailing;
    json_t *root = NULL;
    unsigned char *seen;
    size_t nodes;
    if (!f) { wl_error(err, "cannot open %s", path); return NULL; }
    if (!yaml_parser_initialize(&p)) { fclose(f); return NULL; }
    yaml_parser_set_input_file(&p, f);
    if (!yaml_parser_load(&p, &d)) {
        wl_error(err, "%s:%zu: %s", path, p.problem_mark.line + 1,
                 p.problem ? p.problem : "invalid YAML");
        goto done;
    }
    nodes = (size_t)(d.nodes.top - d.nodes.start);
    if (nodes > 100000 || nodes == 0) {
        wl_error(err, "empty or oversized YAML");
        yaml_document_delete(&d);
        goto done;
    }
    seen = calloc(nodes + 1, 1);
    if (seen) root = yaml_value(&d, 1, seen, nodes, 0, err);
    free(seen);
    yaml_document_delete(&d);
    if (!yaml_parser_load(&p, &trailing)) {
        json_decref(root); root = NULL;
        wl_error(err, "invalid trailing YAML document");
    } else {
        if (yaml_document_get_root_node(&trailing) != NULL) {
            json_decref(root); root = NULL;
            wl_error(err, "exactly one YAML document is required");
        }
        yaml_document_delete(&trailing);
    }
done:
    yaml_parser_delete(&p);
    fclose(f);
    return root;
}

int wl_object_keys(json_t *o, const char *allowed, ulab_error_t *err) {
    const char *k;
    json_t *v;
    if (!json_is_object(o)) return wl_error(err, "expected a mapping");
    json_object_foreach(o, k, v) {
        char needle[256];
        (void)v;
        if (snprintf(needle, sizeof(needle), "|%s|", k) >= (int)sizeof(needle) ||
            strstr(allowed, needle) == NULL)
            return wl_error(err, "unknown field: %s", k);
    }
    return ULAB_OK;
}
static int number(json_t *o, const char *key, double *out, double lo,
                   double hi, int integer, ulab_error_t *err) {
    json_t *v = json_object_get(o, key);
    double d;
    if (!v) return ULAB_OK;
    d = json_number_value(v);
    if (!json_is_number(v) || !isfinite(d) || d < lo || d > hi ||
        (integer && floor(d) != d))
        return wl_error(err, "%s must be %s in [%g, %g]", key,
                         integer ? "an integer" : "a number", lo, hi);
    *out = d;
    return ULAB_OK;
}
static int count(json_t *o, const char *key, size_t *out, size_t hi,
                 ulab_error_t *err) {
    double d = (double)*out;
    if (number(o, key, &d, 0, (double)hi, 1, err)) return ULAB_ERR;
    *out = (size_t)d;
    return ULAB_OK;
}
static int string(json_t *o, const char *key, char *out, size_t n,
                   int required, ulab_error_t *err) {
    json_t *v = json_object_get(o, key);
    const char *s = json_string_value(v);
    if (!v && !required) return ULAB_OK;
    if (!s || !*s || strlen(s) >= n) return wl_error(err, "invalid %s", key);
    memcpy(out, s, strlen(s) + 1);
    return ULAB_OK;
}
static int boolean(json_t *o, const char *key, int *out, ulab_error_t *err) {
    json_t *v = json_object_get(o, key);
    if (!v) return ULAB_OK;
    if (!json_is_boolean(v)) return wl_error(err, "%s must be true or false", key);
    *out = json_is_true(v);
    return ULAB_OK;
}
static int targets(json_t *o, wl_targets_t *t, ulab_error_t *err) {
    return count(o,"sites",&t->sites,1000,err) ||
        count(o,"provisioned_ues",&t->provisioned,1000000,err) ||
        count(o,"attached_ues",&t->attached,1000000,err) ||
        count(o,"active_ues",&t->active,1000000,err) ||
        count(o,"console_sessions",&t->sessions,100000,err) ||
        number(o,"requests_per_second",&t->rps,0,1000000,0,err);
}
int workload_is_file(const char *path) {
    /* Version 2 is shared by workloads and browser scenarios. Dispatch by the
     * explicit top-level kind, never by the version or text in a description. */
    ulab_error_t err = {{0}};
    json_t *root = wl_yaml_load(path, &err);
    const char *kind = json_string_value(json_object_get(root, "kind"));
    int result = kind && !strcmp(kind, "workload");
    json_decref(root);
    return result;
}

static int environment(json_t *root, wl_config_t *c, ulab_error_t *err) {
    scenario_t *s = c->environment;
    json_t *w = json_object_get(root,"world"), *packages;
    size_t sites = 0, live = 0, records = 0, i;
    double seed = 7001;
    scenario_init(s);
    s->version = 1;
    strcpy(s->suite,"workload");
    strcpy(s->name,c->name);
    s->phase_count = 1;
    s->setup.create_networks = s->setup.create_sites = 1;
    s->setup.create_nodes = s->setup.create_node_site_links = 1;
    s->setup.create_packages = s->setup.create_subscribers = s->setup.create_sims = 1;
    s->world.networks = 1;
    s->world.sims_per_subscriber = 1;
    s->world.tower_per_site = s->world.amplifier_per_site = s->world.controller_per_site = 1;
    s->runtime.start_nodes = s->runtime.wait_nodes_ready = 1;
    if (wl_object_keys(w,"|sites|ues_per_site|backend_sims|",err) ||
        count(w,"sites",&sites,1000,err) ||
        count(w,"ues_per_site",&live,500,err) ||
        count(w,"backend_sims",&records,1000000,err) ||
        number(root,"seed",&seed,0,4294967295.0,1,err)) return ULAB_ERR;
    if (sites == 0) return wl_error(err,"world.sites must be at least one");
    s->seed = (uint32_t)seed;
    s->world.sites_per_network = (uint32_t)sites;
    s->world.ues_per_site = (uint32_t)live;
    s->world.sims_per_network = (uint32_t)records;
    if (sites * live + records > 1000000)
        return wl_error(err,"maximum planned SIM population is 1000000");
    packages = json_object_get(root,"packages");
    if (packages && !json_is_array(packages)) return wl_error(err,"packages must be a list");
    if (json_array_size(packages) > ULAB_MAX_PACKAGES) return wl_error(err,"too many packages");
    for (i = 0; i < json_array_size(packages); i++) {
        json_t *p = json_array_get(packages,i);
        package_spec_t *spec = &s->packages[s->package_count++];
        double mb = 0, days = 0, pct = 0;
        spec->active = 1;
        strcpy(spec->currency,"USD"); strcpy(spec->country,"USA");
        strcpy(spec->scope,"network");
        if (wl_object_keys(p,"|ref|name|data_mb|duration_days|amount|assign_percent|currency|country|",err) ||
            string(p,"ref",spec->ref,sizeof(spec->ref),1,err) ||
            string(p,"name",spec->name,sizeof(spec->name),1,err) ||
            string(p,"currency",spec->currency,sizeof(spec->currency),0,err) ||
            string(p,"country",spec->country,sizeof(spec->country),0,err) ||
            number(p,"data_mb",&mb,1,1e12,1,err) ||
            number(p,"duration_days",&days,1,3650,1,err) ||
            number(p,"amount",&spec->amount,0,1e8,0,err) ||
            number(p,"assign_percent",&pct,0,100,1,err)) return ULAB_ERR;
        spec->data_mb = (uint64_t)mb;
        spec->duration_days = (uint32_t)days;
        spec->assign_percent = (uint32_t)pct;
        for (size_t j = 0; j < i; j++)
            if (!strcmp(spec->ref,s->packages[j].ref)) return wl_error(err,"duplicate package ref");
    }
    return scenario_validate(s,err);
}

int wl_config_load(const char *path, const runner_opts_t *opts,
                    wl_config_t *c, ulab_error_t *err) {
    json_t *r, *o, *list;
    wl_targets_t previous;
    size_t i;
    memset(c,0,sizeof(*c));
    c->environment = calloc(1,sizeof(scenario_t));
    if (!c->environment) return wl_error(err,"out of memory");
    c->source = r = wl_yaml_load(path,err);
    if (!r) return ULAB_ERR;
    if (wl_object_keys(r,"|version|kind|name|description|seed|suite|priority|status|tags|world|packages|service_enabled|initial|console|ue|rate|limits|phases|keep_environment|",err)) return ULAB_ERR;
    if (json_integer_value(json_object_get(r,"version")) != 2 ||
        !json_is_string(json_object_get(r,"kind")) ||
        strcmp(json_string_value(json_object_get(r,"kind")),"workload"))
        return wl_error(err,"workload requires version: 2 and kind: workload");
    if (string(r,"name",c->name,sizeof(c->name),1,err)) return ULAB_ERR;
    for (const char *p=c->name; *p; p++)
        if (!islower((unsigned char)*p) && !isdigit((unsigned char)*p) && *p != '-')
            return wl_error(err,"workload name must use lowercase letters, digits and hyphens");
    ulab_copy(c->path,sizeof(c->path),path);
    ulab_copy(c->assets,sizeof(c->assets),opts->workload_assets);
    if (environment(r,c,err)) return ULAB_ERR;
    if (opts->has_seed_override) { c->environment->seed = opts->seed_override; json_object_set_new(c->source,"seed",json_integer(opts->seed_override)); }
    if (string(r,"description",c->environment->description,sizeof(c->environment->description),0,err) ||
        string(r,"status",c->environment->status,sizeof(c->environment->status),0,err) ||
        string(r,"priority",c->environment->priority,sizeof(c->environment->priority),0,err) ||
        string(r,"tags",c->environment->tags,sizeof(c->environment->tags),0,err)) return ULAB_ERR;
    if (strcmp(c->environment->status,"active") && strcmp(c->environment->status,"skip"))
        return wl_error(err,"workload status must be active or skip");
    if (json_object_get(r,"suite") && strcmp(json_string_value(json_object_get(r,"suite")) ?: "","workload"))
        return wl_error(err,"workload suite must be workload");
    c->max_inflight=256; c->ue_concurrency=8; c->provision_batch=32;
    c->request_timeout=30; c->setup_timeout=120; c->job_timeout=600;
    c->initial_timeout=3600; c->setup_max_rps=2; c->provision_rate=0.25; c->attach_rate=2;
    c->site_interval=60; c->dwell=120; c->visible_percent=100;
    c->pause_min=20; c->pause_max=60; c->traffic_mb=1; c->observe_interval=120;
    c->max_response_bytes=64*1024*1024;
    c->max_error_percent=1; c->max_missed_percent=1; c->p95_ms=3000;
    c->cdr_wait=60;
    if (boolean(r,"service_enabled",&c->service_enabled,err) ||
        boolean(r,"keep_environment",&c->keep_environment,err)) return ULAB_ERR;
    o=json_object_get(r,"initial");
    if (o && (wl_object_keys(o,"|sites|provisioned_ues|attached_ues|active_ues|console_sessions|requests_per_second|",err) || targets(o,&c->initial,err))) return ULAB_ERR;
    if (c->initial.sessions || c->initial.rps || c->initial.active)
        return wl_error(err,"initial cannot start console/rate/traffic; use the first phase");
    o=json_object_get(r,"console");
    if (o) {
        if (wl_object_keys(o,"|pages|dwell_seconds|visible_percent|optional_polling|",err) ||
            number(o,"dwell_seconds",&c->dwell,1,86400,0,err) ||
            number(o,"visible_percent",&c->visible_percent,0,100,0,err) ||
            boolean(o,"optional_polling",&c->optional_polling,err)) return ULAB_ERR;
        list=json_object_get(o,"pages");
        if (!json_is_array(list) || !json_array_size(list) || json_array_size(list)>WL_MAX_PAGES)
            return wl_error(err,"console.pages requires 1..16 profile names");
        for (i=0;i<json_array_size(list);i++) {
            const char *v=json_string_value(json_array_get(list,i));
            if (!v || !*v || strlen(v)>=ULAB_MAX_REF) return wl_error(err,"invalid page profile");
            strcpy(c->pages[c->page_count++],v);
        }
    }
    o=json_object_get(r,"ue");
    if (o) {
        size_t mb=(size_t)c->traffic_mb;
        if (wl_object_keys(o,"|attach_rate|concurrency|traffic_mb|pause_min_seconds|pause_max_seconds|observe_interval_seconds|",err) ||
            number(o,"attach_rate",&c->attach_rate,0.01,1000,0,err) ||
            count(o,"concurrency",&c->ue_concurrency,1024,err) ||
            count(o,"traffic_mb",&mb,10000,err) ||
            number(o,"pause_min_seconds",&c->pause_min,0.1,86400,0,err) ||
            number(o,"pause_max_seconds",&c->pause_max,0.1,86400,0,err) ||
            number(o,"observe_interval_seconds",&c->observe_interval,10,86400,0,err)) return ULAB_ERR;
        c->traffic_mb=mb;
    }
    if (!c->ue_concurrency || !c->traffic_mb || c->pause_max<c->pause_min)
        return wl_error(err,"invalid UE concurrency, traffic amount, or pause interval");
    o=json_object_get(r,"rate");
    if (o) {
        if (wl_object_keys(o,"|operation|variables|",err) ||
            string(o,"operation",c->rate_operation,sizeof(c->rate_operation),1,err)) return ULAB_ERR;
        c->rate_variables=json_incref(json_object_get(o,"variables"));
        if (c->rate_variables && !json_is_object(c->rate_variables)) return wl_error(err,"rate.variables must be a mapping");
    }
    o=json_object_get(r,"limits");
    if (o) {
        size_t mb=c->max_response_bytes/1024/1024, wait=c->cdr_wait;
        if (wl_object_keys(o,"|setup_max_rps|max_inflight|request_timeout_seconds|setup_timeout_seconds|job_timeout_seconds|initial_timeout_seconds|provision_rate|provision_batch|site_interval_seconds|max_response_mb|max_error_percent|max_missed_percent|p95_ms|cdr_wait_seconds|",err) ||
            number(o,"setup_max_rps",&c->setup_max_rps,0.01,1000,0,err) ||
            count(o,"max_inflight",&c->max_inflight,100000,err) ||
            count(o,"provision_batch",&c->provision_batch,1000,err) ||
            count(o,"max_response_mb",&mb,1024,err) ||
            count(o,"cdr_wait_seconds",&wait,3600,err) ||
            number(o,"request_timeout_seconds",&c->request_timeout,0.01,600,0,err) ||
            number(o,"setup_timeout_seconds",&c->setup_timeout,1,600,0,err) ||
            number(o,"job_timeout_seconds",&c->job_timeout,1,7200,0,err) ||
            number(o,"initial_timeout_seconds",&c->initial_timeout,1,604800,0,err) ||
            number(o,"provision_rate",&c->provision_rate,0.01,1000,0,err) ||
            number(o,"site_interval_seconds",&c->site_interval,0,86400,0,err) ||
            number(o,"max_error_percent",&c->max_error_percent,0,100,0,err) ||
            number(o,"max_missed_percent",&c->max_missed_percent,0,100,0,err) ||
            number(o,"p95_ms",&c->p95_ms,0,600000,0,err)) return ULAB_ERR;
        c->max_response_bytes=mb*1024*1024; c->cdr_wait=(unsigned)wait;
    }
    if (!c->max_inflight || !c->provision_batch || !c->max_response_bytes)
        return wl_error(err,"inflight, batch and response limits must be positive");
    list=json_object_get(r,"phases");
    if (!json_is_array(list) || !json_array_size(list) || json_array_size(list)>WL_MAX_PHASES)
        return wl_error(err,"phases requires 1..32 entries");
    previous=c->initial;
    for (i=0;i<json_array_size(list);i++) {
        wl_phase_t *p=&c->phases[c->phase_count++];
        o=json_array_get(list,i); p->target=previous; p->reach=600; p->hold=60;
        if (wl_object_keys(o,"|name|ramp_seconds|reach_seconds|hold_seconds|sites|provisioned_ues|attached_ues|active_ues|console_sessions|requests_per_second|",err) ||
            string(o,"name",p->name,sizeof(p->name),1,err) ||
            number(o,"ramp_seconds",&p->ramp,0,604800,0,err) ||
            number(o,"reach_seconds",&p->reach,0.1,604800,0,err) ||
            number(o,"hold_seconds",&p->hold,0.1,604800,0,err) || targets(o,&p->target,err)) return ULAB_ERR;
        if (p->reach<p->ramp || p->target.sites<previous.sites || p->target.provisioned<previous.provisioned)
            return wl_error(err,"phase %s: reach must cover ramp; sites/SIM inventory cannot decrease",p->name);
        for(size_t j=0;j<i;j++) if(!strcmp(p->name,c->phases[j].name)) return wl_error(err,"duplicate phase name");
        previous=p->target;
    }
    for (i=0;i<=c->phase_count;i++) {
        wl_targets_t *t=i?&c->phases[i-1].target:&c->initial;
        size_t live=c->environment->world.ues_per_site;
        size_t records=c->environment->world.sims_per_network;
        if(t->sites>c->environment->world.sites_per_network ||
           t->provisioned>t->sites*live+records || t->attached>t->provisioned ||
           t->attached>t->sites*live || t->active>t->attached ||
           (t->attached && !c->service_enabled) || (t->sessions && !c->page_count) ||
           (t->rps && !c->rate_operation[0]) || ((t->sessions || t->attached || t->rps) && !t->sites))
            return wl_error(err,"inconsistent population/actor target at %s",i?c->phases[i-1].name:"initial");
        if(t->sessions>c->max_sessions) c->max_sessions=t->sessions;
    }
    return ULAB_OK;
}
void wl_config_free(wl_config_t *c) {
    json_decref(c->source); json_decref(c->rate_variables); free(c->environment);
    memset(c,0,sizeof(*c));
}

/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */
#ifndef ULAB_WORKLOAD_H
#define ULAB_WORKLOAD_H

#include <curl/curl.h>
#include <jansson.h>
#include <pthread.h>
#include <signal.h>
#include <sys/types.h>
#include "runner.h"
#include "bff.h"
#include "runtime.h"

#define WL_MAX_PHASES 32
#define WL_MAX_PAGES 16
#define WL_MAX_OPS 64
#define WL_MAX_PAGE_OPS 16
#define WL_PROGRESS_SECONDS 10.0

typedef struct {
    size_t sites, provisioned, attached, active, sessions;
    double rps;
} wl_targets_t;

typedef struct {
    char name[ULAB_MAX_REF];
    double ramp, reach, hold;
    wl_targets_t target;
} wl_phase_t;

typedef struct {
    scenario_t *environment;
    json_t *source;
    char path[ULAB_MAX_PATH], assets[ULAB_MAX_PATH];
    char name[ULAB_MAX_REF];
    wl_targets_t initial;
    wl_phase_t phases[WL_MAX_PHASES];
    size_t phase_count;
    char pages[WL_MAX_PAGES][ULAB_MAX_REF];
    size_t page_count;
    char rate_operation[ULAB_MAX_REF];
    json_t *rate_variables;
    int service_enabled, optional_polling, keep_environment;
    size_t max_inflight, max_sessions, ue_concurrency, provision_batch;
    double request_timeout, setup_timeout, job_timeout, initial_timeout;
    double provision_rate, attach_rate, site_interval, setup_max_rps;
    double dwell, visible_percent, pause_min, pause_max, observe_interval;
    uint64_t traffic_mb;
    size_t max_response_bytes;
    double max_error_percent, max_missed_percent, p95_ms;
    unsigned cdr_wait;
} wl_config_t;

typedef struct {
    char id[ULAB_MAX_REF], name[ULAB_MAX_REF];
    char *query;
    json_t *variables, *required, *allowed_errors;
} wl_operation_t;
typedef struct {
    size_t operation;
    double poll, ttl;
    int optional_poll, background_poll, network_on_mount;
} wl_page_operation_t;
typedef struct {
    char name[ULAB_MAX_REF];
    wl_page_operation_t operations[WL_MAX_PAGE_OPS];
    size_t count;
} wl_page_t;
typedef struct {
    wl_operation_t operations[WL_MAX_OPS];
    size_t operation_count;
    wl_page_t pages[WL_MAX_PAGES];
    size_t page_count;
    json_t *manifest;
} wl_catalog_t;

typedef struct wl_metrics wl_metrics_t;
typedef struct {
    uint64_t completed, failed, missed;
} wl_read_counts_t;
typedef struct {
    char phase[ULAB_MAX_REF], operation[ULAB_MAX_REF], actor[ULAB_MAX_REF];
    char outcome[40], request_id[ULAB_MAX_ID];
    double scheduled, started, ended, dns_ms, connect_ms, tls_ms, ttfb_ms;
    long http_status;
    size_t bytes;
    unsigned known_gaps;
    uint64_t weight;
    char detail[ULAB_MAX_ERR];
} wl_sample_t;

typedef struct {
    CURL *easy;
    const wl_config_t *config;
    bff_client_t *auth;
    wl_metrics_t *metrics;
    _Atomic sig_atomic_t *stop;
    const char *phase;
    char last_outcome[40];
    uint64_t sequence;
    double next_request, next_progress;
} wl_sync_http_t;

typedef struct wl_http wl_http_t;
typedef struct {
    wl_config_t *config;
    const runner_opts_t *opts;
    world_t world;
    runtime_t runtime;
    bff_client_t bff;
    wl_sync_http_t http;
    wl_metrics_t *metrics;
    _Atomic sig_atomic_t *stop;
    char run_dir[ULAB_MAX_PATH];
    FILE *journal;
    pthread_t thread;
    pthread_mutex_t lock;
    int running, done, result, task_kind;
    size_t site_index, *batch, batch_count, batch_sequence;
    ulab_error_t failure;
    unsigned char *site_ready, *provisioned;
    int initialized, network_started, media_started, uncertain;
    uint64_t script_sequence;
    double last_site;
    char phase[ULAB_MAX_REF];
} wl_environment_t;

double wl_now(void);
int wl_error(ulab_error_t *err, const char *fmt, ...);
int wl_path(char *out, size_t n, const char *dir, const char *name,
            ulab_error_t *err);
int wl_safe_path(const char *value);
json_t *wl_yaml_load(const char *path, ulab_error_t *err);
int wl_object_keys(json_t *obj, const char *allowed, ulab_error_t *err);
int workload_is_file(const char *path);
int wl_config_load(const char *path, const runner_opts_t *opts,
                    wl_config_t *c, ulab_error_t *err);
void wl_config_free(wl_config_t *c);
int workload_cleanup_run(const runner_opts_t *opts);
int workload_run(const runner_opts_t *opts, int plan_only);
int wl_catalog_load(const wl_config_t *c, wl_catalog_t *cat, ulab_error_t *err);
void wl_catalog_free(wl_catalog_t *cat);
int wl_catalog_operation(const wl_catalog_t *cat, const char *id);
int wl_catalog_page(const wl_catalog_t *cat, const char *name);
json_t *wl_bind_variables(json_t *template, const world_t *world,
                          size_t site, size_t node, double mounted_at);
const char *wl_classify(long status, int curl_code, const char *body,
                        size_t len, const wl_operation_t *op,
                        json_t **root, unsigned *gaps, ulab_error_t *err);
wl_metrics_t *wl_metrics_open(const char *dir, const wl_config_t *c,
                              ulab_error_t *err);
void wl_metrics_sample(wl_metrics_t *m, const wl_sample_t *s);
void wl_metrics_read_counts(wl_metrics_t *m, wl_read_counts_t *counts);
void wl_metrics_phase_counts(wl_metrics_t *m, const char *phase, wl_read_counts_t *counts);
/* Index 0 is setup; remaining indices follow the configured phase order. */
void wl_metrics_stage_begin(wl_metrics_t *m, size_t index, int hold, double now);
void wl_metrics_stage_end(wl_metrics_t *m, size_t index, int hold, double now,
                          int complete, const wl_targets_t *actual, const char *failure);
void wl_metrics_stage_fail(wl_metrics_t *m, size_t index, int hold,
                           const char *failure);
void wl_metrics_event(wl_metrics_t *m, const char *kind, const char *phase,
                      json_t *data);
int wl_metrics_tick(wl_metrics_t *m, const char *phase,
                     const wl_targets_t *target, const wl_targets_t *actual,
                     size_t requests, size_t transfers, double now);
int wl_metrics_finish(wl_metrics_t *m, const wl_config_t *c, int result,
                       int cleanup_failed, const char *reason,
                       const char *assets, ulab_error_t *err);
void wl_metrics_close(wl_metrics_t *m);
int wl_report_write(FILE *out, const json_t *summary, size_t max_details);
int wl_sync_transport(void *ctx, const char *op, const char *query,
                       const char *variables, json_t **out, ulab_error_t *err);
wl_http_t *wl_http_open(const wl_config_t *c, bff_client_t *auth,
                        wl_metrics_t *metrics, _Atomic sig_atomic_t *stop,
                        ulab_error_t *err);
int wl_http_submit(wl_http_t *h, const wl_operation_t *op, json_t *variables,
                    const char *phase, const char *actor, double scheduled,
                    unsigned *pending, int *last_ok, ulab_error_t *err);
int wl_http_poll(wl_http_t *h, int wait_ms, ulab_error_t *err);
size_t wl_http_pending(wl_http_t *h);
void wl_http_close(wl_http_t *h);
int wl_environment_open(wl_environment_t *e, wl_config_t *c,
                        const runner_opts_t *opts, wl_metrics_t *metrics,
                        const char *run_dir, _Atomic sig_atomic_t *stop,
                        ulab_error_t *err);
int wl_environment_step(wl_environment_t *e, const wl_targets_t *target,
                        const char *phase, ulab_error_t *err);
void wl_environment_counts(wl_environment_t *e, wl_targets_t *actual);
int wl_environment_verify(wl_environment_t *e,ulab_error_t *err);
int wl_environment_cleanup(wl_environment_t *e, ulab_error_t *err);
void wl_environment_close(wl_environment_t *e);
int wl_environment_script(void *ctx, const char *script,
                           const char *args, ulab_error_t *err);
pid_t wl_process_start(const char *scripts, const char *name, const char *args,
                       const char *log_path, ulab_error_t *err);
int wl_process_wait(pid_t pid, double deadline, _Atomic sig_atomic_t *stop,
                     ulab_error_t *err);
int wl_journal(wl_environment_t *e, const char *kind, size_t index,
                const char *state, ulab_error_t *err);
int wl_schedule(wl_environment_t *e, wl_catalog_t *cat,
                  ulab_error_t *err);

#endif

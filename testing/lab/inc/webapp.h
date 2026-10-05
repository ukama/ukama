/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
#ifndef ULAB_WEBAPP_H_
#define ULAB_WEBAPP_H_

#include <jansson.h>
#include <signal.h>
#include <sys/types.h>
#include "runner.h"

#define WEBAPP_WIRE_MAX (1024u * 1024u)
typedef struct {
    pid_t pid;
    int input;
    int output;
    int broken;
    unsigned int sequence;
    double deadline;
    volatile sig_atomic_t *cancel;
    FILE *journal;
    char run_id[81];
} webapp_client_t;

typedef struct {
    json_t *root;
    world_t *world;
    char path[ULAB_MAX_PATH];
} webapp_journal_t;

/* Hooks are called serially. Runtime/final cleanup jobs are bounded child
 * processes; they must not depend on copying child memory back to the parent.
 * Provision/recover use client/journal in the parent; they must persist
 * mutation receipts before reporting acceptance. */
typedef struct {
    void *ctx;
    int (*provision)(void *ctx, webapp_client_t *client, webapp_journal_t *journal, ulab_error_t *err);
    int (*recover)(void *ctx, webapp_journal_t *journal, ulab_error_t *err);
    int (*prepare_cleanup)(void *ctx, ulab_error_t *err);
    int (*runtime_event)(void *ctx, const event_spec_t *event, ulab_error_t *err);
    int (*cleanup_resource)(void *ctx, const char *kind, const char *id,
                            ulab_error_t *err);
    int (*cleanup_runtime)(void *ctx, ulab_error_t *err);
} webapp_hooks_t;

double webapp_now(void);
int webapp_error(ulab_error_t *err, const char *fmt, ...);
int webapp_valid_run_id(const char *id);
int webapp_path(char *out, size_t len, const char *base, const char *name,
                 ulab_error_t *err);
int webapp_client_start(webapp_client_t *client, const char *executable,
                         const char *run_id, const char *run_dir,
                         double deadline, volatile sig_atomic_t *cancel,
                         ulab_error_t *err);
int webapp_call(webapp_client_t *client, const char *action, json_t *inputs,
                 unsigned int timeout_seconds, json_t **response,
                 ulab_error_t *err);
int webapp_client_stop(webapp_client_t *client, int failed, ulab_error_t *err);

int webapp_journal_open(webapp_journal_t *journal, world_t *world,
                         const char *run_dir, ulab_error_t *err);
/* Provisioning uses created=1 (UI receipt), 2 (mapped runtime claim); navigation 0. */
int webapp_journal_bind(webapp_journal_t *journal, json_t *bindings,
                         int created, unsigned int command_id, ulab_error_t *err);
int webapp_journal_save(webapp_journal_t *journal, ulab_error_t *err);
void webapp_journal_close(webapp_journal_t *journal);
int webapp_journal_cleanup(webapp_journal_t *journal,
                            const webapp_hooks_t *hooks, double deadline,
                            ulab_error_t *err);
typedef int (*webapp_job_fn)(void *ctx, ulab_error_t *err);
int webapp_bounded_job(webapp_job_fn fn, void *ctx, double deadline,
                        volatile sig_atomic_t *cancel, ulab_error_t *err);

int webapp_run_local(const runner_opts_t *opts, const scenario_t *scenario,
                     world_t *world, report_t *report, const char *run_dir, ulab_error_t *err);
int webapp_execute(const runner_opts_t *opts, const scenario_t *scenario,
                    world_t *world, report_t *report, const char *run_dir,
                    const webapp_hooks_t *hooks, ulab_error_t *err);
int webapp_event_inputs(const event_spec_t *event, world_t *world,
                         json_t **inputs, ulab_error_t *err);
json_t *webapp_check_inputs(const check_spec_t *check);

const char *webapp_commerce_kind(const event_spec_t *event);
int webapp_commerce_inputs(const event_spec_t *event, world_t *world, json_t *inputs, ulab_error_t *err);
int webapp_commerce_check(const check_spec_t *check, check_spec_t *resolved, world_t *world, json_t **inputs, ulab_error_t *err);

int webapp_inventory_event(const event_spec_t *event, world_t *world, json_t **inputs, ulab_error_t *err);
int webapp_inventory_check(const check_spec_t *check, world_t *world, json_t **inputs, ulab_error_t *err);

#endif

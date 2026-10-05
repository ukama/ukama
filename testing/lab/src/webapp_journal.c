/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include "webapp.h"
#include "util.h"
#include <ctype.h>
#include <errno.h>
#include <fcntl.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <time.h>
#include <unistd.h>

int webapp_journal_save(webapp_journal_t *j, ulab_error_t *err) {
    char tmp[ULAB_MAX_PATH + 8];
    char parent[ULAB_MAX_PATH];
    char *slash;
    int fd;
    int rc;
    FILE *f;
    snprintf(tmp, sizeof(tmp), "%s.tmp", j->path);
    fd = open(tmp, O_WRONLY | O_CREAT | O_EXCL | O_CLOEXEC, 0600);
    if (fd < 0) return webapp_error(err, "cannot create resource journal temporary file");
    f = fdopen(fd, "w");
    if (!f) { close(fd); unlink(tmp); return webapp_error(err, "cannot open resource journal"); }
    rc = json_dumpf(j->root, f, JSON_INDENT(2));
    if (fflush(f) || fsync(fd)) rc = -1;
    if (fclose(f)) rc = -1;
    if (!rc) rc = rename(tmp, j->path);
    if (rc) { unlink(tmp); return webapp_error(err, "cannot persist resource journal"); }
    ulab_copy(parent, sizeof(parent), j->path);
    slash = strrchr(parent, '/');
    if (!slash) return webapp_error(err, "resource journal has no parent directory");
    if (slash == parent) slash[1] = '\0'; else *slash = '\0';
    fd = open(parent, O_RDONLY | O_DIRECTORY | O_CLOEXEC);
    if (fd < 0) return webapp_error(err, "cannot open resource journal directory");
    rc = fsync(fd);
    if (close(fd)) rc = -1;
    if (rc) return webapp_error(err, "cannot persist resource journal directory");
    return ULAB_OK;
}
int webapp_journal_open(webapp_journal_t *j, world_t *world,
                         const char *run_dir, ulab_error_t *err) {
    memset(j, 0, sizeof(*j));
    j->world = world;
    if (webapp_path(j->path, sizeof(j->path), run_dir, "webapp-resources.json", err)) return ULAB_ERR;
    if (access(j->path, F_OK) == 0) return webapp_error(err, "resource journal already exists; use a new run ID");
    j->root = json_pack("{s:i,s:s,s:[],s:[],s:s}", "version", 1, "run_id", world->run_id,
                        "resources", "observations", "cleanup", "pending");
    if (!j->root) return webapp_error(err, "cannot allocate resource journal");
    return webapp_journal_save(j, err);
}
static int binding_target(world_t *w, const char *kind, const char *ref,
                            char **id, const char **name, ulab_error_t *err) {
    network_t *network;
    site_t *site;
    node_t *node;
    if (!strcmp(kind, "network")) {
        network = world_network_by_ref(w, ref);
        if (network) { *id = network->bff_id; *name = network->name; return ULAB_OK; }
    } else if (!strcmp(kind, "site")) {
        site = world_site_by_ref(w, ref);
        if (site) { *id = site->bff_id; *name = site->name; return ULAB_OK; }
    } else if (!strcmp(kind, "node")) {
        node = world_node_by_ref(w, ref);
        if (node) { *id = node->bff_id; *name = node->name; return ULAB_OK; }
    }
    return webapp_error(err, "worker binding references an unknown world entity");
}
int webapp_journal_bind(webapp_journal_t *j, json_t *bindings,
                         int created, unsigned int command_id, ulab_error_t *err) {
    size_t i;
    size_t k;
    json_t *binding;
    json_t *record;
    json_t *array;
    json_t *existing;
    const char *kind;
    const char *ref;
    const char *id;
    const char *source;
    const char *name = NULL;
    const char *reported_name;
    const unsigned char *p;
    char *target = NULL;
    int found;
    if (!json_is_array(bindings)) return webapp_error(err, "binding response is not an array");
    array = json_object_get(j->root, created ? "resources" : "observations");
    json_array_foreach(bindings, i, binding) {
        kind = json_string_value(json_object_get(binding, "kind"));
        ref = json_string_value(json_object_get(binding, "ref"));
        id = json_string_value(json_object_get(binding, "id"));
        source = json_string_value(json_object_get(binding, "observed_via"));
        if (!kind || !ref || !id || !source || !*id || strlen(id) >= ULAB_MAX_ID)
            return webapp_error(err, "invalid entity binding");
        for (p = (const unsigned char *)id; *p; p++)
            if (!isalnum(*p) && *p != '-' && *p != '_') return webapp_error(err, "invalid binding ID");
        if (strcmp(source, created == 2 ? "runtime_claim" : created ? "ui_result" : "url")) return webapp_error(err, "binding source does not establish requested ownership");
        if (binding_target(j->world, kind, ref, &target, &name, err) || !target || !name) return ULAB_ERR;
        if (created == 2 && (strcmp(kind, "node") || !*target || strcmp(target, id)))
            return webapp_error(err, "runtime claim requires an established factory node identity");
        if (*target && strcmp(target, id)) return webapp_error(err, "binding conflicts with an established entity ID");
        if (created) {
            reported_name = json_string_value(json_object_get(binding, "name"));
            if (!reported_name || strcmp(reported_name, name)) return webapp_error(err, "created entity name differs from the planned world name");
        } else if (!*target) {
            node_t *node;
            node = !strcmp(kind, "node") ? world_node_by_ref(j->world, ref) : NULL;
            if (!node || strcmp(node->id, id)) return webapp_error(err, "observation cannot create an unprovisioned entity binding");
        }
        found = 0;
        json_array_foreach(array, k, existing) {
            if (ulab_streq(json_string_value(json_object_get(existing, "kind")), kind) &&
                ulab_streq(json_string_value(json_object_get(existing, "ref")), ref)) {
                if (!ulab_streq(json_string_value(json_object_get(existing, "id")), id)) return webapp_error(err, "journal binding conflict");
                found = 1;
            }
        }
        if (!found) {
            record = json_pack("{s:s,s:s,s:s,s:s,s:b,s:i,s:s}", "kind", kind, "ref", ref, "id", id,
                                "name", name, "owned", created, "command_id", command_id,
                                "cleanup", created ? "pending" : "not_owned");
            if (!record || json_array_append_new(array, record)) return webapp_error(err, "cannot retain entity binding");
        }
        if (!found) json_object_set_new(record, "observed_via", json_string(source));
        /* Keep ownership in memory even if disk persistence fails, so the
         * caller can still attempt bounded cleanup before reporting failure. */
        if (created) ulab_copy(target, ULAB_MAX_ID, id);
        if (webapp_journal_save(j, err)) return ULAB_ERR;
    }
    return ULAB_OK;
}
void webapp_journal_close(webapp_journal_t *j) {
    json_decref(j->root);
    j->root = NULL;
}

int webapp_bounded_job(webapp_job_fn fn, void *ctx, double deadline,
                        volatile sig_atomic_t *cancel, ulab_error_t *err) {
    struct { int rc; ulab_error_t error; } result;
    struct timespec pause = {0, 20000000};
    int pipefd[2];
    int status;
    int timed_out = 0;
    int reaped = 0;
    pid_t pid;
    pid_t got;
    double grace;
    ssize_t count;
    if (!fn) return webapp_error(err, "required runtime/cleanup handler is unavailable");
    if (webapp_now() >= deadline || (cancel && *cancel)) return webapp_error(err, "runtime/cleanup budget exhausted");
    if (pipe(pipefd)) return webapp_error(err, "cannot create runtime result pipe");
    fflush(NULL);
    pid = fork();
    if (pid < 0) { close(pipefd[0]); close(pipefd[1]); return webapp_error(err, "cannot start bounded runtime job"); }
    if (pid == 0) {
        close(pipefd[0]);
        setpgid(0, 0);
        signal(SIGTERM, SIG_DFL); signal(SIGINT, SIG_DFL);
        memset(&result, 0, sizeof(result));
        result.rc = fn(ctx, &result.error);
        count = write(pipefd[1], &result, sizeof(result));
        close(pipefd[1]);
        _exit(count == (ssize_t)sizeof(result) ? 0 : 1);
    }
    setpgid(pid, pid);
    close(pipefd[1]);
    fcntl(pipefd[0], F_SETFL, O_NONBLOCK);
    for (;;) {
        got = waitpid(pid, &status, WNOHANG);
        if (got == pid) { reaped = 1; break; }
        if (got < 0 && errno != EINTR) break;
        if (webapp_now() >= deadline || (cancel && *cancel)) { timed_out = 1; break; }
        nanosleep(&pause, NULL);
    }
    if (!reaped) {
        kill(-pid, SIGTERM);
        grace = webapp_now() + 0.3;
        do {
            if (waitpid(pid, &status, WNOHANG) == pid) { reaped = 1; break; }
            nanosleep(&pause, NULL);
        } while (webapp_now() < grace);
    }
    kill(-pid, SIGKILL);
    if (!reaped) {
        grace = webapp_now() + 2;
        do {
            if (waitpid(pid, &status, WNOHANG) == pid) { reaped = 1; break; }
            nanosleep(&pause, NULL);
        } while (webapp_now() < grace);
    }
    count = read(pipefd[0], &result, sizeof(result));
    close(pipefd[0]);
    if (timed_out || !reaped || count != (ssize_t)sizeof(result) || !WIFEXITED(status) || WEXITSTATUS(status))
        return webapp_error(err, "runtime/cleanup job failed, timed out or was cancelled");
    if (result.rc) { *err = result.error; if (!err->msg[0]) webapp_error(err, "runtime/cleanup job failed"); }
    return result.rc;
}
typedef struct {
    const webapp_hooks_t *hooks;
    const char *kind;
    const char *id;
} cleanup_job_t;
static int delete_resource(void *arg, ulab_error_t *err) {
    cleanup_job_t *job = arg;
    return job->hooks->cleanup_resource(job->hooks->ctx, job->kind, job->id, err);
}
int webapp_journal_cleanup(webapp_journal_t *j, const webapp_hooks_t *hooks,
                            double deadline, ulab_error_t *err) {
    const char *kinds[] = {"node", "site", "network"};
    json_t *resources;
    json_t *r;
    const char *kind;
    const char *id;
    cleanup_job_t job;
    size_t k;
    size_t i;
    int failed = 0;
    int rc;
    ulab_error_t local;
    if (!j->root) return ULAB_OK;
    if (json_is_true(json_object_get(j->root, "uncertain_creation"))) {
        failed = 1; webapp_error(err, "unresolved creation requires manual reconciliation; cleanup is incomplete");
    }
    resources = json_object_get(j->root, "resources");
    for (k = 0; k < sizeof(kinds) / sizeof(kinds[0]); k++) {
        for (i = json_array_size(resources); i > 0; i--) {
            r = json_array_get(resources, i - 1);
            kind = json_string_value(json_object_get(r, "kind"));
            id = json_string_value(json_object_get(r, "id"));
            if (!json_is_true(json_object_get(r, "owned")) || !ulab_streq(kind, kinds[k]) ||
                ulab_streq(json_string_value(json_object_get(r, "cleanup")), "deleted")) continue;
            memset(&local, 0, sizeof(local));
            job.hooks = hooks; job.kind = kind; job.id = id;
            rc = !hooks || !hooks->cleanup_resource ? webapp_error(&local, "owned resource has no cleanup handler") :
                webapp_bounded_job(delete_resource, &job, deadline, NULL, &local);
            json_object_set_new(r, "cleanup", json_string(rc ? "failed" : "deleted"));
            if (rc) { failed = 1; *err = local; json_object_set_new(r, "cleanup_error", json_string(local.msg)); }
            if (webapp_journal_save(j, &local)) { failed = 1; *err = local; }
        }
    }
    json_object_set_new(j->root, "cleanup", json_string(failed ? "failed" : "complete"));
    if (webapp_journal_save(j, &local)) { failed = 1; *err = local; }
    return failed ? ULAB_ERR : ULAB_OK;
}

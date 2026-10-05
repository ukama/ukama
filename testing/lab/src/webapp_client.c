/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include "webapp.h"
#include "util.h"
#include <ctype.h>
#include <errno.h>
#include <fcntl.h>
#include <poll.h>
#include <spawn.h>
#include <stdarg.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <time.h>
#include <unistd.h>

extern char **environ;

double webapp_now(void) {
    struct timespec t;
    clock_gettime(CLOCK_MONOTONIC, &t);
    return t.tv_sec + t.tv_nsec / 1000000000.0;
}
int webapp_error(ulab_error_t *err, const char *fmt, ...) {
    va_list args;
    va_start(args, fmt);
    vsnprintf(err->msg, sizeof(err->msg), fmt, args);
    va_end(args);
    return ULAB_ERR;
}
int webapp_valid_run_id(const char *id) {
    const unsigned char *p;
    if (!id || !*id || strlen(id) > 80 || !isalnum((unsigned char)*id)) return 0;
    for (p = (const unsigned char *)id; *p; p++)
        if (!isalnum(*p) && *p != '-' && *p != '_') return 0;
    return 1;
}
int webapp_path(char *out, size_t len, const char *base, const char *name,
                 ulab_error_t *err) {
    int n;
    n = snprintf(out, len, "%s/%s", base, name);
    return n < 0 || (size_t)n >= len ? webapp_error(err, "web-app path too long") : ULAB_OK;
}
static void pause_ms(void) {
    struct timespec t = {0, 20000000};
    nanosleep(&t, NULL);
}
static int wait_fd(int fd, short events, double deadline,
                    volatile sig_atomic_t *cancel, ulab_error_t *err) {
    struct pollfd p;
    int rc;
    p.fd = fd;
    p.events = events;
    for (;;) {
        if (cancel && *cancel) return webapp_error(err, "web-app run cancelled");
        if (webapp_now() >= deadline) return webapp_error(err, "web-app command timed out");
        p.revents = 0;
        rc = poll(&p, 1, 50);
        if (rc < 0 && errno == EINTR) continue;
        if (rc < 0) return webapp_error(err, "web-app pipe poll failed");
        if (rc && (p.revents & events)) return ULAB_OK;
        if (rc && (p.revents & (POLLHUP | POLLERR | POLLNVAL)))
            return webapp_error(err, "web-app worker closed its pipe");
    }
}
static int log_message(webapp_client_t *c, const char *direction, json_t *value,
                        ulab_error_t *err) {
    json_t *record;
    char *line;
    int rc;
    record = json_pack("{s:s,s:O}", "direction", direction, "message", value);
    line = record ? json_dumps(record, JSON_COMPACT) : NULL;
    json_decref(record);
    if (!line) return webapp_error(err, "cannot encode web-app command journal");
    rc = fprintf(c->journal, "%s\n", line) < 0;
    free(line);
    if (rc || fflush(c->journal) || fsync(fileno(c->journal)))
        return webapp_error(err, "cannot persist web-app command journal");
    return ULAB_OK;
}
static int fd_prepare(int *fd) {
    int replacement;
    if (*fd < 3) {
        replacement = fcntl(*fd, F_DUPFD, 3);
        close(*fd);
        *fd = replacement;
    }
    return *fd < 0 || fcntl(*fd, F_SETFD, FD_CLOEXEC) < 0;
}
int webapp_client_start(webapp_client_t *c, const char *executable,
                         const char *run_id, const char *run_dir,
                         double deadline, volatile sig_atomic_t *cancel,
                         ulab_error_t *err) {
    int in[2] = {-1, -1};
    int out[2] = {-1, -1};
    int rc;
    int fd;
    char path[ULAB_MAX_PATH];
    posix_spawn_file_actions_t actions;
    posix_spawnattr_t attr;
    sigset_t empty;
    char *argv[] = {(char *)executable, "worker", NULL};

    memset(c, 0, sizeof(*c));
    c->input = c->output = -1;
    c->cancel = cancel;
    c->deadline = deadline;
    if (!webapp_valid_run_id(run_id)) return webapp_error(err, "invalid web-app run ID (1..80 letters/digits/hyphens/underscores)");
    ulab_copy(c->run_id, sizeof(c->run_id), run_id);
    if (webapp_path(path, sizeof(path), run_dir, "webapp-commands.jsonl", err)) return ULAB_ERR;
    fd = open(path, O_WRONLY | O_CREAT | O_EXCL | O_CLOEXEC, 0600);
    if (fd < 0) return webapp_error(err, "cannot create web-app command journal");
    c->journal = fdopen(fd, "w");
    if (!c->journal) { close(fd); return webapp_error(err, "cannot open web-app command journal"); }
    if (socketpair(AF_UNIX, SOCK_STREAM, 0, in) || pipe(out) ||
        fd_prepare(&in[0]) || fd_prepare(&in[1]) || fd_prepare(&out[0]) || fd_prepare(&out[1])) goto fail;
    if (webapp_path(path, sizeof(path), run_dir, "webapp-worker.log", err)) goto fail;
    if (posix_spawn_file_actions_init(&actions)) goto fail;
    rc = posix_spawn_file_actions_adddup2(&actions, in[1], STDIN_FILENO);
    rc |= posix_spawn_file_actions_adddup2(&actions, out[1], STDOUT_FILENO);
    rc |= posix_spawn_file_actions_addopen(&actions, STDERR_FILENO, path, O_WRONLY | O_CREAT | O_EXCL, 0600);
    rc |= posix_spawn_file_actions_addclose(&actions, in[0]);
    rc |= posix_spawn_file_actions_addclose(&actions, in[1]);
    rc |= posix_spawn_file_actions_addclose(&actions, out[0]);
    rc |= posix_spawn_file_actions_addclose(&actions, out[1]);
    if (rc || posix_spawnattr_init(&attr)) { posix_spawn_file_actions_destroy(&actions); goto fail; }
    sigemptyset(&empty);
    rc = posix_spawnattr_setsigmask(&attr, &empty);
    rc |= posix_spawnattr_setflags(&attr, POSIX_SPAWN_SETPGROUP | POSIX_SPAWN_SETSIGMASK);
    rc |= posix_spawnattr_setpgroup(&attr, 0);
    if (!rc) rc = posix_spawnp(&c->pid, executable, &actions, &attr, argv, environ);
    posix_spawnattr_destroy(&attr);
    posix_spawn_file_actions_destroy(&actions);
    if (rc) goto fail;
    close(in[1]); close(out[1]);
    c->input = in[0]; c->output = out[0];
    if (fcntl(c->input, F_SETFL, O_NONBLOCK) < 0 || fcntl(c->output, F_SETFL, O_NONBLOCK) < 0) {
        c->broken = 1;
        webapp_client_stop(c, 1, err);
        return webapp_error(err, "cannot configure worker pipes");
    }
    return ULAB_OK;
fail:
    if (in[0] >= 0) close(in[0]);
    if (in[1] >= 0) close(in[1]);
    if (out[0] >= 0) close(out[0]);
    if (out[1] >= 0) close(out[1]);
    fclose(c->journal); c->journal = NULL;
    return webapp_error(err, "cannot launch web-app worker; check --webapp-worker and its build");
}
static int send_line(webapp_client_t *c, const char *line, double deadline,
                       ulab_error_t *err) {
    size_t offset = 0;
    size_t len = strlen(line);
    ssize_t n;
    while (offset < len + 1) {
        if (wait_fd(c->input, POLLOUT, deadline, c->cancel, err)) return ULAB_ERR;
        n = send(c->input, offset < len ? line + offset : "\n",
                  offset < len ? len - offset : 1, MSG_NOSIGNAL);
        if (n < 0 && (errno == EAGAIN || errno == EINTR)) continue;
        if (n <= 0) return webapp_error(err, "cannot write web-app command");
        offset += (size_t)n;
    }
    return ULAB_OK;
}
static int read_response(webapp_client_t *c, double deadline, json_t **reply,
                          ulab_error_t *err) {
    char *buffer;
    size_t used = 0;
    size_t i;
    ssize_t n;
    json_error_t jerr;
    int rc = ULAB_ERR;
    buffer = malloc(WEBAPP_WIRE_MAX + 1);
    if (!buffer) return webapp_error(err, "cannot allocate web-app response buffer");
    while (used < WEBAPP_WIRE_MAX) {
        if (wait_fd(c->output, POLLIN, deadline, c->cancel, err)) goto done;
        n = read(c->output, buffer + used, WEBAPP_WIRE_MAX - used);
        if (n < 0 && (errno == EINTR || errno == EAGAIN)) continue;
        if (n <= 0) { webapp_error(err, "worker exited before replying"); goto done; }
        for (i = used; i < used + (size_t)n; i++) {
            if (buffer[i] != '\n') continue;
            if (i + 1 != used + (size_t)n) { webapp_error(err, "unsolicited worker output after JSONL response"); goto done; }
            *reply = json_loadb(buffer, i, JSON_REJECT_DUPLICATES, &jerr);
            if (!json_is_object(*reply)) { json_decref(*reply); *reply = NULL; webapp_error(err, "malformed worker JSON response"); goto done; }
            rc = ULAB_OK;
            goto done;
        }
        used += (size_t)n;
    }
    webapp_error(err, "worker response exceeds 1 MiB");
done:
    free(buffer);
    return rc;
}
static int same_string(json_t *v, const char *s) {
    return json_is_string(v) && !strcmp(json_string_value(v), s);
}
int webapp_call(webapp_client_t *c, const char *action, json_t *inputs,
                 unsigned int seconds, json_t **response, ulab_error_t *err) {
    double now;
    double remaining;
    double end;
    struct timespec utc;
    json_int_t epoch;
    json_t *command;
    json_t *status;
    json_t *run_status;
    const char *message;
    char *line;
    int rc;
    *response = NULL;
    if (!c->pid || c->broken) return webapp_error(err, "web-app worker is unavailable");
    now = webapp_now(); remaining = c->deadline - now;
    if (remaining > seconds) remaining = seconds;
    if (remaining <= 0 || (c->cancel && *c->cancel)) return webapp_error(err, "web-app run cancelled or scenario deadline exceeded");
    if (seconds == 0 || seconds > 900 || ++c->sequence > 10000) return webapp_error(err, "invalid web-app command budget/sequence");
    end = now + remaining;
    clock_gettime(CLOCK_REALTIME, &utc);
    epoch = (json_int_t)utc.tv_sec * 1000 + utc.tv_nsec / 1000000 + (json_int_t)(remaining * 1000);
    command = json_pack("{s:i,s:s,s:i,s:s,s:I,s:O}", "protocol", 1, "run_id", c->run_id,
                        "command_id", c->sequence, "action", action, "deadline_ms", epoch, "inputs", inputs);
    line = command ? json_dumps(command, JSON_COMPACT) : NULL;
    if (!line || strlen(line) >= WEBAPP_WIRE_MAX) {
        free(line); json_decref(command);
        return webapp_error(err, "cannot encode bounded web-app command");
    }
    rc = log_message(c, "request", command, err);
    json_decref(command);
    if (!rc) rc = send_line(c, line, end, err);
    free(line);
    /* A late successful assertion is never accepted. The extra wait allows
     * the worker to retain failure artifacts after its acceptance deadline. */
    if (!rc) rc = read_response(c, end + 10, response, err);
    if (rc) { c->broken = 1; return rc; }
    if (log_message(c, "response", *response, err)) { c->broken = 1; return ULAB_ERR; }
    status = json_object_get(*response, "status");
    run_status = json_object_get(*response, "run_status");
    if (json_integer_value(json_object_get(*response, "protocol")) != 1 ||
        !same_string(json_object_get(*response, "run_id"), c->run_id) ||
        json_integer_value(json_object_get(*response, "command_id")) != c->sequence ||
        !same_string(json_object_get(*response, "action"), action) ||
        (!same_string(status, "ok") && !same_string(status, "error")) ||
        !json_is_array(json_object_get(*response, "bindings")) ||
        !json_is_array(json_object_get(*response, "artifacts")) ||
        !json_is_integer(json_object_get(*response, "duration_ms")) ||
        json_integer_value(json_object_get(*response, "duration_ms")) < 0 ||
        !json_object_get(*response, "expected") || !json_object_get(*response, "actual")) {
        c->broken = 1;
        return webapp_error(err, "worker response violates protocol/correlation contract");
    }
    if (same_string(status, "error")) {
        if (!same_string(run_status, "failed")) { c->broken = 1; return webapp_error(err, "failed worker response has invalid run status"); }
        message = json_string_value(json_object_get(json_object_get(*response, "error"), "message"));
        return webapp_error(err, "web-app: %.900s", message ? message : "worker operation failed");
    }
    if ((!strcmp(action, "close") && !same_string(run_status, "passed") && !same_string(run_status, "failed")) ||
        (strcmp(action, "close") && !same_string(run_status, "running"))) {
        c->broken = 1;
        return webapp_error(err, "worker returned an invalid successful run status");
    }
    if (webapp_now() > end && strcmp(action, "close")) return webapp_error(err, "worker success arrived after command deadline");
    return ULAB_OK;
}
int webapp_client_stop(webapp_client_t *c, int failed, ulab_error_t *err) {
    double end;
    int status = 0;
    int rc = ULAB_OK;
    int reaped = 0;
    pid_t got;
    json_t *reply = NULL;
    json_t *inputs;
    if (c->pid <= 0) goto finish;
    if (!c->broken) {
        c->cancel = NULL;
        c->deadline = webapp_now() + 10;
        inputs = json_pack("{s:b,s:s}", "failed", failed, "reason", failed ? "C_RUNNER_FAILED" : "EXPLICIT_CLOSE");
        if (webapp_call(c, "close", inputs, 10, &reply, err)) rc = ULAB_ERR;
        if (!failed && reply && same_string(json_object_get(reply, "run_status"), "failed"))
            rc = webapp_error(err, "worker reported failure while closing");
        json_decref(inputs); json_decref(reply);
    } else {
        /* Let the worker retain evidence before terminating browser helpers. */
        kill(c->pid, SIGTERM);
        rc = webapp_error(err, "worker protocol lost; shutdown acknowledgement unavailable");
    }
    if (c->input >= 0) { close(c->input); c->input = -1; }
    end = webapp_now() + (c->broken ? 10 : 2);
    do {
        got = waitpid(c->pid, &status, WNOHANG);
        if (got == c->pid) { reaped = 1; break; }
        if (got < 0 && errno != EINTR) break;
        pause_ms();
    } while (webapp_now() < end);
    /* Reap and terminate any same-group helpers even if the leader exited. */
    kill(-c->pid, SIGKILL);
    if (!reaped) {
        end = webapp_now() + 2;
        do {
            got = waitpid(c->pid, &status, WNOHANG);
            if (got == c->pid) { reaped = 1; break; }
            pause_ms();
        } while (webapp_now() < end);
        rc = webapp_error(err, "worker required forced termination");
    }
    if (reaped && (!WIFEXITED(status) || (WEXITSTATUS(status) != 0 && !(failed && WEXITSTATUS(status) == 1))))
        rc = webapp_error(err, "web-app worker exited abnormally (status=%d)", status);
    c->pid = 0;
finish:
    if (c->input >= 0) close(c->input);
    if (c->output >= 0) close(c->output);
    c->input = c->output = -1;
    if (c->journal) { if (fclose(c->journal)) rc = webapp_error(err, "cannot close web-app journal"); c->journal = NULL; }
    return rc;
}

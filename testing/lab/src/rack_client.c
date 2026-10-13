/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */
#define _GNU_SOURCE
#include <ctype.h>
#include <errno.h>
#include <fcntl.h>
#include <poll.h>
#include <pthread.h>
#include <signal.h>
#include <spawn.h>
#include <stdatomic.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <time.h>
#include <unistd.h>

#include "rack_client.h"
#include "util.h"

extern char **environ;

struct rack_client {
    hardware_config_t config;
    char session[65];
    char prefix[40];
    char **environment;
    char ssh_log[ULAB_MAX_PATH];
    FILE *log;
    atomic_ulong sequence;
    atomic_int stop;
    atomic_int lost;
    pthread_t thread;
    int thread_started;
    unsigned int interval_ms;
};

static uint64_t now_ms(void) {
    struct timespec t;
    clock_gettime(CLOCK_MONOTONIC, &t);
    return (uint64_t)t.tv_sec * 1000u + (uint64_t)t.tv_nsec / 1000000u;
}

static int safe_token(const char *s) {
    size_t i;
    if (s == NULL || s[0] == '\0' || strlen(s) > 64)
        return 0;
    for (i = 0; s[i]; i++) {
        if (!isalnum((unsigned char)s[i]) && s[i] != '-' && s[i] != '_')
            return 0;
    }
    return 1;
}

static int ssh_call(rack_client_t *c, const char *command, unsigned int timeout, char *output,
                    size_t capacity, ulab_error_t *err) {
    char target[ULAB_MAX_NAME + ULAB_MAX_REF + 2];
    char known[ULAB_MAX_PATH + 32];
    char port[16];
    char *argv[28];
    posix_spawn_file_actions_t actions;
    posix_spawnattr_t attributes;
    int pipefd[2];
    int status;
    int rc;
    int done;
    int eof;
    size_t used;
    size_t argc;
    pid_t pid;
    uint64_t deadline;

    snprintf(target, sizeof(target), "%s@%s", c->config.user, c->config.host);
    snprintf(port, sizeof(port), "%u", c->config.ssh_port);
    snprintf(known, sizeof(known), "UserKnownHostsFile=\"%s\"", c->config.known_hosts_file);
    argc = 0;
    argv[argc++] = "ssh";
    argv[argc++] = "-F";
    argv[argc++] = "/dev/null";
    argv[argc++] = "-T";
    argv[argc++] = "-o";
    argv[argc++] = "BatchMode=yes";
    argv[argc++] = "-o";
    argv[argc++] = "StrictHostKeyChecking=yes";
    argv[argc++] = "-o";
    argv[argc++] = "IdentitiesOnly=yes";
    argv[argc++] = "-o";
    argv[argc++] = "ConnectTimeout=5";
    argv[argc++] = "-o";
    argv[argc++] = "GlobalKnownHostsFile=/dev/null";
    argv[argc++] = "-o";
    argv[argc++] = known;
    argv[argc++] = "-i";
    argv[argc++] = c->config.identity_file;
    argv[argc++] = "-p";
    argv[argc++] = port;
    argv[argc++] = target;
    argv[argc++] = (char *)command;
    argv[argc] = NULL;
    output[0] = '\0';
    if (pipe2(pipefd, O_CLOEXEC) != 0)
        goto setup_error;
    if (posix_spawn_file_actions_init(&actions) != 0) {
        close(pipefd[0]);
        close(pipefd[1]);
        goto setup_error;
    }
    if (posix_spawnattr_init(&attributes) != 0) {
        posix_spawn_file_actions_destroy(&actions);
        close(pipefd[0]);
        close(pipefd[1]);
        goto setup_error;
    }
    rc = posix_spawn_file_actions_addopen(&actions, STDIN_FILENO, "/dev/null", O_RDONLY, 0);
    rc |= posix_spawn_file_actions_adddup2(&actions, pipefd[1], STDOUT_FILENO);
    rc |= posix_spawn_file_actions_addopen(&actions, STDERR_FILENO, c->ssh_log,
                                           O_WRONLY | O_CREAT | O_APPEND, 0600);
    rc |= posix_spawn_file_actions_addclose(&actions, pipefd[0]);
    rc |= posix_spawn_file_actions_addclose(&actions, pipefd[1]);
    rc |= posix_spawnattr_setflags(&attributes, POSIX_SPAWN_SETPGROUP);
    rc |= posix_spawnattr_setpgroup(&attributes, 0);
    if (rc == 0)
        rc = posix_spawn(&pid, c->config.ssh_path, &actions, &attributes, argv, c->environment);
    posix_spawn_file_actions_destroy(&actions);
    posix_spawnattr_destroy(&attributes);
    close(pipefd[1]);
    if (rc != 0) {
        close(pipefd[0]);
        goto setup_error;
    }
    if (fcntl(pipefd[0], F_SETFL, O_NONBLOCK) < 0) {
        close(pipefd[0]);
        kill(-pid, SIGKILL);
        while (waitpid(pid, NULL, 0) < 0 && errno == EINTR) {
        }
        goto setup_error;
    }
    deadline = now_ms() + timeout;
    used = 0;
    done = eof = status = 0;
    rc = ULAB_OK;
    while (!done || !eof) {
        struct pollfd pollfd;
        ssize_t count;
        pid_t waited;

        if (now_ms() >= deadline) {
            snprintf(err->msg, sizeof(err->msg),
                     "RF-rack command timed out; outcome may be unknown");
            rc = ULAB_ERR;
            break;
        }
        pollfd.fd = pipefd[0];
        pollfd.events = POLLIN;
        pollfd.revents = 0;
        if (!eof)
            (void)poll(&pollfd, 1, 50);
        else {
            struct timespec pause = {0, 10000000};
            nanosleep(&pause, NULL);
        }
        if (!eof) {
            count = read(pipefd[0], output + used, capacity - used - 1);
            if (count > 0) {
                used += (size_t)count;
                if (used == capacity - 1) {
                    snprintf(err->msg, sizeof(err->msg), "RF-rack response exceeds buffer");
                    rc = ULAB_ERR;
                    break;
                }
            } else if (count == 0)
                eof = 1;
            else if (errno != EAGAIN && errno != EINTR) {
                snprintf(err->msg, sizeof(err->msg), "RF-rack response read failed");
                rc = ULAB_ERR;
                break;
            }
        }
        if (!done) {
            waited = waitpid(pid, &status, WNOHANG);
            if (waited == pid)
                done = 1;
            else if (waited < 0 && errno != EINTR) {
                snprintf(err->msg, sizeof(err->msg), "RF-rack command wait failed");
                rc = ULAB_ERR;
                break;
            }
        }
    }
    output[used] = '\0';
    close(pipefd[0]);
    if (rc != ULAB_OK)
        kill(-pid, SIGKILL);
    if (!done)
        while (waitpid(pid, &status, 0) < 0 && errno == EINTR) {
        }
    if (rc == ULAB_OK && (!WIFEXITED(status) || WEXITSTATUS(status) != 0)) {
        snprintf(err->msg, sizeof(err->msg), "RF-rack SSH/rackctl failed; see rack-ssh.log");
        rc = ULAB_ERR;
    }
    return rc;
setup_error:
    snprintf(err->msg, sizeof(err->msg), "cannot start RF-rack SSH command");
    return ULAB_ERR;
}

static int call(rack_client_t *c, const char *command, int authenticated, unsigned int timeout,
                char request_id[65], json_t **result, ulab_error_t *err) {
    char remote[1024];
    char output[65536];
    json_t *root;
    json_t *value;
    const char *returned_id;
    const char *message;
    size_t i;
    int n;
    int rc;

    *result = NULL;
    /* Remote ssh commands go through the remote login shell. Accept only the
     * fixed rackctl grammar, generated IDs, numeric slots and numeric amounts. */
    for (i = 0; command[i]; i++) {
        if (!isalnum((unsigned char)command[i]) && strchr(" -_", command[i]) == NULL) {
            snprintf(err->msg, sizeof(err->msg), "invalid RF-rack command token");
            return ULAB_ERR;
        }
    }
    snprintf(request_id, 65, "%s-%lu", c->prefix, atomic_fetch_add(&c->sequence, 1));
    n = snprintf(remote, sizeof(remote), "rackctl --json --request-id %s %s%s %s", request_id,
                 authenticated ? "--session " : "", authenticated ? c->session : "", command);
    if (n < 0 || (size_t)n >= sizeof(remote)) {
        snprintf(err->msg, sizeof(err->msg), "RF-rack command too long");
        return ULAB_ERR;
    }
    rc = ssh_call(c, remote, timeout, output, sizeof(output), err);
    root = json_loads(output, JSON_REJECT_DUPLICATES, NULL);
    returned_id = json_string_value(json_object_get(root, "request_id"));
    if (json_integer_value(json_object_get(root, "version")) != 1 || returned_id == NULL ||
        strcmp(returned_id, request_id) != 0) {
        if (rc == ULAB_OK)
            snprintf(err->msg, sizeof(err->msg), "invalid/stale RF-rack response");
        rc = ULAB_ERR;
    } else if (!json_is_true(json_object_get(root, "ok"))) {
        message = json_string_value(json_object_get(json_object_get(root, "error"), "message"));
        snprintf(err->msg, sizeof(err->msg), "RF-rack rejected %.128s: %.700s", command,
                 message ? message : "invalid error response");
        rc = ULAB_ERR;
    } else if (rc == ULAB_OK) {
        value = json_object_get(root, "result");
        if (value == NULL) {
            snprintf(err->msg, sizeof(err->msg), "RF-rack response has no result");
            rc = ULAB_ERR;
        } else
            *result = json_incref(value);
    }
    json_decref(root);
    if (c->log) {
        fprintf(c->log, "request=%s command=%s rc=%d\n", request_id, command, rc);
        fflush(c->log);
    }
    return rc;
}

int rack_client_healthy(rack_client_t *c, ulab_error_t *err) {
    if (c == NULL || atomic_load(&c->lost)) {
        snprintf(err->msg, sizeof(err->msg), "RF-rack session unavailable or keepalive failed");
        return ULAB_ERR;
    }
    return ULAB_OK;
}

int rack_client_call(rack_client_t *c, const char *command, char request_id[65], json_t **result,
                     ulab_error_t *err) {
    *result = NULL;
    if (rack_client_healthy(c, err))
        return ULAB_ERR;
    return call(c, command, 1, c->config.command_timeout_seconds * 1000u, request_id, result, err);
}

static void *keepalive(void *arg) {
    rack_client_t *c = arg;
    uint64_t next = now_ms() + c->interval_ms;
    struct timespec pause = {0, 100000000};
    while (!atomic_load(&c->stop)) {
        if (now_ms() >= next) {
            char id[65];
            json_t *result = NULL;
            ulab_error_t err = {{0}};
            if (call(c, "session keepalive", 1, c->interval_ms, id, &result, &err) ||
                json_integer_value(json_object_get(result, "expires_in_ms")) <= c->interval_ms) {
                json_decref(result);
                atomic_store(&c->lost, 1);
                return NULL;
            }
            json_decref(result);
            next = now_ms() + c->interval_ms;
        }
        nanosleep(&pause, NULL);
    }
    return NULL;
}

int rack_client_open(rack_client_t **out, const hardware_config_t *config, const char *run_dir,
                     ulab_error_t *err) {
    rack_client_t *c;
    struct timespec now;
    char command[96];
    char path[ULAB_MAX_PATH];
    char id[65];
    json_t *result = NULL;
    const char *session;
    json_int_t ttl;
    int rc;
    size_t i;
    size_t count;

    *out = NULL;
    c = calloc(1, sizeof(*c));
    if (c == NULL)
        return ULAB_EINTERNAL;
    c->config = *config;
    /* Event timeout overrides use setenv. Snapshot once before the keepalive
     * thread starts, so concurrent child launches never read mutable environ. */
    for (count = 0; environ[count]; count++) {
    }
    c->environment = calloc(count + 1, sizeof(*c->environment));
    if (c->environment == NULL) {
        free(c);
        return ULAB_EINTERNAL;
    }
    for (i = 0; i < count; i++) {
        c->environment[i] = strdup(environ[i]);
        if (c->environment[i] == NULL) {
            while (i)
                free(c->environment[--i]);
            free(c->environment);
            free(c);
            return ULAB_EINTERNAL;
        }
    }
    atomic_init(&c->stop, 0);
    atomic_init(&c->lost, 0);
    atomic_init(&c->sequence, 1);
    clock_gettime(CLOCK_REALTIME, &now);
    snprintf(c->prefix, sizeof(c->prefix), "lab-%x-%llx", (unsigned int)getpid(),
             (unsigned long long)now.tv_sec * 1000000000u + (unsigned long long)now.tv_nsec);
    if (snprintf(c->ssh_log, sizeof(c->ssh_log), "%s/rack-ssh.log", run_dir) >=
            (int)sizeof(c->ssh_log) ||
        snprintf(path, sizeof(path), "%s/rack.log", run_dir) >= (int)sizeof(path)) {
        rack_client_close(&c, err);
        snprintf(err->msg, sizeof(err->msg), "rack log path too long");
        return ULAB_ERR;
    }
    c->log = fopen(path, "ae");
    if (c->log == NULL) {
        rack_client_close(&c, err);
        snprintf(err->msg, sizeof(err->msg), "cannot open rack log");
        return ULAB_ERR;
    }
    snprintf(command, sizeof(command), "session acquire %s", c->prefix);
    rc = call(c, command, 0, config->command_timeout_seconds * 1000u, id, &result, err);
    session = json_string_value(json_object_get(result, "session_id"));
    ttl = json_integer_value(json_object_get(result, "ttl_ms"));
    if (rc == ULAB_OK && (!safe_token(session) || ttl < 3000 || ttl > 600000)) {
        snprintf(err->msg, sizeof(err->msg),
                 "RF-rack requires a valid session and TTL of at least 3000ms");
        rc = ULAB_ERR;
    }
    if (safe_token(session))
        ulab_copy(c->session, sizeof(c->session), session);
    json_decref(result);
    if (rc == ULAB_OK) {
        c->interval_ms = (unsigned int)ttl / 3u;
        if (c->interval_ms > 2000u)
            c->interval_ms = 2000u;
        if (pthread_create(&c->thread, NULL, keepalive, c) != 0) {
            snprintf(err->msg, sizeof(err->msg), "cannot start RF-rack keepalive");
            rc = ULAB_ERR;
        } else
            c->thread_started = 1;
    }
    if (rc != ULAB_OK) {
        ulab_error_t ignored = {{0}};
        rack_client_close(&c, &ignored);
        return rc;
    }
    *out = c;
    return ULAB_OK;
}

int rack_client_close(rack_client_t **client, ulab_error_t *err) {
    rack_client_t *c = *client;
    char id[65];
    json_t *result = NULL;
    int rc = ULAB_OK;
    size_t i;
    if (c == NULL)
        return ULAB_OK;
    atomic_store(&c->stop, 1);
    if (c->thread_started)
        pthread_join(c->thread, NULL);
    if (c->session[0])
        rc = call(c, "session release", 1, c->config.command_timeout_seconds * 1000u, id, &result,
                  err);
    json_decref(result);
    if (c->log)
        fclose(c->log);
    for (i = 0; c->environment && c->environment[i]; i++)
        free(c->environment[i]);
    free(c->environment);
    free(c);
    *client = NULL;
    return rc;
}

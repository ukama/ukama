/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <errno.h>
#include <getopt.h>
#include <pthread.h>
#include <signal.h>
#include <sqlite3.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <ulfius.h>
#include <unistd.h>

#include "agent.h"
#include "common/utils.h"
#include "db.h"
#include "log.h"
#include "network.h"
#include "package_cache.h"
#include "tasks.h"
#include "wimc.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_services.h"
#include "usys_string.h"
#include "usys_types.h"

#include "version.h"

#define AGENT_EXEC_NAME "agent"
#define AGENT_STOP_TIMEOUT_SEC 5
#define AGENT_RESTART_LIMIT 3

typedef struct {
    char method[WIMC_MAX_NAME_LEN];
    char service[WIMC_MAX_NAME_LEN];
    char execPath[WIMC_MAX_PATH_LEN];
    pid_t pid;
    int port;
    int running;
    int restartCount;
} ManagedAgent;

typedef struct {
    ManagedAgent agents[MAX_AGENTS];
    int count;
} AgentManager;

static volatile sig_atomic_t gTerminate = 0;
static volatile sig_atomic_t gChildEvent = 0;

static void handle_signal(int signum) {

    (void)signum;
    gTerminate = 1;
}

static void handle_child_signal(int signum) {

    (void)signum;
    gChildEvent = 1;
}

static UsysOption longOptions[] = {
    { "logs",    required_argument, 0, 'l' },
    { "url",     required_argument, 0, 'u' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

static void set_log_level(char *slevel) {

    int ilevel;

    ilevel = USYS_LOG_TRACE;

    if (slevel == NULL) {
        return;
    }

    if (!strcmp(slevel, "TRACE")) {
        ilevel = USYS_LOG_TRACE;
    } else if (!strcmp(slevel, "DEBUG")) {
        ilevel = USYS_LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = USYS_LOG_INFO;
    }

    usys_log_set_level(ilevel);
}

static void usage(void) {

    usys_puts("Usage: wimc.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-u, --url                     Hub URL");
    usys_puts("-v, --version                 Software version");
}

static char *trim_token(char *value) {

    char *end;

    if (value == NULL) {
        return NULL;
    }

    while (*value == ' ' || *value == '\t' || *value == '\n' ||
           *value == '\r') {
        value++;
    }

    end = value + strlen(value);
    while (end > value && (*(end - 1) == ' ' ||
           *(end - 1) == '\t' || *(end - 1) == '\n' ||
           *(end - 1) == '\r')) {
        end--;
    }

    *end = '\0';
    return value;
}

static int is_supported_method(const char *method) {

    if (method == NULL || *method == '\0') {
        return USYS_FALSE;
    }

    if (strcmp(method, WIMC_METHOD_CHUNK_STR) == 0 ||
        strcmp(method, WIMC_METHOD_TARGZ_STR) == 0) {
        return USYS_TRUE;
    }

    return USYS_FALSE;
}

static int build_agent_service_name(const char *method,
                                    char *service,
                                    size_t serviceLen) {

    if (method == NULL || service == NULL || serviceLen == 0) {
        return -1;
    }

    if (snprintf(service, serviceLen, "wimc-agent-%s", method) >=
        (int)serviceLen) {
        return -1;
    }

    return 0;
}

static int agent_manager_has_method(AgentManager *mgr,
                                    const char *method) {

    int i;

    if (mgr == NULL || method == NULL) {
        return USYS_FALSE;
    }

    for (i = 0; i < mgr->count; i++) {
        if (strcmp(mgr->agents[i].method, method) == 0) {
            return USYS_TRUE;
        }
    }

    return USYS_FALSE;
}

static int find_agent_from_env(char *execPath, size_t execPathLen) {

    const char *envPath;

    envPath = getenv(WIMC_AGENT_PATH_ENV);
    if (envPath == NULL || *envPath == '\0') {
        return -1;
    }

    if (snprintf(execPath, execPathLen, "%s", envPath) >=
        (int)execPathLen) {
        return -1;
    }

    if (access(execPath, X_OK) == 0) {
        return 0;
    }

    usys_log_error("%s is set but not executable: %s",
                   WIMC_AGENT_PATH_ENV, execPath);

    return -1;
}

static int find_agent_from_self(char *execPath, size_t execPathLen) {

    char selfPath[WIMC_MAX_PATH_LEN];
    char *slash;
    ssize_t len;

    memset(selfPath, 0, sizeof(selfPath));

    len = readlink("/proc/self/exe", selfPath, sizeof(selfPath) - 1);
    if (len <= 0 || len >= (ssize_t)sizeof(selfPath)) {
        return -1;
    }
    selfPath[len] = '\0';

    slash = strrchr(selfPath, '/');
    if (slash == NULL) {
        return -1;
    }

    *(slash + 1) = '\0';

    if (snprintf(execPath, execPathLen, "%s%s", selfPath,
                 AGENT_EXEC_NAME) >= (int)execPathLen) {
        return -1;
    }

    if (access(execPath, X_OK) == 0) {
        return 0;
    }

    return -1;
}

static int find_agent_exec(char *execPath, size_t execPathLen) {

    if (execPath == NULL || execPathLen == 0) {
        return -1;
    }

    if (find_agent_from_env(execPath, execPathLen) == 0) {
        return 0;
    }

    if (getenv(WIMC_AGENT_PATH_ENV) != NULL) {
        return -1;
    }

    if (find_agent_from_self(execPath, execPathLen) == 0) {
        return 0;
    }

    usys_log_error("Unable to find sibling WIMC agent. Set %s if needed",
                   WIMC_AGENT_PATH_ENV);

    return -1;
}

static int agent_manager_add(AgentManager *mgr,
                             const char *method,
                             const char *execPath) {

    ManagedAgent *agent;
    char service[WIMC_MAX_NAME_LEN];
    int port;

    if (mgr == NULL || method == NULL || execPath == NULL) {
        return -1;
    }

    if (!is_supported_method(method)) {
        usys_log_error("Unsupported WIMC agent method: %s", method);
        return -1;
    }

    if (agent_manager_has_method(mgr, method)) {
        return 0;
    }

    if (mgr->count >= MAX_AGENTS) {
        usys_log_error("Too many managed WIMC agents");
        return -1;
    }

    if (build_agent_service_name(method, service, sizeof(service)) != 0) {
        usys_log_error("Unable to build service name for method: %s",
                       method);
        return -1;
    }

    port = usys_find_service_port(service);
    if (port <= 0) {
        usys_log_error("Unable to find service port for %s", service);
        return -1;
    }

    agent = &mgr->agents[mgr->count];
    memset(agent, 0, sizeof(*agent));

    snprintf(agent->method, sizeof(agent->method), "%s", method);
    snprintf(agent->service, sizeof(agent->service), "%s", service);
    snprintf(agent->execPath, sizeof(agent->execPath), "%s", execPath);
    agent->port = port;
    agent->pid = -1;
    agent->running = USYS_FALSE;
    agent->restartCount = 0;

    mgr->count++;
    return 0;
}

static int agent_manager_load(AgentManager *mgr) {

    char priority[WIMC_MAX_ARGS_LEN];
    char execPath[WIMC_MAX_PATH_LEN];
    char *savePtr;
    char *token;
    char *method;
    const char *env;

    if (mgr == NULL) {
        return -1;
    }

    memset(mgr, 0, sizeof(*mgr));
    memset(priority, 0, sizeof(priority));
    memset(execPath, 0, sizeof(execPath));

    if (find_agent_exec(execPath, sizeof(execPath)) != 0) {
        return -1;
    }

    env = getenv(WIMC_METHOD_PRIORITY_ENV);
    if (env == NULL || *env == '\0') {
        env = WIMC_METHOD_TARGZ_STR "," WIMC_METHOD_CHUNK_STR;
    }

    snprintf(priority, sizeof(priority), "%s", env);

    savePtr = NULL;
    token = strtok_r(priority, ",", &savePtr);
    while (token != NULL) {
        method = trim_token(token);

        if (method != NULL && *method != '\0') {
            if (agent_manager_add(mgr, method, execPath) != 0) {
                usys_log_error("Skipping WIMC agent method: %s", method);
            }
        }

        token = strtok_r(NULL, ",", &savePtr);
    }

    if (mgr->count <= 0) {
        usys_log_error("No WIMC agents configured");
        return -1;
    }

    return 0;
}

static int agent_manager_start_one(ManagedAgent *agent, char *logLevel) {

    pid_t pid;

    if (agent == NULL || agent->running) {
        return 0;
    }

    pid = fork();
    if (pid < 0) {
        usys_log_error("Failed to fork WIMC agent %s: %s",
                       agent->method, strerror(errno));
        return -1;
    }

    if (pid == 0) {
        execl(agent->execPath,
              AGENT_EXEC_NAME,
              "-m", agent->method,
              "-l", logLevel ? logLevel : DEF_LOG_LEVEL,
              NULL);
        _exit(127);
    }

    agent->pid = pid;
    agent->running = USYS_TRUE;

    usys_log_debug("Started WIMC agent method=%s service=%s port=%d pid=%d",
                   agent->method, agent->service, agent->port, pid);

    return 0;
}

static int agent_manager_start(AgentManager *mgr, char *logLevel) {

    int i;
    int started;

    if (mgr == NULL) {
        return -1;
    }

    started = 0;

    for (i = 0; i < mgr->count; i++) {
        if (agent_manager_start_one(&mgr->agents[i], logLevel) == 0) {
            started++;
        }
    }

    if (started <= 0) {
        return -1;
    }

    return 0;
}

static ManagedAgent *agent_manager_find_by_pid(AgentManager *mgr,
                                               pid_t pid) {

    int i;

    if (mgr == NULL || pid <= 0) {
        return NULL;
    }

    for (i = 0; i < mgr->count; i++) {
        if (mgr->agents[i].pid == pid) {
            return &mgr->agents[i];
        }
    }

    return NULL;
}

static void agent_manager_reap(AgentManager *mgr, char *logLevel) {

    ManagedAgent *agent;
    int status;
    pid_t pid;

    if (mgr == NULL) {
        return;
    }

    while (1) {
        pid = waitpid(-1, &status, WNOHANG);
        if (pid <= 0) {
            break;
        }

        agent = agent_manager_find_by_pid(mgr, pid);
        if (agent == NULL) {
            continue;
        }

        agent->running = USYS_FALSE;
        agent->pid = -1;

        if (WIFEXITED(status)) {
            usys_log_error("WIMC agent %s exited with code %d",
                           agent->method, WEXITSTATUS(status));
        } else if (WIFSIGNALED(status)) {
            usys_log_error("WIMC agent %s killed by signal %d",
                           agent->method, WTERMSIG(status));
        } else {
            usys_log_error("WIMC agent %s exited", agent->method);
        }

        if (gTerminate) {
            continue;
        }

        if (agent->restartCount >= AGENT_RESTART_LIMIT) {
            usys_log_error("WIMC agent %s restart limit reached",
                           agent->method);
            continue;
        }

        agent->restartCount++;
        usys_log_error("Restarting WIMC agent %s attempt %d",
                       agent->method, agent->restartCount);
        agent_manager_start_one(agent, logLevel);
    }

    gChildEvent = 0;
}

static void agent_manager_stop(AgentManager *mgr) {

    int i;
    int status;
    int elapsed;
    int remaining;

    if (mgr == NULL) {
        return;
    }

    remaining = 0;

    for (i = 0; i < mgr->count; i++) {
        if (mgr->agents[i].running && mgr->agents[i].pid > 0) {
            kill(mgr->agents[i].pid, SIGTERM);
            remaining++;
        }
    }

    elapsed = 0;
    while (remaining > 0 && elapsed < AGENT_STOP_TIMEOUT_SEC) {
        remaining = 0;

        for (i = 0; i < mgr->count; i++) {
            if (!mgr->agents[i].running || mgr->agents[i].pid <= 0) {
                continue;
            }

            if (waitpid(mgr->agents[i].pid, &status, WNOHANG) ==
                mgr->agents[i].pid) {
                mgr->agents[i].running = USYS_FALSE;
                mgr->agents[i].pid = -1;
            } else {
                remaining++;
            }
        }

        if (remaining > 0) {
            sleep(1);
            elapsed++;
        }
    }

    for (i = 0; i < mgr->count; i++) {
        if (mgr->agents[i].running && mgr->agents[i].pid > 0) {
            usys_log_error("Force killing WIMC agent %s pid=%d",
                           mgr->agents[i].method, mgr->agents[i].pid);
            kill(mgr->agents[i].pid, SIGKILL);
            waitpid(mgr->agents[i].pid, &status, 0);
            mgr->agents[i].running = USYS_FALSE;
            mgr->agents[i].pid = -1;
        }
    }
}

int main(int argc, char **argv) {

    int opt;
    int optIdx;
    int rc;
    int taskMutexInit;
    int dbMutexInit;
    int curlInit;
    int webStarted;
    int agentsStarted;
    int wimcPort;
    int ukamaPort;
    Agent *agents;
    WTasks *tasks;
    char *debug;
    char hubURL[WIMC_MAX_URL_LEN];
    UInst serviceInst;
    Config serviceConfig;
    AgentManager agentManager;

    rc = EXIT_FAILURE;
    taskMutexInit = 0;
    dbMutexInit = 0;
    curlInit = 0;
    webStarted = 0;
    agentsStarted = 0;
    agents = NULL;
    tasks = NULL;
    debug = DEF_LOG_LEVEL;

    memset(hubURL, 0, sizeof(hubURL));
    memset(&serviceInst, 0, sizeof(serviceInst));
    memset(&serviceConfig, 0, sizeof(serviceConfig));
    memset(&agentManager, 0, sizeof(agentManager));

    usys_log_set_service(SERVICE_NAME);
//    usys_log_remote_init(SERVICE_NAME);

    wimcPort = usys_find_service_port(SERVICE_NAME);
    if (wimcPort == 0) {
        usys_log_error("Unable to find service port for %s", SERVICE_NAME);
        goto cleanup;
    }

    ukamaPort = usys_find_service_port(SERVICE_UKAMA);
    if (ukamaPort == 0) {
        usys_log_error("Unable to find service port for %s", SERVICE_UKAMA);
        goto cleanup;
    }

    if (snprintf(hubURL, sizeof(hubURL), "http://localhost:%d",
                 ukamaPort) >= (int)sizeof(hubURL)) {
        usys_log_error("Hub URL too long");
        goto cleanup;
    }

    while (USYS_TRUE) {
        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "hvl:u:", longOptions, &optIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            rc = EXIT_SUCCESS;
            goto cleanup;

        case 'v':
            usys_puts(VERSION);
            rc = EXIT_SUCCESS;
            goto cleanup;

        case 'l':
            if (optarg != NULL && *optarg != '\0') {
                debug = optarg;
                set_log_level(debug);
            }
            break;

        case 'u':
            if (optarg == NULL || *optarg == '\0') {
                usage();
                goto cleanup;
            }

            if (snprintf(hubURL, sizeof(hubURL), "%s", optarg) >=
                (int)sizeof(hubURL)) {
                usys_log_error("Hub URL too long");
                goto cleanup;
            }
            break;

        default:
            usage();
            goto cleanup;
        }
    }

    serviceConfig.servicePort = wimcPort;
    serviceConfig.dbFile = strdup(WIMC_DB_PATH);
    serviceConfig.hubURL = strdup(hubURL);

    if (serviceConfig.dbFile == NULL || serviceConfig.hubURL == NULL) {
        usys_log_error("Memory allocation failure");
        goto cleanup;
    }

    signal(SIGINT, handle_signal);
    signal(SIGTERM, handle_signal);
    signal(SIGCHLD, handle_child_signal);

    usys_log_debug("Starting %s", SERVICE_NAME);

    agents = (Agent *)calloc(MAX_AGENTS, sizeof(Agent));
    if (agents == NULL) {
        usys_log_error("Memory failure");
        goto cleanup;
    }

    serviceConfig.agents = &agents;
    serviceConfig.tasks = &tasks;

    if (db_open_or_create(serviceConfig.dbFile, &serviceConfig.db) != 0) {
        usys_log_error("Unable to open/create DB file: %s",
                       serviceConfig.dbFile);
        goto cleanup;
    }

    if (pthread_mutex_init(&serviceConfig.taskMutex, NULL) != 0) {
        usys_log_error("taskMutex init failed");
        goto cleanup;
    }
    taskMutexInit = 1;

    if (pthread_mutex_init(&serviceConfig.dbMutex, NULL) != 0) {
        usys_log_error("dbMutex init failed");
        goto cleanup;
    }
    dbMutexInit = 1;

    if (curl_global_init(CURL_GLOBAL_ALL) != 0) {
        usys_log_error("curl_global_init failed");
        goto cleanup;
    }
    curlInit = 1;

    if (pkg_reconcile_startup(serviceConfig.db, DEFAULT_APPS_PKGS_PATH) != 0) {
        usys_log_error("Package cache startup reconcile failed");
        goto cleanup;
    }

    if (start_web_service(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup");
        goto cleanup;
    }
    webStarted = 1;

    if (agent_manager_load(&agentManager) != 0) {
        usys_log_error("Failed to load WIMC agent configuration");
        goto cleanup;
    }

    if (agent_manager_start(&agentManager, debug) != 0) {
        usys_log_error("Failed to start WIMC agents");
        goto cleanup;
    }
    agentsStarted = 1;

    while (!gTerminate) {
        sleep(1);

        if (gChildEvent) {
            agent_manager_reap(&agentManager, debug);
        }
    }

    rc = EXIT_SUCCESS;

cleanup:
    if (agentsStarted) {
        agent_manager_stop(&agentManager);
    }

    if (webStarted) {
        ulfius_stop_framework(&serviceInst);
        ulfius_clean_instance(&serviceInst);
    }

    if (serviceConfig.db != NULL) {
        sqlite3_close(serviceConfig.db);
        serviceConfig.db = NULL;
    }

    if (dbMutexInit) {
        pthread_mutex_destroy(&serviceConfig.dbMutex);
    }

    if (taskMutexInit) {
        pthread_mutex_destroy(&serviceConfig.taskMutex);
    }

    if (curlInit) {
        curl_global_cleanup();
    }

    clear_tasks(&tasks);
    clear_agents(agents);

    free(agents);
    free(serviceConfig.dbFile);
    free(serviceConfig.hubURL);

    return rc;
}

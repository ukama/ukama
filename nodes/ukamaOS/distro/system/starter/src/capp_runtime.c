/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <jansson.h>
#include <pthread.h>

#include "starter.h"
#include "config.h"

#include "usys_types.h"
#include "usys_log.h"

extern char **environ;

static int token_count(char *str) {

    int count=0, i;

    if (str == NULL)  return count;
    if (!strlen(str)) return count;

    for (i=0; i<strlen(str); i++) {
        if (str[i] == ' ') {
            count++;
        }
    }

    return count++;
}

static void log_runtime(CappRuntime *runtime, int argc, int envc) {

    int i;

    usys_log_debug("Executing capp: binary: %s argc: %d env: %d",
                   runtime->cmd, argc, envc);

    for (i=0; i<argc; i++) {
        if (runtime->argv[i]) {
            usys_log_debug("\t argc-%d: %s", i, runtime->argv[i]);
        }
    }

    for (i=0; i<envc; i++) {
        if (runtime->env[i]) {
            usys_log_debug("\t envc-%d: %s", i, runtime->env[i]);
        }
    }
}

static int create_args_list(char *exec, char *arg, char ***argv) {

    int count=0, i=0;
    char *str, *token;

    if (arg == NULL) return 0;

    count = token_count(arg);
    if (exec) {
        count += 2; /* for binary and NULL, see exec(3) */
        usys_log_debug("Creating args list for capp: %s", exec);
    } else {
        count += 1; /* for env variables. */
        usys_log_debug("Creating env list for capp");
    }

    *argv = (char **)calloc(count, sizeof(char *));
    if (argv == NULL) {
        usys_log_error("Error allocating memory of size: %d",
                       count * sizeof(char *));
        return USYS_FALSE;
    }

    str = strdup(arg);

    if (exec) { /* only need for runtime arguments. */
        (*argv)[0] = strdup(exec);
        i++;
    }

    token = strtok(str, " ");
    while (token != NULL) {
        (*argv)[i] = strdup(token);
        token = strtok(NULL, " ");
        i++;
    }
    (*argv)[i] = (char *) NULL;

    return count;
}

static bool execute_capp(void *arg) {

    Capp *capp = (Capp *)arg;
    CappRuntime *runtime;
    pid_t pid;

    runtime = capp->runtime;

    usys_log_debug("Executing command: %s", runtime->cmd);

    runtime->status = CAPP_RUNTIME_EXEC;

    pid = fork();

    if (pid == -1) {
        /* error executing fork */
        usys_log_error("Error executing fork for capp: %s:%s",
                       capp->name, capp->tag);
        return USYS_FALSE;
    } else if (pid == 0) {

        execvpe(runtime->cmd, runtime->argv, runtime->env);
        exit(0);
    } else {

        int status;
        runtime->pid = pid;
        waitpid(pid, &status, 0);
        runtime->status = CAPP_RUNTIME_DONE;
    }

    return USYS_TRUE;
}

static bool setup_and_execute_capp(Capp *capp, int *error) {

    int argc=0, envc=0;
    CappProc *process=NULL;
    CappRuntime *runtime=NULL;
    pthread_t thread;

    if (!capp || !capp->config) return USYS_FALSE;

    runtime      = (CappRuntime *)calloc(1, sizeof(CappRuntime));
    process      = capp->config->process;
    runtime->cmd = strdup(process->exec);

    /* Runtime arguments list */
    if ((argc = create_args_list(process->exec,
                                 process->argv,
                                 &runtime->argv)) == USYS_FALSE) {
        log_error("Error creating argument list for capp execution: %s",
                  process->exec);
        return USYS_FALSE;
    }

    /* Environment varaibles list */
    if (process->env) {
        if ((envc = create_args_list(NULL,
                                     process->env,
                                     &runtime->env)) == USYS_FALSE) {
            usys_log_error("Error creating env list for capp execution: %s",
                           process->exec);
            return USYS_FALSE;
        }
    }

    capp->runtime = runtime;
    log_runtime(capp->runtime, argc, envc);

    /* thread up */
    pthread_create(&thread,
                   NULL,
                   execute_capp,
                   (void *)capp);

    pthread_join(thread, NULL);

    return USYS_TRUE;
}

static bool create_and_run_capps(Capp *capp, int *error) {

    char *configFileName = NULL;

    if (capp->rootfs == NULL) return USYS_FALSE;

    configFileName = (char *)calloc(1, strlen(capp->rootfs) +
                                    strlen(DEF_CAPP_CONFIG_FILE) + 2);
    sprintf(configFileName, "%s/%s", capp->rootfs, DEF_CAPP_CONFIG_FILE);

    if (!process_capp_config_file(&capp->config, configFileName)) {
        free(configFileName);
        return USYS_FALSE;
    }

    return setup_and_execute_capp(capp, error);
}

void run_space_all_capps(Space *space) {

    CappList *cappList = NULL;
    int error=0;
    
    for (cappList=space->cappList;
         cappList;
         cappList=cappList->next) {

        if (cappList->capp->fetch == CAPP_PKG_NOT_FOUND) continue;

        /* skip if already running or done */
        if (cappList->capp->runtime != NULL) {
            if (cappList->capp->runtime->status == CAPP_RUNTIME_EXEC ||
                cappList->capp->runtime->status == CAPP_RUNTIME_DONE)
                continue;
        }

        if (create_and_run_capps(cappList->capp, &error)) {
            usys_log_error("Unable to execute capp: %s:%s Error: %d",
                           cappList->capp->name,
                           cappList->capp->tag,
                           error);
            error = 0;
            continue;
        }

        usys_log_debug("Executing capp: %s:%s",
                       cappList->capp->name,
                       cappList->capp->tag);
    }
}

void fetch_unpack_run(Space *space, Config *config) {

    CappList *cappList = NULL;
    Capp     *capp = NULL;
    char     *path = NULL;
    int      ret=0;
    int      httpStatus=0;

    char runMe[MAX_BUFFER]      = {0};

    for (cappList=space->cappList;
         cappList;
         cappList=cappList->next) {

        if (cappList->capp->fetch == CAPP_PKG_FOUND) continue;

        capp = cappList->capp;

        /* get the file from wimc.d */
        if (get_capp_path(config, capp->name, capp->tag,
                          &path, &httpStatus) == USYS_NOK) {
            log_error("Error getting path for capp: %s:%s",
                      capp->name, capp->tag);
            continue;
        }

        /* set the fetch flag to avoid fetching the pkg again */
        cappList->capp->fetch = CAPP_PKG_FOUND;

        /* Move file from path to DEF_CAPP_PATH */
        sprintf(runMe, "/bin/cp %s/%s_%s.tar.gz %s",
                path,
                capp->name,
                capp->tag,
                DEF_CAPP_PATH);
        log_debug("Running command: %s", runMe);
        if ((ret = system(runMe)) < 0) {
            usys_log_error("Unable to execute cmd %s for space: %s Code: %d",
                           runMe, space->name, ret);
            continue;
        }

        /* copy the capp file to the space rootfs */
        copy_capp_to_space_rootfs(space->name,
                                  cappList->capp,
                                  DEF_CAPP_PATH,
                                  DEF_SPACE_ROOTFS_PATH);
    }

    /* Now run them all */
    run_space_all_capps(space);
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <jansson.h>
#include <pthread.h>
#include <signal.h>

#include "starter.h"
#include "config.h"
#include "web_client.h"

#include "usys_types.h"
#include "usys_log.h"

/* space.c */
void copy_capp_to_space_rootfs(char *spaceName,
                               Capp *capp,
                               char *sPath,
                               char *dPath);

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

static void* execute_capp(void *arg) {

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

        if (execvpe(runtime->cmd, runtime->argv, runtime->env) == -1) {
            perror("execvpe error");
            exit(EXIT_FAILURE);
        }
        exit(EXIT_SUCCESS);
    } else {

        int status;

        runtime->pid = pid;
        waitpid(pid, &status, 0);

        if (WIFEXITED(status)) {
            usys_log_debug("%s: exited normally with status %d",
                           capp->name, WEXITSTATUS(status));
            runtime->status = CAPP_RUNTIME_DONE;
        } else if (WIFSIGNALED(status)) {
            usys_log_debug("%s: terminated by signal: %d",
                           capp->name, WTERMSIG(status));
            runtime->status = CAPP_RUNTIME_FAILURE;
        } else if (WIFSTOPPED(status)) {
            usys_log_debug("%s: stopped by signal: %d",
                           capp->name, WSTOPSIG(status));
            runtime->status = CAPP_RUNTIME_FAILURE;
        } else {
            usys_log_debug("%s: uknown exit status: %d",
                           capp->name, status);
            runtime->status = CAPP_RUNTIME_UNKNOWN;
        }
    }

    return NULL;
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

    pthread_create(&thread,
                   NULL,
                   execute_capp,
                   (void *)capp);

    return USYS_TRUE;
}

static bool copy_file(const char *srcPath, const char *destPath) {

    char runMe[MAX_BUFFER] = {0};

    usys_log_debug("Copying from %s to %s", srcPath, destPath);

    sprintf(runMe, "/bin/cp -p %s %s", srcPath, destPath);
    if (system(runMe) != 0) {
        usys_log_error("Unable to cp from: %s to: %s", srcPath, destPath);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static bool copy_folder(char *srcFolder, char *destFolder) {

    DIR *dir;
    struct dirent *entry;
    struct stat destStat;

    dir = opendir(srcFolder);
    if (!dir) {
        usys_log_error("Unable to open folder: %s", srcFolder);
        return USYS_FALSE;
    }

    /* check if destination folder exist */
    if (stat(destFolder, &destStat) != 0) {
        if (mkdir(destFolder, 0777) != 0) {
            usys_log_error("Unable to create dest folder: %s", destFolder);
            return USYS_FALSE;
        }
    }

    while ((entry = readdir(dir))) {
        if (entry->d_type == DT_REG) {
            char srcPath[1024];
            char destPath[1024];

            snprintf(srcPath, sizeof(srcPath), "%s/%s",
                     srcFolder, entry->d_name);
            snprintf(destPath, sizeof(destPath), "%s/%s",
                     destFolder, entry->d_name);

            if (!copy_file(srcPath, destPath)) {
                usys_log_error("Unable to copy from: %s to: %s",
                               srcPath, destPath);
                return USYS_FALSE;
            }
        }
    }

    closedir(dir);

    return USYS_TRUE;
}

static bool install_capp(char *appName, char *rootPath) {

    char srcConfigFolder[MAX_BUFFER]  = {0};
    char destConfigFolder[MAX_BUFFER] = {0};
    char sbinFolder[MAX_BUFFER]       = {0};
    char libFolder[MAX_BUFFER]        = {0};

    sprintf(srcConfigFolder,  "%s/conf/", rootPath);
    sprintf(destConfigFolder, "/ukama/configs/%s/", appName);
    sprintf(sbinFolder,       "%s/sbin/", rootPath);
    sprintf(libFolder,        "%s/lib/", rootPath);

    if (copy_folder(srcConfigFolder, destConfigFolder) == USYS_FALSE) {
        usys_log_debug("No config for %s. Skipping", rootPath);
    }

    if (copy_folder(sbinFolder, "/sbin") == USYS_FALSE) {
        usys_log_error("No binary files for %s", rootPath);
        return USYS_FALSE;
    }

    if (copy_folder(libFolder, "/ukama/apps/lib") == USYS_FALSE) {
        usys_log_error("No lib files for %s", rootPath);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static bool create_and_run_capps(Capp *capp, int *error) {

    char *configFileName = NULL;

    if (capp->rootfs == NULL) return USYS_FALSE;

    configFileName = (char *)calloc(1, strlen(capp->rootfs) +
                                    strlen(DEF_CAPP_CONFIG_FILE) + 2);
    sprintf(configFileName, "%s/%s", capp->rootfs, DEF_CAPP_CONFIG_FILE);

    if (!install_capp(capp->name, capp->rootfs)) {
        usys_log_error("Unable to install app at config and sbin");
        free(configFileName);
        return USYS_FALSE;
    }

    if (!process_capp_config_file(&capp->config, configFileName)) {
        free(configFileName);
        return USYS_FALSE;
    }

    return setup_and_execute_capp(capp, error);
}

static CappList *reverse_link_list(CappList *head) {

    CappList *prev = NULL;
    CappList *current = head;
    CappList *next = NULL;

    while (current != NULL) {
        next = current->next;
        current->next = prev;
        prev = current;
        current = next;
    }

    return prev;
}

static bool is_capp_running(Capp *capp, int wait, int maxRetry) {

    int count = 0;

    do {
        if (ping_capp(capp->name)) {
            usys_log_debug("%s capp state is RUN/Active", capp->name);
            return USYS_TRUE;
        }

        sleep(wait);
        count ++;
    } while (count < maxRetry);

    usys_log_debug("%s capp run max out. Retried: %d timeout: %d",
                   capp->name, maxRetry, wait);
    return USYS_FALSE;
}

static void terminate_capp(Capp *capp) {

    CappRuntime *runtime;

    runtime = capp->runtime;

    if (runtime == NULL)   return;
    if (runtime->pid == 0) return;

    kill(runtime->pid, SIGKILL);
    usys_log_debug("%s capp killed via SIGKILL", capp->name);

    /* reset runtime */
    runtime->status = CAPP_RUNTIME_PEND;
    runtime->pid    = 0;

}

static bool check_for_dependant_capps(Space *space, char *name) {

    CappList *cappList = NULL;

    for (cappList=space->cappList; cappList; cappList=cappList->next) {

        /* ignore itself */
        if (strcasecmp(cappList->capp->name, name) == 0) continue;

        if (cappList->capp->depend != NULL) {
            if (strcmp(cappList->capp->depend->name, name) == 0 &&
                strcmp(cappList->capp->depend->state, STATE_DONE) == 0) {
                return USYS_TRUE;
            }
        }
    }

    return USYS_FALSE;
}

void run_space_all_capps(Space *space) {

    CappList *cappList = NULL;
    int error=0;

    space->cappList = reverse_link_list(space->cappList);

    for (cappList=space->cappList; cappList; cappList=cappList->next) {

    retry:
        if (cappList->capp->fetch == CAPP_PKG_NOT_FOUND) continue;

        /* skip if already running or done */
        if (cappList->capp->runtime != NULL) {
            if (cappList->capp->runtime->status == CAPP_RUNTIME_EXEC ||
                cappList->capp->runtime->status == CAPP_RUNTIME_DONE)
                continue;
        }
        
        if (create_and_run_capps(cappList->capp, &error) == USYS_FALSE) {
            usys_log_error("Unable to execute capp: %s:%s Error: %d",
                           cappList->capp->name,
                           cappList->capp->tag,
                           error);
            error = 0;
            continue;
        }

        /* for 'boot' space, each capp must be in run (or done) state
         * before moving we move to execute the next one */
        if (strcasecmp(space->name, SPACE_BOOT) == 0) {

            int status;

            /* any capp waiting on this to be done? */
            if (check_for_dependant_capps(space, cappList->capp->name)) {

                waitpid(cappList->capp->runtime->pid, &status, 0);
                if (cappList->capp->runtime->status != CAPP_RUNTIME_DONE) {
                    /* retry */
                    cappList->capp->runtime->status = CAPP_RUNTIME_PEND;
                    cappList->capp->runtime->pid    = 0;

                    usys_log_error("%s capp failed. Retrying",
                                   cappList->capp->name);

                    goto retry;
                } else {
                    continue;
                }
            }
            
            /* make sure capp is running - 200 on /v1/ping or its done */
            if (is_capp_running(cappList->capp, WAIT_TIME, MAX_RETRIES)) {
                usys_log_debug("Executing capp: %s:%s",
                               cappList->capp->name,
                               cappList->capp->tag);
                continue;
            }

            /* Unable to get the capp running, kill and retry.
             * Could get stuck here forever */
            terminate_capp(cappList->capp);
            goto retry;
        }
    }
}

void fetch_unpack_run(Space *space, Config *config) {

    CappList *cappList = NULL;
    Capp     *capp = NULL;
    char     *path = NULL;
    int      ret=0;
    int      httpStatus=0;

    char runMe[MAX_BUFFER] = {0};

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

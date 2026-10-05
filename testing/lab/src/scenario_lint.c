/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "scenario.h"
#include "util.h"
#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

typedef struct {
    scenario_t *scenario;
    size_t files;
    size_t failed;
} lint_t;

static void lint_path(lint_t *lint, const char *path, unsigned int depth) {
    struct stat st;
    struct dirent **entries;
    int count;
    int i;
    char child[ULAB_MAX_PATH];
    ulab_error_t err;

    if (depth > 32 || lstat(path, &st) != 0) {
        fprintf(stderr, "LINT FAIL %s: %s\n", path,
                depth > 32 ? "directory nesting exceeds 32" : strerror(errno));
        lint->failed++;
        return;
    }
    if (S_ISDIR(st.st_mode)) {
        count = scandir(path, &entries, NULL, alphasort);
        if (count < 0) {
            fprintf(stderr, "LINT FAIL %s: %s\n", path, strerror(errno));
            lint->failed++;
            return;
        }
        for (i = 0; i < count; i++) {
            if (entries[i]->d_name[0] != '.') {
                if (snprintf(child, sizeof(child), "%s/%s", path,
                             entries[i]->d_name) >= (int)sizeof(child)) {
                    fprintf(stderr, "LINT FAIL %s: child path too long\n", path);
                    lint->failed++;
                } else {
                    lint_path(lint, child, depth + 1);
                }
            }
            free(entries[i]);
        }
        free(entries);
        return;
    }
    if (!S_ISREG(st.st_mode)) {
        fprintf(stderr, "LINT FAIL %s: expected a regular file/directory\n", path);
        lint->failed++;
        return;
    }
    if (!ulab_ends(path, ".yaml") && !ulab_ends(path, ".yml")) {
        if (depth == 0) {
            fprintf(stderr, "LINT FAIL %s: expected a .yaml/.yml scenario\n", path);
            lint->failed++;
        }
        return;
    }
    lint->files++;
    memset(&err, 0, sizeof(err));
    if (scenario_load(path, lint->scenario, &err) ||
        scenario_validate(lint->scenario, &err)) {
        fprintf(stderr, "LINT FAIL %s: %s\n", path, err.msg);
        lint->failed++;
    } else {
        printf("LINT OK   %s (v%u; syntax/contract only)\n", path,
               lint->scenario->version);
    }
    fflush(stdout);
}

int scenario_lint_main(int argc, char **argv) {
    lint_t lint;
    int i;
    int rc;
    memset(&lint, 0, sizeof(lint));
    if (argc < 1) {
        fprintf(stderr, "usage: ukama-lab lint <scenario.yaml|dir> [...]\n");
        return ULAB_EUSAGE;
    }
    lint.scenario = calloc(1, sizeof(*lint.scenario));
    if (lint.scenario == NULL) return ULAB_EINTERNAL;
    for (i = 0; i < argc; i++) lint_path(&lint, argv[i], 0);
    printf("LINT      files=%zu errors=%zu (no execution or coverage credit)\n",
           lint.files, lint.failed);
    rc = lint.files > 0 && lint.failed == 0 ? ULAB_OK : ULAB_ESCENARIO;
    free(lint.scenario);
    return rc;
}

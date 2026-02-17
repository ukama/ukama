/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

#include "board_config.h"

static void trim(char *s) {

    char *p = s;
    int l = strlen(p);

    while (l > 0 && isspace(p[l-1])) {
        p[--l] = 0;
    }

    while (*p && isspace(*p)) ++p, --l;

    memmove(s, p, l+1);
}

static int parse_yesno(const char *val) {

    if (!val) return 0;
    if (strcasecmp(val, "yes") == 0)  return 1;
    if (strcasecmp(val, "true") == 0) return 1;

    return 0;
}

static int find_entry(BoardConfig *cfg, const char *key) {
    for (int i = 0; i < cfg->count; i++) {
        if (strcmp(cfg->entries[i].key, key) == 0)
            return i;
    }

    return -1;
}

NodeType detect_node_type(const char *nodeID) {

    if (!nodeID) return NODE_UNKNOWN;

    if (strstr(nodeID, "tnode")) {
        return NODE_TOWER;
    }

    if (strstr(nodeID, "anode")) {
        return NODE_AMPLIFIER;
    }

    return NODE_UNKNOWN;
}

static int load_board_file(BoardConfig *cfg, const char *filePath) {

    char line[MAX_SIZE];

    FILE *fp = fopen(filePath, "r");
    if (!fp) {
        return 0;
    }

    while (fgets(line, sizeof(line), fp)) {

        char key[128];
        char val[128];
        int enabled, idx;

        if (line[0] == '#' || strlen(line) < 3)
            continue;

        char *eq = strchr(line, '=');
        if (!eq)
            continue;

        *eq = 0;

        strncpy(key, line,   sizeof(key)-1);
        strncpy(val, eq + 1, sizeof(val)-1);

        key[sizeof(key)-1] = 0;
        val[sizeof(val)-1] = 0;

        trim(key);
        trim(val);

        enabled = parse_yesno(val);
        idx = find_entry(cfg, key);

        if (idx >= 0) {
            /* override existing (board overrides common) */
            cfg->entries[idx].enabled = enabled;
        } else {
            if (cfg->count < MAX_BOARD_APPS) {
                strncpy(cfg->entries[cfg->count].key,
                        key,
                        MAX_APP_KEY-1);

                cfg->entries[cfg->count].enabled = enabled;
                cfg->count++;
            }
        }
    }

    fclose(fp);
    return 1;
}


int board_config_load(BoardConfig *cfg,
                      const char *boardsDir,
                      NodeType type) {

    char path[MAX_SIZE];
    const char *filename = NULL;

    if (!cfg || !boardsDir) {
        return 0;
    }

    memset(cfg, 0, sizeof(*cfg));
    cfg->type = type;

    switch (type) {
    case NODE_TOWER:
        filename = NODE_TOWER_CONFIG;
        break;

    case NODE_AMPLIFIER:
        filename = NODE_AMPLIFIER_CONFIG;
        break;

    default:
        return 0;
    }

    /* Construct full path: boardsDir/filename */
    if (snprintf(path, sizeof(path), "%s/%s", boardsDir, filename) >= sizeof(path)) {
        /* truncated */
        return 0;
    }

    return load_board_file(cfg, path);
}

int board_is_app_enabled(BoardConfig *cfg,
                         const char *appName) {

    char key[MAX_SIZE];
    
    if (!cfg || !appName)
        return 0;

    snprintf(key, sizeof(key), "%s_APP", appName);
    for (int i = 0; i < cfg->count; i++) {
        if (strcasecmp(cfg->entries[i].key, key) == 0)
            return cfg->entries[i].enabled;
    }

    return 0;
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

#include "usys_log.h"

#define MAX_BUFFER_SIZE  1024
#define MAX_LINE_LENGTH   256
#define MAX_KEY_LENGTH    128
#define MAX_VALUE_LENGTH  128

typedef struct BoardConfigEntry {

    char key[MAX_KEY_LENGTH];
    char value[MAX_VALUE_LENGTH];
} BoardConfigEntry;

typedef struct BoardConfig {

    BoardConfigEntry *entries;
    size_t           count;
} BoardConfig;

static void freeBoardConfig(BoardConfig *boardConfig) {

    free(boardConfig->entries);

    boardConfig->entries = NULL;
    boardConfig->count   = 0;
}

static char *trimWhitespace(char *str) {

    char *end;

    while (isspace((unsigned char)*str)) str++;
    if (*str == 0) return str;

    end = str + strlen(str) - 1;
    while (end > str && isspace((unsigned char)*end)) end--;

    *(end + 1) = '\0';

    return str;
}

static BoardConfig parseBoardConfigFile(const char *fileName) {

    BoardConfig boardConfig;
    FILE *file;
    char line[MAX_LINE_LENGTH];
    char *trimmedLine;
    char *delimiter;
    char *key;
    char *value;

    boardConfig.entries = NULL;
    boardConfig.count = 0;

    file = fopen(fileName, "r");
    if (!file) {
        usys_log_error("Error opening file: %s", fileName);
        return boardConfig;
    }

    while (fgets(line, sizeof(line), file)) {
        trimmedLine = trimWhitespace(line);
        if (trimmedLine[0] == '\0' || trimmedLine[0] == '#') continue;

        delimiter = strchr(trimmedLine, '=');
        if (!delimiter) {
            usys_log_error("Invalid line in config file: %s\n",
                           trimmedLine, fileName);
            continue;
        }

        *delimiter = '\0';
        key   = trimWhitespace(trimmedLine);
        value = trimWhitespace(delimiter + 1);

        boardConfig.entries = realloc(boardConfig.entries,
                                      (boardConfig.count + 1) *
                                      sizeof(BoardConfigEntry));
        if (!boardConfig.entries) {
            usys_log_error("Memory allocation failed. Size: %d",
                           (boardConfig.count + 1) * sizeof(BoardConfigEntry));
            fclose(file);
            return boardConfig;
        }

        strncpy(boardConfig.entries[boardConfig.count].key, key, MAX_KEY_LENGTH);
        strncpy(boardConfig.entries[boardConfig.count].value, value, MAX_VALUE_LENGTH);
        boardConfig.count++;
    }

    fclose(file);
    return boardConfig;
}

static BoardConfig mergeBoardConfigs(const BoardConfig *commonConfig,
                                     const BoardConfig *boardConfig) {

    BoardConfig mergedConfig;
    size_t i, j;
    int found;

    mergedConfig.entries = malloc(commonConfig->count * sizeof(BoardConfigEntry));
    if (!mergedConfig.entries) {
        usys_log_error("Memory allocation failed. Size: %d",
                       commonConfig->count * sizeof(BoardConfigEntry));
        exit(EXIT_FAILURE);
    }
    memcpy(mergedConfig.entries,
           commonConfig->entries,
           commonConfig->count * sizeof(BoardConfigEntry));
    mergedConfig.count = commonConfig->count;

    for (i = 0; i < boardConfig->count; i++) {
        found = 0;
        for (j = 0; j < mergedConfig.count; j++) {
            if (strcmp(boardConfig->entries[i].key,
                       mergedConfig.entries[j].key) == 0) {
                strncpy(mergedConfig.entries[j].value,
                        boardConfig->entries[i].value, MAX_VALUE_LENGTH);
                found = 1;
                break;
            }
        }
        if (!found) {
            mergedConfig.entries = realloc(mergedConfig.entries,
                                           (mergedConfig.count + 1) *
                                           sizeof(BoardConfigEntry));
            if (!mergedConfig.entries) {
                usys_log_error("Memory allocation failed: %d",
                               (mergedConfig.count + 1)*sizeof(BoardConfigEntry));
                exit(EXIT_FAILURE);
            }
            strncpy(mergedConfig.entries[mergedConfig.count].key,
                    boardConfig->entries[i].key, MAX_KEY_LENGTH);
            strncpy(mergedConfig.entries[mergedConfig.count].value,
                    boardConfig->entries[i].value, MAX_VALUE_LENGTH);
            mergedConfig.count++;
        }
    }

    return mergedConfig;
}

char *getAppsFromBoardConfigs(const char *commonFile,
                              const char *boardFile) {

    BoardConfig commonConfig, boardConfig, mergedConfig;
    size_t i=0, bufferSize=MAX_BUFFER_SIZE;
    char *result = NULL;

    commonConfig = parseBoardConfigFile(commonFile);
    boardConfig  = parseBoardConfigFile(boardFile);
    mergedConfig = mergeBoardConfigs(&commonConfig, &boardConfig);

    result = malloc(bufferSize);
    if (!result) {
        usys_log_error("Memory allocation failed Size: %d", bufferSize);
        return NULL;
    }
    result[0] = '\0';

    for (i = 0; i < mergedConfig.count; i++) {
        if (strcmp(mergedConfig.entries[i].value, "yes") == 0) {
            if (strlen(result) + strlen(mergedConfig.entries[i].key) + 2 >
                bufferSize) {
                bufferSize *= 2;
                result = realloc(result, bufferSize);
                if (!result) {
                    usys_log_error("Memory allocation failed. Size: %d",
                                   bufferSize);
                    return NULL;
                }
            }

            if (strlen(result) > 0) {
                strcat(result, ",");
            }

            strcat(result, mergedConfig.entries[i].key);
        }
    }

    freeBoardConfig(&commonConfig);
    freeBoardConfig(&boardConfig);
    freeBoardConfig(&mergedConfig);

    return result;
}

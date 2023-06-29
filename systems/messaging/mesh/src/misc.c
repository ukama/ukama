/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string.h>

/*
 * split_strings --
 *
 */
void split_strings(char *input, char **str1, char **str2, char *delimiter) {

    char *token=NULL;

    token = strtok(input, delimiter);

    if (token != NULL && str1) {
        *str1 = strdup(token);

        token = strtok(NULL, delimiter);
        if (token != NULL && str2) {
            *str2 = strdup(token);
        }
    }
}

/*
 * get_substring_after_index --
 */
const int get_substring_after_index(char **ptr, char *str, int after, char ch) {

    int i=0, j=0;

    if (str == NULL) return 0;

    for (i=0; i<strlen(str); i++) {
        if (str[i] == ch) {
            j++;
            if (j == after) {
                *ptr = str + i + 1;
                break;
            }
        }
    }

    return 1;
}

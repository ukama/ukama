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
#include <unistd.h>

#include "usys_log.h"
#include "usys_services.h"

#define SERVICE_NAME "testApp"
#define MSG_LENGTH   64

void generate_random_message(char* message) {
    char charset[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    int charsetSize = sizeof(charset) - 1;
    
    srand(time(NULL));
    
    for (int i = 0; i < MSG_LENGTH; ++i) {
        int key = rand() % charsetSize;
        message[i] = charset[key];
    }
    
    message[MSG_LENGTH] = '\0';
}

int main(int argc, char **argv) {

    char logMessage[MSG_LENGTH+1] = {0};

    usys_log_set_service(SERVICE_NAME);
    usys_log_remote_init(SERVICE_NAME);

    while (1) {

        memset(&logMessage[0], 0, MSG_LENGTH+1);
        generate_random_message(logMessage);

        usys_log_error(&logMessage[0]);
        usys_log_debug(&logMessage[0]);
        usys_log_info(&logMessage[0]);
    }
    
	return 0;
}

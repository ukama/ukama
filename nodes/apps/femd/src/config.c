/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "config.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

int config_init(Config *config) {
    if (!config) {
        printf("[ERROR] Null config pointer\n");
        return STATUS_NOK;
    }
    
    // Initialize with default values
    memset(config, 0, sizeof(Config));
    
    config->serviceName = strdup(DEF_SERVICE_NAME);
    config->servicePort = DEF_SERVICE_PORT;
    config->logLevel = strdup(DEF_LOG_LEVEL);
    config->configFile = strdup(DEF_CONFIG_FILE);
    
    if (!config->serviceName || !config->logLevel || !config->configFile) {
        printf("[ERROR] Failed to allocate memory for config\n");
        config_free(config);
        return STATUS_NOK;
    }
    
    printf("[DEBUG] Configuration initialized with defaults\n");
    return STATUS_OK;
}

void config_free(Config *config) {
    if (!config) {
        return;
    }
    
    if (config->serviceName) {
        free(config->serviceName);
        config->serviceName = NULL;
    }
    
    if (config->logLevel) {
        free(config->logLevel);
        config->logLevel = NULL;
    }
    
    if (config->configFile) {
        free(config->configFile);
        config->configFile = NULL;
    }
    
    printf("[DEBUG] Configuration freed\n");
}

int config_load_from_file(Config *config, const char *filename) {
    if (!config || !filename) {
        printf("[ERROR] Null config or filename\n");
        return STATUS_NOK;
    }
    
    // Check if file exists
    if (access(filename, F_OK) != 0) {
        printf("[WARN] Config file does not exist: %s\n", filename);
        return STATUS_NOK;
    }
    
    // Simple file reading
    FILE *file = fopen(filename, "r");
    if (!file) {
        printf("[ERROR] Failed to open config file: %s\n", filename);
        return STATUS_NOK;
    }
    
    char buffer[1024];
    if (!fgets(buffer, sizeof(buffer), file)) {
        printf("[ERROR] Failed to read config file: %s\n", filename);
        fclose(file);
        return STATUS_NOK;
    }
    fclose(file);
    
    printf("[DEBUG] Read config file: %s\n", filename);
    
    // Simple parsing - look for key=value pairs
    char *line = strtok(buffer, "\n");
    while (line != NULL) {
        // Skip comments and empty lines
        if (line[0] == '#' || line[0] == '\0' || line[0] == '\n') {
            line = strtok(NULL, "\n");
            continue;
        }
        
        // Look for key=value
        char *equals = strchr(line, '=');
        if (equals != NULL) {
            *equals = '\0';
            char *key = line;
            char *value = equals + 1;
            
            // Trim whitespace (simple version)
            while (*key == ' ' || *key == '\t') key++;
            while (*value == ' ' || *value == '\t') value++;
            
            // Remove trailing whitespace from value
            char *end = value + strlen(value) - 1;
            while (end > value && (*end == ' ' || *end == '\t' || *end == '\n' || *end == '\r')) {
                *end = '\0';
                end--;
            }
            
            // Update config based on key
            if (strcmp(key, "service_name") == 0) {
                if (config->serviceName) free(config->serviceName);
                config->serviceName = strdup(value);
                printf("[DEBUG] Config: service_name = %s\n", value);
            } else if (strcmp(key, "service_port") == 0) {
                config->servicePort = atoi(value);
                printf("[DEBUG] Config: service_port = %d\n", config->servicePort);
            } else if (strcmp(key, "log_level") == 0) {
                if (config->logLevel) free(config->logLevel);
                config->logLevel = strdup(value);
                printf("[DEBUG] Config: log_level = %s\n", value);
            } else {
                printf("[DEBUG] Unknown config key: %s\n", key);
            }
        }
        
        line = strtok(NULL, "\n");
    }
    
    return STATUS_OK;
}

void config_print(const Config *config) {
    if (!config) {
        printf("[ERROR] Null config pointer\n");
        return;
    }
    
    printf("[INFO] Configuration:\n");
    printf("[INFO]   Service Name: %s\n", config->serviceName ? config->serviceName : "NULL");
    printf("[INFO]   Service Port: %d\n", config->servicePort);
    printf("[INFO]   Log Level: %s\n", config->logLevel ? config->logLevel : "NULL");
    printf("[INFO]   Config File: %s\n", config->configFile ? config->configFile : "NULL");
}
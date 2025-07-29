/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "config.h"
#include "femd.h"
#include "usys_mem.h"
#include "usys_string.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define DEF_SERVICE_NAME         "femd"
#define DEF_LOG_LEVEL            "INFO"
#define DEF_CONFIG_FILE          "./config/femd.conf"

int config_init(Config *config) {
    if (!config) {
        usys_log_error("Null config pointer");
        return STATUS_NOK;
    }
    
    memset(config, 0, sizeof(Config));
    
    config->serviceName = usys_strdup(DEF_SERVICE_NAME);
    config->servicePort = 0; /* Will be set by usys_find_service_port */
    config->logLevel = usys_strdup(DEF_LOG_LEVEL);
    config->configFile = usys_strdup(DEF_CONFIG_FILE);
    
    if (!config->serviceName || !config->logLevel || !config->configFile) {
        usys_log_error("Failed to allocate memory for config");
        config_free(config);
        return STATUS_NOK;
    }
    
    usys_log_debug("Configuration initialized with defaults");
    return STATUS_OK;
}

void config_free(Config *config) {
    if (!config) {
        return;
    }
    
    if (config->serviceName) {
        usys_free(config->serviceName);
        config->serviceName = NULL;
    }
    
    if (config->logLevel) {
        usys_free(config->logLevel);
        config->logLevel = NULL;
    }
    
    if (config->configFile) {
        usys_free(config->configFile);
        config->configFile = NULL;
    }
    
    usys_log_debug("Configuration freed");
}

int config_load_from_file(Config *config, const char *filename) {
    if (!config || !filename) {
        usys_log_error("Null config or filename");
        return STATUS_NOK;
    }
    
    if (access(filename, F_OK) != 0) {
        usys_log_warn("Config file does not exist: %s", filename);
        return STATUS_NOK;
    }
    
    FILE *file = fopen(filename, "r");
    if (!file) {
        usys_log_error("Failed to open config file: %s", filename);
        return STATUS_NOK;
    }
    
    char buffer[1024];
    if (!fgets(buffer, sizeof(buffer), file)) {
        usys_log_error("Failed to read config file: %s", filename);
        fclose(file);
        return STATUS_NOK;
    }
    fclose(file);
    
    usys_log_debug("Read config file: %s", filename);
    
    char *line = strtok(buffer, "\n");
    while (line != NULL) {
        if (line[0] == '#' || line[0] == '\0' || line[0] == '\n') {
            line = strtok(NULL, "\n");
            continue;
        }
        
        char *equals = strchr(line, '=');
        if (equals != NULL) {
            *equals = '\0';
            char *key = line;
            char *value = equals + 1;
            
            while (*key == ' ' || *key == '\t') key++;
            while (*value == ' ' || *value == '\t') value++;
            
            char *end = value + strlen(value) - 1;
            while (end > value && (*end == ' ' || *end == '\t' || *end == '\n' || *end == '\r')) {
                *end = '\0';
                end--;
            }
            
            if (strcmp(key, "service_name") == 0) {
                if (config->serviceName) usys_free(config->serviceName);
                config->serviceName = usys_strdup(value);
                usys_log_debug("Config: service_name = %s", value);
            } else if (strcmp(key, "service_port") == 0) {
                /* Note: service_port from config file is ignored, using usys_find_service_port instead */
                usys_log_debug("Config: service_port = %s (ignored, using services file)", value);
            } else if (strcmp(key, "log_level") == 0) {
                if (config->logLevel) usys_free(config->logLevel);
                config->logLevel = usys_strdup(value);
                usys_log_debug("Config: log_level = %s", value);
            } else {
                usys_log_debug("Unknown config key: %s", key);
            }
        }
        
        line = strtok(NULL, "\n");
    }
    
    return STATUS_OK;
}

void config_print(const Config *config) {
    if (!config) {
        usys_log_error("Null config pointer");
        return;
    }
    
    usys_log_info("Configuration:");
    usys_log_info("  Service Name: %s", config->serviceName ? config->serviceName : "NULL");
    usys_log_info("  Service Port: %d", config->servicePort);
    usys_log_info("  Log Level: %s", config->logLevel ? config->logLevel : "NULL");
    usys_log_info("  Config File: %s", config->configFile ? config->configFile : "NULL");
}
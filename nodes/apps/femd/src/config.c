#include "config.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_file.h"

int config_init(Config *config) {
    if (!config) {
        usys_log_error("Null config pointer");
        return STATUS_NOK;
    }
    
    // Initialize with default values
    memset(config, 0, sizeof(Config));
    
    config->serviceName = usys_strdup(DEF_SERVICE_NAME);
    config->servicePort = DEF_SERVICE_PORT;
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
    
    // Check if file exists
    if (access(filename, F_OK) != 0) {
        usys_log_warn("Config file does not exist: %s", filename);
        return STATUS_NOK;
    }
    
    // For now, we'll do a simple file read and basic parsing
    // In a real implementation, you might use a proper config parser
    char buffer[1024];
    int bytes_read = usys_file_read(filename, buffer, sizeof(buffer) - 1);
    
    if (bytes_read <= 0) {
        usys_log_error("Failed to read config file: %s", filename);
        return STATUS_NOK;
    }
    
    buffer[bytes_read] = '\0';
    usys_log_debug("Read %d bytes from config file", bytes_read);
    
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
                if (config->serviceName) usys_free(config->serviceName);
                config->serviceName = usys_strdup(value);
                usys_log_debug("Config: service_name = %s", value);
            } else if (strcmp(key, "service_port") == 0) {
                config->servicePort = atoi(value);
                usys_log_debug("Config: service_port = %d", config->servicePort);
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
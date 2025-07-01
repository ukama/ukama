#ifndef CONFIG_H
#define CONFIG_H

#include "usys_types.h"

typedef struct {
    char *serviceName;
    int  servicePort;
    char *logLevel;
    char *configFile;
} Config;

// Default configuration values
#define DEF_SERVICE_NAME         "femd"
#define DEF_SERVICE_PORT         8080
#define DEF_LOG_LEVEL            "INFO"
#define DEF_CONFIG_FILE          "./config/femd.conf"

// Function declarations
int config_init(Config *config);
void config_free(Config *config);
int config_load_from_file(Config *config, const char *filename);
void config_print(const Config *config);

#endif /* CONFIG_H */
/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "femd.h"

// Global variable to control main loop
volatile sig_atomic_t g_running = 1;

void handle_sigint(int signum) {
    printf("[INFO] Received signal %d, shutting down...\n", signum);
    g_running = 0;
}

void print_usage(const char *program) {
    printf("Usage: %s [OPTIONS]\n", program);
    printf("FEM Daemon - Front End Module control daemon\n\n");
    printf("Options:\n");
    printf("  -h, --help              Show this help message\n");
    printf("  -v, --version           Show version information\n");
    printf("  -l, --log-level LEVEL   Set log level (DEBUG, INFO, WARN, ERROR)\n");
    printf("  -p, --port PORT         Set service port (default: %d)\n", DEF_SERVICE_PORT);
    printf("  -c, --config FILE       Load configuration from file\n");
    printf("\n");
    printf("Examples:\n");
    printf("  %s                           # Start with default settings\n", program);
    printf("  %s -l DEBUG                  # Start with debug logging\n", program);
    printf("  %s -c /etc/femd/femd.conf    # Use specific config file\n", program);
}

void print_version(void) {
    printf("FEM Daemon version %s\n", FEM_VERSION);
    printf("Built on %s %s\n", __DATE__, __TIME__);
}

int main(int argc, char **argv) {
    int ret = STATUS_NOK;
    Config config = {0};
    
    // Initialize logging (basic printf for now)
    printf("[INFO] Starting FEM daemon v%s\n", FEM_VERSION);
    
    // Parse command line arguments
    static struct option long_options[] = {
        {"help",      no_argument,       0, 'h'},
        {"version",   no_argument,       0, 'v'},
        {"log-level", required_argument, 0, 'l'},
        {"port",      required_argument, 0, 'p'},
        {"config",    required_argument, 0, 'c'},
        {0, 0, 0, 0}
    };
    
    int opt;
    while ((opt = getopt_long(argc, argv, "hvl:p:c:", long_options, NULL)) != -1) {
        switch (opt) {
            case 'h':
                print_usage(argv[0]);
                exit(0);
                break;
            case 'v':
                print_version();
                exit(0);
                break;
            case 'l':
                printf("[INFO] Log level set to: %s\n", optarg);
                break;
            case 'p':
                // We'll store this for later use
                printf("[INFO] Port set to: %s\n", optarg);
                break;
            case 'c':
                printf("[INFO] Config file: %s\n", optarg);
                break;
            default:
                print_usage(argv[0]);
                exit(1);
        }
    }
    
    // Initialize configuration
    if (config_init(&config) != STATUS_OK) {
        printf("[ERROR] Failed to initialize configuration\n");
        goto cleanup;
    }
    
    // Load configuration from file if specified
    const char *config_file = DEF_CONFIG_FILE;
    if (config_load_from_file(&config, config_file) == STATUS_OK) {
        printf("[INFO] Configuration loaded from: %s\n", config_file);
        config_print(&config);
    } else {
        printf("[WARN] Could not load config file %s, using defaults\n", config_file);
    }
    
    // Set up signal handling
    signal(SIGINT, handle_sigint);
    signal(SIGTERM, handle_sigint);
    
    printf("[INFO] FEM daemon started successfully\n");
    printf("[INFO] Service: %s, Port: %d\n", config.serviceName, config.servicePort);
    
    // Main daemon loop
    while (g_running) {
        printf("[DEBUG] Daemon is running...\n");
        sleep(5); // Sleep for 5 seconds
    }
    
    ret = STATUS_OK;
    
cleanup:
    printf("[INFO] Shutting down FEM daemon...\n");
    config_free(&config);
    printf("[INFO] FEM daemon shutdown complete\n");
    
    return ret == STATUS_OK ? 0 : 1;
}
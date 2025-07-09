/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "femd.h"
#include "gpio_controller.h"
#include "safety_monitor.h"

volatile sig_atomic_t g_running = 1;

void handle_sigint(int signum) {
    usys_log_info("Received signal %d, shutting down...", signum);
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
    GpioController gpio_controller = {0};
    I2CController i2c_controller = {0};
    WebAPIServer web_server = {0};
    SafetyMonitor safety_monitor = {0};
    
    usys_log_set_service(SERVICE_NAME);
    usys_log_remote_init(SERVICE_NAME);
    
    usys_log_info("Starting FEM daemon v%s", FEM_VERSION);
    
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
                usys_log_info("Log level set to: %s", optarg);
                break;
            case 'p':
                usys_log_info("Port set to: %s", optarg);
                break;
            case 'c':
                usys_log_info("Config file: %s", optarg);
                break;
            default:
                print_usage(argv[0]);
                exit(1);
        }
    }
    
    if (config_init(&config) != STATUS_OK) {
        usys_log_error("Failed to initialize configuration");
        goto cleanup;
    }
    
    const char *config_file = DEF_CONFIG_FILE;
    if (config_load_from_file(&config, config_file) == STATUS_OK) {
        usys_log_info("Configuration loaded from: %s", config_file);
        config_print(&config);
    } else {
        usys_log_warn("Could not load config file %s, using defaults", config_file);
    }
    
    signal(SIGINT, handle_sigint);
    signal(SIGTERM, handle_sigint);
    
    if (gpio_controller_init(&gpio_controller, NULL) != STATUS_OK) {
        usys_log_error("Failed to initialize GPIO controller");
        goto cleanup;
    }
    
    if (i2c_controller_init(&i2c_controller) != STATUS_OK) {
        usys_log_error("Failed to initialize I2C controller");
        goto cleanup;
    }
    
    if (web_api_init(&web_server, config.servicePort, &gpio_controller, &i2c_controller) != STATUS_OK) {
        usys_log_error("Failed to initialize Web API server");
        goto cleanup;
    }
    
    if (web_api_start(&web_server) != STATUS_OK) {
        usys_log_error("Failed to start Web API server");
        goto cleanup;
    }
    
    if (safety_monitor_init(&safety_monitor, &gpio_controller, &i2c_controller) != STATUS_OK) {
        usys_log_error("Failed to initialize safety monitor");
        goto cleanup;
    }
    
    if (safety_monitor_start(&safety_monitor) != STATUS_OK) {
        usys_log_error("Failed to start safety monitor");
        goto cleanup;
    }
    
    usys_log_info("FEM daemon started successfully");
    usys_log_info("Service: %s, Port: %d", config.serviceName, config.servicePort);
    
    usys_log_info("Testing GPIO functionality...");
    
    GpioStatus status;
    if (gpio_get_all_status(&gpio_controller, FEM_UNIT_1, &status) == STATUS_OK) {
        usys_log_info("FEM1 Status - TX_RF: %s, RX_RF: %s, PA_VDS: %s",
               status.tx_rf_enable ? "ON" : "OFF",
               status.rx_rf_enable ? "ON" : "OFF", 
               status.pa_vds_enable ? "ON" : "OFF");
    }
    
    usys_log_info("Testing I2C device detection...");
    i2c_print_device_scan(FEM_UNIT_1);
    
    usys_log_info("Testing GPIO control - enabling TX_RF for FEM1");
    gpio_set_tx_rf(&gpio_controller, FEM_UNIT_1, true);
    
    usys_log_info("Testing GPIO control - disabling TX_RF for FEM1");
    gpio_set_tx_rf(&gpio_controller, FEM_UNIT_1, false);
    
    while (g_running) {
        usys_log_debug("Daemon is running...");
        sleep(5);
    }
    
    ret = STATUS_OK;
    
cleanup:
    usys_log_info("Shutting down FEM daemon...");
    safety_monitor_cleanup(&safety_monitor);
    web_api_cleanup(&web_server);
    i2c_controller_cleanup(&i2c_controller);
    gpio_controller_cleanup(&gpio_controller);
    config_free(&config);
    usys_log_info("FEM daemon shutdown complete");
    
    return (ret == STATUS_OK) ? 0 : 1;
}
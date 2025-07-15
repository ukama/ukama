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
#include <pthread.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <errno.h>

#include "web_api.h"
#include "femd.h"
#include "jserdes.h"
#include "http_status.h"

static pthread_t server_thread;
static int server_socket = -1;

static void* web_api_server_thread(void *arg);
static void web_api_handle_client(WebAPIServer *server, int client_socket);
static int web_api_parse_request(const char *raw_request, HTTPRequest *request);
static void web_api_send_response(int client_socket, const HTTPResponse *response);
static int web_api_route_request(WebAPIServer *server, const HTTPRequest *request, HTTPResponse *response);

int web_api_init(WebAPIServer *server, int port, GpioController *gpio_ctrl, I2CController *i2c_ctrl) {
    if (!server || !gpio_ctrl || !i2c_ctrl) {
        usys_log_error("Invalid parameters for web API initialization");
        return STATUS_NOK;
    }

    memset(server, 0, sizeof(WebAPIServer));
    server->port = port;
    server->running = false;
    server->gpio_controller = gpio_ctrl;
    server->i2c_controller = i2c_ctrl;

    usys_log_info("Web API initialized on port %d", port);
    return STATUS_OK;
}

int web_api_start(WebAPIServer *server) {
    if (!server) {
        return STATUS_NOK;
    }

    if (server->running) {
        usys_log_warn("Web API server is already running");
        return STATUS_OK;
    }

    if (pthread_create(&server_thread, NULL, web_api_server_thread, server) != 0) {
        usys_log_error("Failed to create web API server thread");
        return STATUS_NOK;
    }

    server->running = true;
    usys_log_info("Web API server started on port %d", server->port);
    return STATUS_OK;
}

void web_api_stop(WebAPIServer *server) {
    if (!server || !server->running) {
        return;
    }

    server->running = false;

    if (server_socket >= 0) {
        close(server_socket);
        server_socket = -1;
    }

    pthread_join(server_thread, NULL);

    usys_log_info("Web API server stopped");
}

void web_api_cleanup(WebAPIServer *server) {
    if (server) {
        web_api_stop(server);
        memset(server, 0, sizeof(WebAPIServer));
    }
}

static void* web_api_server_thread(void *arg) {
    WebAPIServer *server = (WebAPIServer*)arg;
    struct sockaddr_in server_addr, client_addr;
    socklen_t client_len = sizeof(client_addr);
    int client_socket;

    server_socket = socket(AF_INET, SOCK_STREAM, 0);
    if (server_socket < 0) {
        usys_log_error("Failed to create socket: %s", strerror(errno));
        return NULL;
    }

    int opt = 1;
    setsockopt(server_socket, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(server->port);

    if (bind(server_socket, (struct sockaddr*)&server_addr, sizeof(server_addr)) < 0) {
        usys_log_error("Failed to bind socket: %s", strerror(errno));
        close(server_socket);
        return NULL;
    }

    if (listen(server_socket, 5) < 0) {
        usys_log_error("Failed to listen on socket: %s", strerror(errno));
        close(server_socket);
        return NULL;
    }

    usys_log_info("Web API server listening on port %d", server->port);

    while (server->running) {
        client_socket = accept(server_socket, (struct sockaddr*)&client_addr, &client_len);
        
        if (client_socket < 0) {
            if (server->running) {
                usys_log_error("Failed to accept connection: %s", strerror(errno));
            }
            break;
        }

        usys_log_debug("Client connected from %s:%d", 
                       inet_ntoa(client_addr.sin_addr), ntohs(client_addr.sin_port));

        web_api_handle_client(server, client_socket);
        close(client_socket);
    }

    if (server_socket >= 0) {
        close(server_socket);
    }

    return NULL;
}

static void web_api_handle_client(WebAPIServer *server, int client_socket) {
    char buffer[2048];
    HTTPRequest request;
    HTTPResponse response;
    int bytes_received;

    bytes_received = recv(client_socket, buffer, sizeof(buffer) - 1, 0);
    if (bytes_received <= 0) {
        return;
    }

    buffer[bytes_received] = '\0';
    usys_log_debug("Received HTTP request:\n%s", buffer);

    if (web_api_parse_request(buffer, &request) != STATUS_OK) {
        web_api_set_error_response(&response, 400, "Bad Request");
        web_api_send_response(client_socket, &response);
        return;
    }

    if (web_api_route_request(server, &request, &response) != STATUS_OK) {
        web_api_set_error_response(&response, 500, "Internal Server Error");
    }

    web_api_send_response(client_socket, &response);
}

static int web_api_parse_request(const char *raw_request, HTTPRequest *request) {
    if (!raw_request || !request) {
        return STATUS_NOK;
    }

    memset(request, 0, sizeof(HTTPRequest));

    // Create a copy of the request to avoid modifying the original
    size_t request_len = strlen(raw_request);
    char *request_copy = malloc(request_len + 1);
    if (!request_copy) {
        usys_log_error("Failed to allocate memory for request parsing");
        return STATUS_NOK;
    }
    strcpy(request_copy, raw_request);

    char *line = strtok(request_copy, "\r\n");
    if (!line) {
        free(request_copy);
        return STATUS_NOK;
    }

    if (sscanf(line, "%15s %255s", request->method, request->path) != 2) {
        free(request_copy);
        return STATUS_NOK;
    }

    request->content_length = 0;
    char *content_length_line = strstr(raw_request, "Content-Length:");
    if (content_length_line) {
        sscanf(content_length_line, "Content-Length: %d", &request->content_length);
    }

    char *body_start = strstr(raw_request, "\r\n\r\n");
    if (body_start && request->content_length > 0) {
        body_start += 4; // Skip \r\n\r\n
        int body_len = strlen(body_start);
        if (body_len > 0 && body_len < sizeof(request->body)) {
            strncpy(request->body, body_start, sizeof(request->body) - 1);
        }
    }

    free(request_copy);
    return STATUS_OK;
}

static void web_api_send_response(int client_socket, const HTTPResponse *response) {
    char http_response[4096];
    const char *status_text;

    switch (response->status_code) {
        case 200: status_text = "OK"; break;
        case 400: status_text = "Bad Request"; break;
        case 404: status_text = "Not Found"; break;
        case 500: status_text = "Internal Server Error"; break;
        default: status_text = "Unknown"; break;
    }

    snprintf(http_response, sizeof(http_response),
             "HTTP/1.1 %d %s\r\n"
             "Content-Type: %s\r\n"
             "Content-Length: %d\r\n"
             "Access-Control-Allow-Origin: *\r\n"
             "Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS\r\n"
             "Access-Control-Allow-Headers: Content-Type\r\n"
             "\r\n"
             "%s",
             response->status_code, status_text,
             response->content_type,
             response->body_length,
             response->body);

    send(client_socket, http_response, strlen(http_response), 0);
    usys_log_debug("Sent HTTP response: %d %s", response->status_code, status_text);
}

static int web_api_route_request(WebAPIServer *server, const HTTPRequest *request, HTTPResponse *response) {
    usys_log_debug("Routing %s %s", request->method, request->path);

    if (strcmp(request->method, "OPTIONS") == 0) {
        web_api_set_response(response, HttpStatus_OK, "text/plain", "");
        return STATUS_OK;
    }

    if (strncmp(request->path, "/v1/fem/", 8) == 0) {
        int fem_unit = parse_fem_unit(request->path);
        if (fem_unit < 1 || fem_unit > 2) {
            web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid FEM unit");
            return STATUS_OK;
        }

        if (strcmp(request->method, "GET") == 0 && 
            (strcmp(request->path + 8, "1/gpio") == 0 || strcmp(request->path + 8, "2/gpio") == 0)) {
            return api_gpio_get_status(server, fem_unit, response);
        }

        if (strcmp(request->method, "POST") == 0) {
            if (strstr(request->path, "/gpio/tx_rf")) {
                JsonObj *json = parse_request_json(request->body);
                if (json) {
                    bool enable = json_is_true(json_object_get(json, "enable"));
                    int result = api_gpio_set_control(server, fem_unit, "tx_rf", enable, response);
                    json_free(&json);
                    return result;
                } else {
                    web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid JSON");
                    return STATUS_NOK;
                }
            }
            else if (strstr(request->path, "/gpio/rx_rf")) {
                JsonObj *json = parse_request_json(request->body);
                if (json) {
                    bool enable = json_is_true(json_object_get(json, "enable"));
                    int result = api_gpio_set_control(server, fem_unit, "rx_rf", enable, response);
                    json_free(&json);
                    return result;
                } else {
                    web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid JSON");
                    return STATUS_NOK;
                }
            }
            else if (strstr(request->path, "/gpio/pa_vds")) {
                JsonObj *json = parse_request_json(request->body);
                if (json) {
                    bool enable = json_is_true(json_object_get(json, "enable"));
                    int result = api_gpio_set_control(server, fem_unit, "pa_vds", enable, response);
                    json_free(&json);
                    return result;
                } else {
                    web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid JSON");
                    return STATUS_NOK;
                }
            }
            else if (strstr(request->path, "/gpio/tx_rfpal")) {
                JsonObj *json = parse_request_json(request->body);
                if (json) {
                    bool enable = json_is_true(json_object_get(json, "enable"));
                    int result = api_gpio_set_control(server, fem_unit, "tx_rfpal", enable, response);
                    json_free(&json);
                    return result;
                } else {
                    web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid JSON");
                    return STATUS_NOK;
                }
            }
            else if (strstr(request->path, "/gpio/28v_vds")) {
                JsonObj *json = parse_request_json(request->body);
                if (json) {
                    bool enable = json_is_true(json_object_get(json, "enable"));
                    int result = api_gpio_set_control(server, fem_unit, "28v_vds", enable, response);
                    json_free(&json);
                    return result;
                } else {
                    web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid JSON");
                    return STATUS_NOK;
                }
            }
        }

        if (strcmp(request->method, "POST") == 0 && strstr(request->path, "/i2c/dac")) {
            JsonObj *json = parse_request_json(request->body);
            if (json) {
                float carrier_voltage = 0.0f, peak_voltage = 0.0f;
                if (json_deserialize_dac_control(json, &carrier_voltage, &peak_voltage) == JSON_OK) {
                    int result = api_dac_set_voltages(server, fem_unit, carrier_voltage, peak_voltage, response);
                    json_free(&json);
                    return result;
                } else {
                    json_free(&json);
                    web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid DAC control parameters");
                    return STATUS_NOK;
                }
            } else {
                web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid JSON");
                return STATUS_NOK;
            }
        }

        if (strcmp(request->method, "GET") == 0 && strstr(request->path, "/i2c/dac")) {
            return api_dac_get_config(server, fem_unit, response);
        }

        if (strcmp(request->method, "GET") == 0 && strstr(request->path, "/i2c/temperature")) {
            return api_temp_read(server, fem_unit, response);
        }

        if (strcmp(request->method, "GET") == 0 && strstr(request->path, "/i2c/adc")) {
            return api_adc_read_all(server, fem_unit, response);
        }

        if (strcmp(request->method, "GET") == 0 && strstr(request->path, "/i2c/eeprom")) {
            return api_eeprom_read_serial(server, fem_unit, response);
        }
        if (strcmp(request->method, "POST") == 0 && strstr(request->path, "/i2c/eeprom")) {
            JsonObj *json = parse_request_json(request->body);
            if (json) {
                char *serial = NULL;
                if (json_deserialize_eeprom_write(json, &serial) == JSON_OK && serial) {
                    int result = api_eeprom_write_serial(server, fem_unit, serial, response);
                    free(serial);
                    json_free(&json);
                    return result;
                } else {
                    json_free(&json);
                    web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid serial number");
                    return STATUS_NOK;
                }
            } else {
                web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid JSON");
                return STATUS_NOK;
            }
        }
    }

    if (strcmp(request->method, "GET") == 0 && strcmp(request->path, "/health") == 0) {
        JsonObj *health_json = NULL;
        ServiceInfo info = {
            .service_name = "femd",
            .version = FEM_VERSION,
            .uptime = 0  // TODO: Calculate actual uptime
        };
        if (json_serialize_service_info(&health_json, &info) == JSON_OK) {
            web_api_set_json_response(response, HttpStatus_OK, health_json);
            json_free(&health_json);
        } else {
            web_api_set_error_response(response, HttpStatus_InternalServerError, "Failed to serialize health info");
        }
        return STATUS_OK;
    }

    web_api_set_error_response(response, HttpStatus_NotFound, "Endpoint not found");
    return STATUS_OK;
}

void web_api_set_response(HTTPResponse *response, HttpStatusCode status, const char *content_type, const char *body) {
    response->status_code = status;
    strncpy(response->content_type, content_type, sizeof(response->content_type) - 1);
    if (body) {
        strncpy(response->body, body, sizeof(response->body) - 1);
        response->body_length = strlen(response->body);
    } else {
        response->body[0] = '\0';
        response->body_length = 0;
    }
}

void web_api_set_json_response(HTTPResponse *response, HttpStatusCode status, JsonObj *json) {
    if (json) {
        char *json_str = json_dumps(json, JSON_COMPACT);
        if (json_str) {
            web_api_set_response(response, status, "application/json", json_str);
            free(json_str);
        } else {
            web_api_set_response(response, HttpStatus_InternalServerError, "application/json", "{\"error\":\"JSON serialization failed\"}");
        }
    } else {
        web_api_set_response(response, HttpStatus_InternalServerError, "application/json", "{\"error\":\"No JSON data\"}");
    }
}

void web_api_set_error_response(HTTPResponse *response, HttpStatusCode status, const char *error_message) {
    JsonObj *json_error = NULL;
    if (json_serialize_error(&json_error, status, error_message) == JSON_OK) {
        web_api_set_json_response(response, status, json_error);
        json_free(&json_error);
    } else {
        web_api_set_response(response, status, "application/json", "{\"error\":\"Error serialization failed\"}");
    }
}

int parse_fem_unit(const char *path) {
    if (strstr(path, "/fem/1/")) return 1;
    if (strstr(path, "/fem/2/")) return 2;
    return -1;
}

JsonObj *parse_request_json(const char *json_str) {
    if (!json_str || strlen(json_str) == 0) {
        return NULL;
    }
    
    json_error_t error;
    JsonObj *json = json_loads(json_str, 0, &error);
    if (!json) {
        usys_log_error("JSON parse error: %s", error.text);
        return NULL;
    }
    
    return json;
}

bool parse_json_bool(const char *json, const char *key) {
    if (!json || !key) return false;
    
    char search_pattern[64];
    snprintf(search_pattern, sizeof(search_pattern), "\"%s\":", key);
    
    char *pos = strstr(json, search_pattern);
    if (pos) {
        pos += strlen(search_pattern);
        while (*pos == ' ' || *pos == '\t') pos++;
        
        if (strncmp(pos, "true", 4) == 0) return true;
        if (*pos == '1') return true;
    }
    return false;
}

float parse_json_float(const char *json, const char *key) {
    if (!json || !key) return 0.0f;
    
    char search_pattern[64];
    snprintf(search_pattern, sizeof(search_pattern), "\"%s\":", key);
    
    char *pos = strstr(json, search_pattern);
    if (pos) {
        pos += strlen(search_pattern);
        while (*pos == ' ' || *pos == '\t') pos++;
        
        return (float)atof(pos);
    }
    return 0.0f;
}

int parse_json_string(const char *json, const char *key, char *value, size_t max_len) {
    if (!json || !key || !value) return 0;
    
    char search_pattern[64];
    snprintf(search_pattern, sizeof(search_pattern), "\"%s\":", key);
    
    char *pos = strstr(json, search_pattern);
    if (pos) {
        pos += strlen(search_pattern);
        while (*pos == ' ' || *pos == '\t') pos++;
        
        if (*pos == '"') {
            pos++; // Skip opening quote
            char *end = strchr(pos, '"');
            if (end) {
                size_t len = end - pos;
                if (len < max_len) {
                    strncpy(value, pos, len);
                    value[len] = '\0';
                    return (int)len;
                }
            }
        }
    }
    return 0;
}

void create_json_gpio_status(const GpioStatus *status, char *json_buffer, size_t buffer_size) {
    snprintf(json_buffer, buffer_size,
             "{"
             "\"tx_rf_enable\":%s,"
             "\"rx_rf_enable\":%s,"
             "\"pa_vds_enable\":%s,"
             "\"rf_pal_enable\":%s,"
             "\"28v_vds_enable\":%s,"
             "\"psu_pgood\":%s"
             "}",
             status->tx_rf_enable ? "true" : "false",
             status->rx_rf_enable ? "true" : "false", 
             status->pa_vds_enable ? "true" : "false",
             status->rf_pal_enable ? "true" : "false",
             status->pa_disable ? "false" : "true", // Inverted logic
             status->pg_reg_5v ? "true" : "false");
}

void create_json_error(const char *error_message, char *json_buffer, size_t buffer_size) {
    snprintf(json_buffer, buffer_size, "{\"error\":\"%s\"}", error_message);
}

void create_json_success(const char *message, char *json_buffer, size_t buffer_size) {
    snprintf(json_buffer, buffer_size, "{\"status\":\"success\",\"message\":\"%s\"}", message);
}

int api_gpio_get_status(WebAPIServer *server, int fem_unit, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    GpioStatus status;
    
    if (gpio_get_all_status(server->gpio_controller, unit, &status) != STATUS_OK) {
        web_api_set_error_response(response, HttpStatus_InternalServerError, "Failed to read GPIO status");
        return STATUS_NOK;
    }
    
    JsonObj *json_response = NULL;
    if (json_serialize_gpio_status(&json_response, &status, fem_unit) == JSON_OK) {
        web_api_set_json_response(response, HttpStatus_OK, json_response);
        json_free(&json_response);
    } else {
        web_api_set_error_response(response, HttpStatus_InternalServerError, "Failed to serialize GPIO status");
        return STATUS_NOK;
    }
    
    return STATUS_OK;
}

int api_gpio_set_control(WebAPIServer *server, int fem_unit, const char *gpio_name, bool enable, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    int result = STATUS_NOK;
    
    if (strcmp(gpio_name, "tx_rf") == 0) {
        result = gpio_set_tx_rf(server->gpio_controller, unit, enable);
    }
    else if (strcmp(gpio_name, "rx_rf") == 0) {
        result = gpio_set_rx_rf(server->gpio_controller, unit, enable);
    }
    else if (strcmp(gpio_name, "pa_vds") == 0) {
        result = gpio_set_pa_vds(server->gpio_controller, unit, enable);
    }
    else if (strcmp(gpio_name, "tx_rfpal") == 0) {
        result = gpio_set_tx_rfpal(server->gpio_controller, unit, enable);
    }
    else if (strcmp(gpio_name, "28v_vds") == 0) {
        result = gpio_set_28v_vds(server->gpio_controller, unit, enable);
    }
    else {
        web_api_set_error_response(response, HttpStatus_BadRequest, "Invalid GPIO name");
        return STATUS_NOK;
    }
    
    if (result == STATUS_OK) {
        JsonObj *json_response = NULL;
        if (json_serialize_success(&json_response, "GPIO control successful") == JSON_OK) {
            web_api_set_json_response(response, HttpStatus_OK, json_response);
            json_free(&json_response);
        } else {
            web_api_set_error_response(response, HttpStatus_InternalServerError, "Failed to serialize success response");
        }
    } else {
        web_api_set_error_response(response, HttpStatus_InternalServerError, "Failed to set GPIO");
    }
    
    return result;
}

int api_dac_set_voltages(WebAPIServer *server, int fem_unit, float carrier_voltage, float peak_voltage, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    
    if (dac_init(server->i2c_controller, unit) != STATUS_OK) {
        web_api_set_error_response(response, 500, "Failed to initialize DAC");
        return STATUS_NOK;
    }
    
    int result1 = dac_set_carrier_voltage(server->i2c_controller, unit, carrier_voltage);
    int result2 = dac_set_peak_voltage(server->i2c_controller, unit, peak_voltage);
    
    if (result1 == STATUS_OK && result2 == STATUS_OK) {
        char json_response[256];
        snprintf(json_response, sizeof(json_response),
                 "{\"status\":\"success\",\"carrier_voltage\":%.2f,\"peak_voltage\":%.2f,\"fem_unit\":%d}",
                 carrier_voltage, peak_voltage, fem_unit);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 500, "Failed to set DAC voltages");
        return STATUS_NOK;
    }
}

int api_dac_get_config(WebAPIServer *server, int fem_unit, HTTPResponse *response) {
    float carrier_voltage, peak_voltage;
    
    if (dac_get_config(server->i2c_controller, &carrier_voltage, &peak_voltage) == STATUS_OK) {
        char json_response[256];
        snprintf(json_response, sizeof(json_response),
                 "{\"carrier_voltage\":%.2f,\"peak_voltage\":%.2f,\"fem_unit\":%d}",
                 carrier_voltage, peak_voltage, fem_unit);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 500, "Failed to read DAC configuration");
        return STATUS_NOK;
    }
}

int api_temp_read(WebAPIServer *server, int fem_unit, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    float temperature;
    
    if (temp_sensor_init(server->i2c_controller, unit) != STATUS_OK) {
        web_api_set_error_response(response, 500, "Failed to initialize temperature sensor");
        return STATUS_NOK;
    }
    
    if (temp_sensor_read(server->i2c_controller, unit, &temperature) == STATUS_OK) {
        char json_response[256];
        snprintf(json_response, sizeof(json_response),
                 "{\"temperature\":%.1f,\"unit\":\"celsius\",\"fem_unit\":%d}",
                 temperature, fem_unit);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 500, "Failed to read temperature");
        return STATUS_NOK;
    }
}

int api_temp_set_threshold(WebAPIServer *server, int fem_unit, float threshold, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    
    if (temp_sensor_set_threshold(server->i2c_controller, unit, threshold) == STATUS_OK) {
        char json_response[256];
        snprintf(json_response, sizeof(json_response),
                 "{\"status\":\"success\",\"threshold\":%.1f,\"fem_unit\":%d}",
                 threshold, fem_unit);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 500, "Failed to set temperature threshold");
        return STATUS_NOK;
    }
}

int api_adc_read_channel(WebAPIServer *server, int fem_unit, int channel, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    float voltage;
    
    if (adc_init(server->i2c_controller, unit) != STATUS_OK) {
        web_api_set_error_response(response, 500, "Failed to initialize ADC");
        return STATUS_NOK;
    }
    
    if (adc_read_channel(server->i2c_controller, unit, channel, &voltage) == STATUS_OK) {
        char json_response[256];
        snprintf(json_response, sizeof(json_response),
                 "{\"channel\":%d,\"voltage\":%.3f,\"fem_unit\":%d}",
                 channel, voltage, fem_unit);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 500, "Failed to read ADC channel");
        return STATUS_NOK;
    }
}

int api_adc_read_all(WebAPIServer *server, int fem_unit, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    float reverse_power, pa_current;
    
    if (adc_init(server->i2c_controller, unit) != STATUS_OK) {
        web_api_set_error_response(response, 500, "Failed to initialize ADC");
        return STATUS_NOK;
    }
    
    int result1 = adc_read_reverse_power(server->i2c_controller, unit, &reverse_power);
    int result2 = adc_read_pa_current(server->i2c_controller, unit, &pa_current);
    
    if (result1 == STATUS_OK && result2 == STATUS_OK) {
        char json_response[512];
        snprintf(json_response, sizeof(json_response),
                 "{"
                 "\"reverse_power\":%.2f,"
                 "\"pa_current\":%.3f,"
                 "\"fem_unit\":%d,"
                 "\"units\":{\"reverse_power\":\"dBm\",\"pa_current\":\"A\"}"
                 "}",
                 reverse_power, pa_current, fem_unit);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 500, "Failed to read ADC values");
        return STATUS_NOK;
    }
}

int api_adc_set_safety(WebAPIServer *server, float max_reverse_power, float max_current, HTTPResponse *response) {
    if (adc_set_safety_thresholds(server->i2c_controller, max_reverse_power, max_current) == STATUS_OK) {
        char json_response[256];
        snprintf(json_response, sizeof(json_response),
                 "{\"status\":\"success\",\"max_reverse_power\":%.1f,\"max_current\":%.1f}",
                 max_reverse_power, max_current);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 500, "Failed to set safety thresholds");
        return STATUS_NOK;
    }
}

int api_eeprom_write_serial(WebAPIServer *server, int fem_unit, const char *serial, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    
    if (eeprom_write_serial(server->i2c_controller, unit, serial) == STATUS_OK) {
        char json_response[256];
        snprintf(json_response, sizeof(json_response),
                 "{\"status\":\"success\",\"serial\":\"%s\",\"fem_unit\":%d}",
                 serial, fem_unit);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 500, "Failed to write serial number");
        return STATUS_NOK;
    }
}

int api_eeprom_read_serial(WebAPIServer *server, int fem_unit, HTTPResponse *response) {
    FemUnit unit = (fem_unit == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    char serial[32];
    
    if (eeprom_read_serial(server->i2c_controller, unit, serial, sizeof(serial)) == STATUS_OK) {
        char json_response[256];
        snprintf(json_response, sizeof(json_response),
                 "{\"serial\":\"%s\",\"fem_unit\":%d}",
                 serial, fem_unit);
        web_api_set_response(response, 200, "application/json", json_response);
        return STATUS_OK;
    } else {
        web_api_set_error_response(response, 404, "No serial number found");
        return STATUS_NOK;
    }
}
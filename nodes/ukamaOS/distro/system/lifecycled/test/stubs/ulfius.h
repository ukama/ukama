#pragma once

#include <stddef.h>

#include "jansson.h"

#define U_OK 0
#define U_CALLBACK_CONTINUE 0

struct _u_request {
    const unsigned char *binary_body;
    size_t binary_body_length;
};

struct _u_response {
    int unused;
};

struct _u_instance {
    int unused;
};

typedef struct _u_request URequest;
typedef struct _u_response UResponse;
typedef struct _u_instance UInst;

int ulfius_init_instance(struct _u_instance *instance,
                         unsigned int port,
                         const char *bindAddress,
                         const char *prefix);
int ulfius_add_endpoint_by_val(
    struct _u_instance *instance,
    const char *method,
    const char *prefix,
    const char *format,
    unsigned int priority,
    int (*callback)(const struct _u_request *,
                    struct _u_response *,
                    void *),
    void *userData);
int ulfius_start_framework(struct _u_instance *instance);
void ulfius_stop_framework(struct _u_instance *instance);
void ulfius_clean_instance(struct _u_instance *instance);
int ulfius_set_string_body_response(struct _u_response *response,
                                    unsigned int status,
                                    const char *body);
int ulfius_set_json_body_response(struct _u_response *response,
                                  unsigned int status,
                                  const json_t *body);


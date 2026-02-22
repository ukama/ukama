#include <ulfius.h>
#include <jansson.h>
#include <pthread.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

#include "usys_log.h"

/* Simple in-process mock services for noded/notifyd/femd.
 * These are real HTTP servers (Ulfius) bound to the same ports device.d expects.
 */

typedef struct {
    struct _u_instance inst;
    int port;
    pthread_t thread;
    int started;

    /* Counters */
    pthread_mutex_t lock;
    unsigned notify_posts;
    unsigned fem_puts;
    unsigned noded_gets;

    /* Last request info */
    char last_path[256];
    char last_body[2048];
} MockService;

static int cb_store_and_accepted(const struct _u_request *req, struct _u_response *resp, void *user_data) {
    MockService *ms = (MockService *)user_data;
    const char *path = req->http_url ? req->http_url : "";

    pthread_mutex_lock(&ms->lock);
    if (req->http_verb && strcmp(req->http_verb, "POST") == 0) {
        ms->notify_posts++;
    }
    snprintf(ms->last_path, sizeof(ms->last_path), "%s", path);
    ms->last_body[0] = '\0';
    if (req->binary_body && req->binary_body_length > 0) {
        size_t n = req->binary_body_length;
        if (n >= sizeof(ms->last_body)) n = sizeof(ms->last_body) - 1;
        memcpy(ms->last_body, req->binary_body, n);
        ms->last_body[n] = '\0';
    }
    pthread_mutex_unlock(&ms->lock);

    ulfius_set_empty_body_response(resp, 202);
    return U_CALLBACK_CONTINUE;
}

static int cb_femd_gpio(const struct _u_request *req, struct _u_response *resp, void *user_data) {
    MockService *ms = (MockService *)user_data;
    const char *path = req->http_url ? req->http_url : "";

    pthread_mutex_lock(&ms->lock);
    ms->fem_puts++;
    snprintf(ms->last_path, sizeof(ms->last_path), "%s", path);
    ms->last_body[0] = '\0';
    if (req->binary_body && req->binary_body_length > 0) {
        size_t n = req->binary_body_length;
        if (n >= sizeof(ms->last_body)) n = sizeof(ms->last_body) - 1;
        memcpy(ms->last_body, req->binary_body, n);
        ms->last_body[n] = '\0';
    }
    pthread_mutex_unlock(&ms->lock);

    ulfius_set_empty_body_response(resp, 202);
    return U_CALLBACK_CONTINUE;
}

static int cb_noded_nodeinfo(const struct _u_request *req, struct _u_response *resp, void *user_data) {
    MockService *ms = (MockService *)user_data;
    json_t *root = NULL;
    json_t *ni = NULL;

    (void)req;

    pthread_mutex_lock(&ms->lock);
    ms->noded_gets++;
    pthread_mutex_unlock(&ms->lock);

    /* Minimal payload expected by json_deserialize_node_info(..., JTAG_NODE_ID, ...) */
    /* Shape in your codebase commonly: {"nodeInfo":{"nodeId":"..."}} or {"nodeInfo":{"nodeID":"..."}}.
       We'll include multiple common keys to be robust. */
    root = json_object();
    ni = json_object();
    json_object_set_new(ni, "nodeId", json_string("ukama-tnode-0000-1111-2222-3333"));
    json_object_set_new(ni, "nodeID", json_string("ukama-tnode-0000-1111-2222-3333"));
    json_object_set_new(ni, "id",     json_string("ukama-tnode-0000-1111-2222-3333"));
    json_object_set_new(root, "nodeInfo", ni);

    ulfius_set_json_body_response(resp, 200, root);
    json_decref(root);
    return U_CALLBACK_CONTINUE;
}

static void *run_framework(void *arg) {
    MockService *ms = (MockService *)arg;
    if (ulfius_start_framework(&ms->inst) != U_OK) {
        return NULL;
    }
    ms->started = 1;
    /* Ulfius runs its own threads; just sleep forever until stopped */
    while (ms->started) {
        usleep(100 * 1000);
    }
    return NULL;
}

static int start_instance(MockService *ms, int port) {
    memset(ms, 0, sizeof(*ms));
    ms->port = port;
    pthread_mutex_init(&ms->lock, NULL);

    if (ulfius_init_instance(&ms->inst, port, NULL, NULL) != U_OK) {
        return -1;
    }

    /* CORS header like other services */
    u_map_put(ms->inst.default_headers, "Access-Control-Allow-Origin", "*");

    return 0;
}

static int start_mock(MockService *ms) {
    if (pthread_create(&ms->thread, NULL, run_framework, ms) != 0) {
        return -1;
    }

    /* Wait a bit for it to bind */
    for (int i = 0; i < 50; i++) {
        if (ms->started) return 0;
        usleep(20 * 1000);
    }

    return -1;
}

static void stop_mock(MockService *ms) {
    if (!ms) return;
    ms->started = 0;
    ulfius_stop_framework(&ms->inst);
    ulfius_clean_instance(&ms->inst);
    pthread_join(ms->thread, NULL);
    pthread_mutex_destroy(&ms->lock);
}

/* Public API used by tests */
int mock_noded_start(MockService *ms, int port) {
    if (start_instance(ms, port) != 0) return -1;
    ulfius_add_endpoint_by_val(&ms->inst, "GET", NULL, "/v1/nodeinfo", 0, &cb_noded_nodeinfo, ms);
    return start_mock(ms);
}

int mock_notifyd_start(MockService *ms, int port) {
    if (start_instance(ms, port) != 0) return -1;
    /* notify endpoint: /notify/v1/event/<serviceName>  */
    ulfius_add_endpoint_by_val(&ms->inst, "POST", NULL, "/notify/v1/event/*", 0, &cb_store_and_accepted, ms);
    return start_mock(ms);
}

int mock_femd_start(MockService *ms, int port) {
    if (start_instance(ms, port) != 0) return -1;
    /* femd gpio endpoints used by wc_put_gpio_to_femd */
    ulfius_add_endpoint_by_val(&ms->inst, "PUT", NULL, "/v1/fems/*/gpio", 0, &cb_femd_gpio, ms);
    return start_mock(ms);
}

void mock_service_stop(MockService *ms) {
    stop_mock(ms);
}

void mock_service_reset(MockService *ms) {
    pthread_mutex_lock(&ms->lock);
    ms->notify_posts = 0;
    ms->fem_puts = 0;
    ms->noded_gets = 0;
    ms->last_path[0] = '\0';
    ms->last_body[0] = '\0';
    pthread_mutex_unlock(&ms->lock);
}

unsigned mock_notify_posts(MockService *ms) {
    unsigned v;
    pthread_mutex_lock(&ms->lock);
    v = ms->notify_posts;
    pthread_mutex_unlock(&ms->lock);
    return v;
}

unsigned mock_fem_puts(MockService *ms) {
    unsigned v;
    pthread_mutex_lock(&ms->lock);
    v = ms->fem_puts;
    pthread_mutex_unlock(&ms->lock);
    return v;
}

unsigned mock_noded_gets(MockService *ms) {
    unsigned v;
    pthread_mutex_lock(&ms->lock);
    v = ms->noded_gets;
    pthread_mutex_unlock(&ms->lock);
    return v;
}

void mock_last_body(MockService *ms, char *buf, size_t buflen) {
    pthread_mutex_lock(&ms->lock);
    snprintf(buf, buflen, "%s", ms->last_body);
    pthread_mutex_unlock(&ms->lock);
}

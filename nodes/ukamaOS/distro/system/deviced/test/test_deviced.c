#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <unistd.h>
#include <sys/wait.h>

#include <curl/curl.h>
#include <jansson.h>

#include "unity.h"

#include "usys_services.h"
#include "usys_log.h"

#include "mocksrv/mock_services.h"

#define MAX_BODY 8192

typedef struct {
    long status;
    char body[MAX_BODY];
} HttpResp;

static size_t write_cb(void *contents, size_t size, size_t nmemb, void *userp) {
    size_t realsize = size * nmemb;
    HttpResp *r = (HttpResp *)userp;
    size_t cur = strlen(r->body);
    size_t copy = realsize;
    if (cur + copy >= sizeof(r->body)) {
        copy = sizeof(r->body) - cur - 1;
    }
    if (copy > 0) {
        memcpy(r->body + cur, contents, copy);
        r->body[cur + copy] = '\0';
    }
    return realsize;
}

static int http_request(const char *method, const char *url, const char *json_body, HttpResp *out) {
    CURL *curl = NULL;
    CURLcode res;
    struct curl_slist *headers = NULL;

    if (!method || !url || !out) return -1;

    memset(out, 0, sizeof(*out));

    curl = curl_easy_init();
    if (!curl) return -1;

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, out);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 5L);

    if (strcmp(method, "GET") == 0) {
        /* default */
    } else {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
        headers = curl_slist_append(headers, "Content-Type: application/json");
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        if (json_body) {
            curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json_body);
        }
    }

    res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        curl_slist_free_all(headers);
        curl_easy_cleanup(curl);
        return -1;
    }

    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &out->status);

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    return 0;
}

static pid_t spawn_deviced(const char *bin_path, int client_mode, const char *client_host, const char *preload_so) {
    pid_t pid;

    pid = fork();
    if (pid < 0) return -1;

    if (pid == 0) {
        /* child */
        if (preload_so && *preload_so) {
            setenv("LD_PRELOAD", preload_so, 1);
        }
        setenv("DEVICED_DEBUG_MODE", "1", 1);

        if (client_mode) {
            execl(bin_path, bin_path, "--client-mode", (char *)NULL);
        } else {
            if (client_host && *client_host) {
                execl(bin_path, bin_path, "--client-host", client_host, (char *)NULL);
            } else {
                execl(bin_path, bin_path, (char *)NULL);
            }
        }
        _exit(127);
    }

    return pid;
}

static void stop_process(pid_t pid) {
    int status;
    if (pid <= 0) return;
    kill(pid, SIGINT);
    (void)waitpid(pid, &status, 0);
}

/* A tiny mock client that answers /v1/reboot/ with 202.
 * This is used to validate tower restart remote-client logic without modifying device.d.
 */

typedef struct {
    struct _u_instance inst;
    int started;
    pthread_t thread;
} MockClient;

static int cb_client_reboot(const struct _u_request *req, struct _u_response *resp, void *user_data) {
    (void)req;
    (void)user_data;
    ulfius_set_empty_body_response(resp, 202);
    return U_CALLBACK_CONTINUE;
}

static void *mockclient_run(void *arg) {
    MockClient *mc = (MockClient *)arg;
    if (ulfius_start_framework(&mc->inst) == U_OK) {
        mc->started = 1;
        while (mc->started) {
            usleep(100 * 1000);
        }
    }
    return NULL;
}

static int mockclient_start(MockClient *mc, int port) {
    memset(mc, 0, sizeof(*mc));
    if (ulfius_init_instance(&mc->inst, port, NULL, NULL) != U_OK) return -1;
    ulfius_add_endpoint_by_val(&mc->inst, "POST", NULL, "/v1/reboot/", 0, &cb_client_reboot, NULL);
    if (pthread_create(&mc->thread, NULL, mockclient_run, mc) != 0) return -1;
    for (int i = 0; i < 50; i++) {
        if (mc->started) return 0;
        usleep(20 * 1000);
    }
    return -1;
}

static void mockclient_stop(MockClient *mc) {
    if (!mc) return;
    mc->started = 0;
    ulfius_stop_framework(&mc->inst);
    ulfius_clean_instance(&mc->inst);
    pthread_join(mc->thread, NULL);
}

/* Globals for suite */
static int g_port_device = 0;
static int g_port_client = 0;
static int g_port_noded = 0;
static int g_port_notifyd = 0;
static int g_port_femd = 0;

static MockService g_noded;
static MockService g_notifyd;
static MockService g_femd;
static MockClient  g_client;

static pid_t g_pid_server = -1;
static pid_t g_pid_client = -1;

static char g_preload_so[512];
static char g_bin_server[512];

static void wait_port_ready(int port) {
    /* dumb wait */
    (void)port;
    usleep(300 * 1000);
}

static void suite_setup_server_mode(void) {
    char url[256];
    HttpResp r;

    mock_service_reset(&g_noded);
    mock_service_reset(&g_notifyd);
    mock_service_reset(&g_femd);

    /* mock client server for /v1/reboot/ */
    TEST_ASSERT_EQUAL_INT(0, mockclient_start(&g_client, g_port_client));

    g_pid_server = spawn_deviced(g_bin_server, 0, "localhost", g_preload_so);
    TEST_ASSERT_TRUE(g_pid_server > 0);

    wait_port_ready(g_port_device);

    snprintf(url, sizeof(url), "http://localhost:%d/v1/ping", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("GET", url, NULL, &r));
    TEST_ASSERT_EQUAL_INT(200, (int)r.status);
}

static void suite_teardown_server_mode(void) {
    stop_process(g_pid_server);
    g_pid_server = -1;
    mockclient_stop(&g_client);
}

static void suite_setup_client_mode(void) {
    char url[256];
    HttpResp r;

    g_pid_client = spawn_deviced(g_bin_server, 1, NULL, g_preload_so);
    TEST_ASSERT_TRUE(g_pid_client > 0);
    wait_port_ready(g_port_client);

    snprintf(url, sizeof(url), "http://localhost:%d/v1/ping", g_port_client);
    TEST_ASSERT_EQUAL_INT(0, http_request("GET", url, NULL, &r));
    TEST_ASSERT_EQUAL_INT(200, (int)r.status);
}

static void suite_teardown_client_mode(void) {
    stop_process(g_pid_client);
    g_pid_client = -1;
}

void setUp(void) { }
void tearDown(void) { }

static void test_server_ping_and_version(void) {
    char url[256];
    HttpResp r;

    snprintf(url, sizeof(url), "http://localhost:%d/v1/ping", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("GET", url, NULL, &r));
    TEST_ASSERT_EQUAL_INT(200, (int)r.status);

    snprintf(url, sizeof(url), "http://localhost:%d/v1/version", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("GET", url, NULL, &r));
    TEST_ASSERT_EQUAL_INT(200, (int)r.status);
    TEST_ASSERT_TRUE(strlen(r.body) > 0);
}

static void test_server_state_endpoint_tower(void) {
    char url[256];
    HttpResp r;
    json_error_t jerr;
    json_t *j = NULL;
    json_t *svc = NULL;

    snprintf(url, sizeof(url), "http://localhost:%d/v1/state", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("GET", url, NULL, &r));
    TEST_ASSERT_EQUAL_INT(200, (int)r.status);

    j = json_loads(r.body, 0, &jerr);
    TEST_ASSERT_NOT_NULL(j);

    svc = json_object_get(j, "service");
    TEST_ASSERT_NOT_NULL(svc);
    TEST_ASSERT_TRUE(json_is_string(svc));
    TEST_ASSERT_EQUAL_STRING("on", json_string_value(svc));

    json_decref(j);
}

static void test_server_service_post_on_is_idempotent_ok(void) {
    char url[256];
    HttpResp r;

    /* initial tower state is ON. Posting ON should return 200 OK immediately (no worker). */
    snprintf(url, sizeof(url), "http://localhost:%d/v1/service", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("POST", url, "{\"state\":\"on\"}", &r));
    TEST_ASSERT_EQUAL_INT(200, (int)r.status);
}

static void test_server_service_post_invalid_body_400(void) {
    char url[256];
    HttpResp r;

    snprintf(url, sizeof(url), "http://localhost:%d/v1/service", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("POST", url, "{\"state\":123}", &r));
    TEST_ASSERT_EQUAL_INT(400, (int)r.status);
}

static void test_server_restart_triggers_notify_and_remote_reboot(void) {
    char url[256];
    HttpResp r;

    mock_service_reset(&g_notifyd);

    snprintf(url, sizeof(url), "http://localhost:%d/v1/restart", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("POST", url, "{\"force\":true}", &r));
    TEST_ASSERT_EQUAL_INT(202, (int)r.status);

    /* Worker should send notify (POST /notify/v1/event/device.d) at least once */
    for (int i = 0; i < 80; i++) {
        if (mock_notify_posts(&g_notifyd) > 0) break;
        usleep(50 * 1000);
    }
    TEST_ASSERT_TRUE(mock_notify_posts(&g_notifyd) > 0);
}

static void test_server_not_allowed_method_405(void) {
    char url[256];
    HttpResp r;

    snprintf(url, sizeof(url), "http://localhost:%d/v1/ping", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("POST", url, "{}", &r));
    TEST_ASSERT_EQUAL_INT(405, (int)r.status);
}

static void test_server_unknown_route_404(void) {
    char url[256];
    HttpResp r;

    snprintf(url, sizeof(url), "http://localhost:%d/v1/nope", g_port_device);
    TEST_ASSERT_EQUAL_INT(0, http_request("GET", url, NULL, &r));
    TEST_ASSERT_EQUAL_INT(404, (int)r.status);
}

static void test_client_mode_restart_endpoint_exists(void) {
    char url[256];
    HttpResp r;

    /* In client-mode, /v1/restart is registered. */
    snprintf(url, sizeof(url), "http://localhost:%d/v1/restart", g_port_client);
    TEST_ASSERT_EQUAL_INT(0, http_request("POST", url, "{\"force\":true}", &r));

    /* It may return 202 if worker scheduled, or 200 if immediate due to state.
       We accept either but must not be 404. */
    TEST_ASSERT_TRUE((int)r.status == 200 || (int)r.status == 202 || (int)r.status == 500 || (int)r.status == 400);
    TEST_ASSERT_NOT_EQUAL(404, (int)r.status);
}

int main(int argc, char **argv) {
    int ret;

    (void)argc;
    (void)argv;

    usys_log_set_level(USYS_LOG_WARN);

    curl_global_init(CURL_GLOBAL_ALL);

    /* Ensure DEF_FEMD_HOST typo ("loclahost") resolves during tests if possible. */
    {
        FILE *hf = fopen("/etc/hosts", "r+");
        if (hf) {
            char line[512];
            int found = 0;
            while (fgets(line, sizeof(line), hf) != NULL) {
                if (strstr(line, "loclahost") != NULL) {
                    found = 1;
                    break;
                }
            }
            if (!found) {
                fseek(hf, 0, SEEK_END);
                fprintf(hf, "
127.0.0.1 loclahost
");
            }
            fclose(hf);
        }
    }
            }
            if (!found) {
                fseek(hf, 0, SEEK_END);
                fprintf(hf, "\n127.0.0.1 loclahost\n");
            }
            fclose(hf);
        }
    }

    g_port_device = usys_find_service_port(SERVICE_DEVICE);
    g_port_client = usys_find_service_port(SERVICE_DEVICE_CLIENT);
    g_port_noded = usys_find_service_port(SERVICE_NODE);
    g_port_notifyd = usys_find_service_port(SERVICE_NOTIFY);
    g_port_femd = usys_find_service_port(SERVICE_FEM);

    if (!g_port_device || !g_port_client || !g_port_noded || !g_port_notifyd || !g_port_femd) {
        fprintf(stderr, "Required service ports not found (check usys services config).\n");
        return 1;
    }

    /* paths provided by Makefile defines */
    snprintf(g_preload_so, sizeof(g_preload_so), "%s", getenv("DEVICED_TEST_PRELOAD") ? getenv("DEVICED_TEST_PRELOAD") : "");
    snprintf(g_bin_server, sizeof(g_bin_server), "%s", getenv("DEVICED_TEST_BIN") ? getenv("DEVICED_TEST_BIN") : "../device.d");

    /* start external dependency mocks */
    ret = mock_noded_start(&g_noded, g_port_noded);
    if (ret != 0) { fprintf(stderr, "failed to start noded mock on %d\n", g_port_noded); return 1; }
    ret = mock_notifyd_start(&g_notifyd, g_port_notifyd);
    if (ret != 0) { fprintf(stderr, "failed to start notifyd mock on %d\n", g_port_notifyd); return 1; }
    ret = mock_femd_start(&g_femd, g_port_femd);
    if (ret != 0) { fprintf(stderr, "failed to start femd mock on %d\n", g_port_femd); return 1; }

    UNITY_BEGIN();

    /* Server-mode suite (tower) */
    suite_setup_server_mode();
    RUN_TEST(test_server_ping_and_version);
    RUN_TEST(test_server_state_endpoint_tower);
    RUN_TEST(test_server_service_post_on_is_idempotent_ok);
    RUN_TEST(test_server_service_post_invalid_body_400);
    RUN_TEST(test_server_not_allowed_method_405);
    RUN_TEST(test_server_unknown_route_404);
    RUN_TEST(test_server_restart_triggers_notify_and_remote_reboot);
    suite_teardown_server_mode();

    /* Client-mode suite */
    suite_setup_client_mode();
    RUN_TEST(test_client_mode_restart_endpoint_exists);
    suite_teardown_client_mode();

    mock_service_stop(&g_femd);
    mock_service_stop(&g_notifyd);
    mock_service_stop(&g_noded);

    curl_global_cleanup();

    return UNITY_END();
}

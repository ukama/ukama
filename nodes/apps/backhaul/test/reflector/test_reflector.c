/*
 * ukama_test_reflector.c
 *
 * Minimal Ukama-style bootstrap + reflector server for testing backhauld.
 *
 * Endpoints:
 *   1) Bootstrap discovery:
 *      GET /reflector
 *        -> { "reflectorNearUrl": ".../reflector/v1", "reflectorFarUrl": ".../reflector/v1", ... }
 *
 *   2) Reflector API (base is returned in reflectorNearUrl/FarUrl):
 *      GET  /reflector/v1/ping
 *      GET  /reflector/v1/download/:bytes
 *      POST /reflector/v1/upload
 *
 * Optional query params for fault injection:
 *   latency_ms=50   (adds base latency before responding)
 *   jitter_ms=20    (adds +/- jitter around latency)
 *   loss_pct=5      (randomly returns 503 with this probability)
 *   chunk_bytes=16384
 *   chunk_delay_ms=2  (throttle download by sleeping per chunk)
 *
 * Run:
 *   REFLECTOR_PORT=8088 REFLECTOR_HOST=127.0.0.1 REFLECTOR_SCHEME=http ./ukama_test_reflector
 *
 * Test:
 *   curl -s http://127.0.0.1:8088/reflector | jq
 *   curl -i "http://127.0.0.1:8088/reflector/v1/ping?latency_ms=80&jitter_ms=20"
 *   curl -o /dev/null -s "http://127.0.0.1:8088/reflector/v1/download/1048576?chunk_bytes=16384&chunk_delay_ms=1"
 *   dd if=/dev/zero bs=1M count=5 | curl -s -X POST --data-binary @- http://127.0.0.1:8088/reflector/v1/upload | jq
 */

#include <ulfius.h>
#include <jansson.h>

#include <openssl/sha.h>

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#define DEF_PORT    8088
#define DEF_SCHEME  "http"
#define DEF_HOST    "127.0.0.1"

#define BOOTSTRAP_EP       "/reflector"
#define REFLECTOR_API_BASE "/reflector/v1"

#define MAX_UPLOAD_BYTES (64 * 1024 * 1024)

#define MIN(a,b) ((a) < (b) ? (a) : (b))

#define STATUS_OK   (0)
#define STATUS_NOK  (-1)

typedef struct {
    int   port;
    char  *scheme;
    char  *host;

    /* Returned to client */
    char  nearUrl[256];
    char  farUrl[256];

    unsigned int rngSeed;
} AppConfig;

static int env_to_int(const char *name, int def) {
    const char *v = getenv(name);
    if (!v || !*v) return def;
    return atoi(v);
}

static char *env_to_strdup(const char *name, const char *def) {
    const char *v = getenv(name);
    if (!v || !*v) v = def;
    return v ? strdup(v) : NULL;
}

static void ms_sleep(int ms) {
    if (ms <= 0) return;
    usleep((useconds_t)ms * 1000);
}

static int rand_range(int minV, int maxV) {
    if (maxV <= minV) return minV;
    return minV + (rand() % (maxV - minV + 1));
}

/* Ulfius query param helper:
 * In Ulfius, query parameters are generally in request->map_url as well.
 */
static const char *qp_get(const struct _u_request *request, const char *name) {
    const char *v = u_map_get(request->map_url, name);
    return (v && *v) ? v : NULL;
}

static int qp_get_int(const struct _u_request *request,
                      const char *name,
                      int def) {
    const char *v = qp_get(request, name);
    if (!v) return def;
    return atoi(v);
}

static int should_drop(const struct _u_request *request) {
    int lossPct = qp_get_int(request, "loss_pct", 0);
    if (lossPct <= 0) return 0;
    if (lossPct >= 100) return 1;
    return ((rand() % 100) < lossPct) ? 1 : 0;
}

static void apply_latency(const struct _u_request *request) {
    int latencyMs = qp_get_int(request, "latency_ms", 0);
    int jitterMs  = qp_get_int(request, "jitter_ms", 0);

    if (latencyMs <= 0 && jitterMs <= 0) return;

    int addJitter = 0;
    if (jitterMs > 0) addJitter = rand_range(-jitterMs, jitterMs);

    ms_sleep(latencyMs + addJitter);
}

static void sha256_hex(const unsigned char *data,
                       size_t len,
                       char outHex[65]) {

    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256_CTX ctx;

    SHA256_Init(&ctx);
    SHA256_Update(&ctx, data, len);
    SHA256_Final(hash, &ctx);

    for (int i = 0; i < SHA256_DIGEST_LENGTH; i++) {
        sprintf(outHex + (i * 2), "%02x", hash[i]);
    }
    outHex[64] = '\0';
}

/* --------------------- Bootstrap discovery --------------------- */

static int cb_bootstrap(const struct _u_request *request,
                        struct _u_response *response,
                        void *user_data) {

    AppConfig *cfg = (AppConfig *)user_data;

    if (should_drop(request)) {
        ulfius_set_string_body_response(response, 503, "dropped");
        return U_CALLBACK_CONTINUE;
    }

    apply_latency(request);

    json_t *root = json_object();
    json_object_set_new(root, "ok", json_true());
    json_object_set_new(root, "ts", json_integer((json_int_t)time(NULL)));

    /* These two keys map cleanly into your Config fields */
    json_object_set_new(root, "reflectorNearUrl", json_string(cfg->nearUrl));
    json_object_set_new(root, "reflectorFarUrl",  json_string(cfg->farUrl));

    /* Optional extras; safe to ignore */
    json_object_set_new(root, "version", json_string("ukama-test-reflector-1"));

    ulfius_set_json_body_response(response, 200, root);
    json_decref(root);

    return U_CALLBACK_CONTINUE;
}

/* --------------------- Reflector API --------------------- */

static int cb_ping(const struct _u_request *request,
                   struct _u_response *response,
                   void *user_data) {

    (void)user_data;

    if (should_drop(request)) {
        ulfius_set_string_body_response(response, 503, "dropped");
        return U_CALLBACK_CONTINUE;
    }

    apply_latency(request);

    struct timespec ts;
    clock_gettime(CLOCK_REALTIME, &ts);

    char body[256];
    snprintf(body, sizeof(body), "OK ts=%ld.%03ld\n",
             (long)ts.tv_sec, (long)(ts.tv_nsec / 1000000));

    ulfius_set_string_body_response(response, 200, body);
    return U_CALLBACK_CONTINUE;
}

static int cb_download(const struct _u_request *request,
                       struct _u_response *response,
                       void *user_data) {

    (void)user_data;

    if (should_drop(request)) {
        ulfius_set_string_body_response(response, 503, "dropped");
        return U_CALLBACK_CONTINUE;
    }

    apply_latency(request);

    const char *bytesStr = u_map_get(request->map_url, "bytes");
    if (!bytesStr || !*bytesStr) {
        ulfius_set_string_body_response(response, 400, "missing :bytes");
        return U_CALLBACK_CONTINUE;
    }

    long long totalBytes = atoll(bytesStr);
    if (totalBytes <= 0 || totalBytes > (long long)(512 * 1024 * 1024)) {
        ulfius_set_string_body_response(response, 400, "invalid :bytes");
        return U_CALLBACK_CONTINUE;
    }

    int chunkBytes   = qp_get_int(request, "chunk_bytes", 16384);
    int chunkDelayMs = qp_get_int(request, "chunk_delay_ms", 0);

    if (chunkBytes < 1024) chunkBytes = 1024;
    if (chunkBytes > 1024 * 1024) chunkBytes = 1024 * 1024;

    char *chunk = (char *)malloc((size_t)chunkBytes);
    if (!chunk) {
        ulfius_set_string_body_response(response, 500, "oom");
        return U_CALLBACK_CONTINUE;
    }

    for (int i = 0; i < chunkBytes; i++) {
        chunk[i] = (char)('A' + (i % 26));
    }

    /* Simple: allocate whole body (fine for test). */
    char *body = (char *)malloc((size_t)totalBytes);
    if (!body) {
        free(chunk);
        ulfius_set_string_body_response(response, 500, "oom");
        return U_CALLBACK_CONTINUE;
    }

    long long off = 0;
    while (off < totalBytes) {
        long long n = MIN((long long)chunkBytes, totalBytes - off);
        memcpy(body + off, chunk, (size_t)n);
        off += n;
        if (chunkDelayMs > 0) ms_sleep(chunkDelayMs);
    }

    ulfius_set_binary_body_response(response, 200, body, (size_t)totalBytes);

    free(chunk);
    free(body);
    return U_CALLBACK_CONTINUE;
}

static int cb_upload(const struct _u_request *request,
                     struct _u_response *response,
                     void *user_data) {

    (void)user_data;

    if (should_drop(request)) {
        ulfius_set_string_body_response(response, 503, "dropped");
        return U_CALLBACK_CONTINUE;
    }

    apply_latency(request);

    if (!request->binary_body || request->binary_body_length == 0) {
        ulfius_set_string_body_response(response, 400, "missing body");
        return U_CALLBACK_CONTINUE;
    }

    if (request->binary_body_length > MAX_UPLOAD_BYTES) {
        ulfius_set_string_body_response(response, 413, "body too large");
        return U_CALLBACK_CONTINUE;
    }

    char hashHex[65] = {0};
    sha256_hex((const unsigned char *)request->binary_body,
               (size_t)request->binary_body_length,
               hashHex);

    json_t *root = json_object();
    json_object_set_new(root, "ok", json_true());
    json_object_set_new(root, "bytesReceived",
                        json_integer((json_int_t)request->binary_body_length));
    json_object_set_new(root, "sha256", json_string(hashHex));
    json_object_set_new(root, "ts", json_integer((json_int_t)time(NULL)));

    ulfius_set_json_body_response(response, 200, root);
    json_decref(root);

    return U_CALLBACK_CONTINUE;
}

/* --------------------- Setup --------------------- */

static void setup_endpoints(struct _u_instance *inst, AppConfig *cfg) {

    ulfius_add_endpoint_by_val(inst, "GET",  NULL, BOOTSTRAP_EP, 0, &cb_bootstrap, cfg);

    ulfius_add_endpoint_by_val(inst, "GET",  NULL, REFLECTOR_API_BASE "/ping", 0, &cb_ping, cfg);
    ulfius_add_endpoint_by_val(inst, "GET",  NULL, REFLECTOR_API_BASE "/download/:bytes", 0, &cb_download, cfg);
    ulfius_add_endpoint_by_val(inst, "POST", NULL, REFLECTOR_API_BASE "/upload", 0, &cb_upload, cfg);
}

static int config_load(AppConfig *cfg) {

    memset(cfg, 0, sizeof(*cfg));

    cfg->port   = env_to_int("REFLECTOR_PORT", DEF_PORT);
    cfg->scheme = env_to_strdup("REFLECTOR_SCHEME", DEF_SCHEME);
    cfg->host   = env_to_strdup("REFLECTOR_HOST", DEF_HOST);

    cfg->rngSeed = (unsigned int)time(NULL);
    const char *seed = getenv("REFLECTOR_SEED");
    if (seed && *seed) cfg->rngSeed = (unsigned int)atoi(seed);

    /* By default, near/far both point to this same server */
    snprintf(cfg->nearUrl, sizeof(cfg->nearUrl),
             "%s://%s:%d%s",
             cfg->scheme, cfg->host, cfg->port, REFLECTOR_API_BASE);

    snprintf(cfg->farUrl, sizeof(cfg->farUrl),
             "%s://%s:%d%s",
             cfg->scheme, cfg->host, cfg->port, REFLECTOR_API_BASE);

    return STATUS_OK;
}

static void config_free(AppConfig *cfg) {
    if (!cfg) return;
    free(cfg->scheme);
    free(cfg->host);
    memset(cfg, 0, sizeof(*cfg));
}

int main(void) {

    AppConfig cfg;
    struct _u_instance inst;

    if (config_load(&cfg) != STATUS_OK) {
        fprintf(stderr, "config_load failed\n");
        return 1;
    }

    srand(cfg.rngSeed);

    if (ulfius_init_instance(&inst, cfg.port, NULL, NULL) != U_OK) {
        fprintf(stderr, "ulfius_init_instance failed\n");
        config_free(&cfg);
        return 1;
    }

    u_map_put(inst.default_headers, "Access-Control-Allow-Origin", "*");
    u_map_put(inst.default_headers, "Server", "ukama-test-reflector");

    setup_endpoints(&inst, &cfg);

    printf("Test reflector listening on %s://%s:%d\n", cfg.scheme, cfg.host, cfg.port);
    printf("Bootstrap (discovery):  %s://%s:%d%s\n", cfg.scheme, cfg.host, cfg.port, BOOTSTRAP_EP);
    printf("Near URL returned:      %s\n", cfg.nearUrl);
    printf("Far URL returned:       %s\n", cfg.farUrl);
    printf("\n");

    if (ulfius_start_framework(&inst) != U_OK) {
        fprintf(stderr, "ulfius_start_framework failed\n");
        ulfius_clean_instance(&inst);
        config_free(&cfg);
        return 1;
    }

    pause();

    ulfius_stop_framework(&inst);
    ulfius_clean_instance(&inst);
    config_free(&cfg);

    return 0;
}

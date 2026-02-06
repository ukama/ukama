/*
 * test_reflector.c
 *
 * Ukama-style bootstrap + reflector server for testing backhauld/backhaul-monitor.
 *
 * Endpoints:
 *   GET  /reflector
 *     -> JSON: { "reflectorNearUrl": ".../reflector/v1", "reflectorFarUrl": ".../reflector/v1", ... }
 *
 *   GET  /reflector/v1/ping
 *   GET  /reflector/v1/download/:bytes
 *   POST /reflector/v1/upload
 *
 * Fault injection (optional query params):
 *   latency_ms=50
 *   jitter_ms=20
 *   loss_pct=5
 *   chunk_bytes=16384
 *   chunk_delay_ms=2
 *
 * Env:
 *   REFLECTOR_PORT=8088
 *   REFLECTOR_HOST=127.0.0.1
 *   REFLECTOR_SCHEME=http
 *   REFLECTOR_SEED=123
 *   REFLECTOR_MAX_UPLOAD_BYTES=67108864
 *   REFLECTOR_MAX_DOWNLOAD_BYTES=67108864
 *
 * Run: 
 *
 *   REFLECTOR_SCHEME=http \
 *   REFLECTOR_HOST=127.0.0.1 \
 *   REFLECTOR_PORT=8088       \
 *   REFLECTOR_SEED=1                  \
 *   REFLECTOR_MAX_UPLOAD_BYTES=67108864 \
 *   REFLECTOR_MAX_DOWNLOAD_BYTES=67108864         \
 *   ./ukama_test_reflector 
 *
 */

#include <ulfius.h>
#include <jansson.h>

#include <openssl/evp.h>

#include <errno.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include <arpa/inet.h>
#include <netinet/in.h>
#include <sys/socket.h>

#define DEF_PORT    8088
#define DEF_SCHEME  "http"
#define DEF_HOST    "127.0.0.1"

#define BOOTSTRAP_EP       "/reflector"
#define REFLECTOR_API_BASE "/reflector/v1"

#define DEF_MAX_UPLOAD_BYTES   (64 * 1024 * 1024)
#define DEF_MAX_DOWNLOAD_BYTES (64 * 1024 * 1024)

#define MIN(a,b) ((a) < (b) ? (a) : (b))

typedef struct {
	int   port;
	char  *scheme;
	char  *host;

	char  nearUrl[256];
	char  farUrl[256];

	unsigned int rngSeed;

	long maxUploadBytes;
	long maxDownloadBytes;

	unsigned long reqId;
} AppConfig;

static const char *safe_str(const char *s) {
	return (s && *s) ? s : "-";
}

static int env_to_int(const char *name, int def) {
	const char *v = getenv(name);
	if (!v || !*v) return def;
	return atoi(v);
}

static long env_to_long(const char *name, long def) {
	const char *v = getenv(name);
	if (!v || !*v) return def;
	errno = 0;
	long x = strtol(v, NULL, 10);
	if (errno != 0) return def;
	return x;
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

/* Ulfius: query + path params are typically visible via request->map_url */
static const char *qp_get(const struct _u_request *request, const char *name) {
	const char *v = u_map_get(request->map_url, name);
	return (v && *v) ? v : NULL;
}

static int qp_get_int(const struct _u_request *request, const char *name, int def) {
	const char *v = qp_get(request, name);
	if (!v) return def;
	return atoi(v);
}

static unsigned long next_req_id(AppConfig *cfg) {
	cfg->reqId += 1;
	return cfg->reqId;
}

static void sockaddr_to_str(const struct sockaddr *sa, char *out, size_t outLen) {
	if (!out || outLen == 0) return;
	out[0] = '\0';

	if (!sa) {
		snprintf(out, outLen, "-");
		return;
	}

	if (sa->sa_family == AF_INET) {
		const struct sockaddr_in *s = (const struct sockaddr_in *)sa;
		char ip[INET_ADDRSTRLEN] = {0};
		if (inet_ntop(AF_INET, &s->sin_addr, ip, sizeof(ip))) {
			snprintf(out, outLen, "%s:%u", ip, (unsigned)ntohs(s->sin_port));
		} else {
			snprintf(out, outLen, "ipv4:?"); /* best-effort */
		}
		return;
	}

	if (sa->sa_family == AF_INET6) {
		const struct sockaddr_in6 *s = (const struct sockaddr_in6 *)sa;
		char ip[INET6_ADDRSTRLEN] = {0};
		if (inet_ntop(AF_INET6, &s->sin6_addr, ip, sizeof(ip))) {
			snprintf(out, outLen, "[%s]:%u", ip, (unsigned)ntohs(s->sin6_port));
		} else {
			snprintf(out, outLen, "ipv6:?");
		}
		return;
	}

	snprintf(out, outLen, "af=%d", (int)sa->sa_family);
}

static void dbg_print_map(const char *label, const struct _u_map *m) {
	if (!m) return;

	int count = u_map_count(m);
	if (count <= 0) return;

	const char **keys = u_map_enum_keys(m);
	if (!keys) return;

	fprintf(stderr, "  %s:\n", safe_str(label));

	/* In this Ulfius API, keys are owned by the map. Do NOT free(). */
	for (int i = 0; i < count; i++) {
		const char *k = keys[i];
		if (!k || !*k) continue;
		const char *v = u_map_get(m, k);
		fprintf(stderr, "    %s=%s\n", safe_str(k), safe_str(v));
	}
}

static void dbg_log_request(unsigned long id, const char *ep, const struct _u_request *request) {
	const char *method = request ? request->http_verb : NULL;
	const char *path   = request ? request->http_url  : NULL;

	char remote[128];
	remote[0] = '\0';
	if (request) {
		sockaddr_to_str((const struct sockaddr *)request->client_address, remote, sizeof(remote));
	}

	fprintf(stderr, "[REFLECTOR] #%lu IN  %s\n", id, safe_str(ep));
	fprintf(stderr, "  method=%s path=%s remote=%s body_len=%zu\n",
	        safe_str(method),
	        safe_str(path),
	        safe_str(remote),
	        request ? (size_t)request->binary_body_length : 0UL);

	if (request) {
		const struct _u_map *params  = request->map_url;
		const struct _u_map *headers = request->map_header;

		dbg_print_map("params",  params);
		dbg_print_map("headers", headers);
	}
}

static void dbg_log_response(unsigned long id, const char *ep, int status, size_t bytes) {
	fprintf(stderr, "[REFLECTOR] #%lu OUT %s status=%d bytes=%zu\n",
	        id, safe_str(ep), status, bytes);
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

static int sha256_hex(const unsigned char *data, size_t len, char outHex[65]) {
	unsigned char md[EVP_MAX_MD_SIZE];
	unsigned int mdLen = 0;

	if (!outHex) return -1;
	outHex[0] = '\0';

	EVP_MD_CTX *ctx = EVP_MD_CTX_new();
	if (!ctx) return -1;

	if (EVP_DigestInit_ex(ctx, EVP_sha256(), NULL) != 1 ||
	    EVP_DigestUpdate(ctx, data, len) != 1 ||
	    EVP_DigestFinal_ex(ctx, md, &mdLen) != 1) {
		EVP_MD_CTX_free(ctx);
		return -1;
	}

	EVP_MD_CTX_free(ctx);

	if (mdLen != 32) return -1;

	for (unsigned int i = 0; i < mdLen; i++) {
		sprintf(outHex + (i * 2), "%02x", md[i]);
	}
	outHex[64] = '\0';
	return 0;
}

static void build_urls(AppConfig *cfg) {
	snprintf(cfg->nearUrl, sizeof(cfg->nearUrl),
	         "%s://%s:%d%s",
	         cfg->scheme, cfg->host, cfg->port, REFLECTOR_API_BASE);

	snprintf(cfg->farUrl, sizeof(cfg->farUrl),
	         "%s://%s:%d%s",
	         cfg->scheme, cfg->host, cfg->port, REFLECTOR_API_BASE);
}

static void print_config(const AppConfig *cfg) {
	fprintf(stderr, "[REFLECTOR] listen=%s://%s:%d\n",
	        safe_str(cfg->scheme), safe_str(cfg->host), cfg->port);
	fprintf(stderr, "[REFLECTOR] bootstrap=%s://%s:%d%s\n",
	        safe_str(cfg->scheme), safe_str(cfg->host), cfg->port, BOOTSTRAP_EP);
	fprintf(stderr, "[REFLECTOR] nearUrl=%s\n", safe_str(cfg->nearUrl));
	fprintf(stderr, "[REFLECTOR] farUrl=%s\n", safe_str(cfg->farUrl));
	fprintf(stderr, "[REFLECTOR] maxUploadBytes=%ld maxDownloadBytes=%ld\n",
	        cfg->maxUploadBytes, cfg->maxDownloadBytes);
}

static int cb_bootstrap(const struct _u_request *request,
                        struct _u_response *response,
                        void *user_data) {

	AppConfig *cfg = (AppConfig *)user_data;
	unsigned long id = next_req_id(cfg);

	dbg_log_request(id, "GET /reflector", request);

	if (should_drop(request)) {
		const char *msg = "dropped";
		ulfius_set_string_body_response(response, 503, msg);
		dbg_log_response(id, "GET /reflector", 503, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	apply_latency(request);

	json_t *root = json_object();
	json_object_set_new(root, "ok", json_true());
	json_object_set_new(root, "ts", json_integer((json_int_t)time(NULL)));

	/* Keys expected by backhaul-monitor Config */
	json_object_set_new(root, "reflectorNearUrl", json_string(cfg->nearUrl));
	json_object_set_new(root, "reflectorFarUrl",  json_string(cfg->farUrl));

	json_object_set_new(root, "version", json_string("ukama-test-reflector-1"));

	ulfius_set_json_body_response(response, 200, root);

	/* log size without leaking full body */
	char *dump = json_dumps(root, JSON_COMPACT);
	dbg_log_response(id, "GET /reflector", 200, dump ? strlen(dump) : 0);
	if (dump) free(dump);

	json_decref(root);
	return U_CALLBACK_CONTINUE;
}

static int cb_ping(const struct _u_request *request,
                   struct _u_response *response,
                   void *user_data) {

	AppConfig *cfg = (AppConfig *)user_data;
	unsigned long id = next_req_id(cfg);

	dbg_log_request(id, "GET /reflector/v1/ping", request);

	if (should_drop(request)) {
		const char *msg = "dropped";
		ulfius_set_string_body_response(response, 503, msg);
		dbg_log_response(id, "GET /reflector/v1/ping", 503, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	apply_latency(request);

	struct timespec ts;
	clock_gettime(CLOCK_REALTIME, &ts);

	char body[128];
	snprintf(body, sizeof(body), "OK ts=%ld.%03ld\n",
	         (long)ts.tv_sec,
	         (long)(ts.tv_nsec / 1000000L));

	ulfius_set_string_body_response(response, 200, body);
	dbg_log_response(id, "GET /reflector/v1/ping", 200, strlen(body));
	return U_CALLBACK_CONTINUE;
}

static int cb_download(const struct _u_request *request,
                       struct _u_response *response,
                       void *user_data) {

	AppConfig *cfg = (AppConfig *)user_data;
	unsigned long id = next_req_id(cfg);

	dbg_log_request(id, "GET /reflector/v1/download/:bytes", request);

	if (should_drop(request)) {
		const char *msg = "dropped";
		ulfius_set_string_body_response(response, 503, msg);
		dbg_log_response(id, "GET /reflector/v1/download/:bytes", 503, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	apply_latency(request);

	/* Ulfius puts :bytes into map_url under "bytes" */
	const char *bytesStr = u_map_get(request->map_url, "bytes");
	if (!bytesStr || !*bytesStr) {
		const char *msg = "missing :bytes";
		ulfius_set_string_body_response(response, 400, msg);
		dbg_log_response(id, "GET /reflector/v1/download/:bytes", 400, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	errno = 0;
	long long totalBytes = strtoll(bytesStr, NULL, 10);
	if (errno != 0 || totalBytes <= 0) {
		const char *msg = "invalid :bytes";
		ulfius_set_string_body_response(response, 400, msg);
		dbg_log_response(id, "GET /reflector/v1/download/:bytes", 400, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	if (totalBytes > (long long)cfg->maxDownloadBytes) {
		char msg[128];
		snprintf(msg, sizeof(msg), "too large (max=%ld)", cfg->maxDownloadBytes);
		ulfius_set_string_body_response(response, 413, msg);
		dbg_log_response(id, "GET /reflector/v1/download/:bytes", 413, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	int chunkBytes   = qp_get_int(request, "chunk_bytes", 16384);
	int chunkDelayMs = qp_get_int(request, "chunk_delay_ms", 0);

	if (chunkBytes < 1024) chunkBytes = 1024;
	if (chunkBytes > 1024 * 1024) chunkBytes = 1024 * 1024;

	/* test/simple: allocate body in memory (bounded by maxDownloadBytes) */
	char *body = (char *)malloc((size_t)totalBytes);
	if (!body) {
		const char *msg = "oom";
		ulfius_set_string_body_response(response, 500, msg);
		dbg_log_response(id, "GET /reflector/v1/download/:bytes", 500, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	/* fill deterministic pattern, with optional chunk delay */
	long long off = 0;
	while (off < totalBytes) {
		long long n = MIN((long long)chunkBytes, totalBytes - off);
		for (long long i = 0; i < n; i++) {
			body[off + i] = (char)('A' + ((off + i) % 26));
		}
		off += n;
		if (chunkDelayMs > 0) ms_sleep(chunkDelayMs);
	}

	ulfius_set_binary_body_response(response, 200, body, (size_t)totalBytes);
	dbg_log_response(id, "GET /reflector/v1/download/:bytes", 200, (size_t)totalBytes);

	free(body);
	return U_CALLBACK_CONTINUE;
}

static int cb_upload(const struct _u_request *request,
                     struct _u_response *response,
                     void *user_data) {

	AppConfig *cfg = (AppConfig *)user_data;
	unsigned long id = next_req_id(cfg);

	dbg_log_request(id, "POST /reflector/v1/upload", request);

	if (should_drop(request)) {
		const char *msg = "dropped";
		ulfius_set_string_body_response(response, 503, msg);
		dbg_log_response(id, "POST /reflector/v1/upload", 503, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	apply_latency(request);

	if (!request->binary_body || request->binary_body_length == 0) {
		const char *msg = "missing body";
		ulfius_set_string_body_response(response, 400, msg);
		dbg_log_response(id, "POST /reflector/v1/upload", 400, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	if ((long)request->binary_body_length > cfg->maxUploadBytes) {
		const char *msg = "body too large";
		ulfius_set_string_body_response(response, 413, msg);
		dbg_log_response(id, "POST /reflector/v1/upload", 413, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	char hashHex[65];
	if (sha256_hex((const unsigned char *)request->binary_body,
	               (size_t)request->binary_body_length,
	               hashHex) != 0) {
		const char *msg = "sha256 failed";
		ulfius_set_string_body_response(response, 500, msg);
		dbg_log_response(id, "POST /reflector/v1/upload", 500, strlen(msg));
		return U_CALLBACK_CONTINUE;
	}

	json_t *root = json_object();
	json_object_set_new(root, "ok", json_true());
	json_object_set_new(root, "bytesReceived",
	                    json_integer((json_int_t)request->binary_body_length));
	json_object_set_new(root, "sha256", json_string(hashHex));
	json_object_set_new(root, "ts", json_integer((json_int_t)time(NULL)));

	ulfius_set_json_body_response(response, 200, root);

	char *dump = json_dumps(root, JSON_COMPACT);
	dbg_log_response(id, "POST /reflector/v1/upload", 200, dump ? strlen(dump) : 0);
	if (dump) free(dump);

	json_decref(root);
	return U_CALLBACK_CONTINUE;
}

static void setup_endpoints(struct _u_instance *inst, AppConfig *cfg) {

	ulfius_add_endpoint_by_val(inst, "GET",  NULL, BOOTSTRAP_EP, 0, &cb_bootstrap, cfg);

	ulfius_add_endpoint_by_val(inst, "GET",  NULL, REFLECTOR_API_BASE "/ping", 0, &cb_ping, cfg);
	ulfius_add_endpoint_by_val(inst, "GET",  NULL, REFLECTOR_API_BASE "/download/:bytes",
                               0, &cb_download, cfg);
	ulfius_add_endpoint_by_val(inst, "POST", NULL, REFLECTOR_API_BASE "/upload", 0, &cb_upload, cfg);
}

static int config_load(AppConfig *cfg) {
	memset(cfg, 0, sizeof(*cfg));

	cfg->port   = env_to_int("REFLECTOR_PORT", DEF_PORT);
	cfg->scheme = env_to_strdup("REFLECTOR_SCHEME", DEF_SCHEME);
	cfg->host   = env_to_strdup("REFLECTOR_HOST", DEF_HOST);

	cfg->rngSeed = (unsigned int)time(NULL);
	{
		const char *seed = getenv("REFLECTOR_SEED");
		if (seed && *seed) cfg->rngSeed = (unsigned int)atoi(seed);
	}

	cfg->maxUploadBytes   = env_to_long("REFLECTOR_MAX_UPLOAD_BYTES", DEF_MAX_UPLOAD_BYTES);
	cfg->maxDownloadBytes = env_to_long("REFLECTOR_MAX_DOWNLOAD_BYTES", DEF_MAX_DOWNLOAD_BYTES);

	if (cfg->maxUploadBytes <= 0) cfg->maxUploadBytes = DEF_MAX_UPLOAD_BYTES;
	if (cfg->maxDownloadBytes <= 0) cfg->maxDownloadBytes = DEF_MAX_DOWNLOAD_BYTES;

	build_urls(cfg);
	return 0;
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

	if (config_load(&cfg) != 0) {
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

	print_config(&cfg);
	fprintf(stderr, "\n");

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

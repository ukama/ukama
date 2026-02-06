/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include <curl/curl.h>
#include "jansson.h"

#include "backhauld.h"
#include "web_client.h"
#include "json_types.h"
#include "json_serdes.h"
#include "usys_log.h"
#include "usys_mem.h"

#define REF_PING_EP   "/v1/ping"
#define REF_DL_EP     "/v1/download"
#define REF_UL_EP     "/v1/upload"

#ifndef REFLECTOR_NEAR_KEY
#define REFLECTOR_NEAR_KEY "reflectorNearUrl"
#endif

#ifndef REFLECTOR_FAR_KEY
#define REFLECTOR_FAR_KEY  "reflectorFarUrl"
#endif

typedef struct {
	size_t	bytes;
	size_t	limit;
	int		stop;
} Sink;

static size_t sink_write_cb(void *ptr, size_t size, size_t nmemb, void *userdata) {

	Sink *s = (Sink *)userdata;
	size_t n = size * nmemb;

	if (!s) return n;

	s->bytes += n;

	if (s->limit > 0 && s->bytes >= s->limit) {
		/* stop early by returning 0 -> causes CURLE_WRITE_ERROR,
		   but we avoid early stop in throughput runs, keep for future */
	}

	return n;
}

static size_t str_write_cb(void *ptr, size_t size, size_t nmemb, void *userdata) {

	size_t n = size * nmemb;
	char **buf = (char **)userdata;

	if (!buf) return n;

	size_t curLen = (*buf) ? strlen(*buf) : 0;
	char *nb = (char *)realloc(*buf, curLen + n + 1);
	if (!nb) return 0;

	memcpy(nb + curLen, ptr, n);
	nb[curLen + n] = 0;
	*buf = nb;

	return n;
}

static const char *curl_rc_str(CURLcode rc) {
	return curl_easy_strerror(rc);
}

static void log_curl_failure(const char *tag,
                             const char *url,
                             CURLcode rc,
                             long httpCode,
                             double ttfb,
                             double total,
                             size_t bytes,
                             const char *respStr /* optional, may be NULL */) {

	/* Keep snippet small to avoid log spam */
	char snippet[161];
	snippet[0] = '\0';

	if (respStr && *respStr) {
		size_t n = strlen(respStr);
		if (n > 160) n = 160;
		memcpy(snippet, respStr, n);
		snippet[n] = '\0';
	}

	/* Loud + actionable */
	usys_log_error("[web_client] %s FAILED: url=%s curl_rc=%d(%s) http=%ld ttfb=%.2fms total=%.2fms bytes=%zu%s%s",
	               tag,
	               (url && *url) ? url : "-",
	               (int)rc,
	               curl_rc_str(rc),
	               httpCode,
	               ttfb * 1000.0,
	               total * 1000.0,
	               bytes,
	               (snippet[0] ? " resp_snip=\"" : ""),
	               (snippet[0] ? snippet : ""));

	if (snippet[0]) {
		usys_log_error("[web_client] %s resp_snip_end", tag);
	}
}

static int http_2xx(long code) {
	return (code >= 200 && code < 300);
}

static int do_get(Config *config, const char *url, ProbeResult *out, char **respStr) {

	CURL *curl = NULL;
	CURLcode rc = CURLE_OK;

	Sink sink = {0};
	double ttfb = 0.0, total = 0.0;
	long httpCode = 0;

	if (!config || !url || !*url) return STATUS_NOK;

	if (out) memset(out, 0, sizeof(*out));
	if (respStr) *respStr = NULL;

	curl = curl_easy_init();
	if (!curl) {
		usys_log_error("[web_client] do_get: curl_easy_init failed url=%s", url);
		return STATUS_NOK;
	}

	curl_easy_setopt(curl, CURLOPT_URL, url);
	curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, (long)config->connectTimeoutMs);
	curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, (long)config->totalTimeoutMs);
	curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);

	/* Optional: make DNS/proxy surprises visible quickly */
	/* curl_easy_setopt(curl, CURLOPT_NOPROXY, "*"); */

	if (respStr) {
		curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, str_write_cb);
		curl_easy_setopt(curl, CURLOPT_WRITEDATA, respStr);
	} else {
		curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, sink_write_cb);
		curl_easy_setopt(curl, CURLOPT_WRITEDATA, &sink);
	}

	rc = curl_easy_perform(curl);

	/* Always try to read metadata, even on failure */
	curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &httpCode);
	curl_easy_getinfo(curl, CURLINFO_STARTTRANSFER_TIME, &ttfb);
	curl_easy_getinfo(curl, CURLINFO_TOTAL_TIME, &total);

	int ok = (rc == CURLE_OK && http_2xx(httpCode)) ? 1 : 0;

	if (out) {
		out->httpCode = httpCode;
		out->ttfbMs   = ttfb * 1000.0;
		out->totalMs  = total * 1000.0;
		out->ok       = ok;
		out->stalled  = (ok && out->ttfbMs >= (double)config->stallThresholdMs) ? 1 : 0;
	}

	if (!ok) {
		log_curl_failure("GET", url, rc, httpCode, ttfb, total,
		                 respStr && *respStr ? strlen(*respStr) : sink.bytes,
		                 (respStr ? (const char *)(*respStr) : NULL));
		/* Clean up before returning */
		curl_easy_cleanup(curl);
		/* If we allocated respStr but call failed, free it to avoid leaks */
		if (respStr && *respStr) {
			free(*respStr);
			*respStr = NULL;
		}
		return STATUS_NOK;
	}

	curl_easy_cleanup(curl);
	return STATUS_OK;
}

static int do_post(Config *config,
                   const char *url,
                   const void *body,
                   size_t bodyLen,
                   TransferResult *out) {

	CURL *curl = NULL;
	CURLcode rc = CURLE_OK;

	Sink sink = {0};
	double total = 0.0;
	long httpCode = 0;

	if (!config || !url || !*url || !body) return STATUS_NOK;
	if (bodyLen == 0) return STATUS_NOK;

	if (out) memset(out, 0, sizeof(*out));

	curl = curl_easy_init();
	if (!curl) {
		usys_log_error("[web_client] do_post: curl_easy_init failed url=%s", url);
		return STATUS_NOK;
	}

	curl_easy_setopt(curl, CURLOPT_URL, url);
	curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, (long)config->connectTimeoutMs);
	curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, (long)config->totalTimeoutMs);

	curl_easy_setopt(curl, CURLOPT_POST, 1L);
	curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
	curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, (long)bodyLen);

	/* Read response to complete transfer */
	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, sink_write_cb);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, &sink);

	rc = curl_easy_perform(curl);

	curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &httpCode);
	curl_easy_getinfo(curl, CURLINFO_TOTAL_TIME, &total);

	int ok = (rc == CURLE_OK && http_2xx(httpCode)) ? 1 : 0;

	if (out) {
		out->httpCode = httpCode;
		out->seconds  = total;
		out->ok       = ok;

		/* compute mbps only if ok and timing valid */
		if (ok && total > 0.0) {
			double mbits = ((double)bodyLen * 8.0) / 1000000.0;
			out->mbps = mbits / total;
		} else {
			out->mbps = 0.0;
		}
	}

	if (!ok) {
		log_curl_failure("POST", url, rc, httpCode, 0.0, total, sink.bytes, NULL);
		curl_easy_cleanup(curl);
		return STATUS_NOK;
	}

	curl_easy_cleanup(curl);
	return STATUS_OK;
}

int wc_init(void) {
	return (curl_global_init(CURL_GLOBAL_DEFAULT) == 0) ? USYS_TRUE : USYS_FALSE;
}

void wc_cleanup(void) {
	curl_global_cleanup();
}

static void join_url(char *dst, size_t dstLen, const char *baseUrl, const char *ep) {

	if (!dst || dstLen == 0) return;
	snprintf(dst, dstLen, "%s%s", baseUrl, ep);
}

int wc_fetch_reflectors(Config *config, ReflectorSet *set) {

	char url[512] = {0};
	char *resp = NULL;
	ProbeResult pr;

	if (!config || !set) return STATUS_NOK;

	memset(set, 0, sizeof(*set));

	/* ENV overrides win */
	if (config->reflectorNearUrl && *config->reflectorNearUrl &&
		config->reflectorFarUrl && *config->reflectorFarUrl) {
		strncpy(set->nearUrl, config->reflectorNearUrl, sizeof(set->nearUrl)-1);
		strncpy(set->farUrl,  config->reflectorFarUrl,  sizeof(set->farUrl)-1);
		set->ts = time(NULL);
		return STATUS_OK;
	}

	snprintf(url, sizeof(url), "%s://%s%s",
			 config->bootstrapScheme,
			 config->bootstrapHost,
			 config->bootstrapEp);

	if (do_get(config, url, &pr, &resp) != STATUS_OK || !pr.ok || !resp) {
		if (resp) free(resp);
		return STATUS_NOK;
	}

	json_error_t jerr;
	json_t *root = json_loads(resp, 0, &jerr);
	free(resp);

	if (!root) return STATUS_NOK;

    json_t *jn = json_object_get(root, REFLECTOR_NEAR_KEY);
    json_t *jf = json_object_get(root, REFLECTOR_FAR_KEY);

	if (!json_is_string(jn) || !json_is_string(jf)) {
		json_decref(root);
		return STATUS_NOK;
	}

	strncpy(set->nearUrl, json_string_value(jn), sizeof(set->nearUrl)-1);
	strncpy(set->farUrl,  json_string_value(jf), sizeof(set->farUrl)-1);
	set->ts = time(NULL);

	json_decref(root);
	return STATUS_OK;
}

int wc_probe_ping(Config *config, const char *baseUrl, ProbeResult *out) {

	char url[512] = {0};

	if (!config || !baseUrl || !*baseUrl || !out) return STATUS_NOK;

	join_url(url, sizeof(url), baseUrl, REF_PING_EP);
	return do_get(config, url, out, NULL);
}

int wc_download_blob(Config *config, const char *baseUrl, int bytes, TransferResult *out) {

	char url[512] = {0};
	CURL *curl = NULL;
	CURLcode rc;

	Sink sink = {0};
	double total=0;
	long httpCode=0;
	double dlSpeed=0;

	if (!config || !baseUrl || !out) return STATUS_NOK;
	memset(out, 0, sizeof(*out));

    snprintf(url, sizeof(url), "%s%s/%d", baseUrl, REF_DL_EP, bytes);

	curl = curl_easy_init();
	if (!curl) return STATUS_NOK;

	curl_easy_setopt(curl, CURLOPT_URL, url);
	curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, (long)config->connectTimeoutMs);
	curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, (long)config->totalTimeoutMs);
	curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);

	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, sink_write_cb);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, &sink);

	rc = curl_easy_perform(curl);

	curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &httpCode);
	curl_easy_getinfo(curl, CURLINFO_TOTAL_TIME, &total);
	curl_easy_getinfo(curl, CURLINFO_SPEED_DOWNLOAD, &dlSpeed); /* bytes/sec */

	out->httpCode = httpCode;
	out->seconds = total;
	out->ok = (rc == CURLE_OK && httpCode >= 200 && httpCode < 300) ? 1 : 0;

	if (out->ok && dlSpeed > 0.0) {
		double mbps = (dlSpeed * 8.0) / 1000000.0;
		out->mbps = mbps;
	} else if (out->ok && total > 0.0) {
		double mbits = ((double)sink.bytes * 8.0) / 1000000.0;
		out->mbps = mbits / total;
	}

	curl_easy_cleanup(curl);
	return STATUS_OK;
}

int wc_upload_echo(Config *config, const char *baseUrl, int bytes, TransferResult *out) {

	char url[512] = {0};

	if (!config || !baseUrl || !out) return STATUS_NOK;

    join_url(url, sizeof(url), baseUrl, REF_UL_EP);

	void *buf = usys_calloc(1, bytes);
	if (!buf) return STATUS_NOK;

	/* no need to randomize; reflector should accept any bytes */
	int rc = do_post(config, url, buf, (size_t)bytes, out);

	usys_free(buf);
	return rc;
}

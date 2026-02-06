/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include <curl/curl.h>
#include <jansson.h>

#include "agents.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_api.h"
#include "usys_file.h"
#include "usys_services.h"

#ifndef RETURN_OK
#define RETURN_OK 0
#endif

#ifndef RETURN_NOTOK
#define RETURN_NOTOK -1
#endif

typedef struct {
	char	*data;
	size_t	len;
} HttpBuf;

typedef struct {
	int		stateCode;
	int		linkGuessCode;
	double	confidence;
	double	dlGoodputMbps;
	double	ulGoodputMbps;
	double	bufferbloatInflationFactor;
	double	capDetectedMbps;
	double	nearTtfbMedianMs;
	double	nearTtfbP95Ms;
	double	nearTtfbP99Ms;
	double	farTtfbMedianMs;
	double	farTtfbP95Ms;
	double	farTtfbP99Ms;
	double	probeSuccessRatePct;
	double	stallRatePct;
	double	lastMicroTs;
	double	lastMultiTs;
	double	lastChgTs;
	double	lastClassifyTs;
	double	lastDiagTs;
	double	diagPresent;
} BackhaulMetrics;

static int str_ieq(const char *a, const char *b) {
	if (!a || !b) return 0;
	while (*a && *b) {
		char ca = *a;
		char cb = *b;
		if (ca >= 'A' && ca <= 'Z') ca = (char)(ca - 'A' + 'a');
		if (cb >= 'A' && cb <= 'Z') cb = (char)(cb - 'A' + 'a');
		if (ca != cb) return 0;
		a++;
		b++;
	}
	return (*a == '\0' && *b == '\0');
}

/* very small helper: “contains”, but case-insensitive */
static int str_icontains(const char *hay, const char *needle) {
	size_t i = 0, nlen = 0, hlen = 0;

	if (!hay || !needle) return 0;

	hlen = strlen(hay);
	nlen = strlen(needle);
	if (nlen == 0 || hlen < nlen) return 0;

	for (i = 0; i + nlen <= hlen; i++) {
		size_t j = 0;
		for (j = 0; j < nlen; j++) {
			char c1 = hay[i + j];
			char c2 = needle[j];
			if (c1 >= 'A' && c1 <= 'Z') c1 = (char)(c1 - 'A' + 'a');
			if (c2 >= 'A' && c2 <= 'Z') c2 = (char)(c2 - 'A' + 'a');
			if (c1 != c2) break;
		}
		if (j == nlen) return 1;
	}
	return 0;
}

static size_t curl_write_cb(void *contents, size_t size, size_t nmemb, void *userp) {
	size_t total = size * nmemb;
	HttpBuf *buf = (HttpBuf *)userp;

	char *newMem = realloc(buf->data, buf->len + total + 1);
	if (!newMem) {
		return 0;
	}

	buf->data = newMem;
	memcpy(buf->data + buf->len, contents, total);
	buf->len += total;
	buf->data[buf->len] = '\0';

	return total;
}

static int http_get_json(const char *url, json_t **outRoot) {
	int ret = RETURN_NOTOK;
	CURL *curl = NULL;
	CURLcode cres = CURLE_OK;
	HttpBuf buf = {0};
	json_error_t jerr;

	if (!url || !outRoot) return RETURN_NOTOK;

	curl = curl_easy_init();
	if (!curl) {
		usys_log_error("backhaul_agent: curl_easy_init failed");
		return RETURN_NOTOK;
	}

	curl_easy_setopt(curl, CURLOPT_URL, url);
	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, curl_write_cb);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, &buf);
	curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, 1500L);
	curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, 700L);
	curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);

	cres = curl_easy_perform(curl);
	if (cres != CURLE_OK) {
		usys_log_error("backhaul_agent: curl failed url=%s err=%s",
		               url, curl_easy_strerror(cres));
		goto done;
	}

	*outRoot = json_loads(buf.data ? buf.data : "{}", 0, &jerr);
	if (!*outRoot) {
		usys_log_error("backhaul_agent: json parse failed url=%s line=%d: %s",
		               url, jerr.line, jerr.text);
		goto done;
	}

	ret = RETURN_OK;

done:
	if (curl) curl_easy_cleanup(curl);
	if (buf.data) free(buf.data);
	return ret;
}

static int map_backhaul_state(const char *s) {
	/* DOWN=0, LIMITED=1, GOOD=2, UNKNOWN=3 */
	if (!s) return 3;
	if (str_ieq(s, "DOWN")) return 0;
	if (str_ieq(s, "LIMITED")) return 1;
	if (str_ieq(s, "GOOD")) return 2;
	return 3;
}

static int map_link_guess(const char *s) {
	/* UNKNOWN=0, TERRESTRIAL_LIKE=1, SAT_LIKE=2, CELLULAR_LIKE=3 */
	if (!s) return 0;
	if (str_ieq(s, "TERRESTRIAL_LIKE")) return 1;
	if (str_ieq(s, "SAT_LIKE")) return 2;
	if (str_ieq(s, "CELLULAR_LIKE")) return 3;
	return 0;
}

static double j_get_num(json_t *root, const char *key) {
	json_t *v = NULL;
	if (!root || !key) return 0;

	v = json_object_get(root, key);
	if (!v) return 0;

	if (json_is_real(v)) return json_real_value(v);
	if (json_is_integer(v)) return (double)json_integer_value(v);

	return 0;
}

static const char *j_get_str(json_t *root, const char *key) {
	json_t *v = NULL;
	if (!root || !key) return NULL;

	v = json_object_get(root, key);
	if (!v) return NULL;
	if (!json_is_string(v)) return NULL;

	return json_string_value(v);
}

static int backhaul_read_metrics(const char *url, BackhaulMetrics *out) {
	int ret = RETURN_NOTOK;
	json_t *root = NULL;

	if (!out) return RETURN_NOTOK;
	memset(out, 0, sizeof(*out));

	if (http_get_json(url, &root) != RETURN_OK) {
		return RETURN_NOTOK;
	}

	out->stateCode = map_backhaul_state(j_get_str(root, "backhaulState"));
	out->linkGuessCode = map_link_guess(j_get_str(root, "linkGuess"));
	out->confidence = j_get_num(root, "confidence");
	out->dlGoodputMbps = j_get_num(root, "dlGoodputMbps");
	out->ulGoodputMbps = j_get_num(root, "ulGoodputMbps");
	out->bufferbloatInflationFactor = j_get_num(root, "bufferbloatInflationFactor");
	out->capDetectedMbps = j_get_num(root, "capDetectedMbps");

	out->nearTtfbMedianMs = j_get_num(root, "nearTtfbMedianMs");
	out->nearTtfbP95Ms = j_get_num(root, "nearTtfbP95Ms");
	out->nearTtfbP99Ms = j_get_num(root, "nearTtfbP99Ms");

	out->farTtfbMedianMs = j_get_num(root, "farTtfbMedianMs");
	out->farTtfbP95Ms = j_get_num(root, "farTtfbP95Ms");
	out->farTtfbP99Ms = j_get_num(root, "farTtfbP99Ms");

	out->probeSuccessRatePct = j_get_num(root, "probeSuccessRatePct");
	out->stallRatePct = j_get_num(root, "stallRatePct");

	out->lastMicroTs = j_get_num(root, "lastMicroTs");
	out->lastMultiTs = j_get_num(root, "lastMultiTs");
	out->lastChgTs = j_get_num(root, "lastChgTs");
	out->lastClassifyTs = j_get_num(root, "lastClassifyTs");
	out->lastDiagTs = j_get_num(root, "lastDiagTs");

	{
		const char *diag = j_get_str(root, "lastDiagName");
		out->diagPresent = (diag && diag[0]) ? 1.0 : 0.0;
	}

	ret = RETURN_OK;

	if (root) json_decref(root);
	return ret;
}

static int backhaul_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                               BackhaulMetrics *m,
                                               metricAddFunc addFunc) {
	int ret = RETURN_OK;

	for (int idx = 0; idx < cfgStat->kpiCount; idx++) {
		KPIConfig *kpi = &(cfgStat->kpi[idx]);
		double val = 0;

		if (!kpi || !kpi->fqname) continue;

		/* Match on fqname (same style you already use elsewhere) */
		if (str_icontains(kpi->fqname, "state")) {
			val = (double)m->stateCode;

		} else if (str_icontains(kpi->fqname, "link_guess")) {
			val = (double)m->linkGuessCode;

		} else if (str_icontains(kpi->fqname, "confidence")) {
			val = m->confidence;

		} else if (str_icontains(kpi->fqname, "dl_goodput_mbps")) {
			val = m->dlGoodputMbps;

		} else if (str_icontains(kpi->fqname, "ul_goodput_mbps")) {
			val = m->ulGoodputMbps;

		} else if (str_icontains(kpi->fqname, "bufferbloat_inflation_factor")) {
			val = m->bufferbloatInflationFactor;

		} else if (str_icontains(kpi->fqname, "cap_detected_mbps")) {
			val = m->capDetectedMbps;

		} else if (str_icontains(kpi->fqname, "near_ttfb_median_ms")) {
			val = m->nearTtfbMedianMs;

		} else if (str_icontains(kpi->fqname, "near_ttfb_p95_ms")) {
			val = m->nearTtfbP95Ms;

		} else if (str_icontains(kpi->fqname, "near_ttfb_p99_ms")) {
			val = m->nearTtfbP99Ms;

		} else if (str_icontains(kpi->fqname, "far_ttfb_median_ms")) {
			val = m->farTtfbMedianMs;

		} else if (str_icontains(kpi->fqname, "far_ttfb_p95_ms")) {
			val = m->farTtfbP95Ms;

		} else if (str_icontains(kpi->fqname, "far_ttfb_p99_ms")) {
			val = m->farTtfbP99Ms;

		} else if (str_icontains(kpi->fqname, "probe_success_rate_pct")) {
			val = m->probeSuccessRatePct;

		} else if (str_icontains(kpi->fqname, "stall_rate_pct")) {
			val = m->stallRatePct;

		} else if (str_icontains(kpi->fqname, "last_micro_ts")) {
			val = m->lastMicroTs;

		} else if (str_icontains(kpi->fqname, "last_multi_ts")) {
			val = m->lastMultiTs;

		} else if (str_icontains(kpi->fqname, "last_change_ts")) {
			val = m->lastChgTs;

		} else if (str_icontains(kpi->fqname, "last_classify_ts")) {
			val = m->lastClassifyTs;

		} else if (str_icontains(kpi->fqname, "last_diag_ts")) {
			val = m->lastDiagTs;

		} else if (str_icontains(kpi->fqname, "diag_present")) {
			val = m->diagPresent;

		} else {
			continue;
		}

		addFunc(kpi, &val);
	}

	return ret;
}

int backhaul_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {
	int ret = RETURN_OK;
	int port = 0;
	char urlBuf[256] = {0};
	const char *path = NULL;

	BackhaulMetrics *m = calloc(1, sizeof(BackhaulMetrics));
	if (!m) {
		usys_log_error("backhaul_agent: oom allocating metrics");
		return RETURN_NOTOK;
	}

	/* Resolve backhaul.d port from /etc/services */
	port = usys_find_service_port(SERVICE_BACKHAUL);
	if (port <= 0) {
		usys_log_error("backhaul_agent: could not resolve service port for '%s'",
		               SERVICE_BACKHAUL);
		free(m);
		return RETURN_NOTOK;
	}

	/* cfgStat->url is treated as endpoint path, default /v1/status */
	path = (cfgStat && cfgStat->url && cfgStat->url[0]) ? cfgStat->url : "/v1/status";

	/* Build full URL */
	snprintf(urlBuf, sizeof(urlBuf), "http://127.0.0.1:%d%s", port, path);

	if (backhaul_read_metrics(urlBuf, m) != RETURN_OK) {
		usys_log_error("backhaul_agent: failed to read %s", urlBuf);
		free(m);
		return RETURN_NOTOK;
	}

	if (backhaul_push_stat_to_metric_server(cfgStat, m, addFunc) != RETURN_OK) {
		usys_log_error("backhaul_agent: failed to push metrics to server");
		ret = RETURN_NOTOK;
	}

	free(m);
	return ret;
}

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
	double	poeTotalPowerWatts;
	double	poeMaxPowerWatts;
	double	systemTemperatureC;
	double	ambientTemperatureC;
	double	systemPowerWatts;
	double	inputVoltage;
	double	systemCurrentAmps;
	double	inputLinkFailureAlarm;
	double	inputPoeFailureAlarm;
} SwitchMetrics;

/* very small helper: case-insensitive substring match */
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

static size_t curl_write_cb(void *contents, size_t size, size_t nmemb,
                            void *userp) {
	size_t total = size * nmemb;
	HttpBuf *buf = (HttpBuf *)userp;

	char *newMem = realloc(buf->data, buf->len + total + 1);
	if (!newMem) return 0;

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
		usys_log_error("switchd_agent: curl_easy_init failed");
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
		usys_log_error("switchd_agent: curl failed url=%s err=%s",
		               url, curl_easy_strerror(cres));
		goto done;
	}

	*outRoot = json_loads(buf.data ? buf.data : "{}", 0, &jerr);
	if (!*outRoot) {
		usys_log_error("switchd_agent: json parse failed url=%s line=%d: %s",
		               url, jerr.line, jerr.text);
		goto done;
	}

	ret = RETURN_OK;

done:
	if (curl) curl_easy_cleanup(curl);
	if (buf.data) free(buf.data);
	return ret;
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

static int j_get_bool(json_t *root, const char *key, int defVal) {
	json_t *v = NULL;
	if (!root || !key) return defVal;

	v = json_object_get(root, key);
	if (!v) return defVal;
	if (!json_is_boolean(v)) return defVal;

	return json_is_true(v) ? 1 : 0;
}

static int switchd_read_metrics(const char *url, SwitchMetrics *out) {
	int ret = RETURN_NOTOK;
	json_t *root = NULL;

	if (!out) return RETURN_NOTOK;
	memset(out, 0, sizeof(*out));

	if (http_get_json(url, &root) != RETURN_OK) {
		return RETURN_NOTOK;
	}

	out->poeTotalPowerWatts  = j_get_num(root, "poeTotalPowerWatts");
	out->poeMaxPowerWatts    = j_get_num(root, "poeMaxPowerWatts");
	out->systemTemperatureC  = j_get_num(root, "systemTemperatureC");
	out->ambientTemperatureC = j_get_num(root, "ambientTemperatureC");
	out->systemPowerWatts    = j_get_num(root, "systemPowerWatts");
	out->inputVoltage        = j_get_num(root, "inputVoltage");
	out->systemCurrentAmps   = j_get_num(root, "systemCurrentAmps");

	out->inputLinkFailureAlarm =
		(double)j_get_bool(root, "inputLinkFailureAlarm", 0);
	out->inputPoeFailureAlarm =
		(double)j_get_bool(root, "inputPoeFailureAlarm", 0);

	ret = RETURN_OK;

	if (root) json_decref(root);
	return ret;
}

static int switchd_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                              SwitchMetrics *m,
                                              metricAddFunc addFunc) {
	for (int idx = 0; idx < cfgStat->kpiCount; idx++) {
		KPIConfig *kpi = &(cfgStat->kpi[idx]);
		double val = 0;

		if (!kpi || !kpi->fqname) continue;

		if (str_icontains(kpi->fqname, "poe_total_power_watts")) {
			val = m->poeTotalPowerWatts;

		} else if (str_icontains(kpi->fqname, "poe_max_power_watts")) {
			val = m->poeMaxPowerWatts;

		} else if (str_icontains(kpi->fqname, "system_temperature_c")) {
			val = m->systemTemperatureC;

		} else if (str_icontains(kpi->fqname, "ambient_temperature_c")) {
			val = m->ambientTemperatureC;

		} else if (str_icontains(kpi->fqname, "system_power_watts")) {
			val = m->systemPowerWatts;

		} else if (str_icontains(kpi->fqname, "input_voltage")) {
			val = m->inputVoltage;

		} else if (str_icontains(kpi->fqname, "system_current_amps")) {
			val = m->systemCurrentAmps;

		} else if (str_icontains(kpi->fqname, "input_link_failure_alarm")) {
			val = m->inputLinkFailureAlarm;

		} else if (str_icontains(kpi->fqname, "input_poe_failure_alarm")) {
			val = m->inputPoeFailureAlarm;

		} else {
			continue;
		}

		addFunc(kpi, &val);
	}

	return RETURN_OK;
}

int switchd_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {
	int port = 0;
	char urlBuf[256] = {0};
	const char *path = NULL;

	SwitchMetrics *m = calloc(1, sizeof(SwitchMetrics));
	if (!m) {
		usys_log_error("switchd_agent: oom allocating metrics");
		return RETURN_NOTOK;
	}

	port = usys_find_service_port(SERVICE_SWITCH);
	if (port <= 0) {
		usys_log_error("switchd_agent: could not resolve service port for '%s'",
		               SERVICE_SWITCH);
		free(m);
		return RETURN_NOTOK;
	}

	path = (cfgStat && cfgStat->url && cfgStat->url[0]) ?
	       cfgStat->url : "/v1/metrics/switch";

	snprintf(urlBuf, sizeof(urlBuf), "http://127.0.0.1:%d%s", port, path);

	if (switchd_read_metrics(urlBuf, m) != RETURN_OK) {
		usys_log_error("switchd_agent: failed to read %s", urlBuf);
		free(m);
		return RETURN_NOTOK;
	}

	if (switchd_push_stat_to_metric_server(cfgStat, m, addFunc) != RETURN_OK) {
		usys_log_error("switchd_agent: failed to push metrics to server");
		free(m);
		return RETURN_NOTOK;
	}

	free(m);
	return RETURN_OK;
}

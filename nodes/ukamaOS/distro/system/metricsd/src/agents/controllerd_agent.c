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
	double	solarPanelVoltage;
	double	solarPanelCurrent;
	double	solarPanelPower;
	double	solarYieldToday;
	double	solarYieldTotal;
	double	batteryVoltage;
	double	batteryCurrent;
	double	mpptEfficiency;
	double	batteryChargePercentage;
	double	controllerTemperature;
	double	loadCurrent;
	int	hasBatteryChargePercentage;
	int	hasControllerTemperature;
	int	hasLoadCurrent;
} ControllerMetrics;

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
		usys_log_error("controllerd_agent: curl_easy_init failed");
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
		usys_log_error("controllerd_agent: curl failed url=%s err=%s",
		               url, curl_easy_strerror(cres));
		goto done;
	}

	*outRoot = json_loads(buf.data ? buf.data : "{}", 0, &jerr);
	if (!*outRoot) {
		usys_log_error("controllerd_agent: json parse failed url=%s line=%d: %s",
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

static const char *j_get_str(json_t *root, const char *key) {
	json_t *v = NULL;
	if (!root || !key) return NULL;

	v = json_object_get(root, key);
	if (!v || !json_is_string(v)) return NULL;

	return json_string_value(v);
}

static json_t *j_get_arr(json_t *root, const char *key) {
	json_t *v = NULL;
	if (!root || !key) return NULL;

	v = json_object_get(root, key);
	if (!v || !json_is_array(v)) return NULL;

	return v;
}

static int controllerd_read_metrics(const char *url, ControllerMetrics *out) {
	int ret = RETURN_NOTOK;
	json_t *root = NULL;
	json_t *metrics = NULL;

	if (!out) return RETURN_NOTOK;
	memset(out, 0, sizeof(*out));

	if (http_get_json(url, &root) != RETURN_OK) {
		return RETURN_NOTOK;
	}

	metrics = j_get_arr(root, "metrics");
	if (!metrics) {
		usys_log_error("controllerd_agent: missing metrics array in %s", url);
		goto done;
	}

	for (size_t i = 0; i < json_array_size(metrics); i++) {
		json_t *metric = json_array_get(metrics, i);
		const char *name = NULL;
		double value = 0;

		if (!metric || !json_is_object(metric)) continue;

		name = j_get_str(metric, "name");
		if (!name) continue;

		value = j_get_num(metric, "value");

		if (str_ieq(name, "solar_panel_voltage")) {
			out->solarPanelVoltage = value;
		} else if (str_ieq(name, "solar_panel_current")) {
			out->solarPanelCurrent = value;
		} else if (str_ieq(name, "solar_panel_power")) {
			out->solarPanelPower = value;
		} else if (str_ieq(name, "solar_yield_today")) {
			out->solarYieldToday = value;
		} else if (str_ieq(name, "solar_yield_total")) {
			out->solarYieldTotal = value;
		} else if (str_ieq(name, "battery_voltage")) {
			out->batteryVoltage = value;
		} else if (str_ieq(name, "battery_current")) {
			out->batteryCurrent = value;
		} else if (str_ieq(name, "mppt_efficiency")) {
			out->mpptEfficiency = value;
		} else if (str_ieq(name, "battery_charge_percentage")) {
			out->batteryChargePercentage = value;
			out->hasBatteryChargePercentage = 1;
		} else if (str_ieq(name, "controller_temperature")) {
			out->controllerTemperature = value;
			out->hasControllerTemperature = 1;
		} else if (str_ieq(name, "load_current")) {
			out->loadCurrent = value;
			out->hasLoadCurrent = 1;
		}
	}

	ret = RETURN_OK;

done:
	if (root) json_decref(root);
	return ret;
}

static int controllerd_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                                  ControllerMetrics *m,
                                                  metricAddFunc addFunc) {
	for (int idx = 0; idx < cfgStat->kpiCount; idx++) {
		KPIConfig *kpi = &(cfgStat->kpi[idx]);
		double val = 0;

		if (!kpi || !kpi->fqname) continue;

		if (str_icontains(kpi->fqname, "solar_panel_voltage")) {
			val = m->solarPanelVoltage;
		} else if (str_icontains(kpi->fqname, "solar_panel_current")) {
			val = m->solarPanelCurrent;
		} else if (str_icontains(kpi->fqname, "solar_panel_power")) {
			val = m->solarPanelPower;
		} else if (str_icontains(kpi->fqname, "solar_yield_today")) {
			val = m->solarYieldToday;
		} else if (str_icontains(kpi->fqname, "solar_yield_total")) {
			val = m->solarYieldTotal;
		} else if (str_icontains(kpi->fqname, "battery_voltage")) {
			val = m->batteryVoltage;
		} else if (str_icontains(kpi->fqname, "battery_current")) {
			val = m->batteryCurrent;
		} else if (str_icontains(kpi->fqname, "mppt_efficiency")) {
			val = m->mpptEfficiency;
		} else if (str_icontains(kpi->fqname, "battery_charge_percentage")) {
			if (!m->hasBatteryChargePercentage) continue;
			val = m->batteryChargePercentage;
		} else if (str_icontains(kpi->fqname, "controller_temperature")) {
			if (!m->hasControllerTemperature) continue;
			val = m->controllerTemperature;
		} else if (str_icontains(kpi->fqname, "load_current")) {
			if (!m->hasLoadCurrent) continue;
			val = m->loadCurrent;
		} else {
			continue;
		}

		addFunc(kpi, &val);
	}

	return RETURN_OK;
}

int controllerd_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {
	int port = 0;
	char urlBuf[256] = {0};
	const char *path = NULL;

	ControllerMetrics *m = calloc(1, sizeof(ControllerMetrics));
	if (!m) {
		usys_log_error("controllerd_agent: oom allocating metrics");
		return RETURN_NOTOK;
	}

	port = usys_find_service_port(SERVICE_CONTROLLER);
	if (port <= 0) {
		usys_log_error("controllerd_agent: could not resolve service port for '%s'",
		               SERVICE_CONTROLLER);
		free(m);
		return RETURN_NOTOK;
	}

	path = (cfgStat && cfgStat->url && cfgStat->url[0]) ?
	       cfgStat->url : "/v1/controller/metrics";

	snprintf(urlBuf, sizeof(urlBuf), "http://127.0.0.1:%d%s", port, path);

	if (controllerd_read_metrics(urlBuf, m) != RETURN_OK) {
		usys_log_error("controllerd_agent: failed to read %s", urlBuf);
		free(m);
		return RETURN_NOTOK;
	}

	if (controllerd_push_stat_to_metric_server(cfgStat, m, addFunc) != RETURN_OK) {
		usys_log_error("controllerd_agent: failed to push metrics to server");
		free(m);
		return RETURN_NOTOK;
	}

	free(m);
	return RETURN_OK;
}

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
	int     present;             /* 1 if we found this fem in JSON */
	int     ok;                  /* fem.ok */
	double  temperatureC;        /* fem.temperature.temperature */
	double  forwardPower;        /* fem.adc.forward_power */
	double  reversePower;        /* fem.adc.reverse_power */
	double  paCurrent;           /* fem.adc.pa_current */
	double  gpio28vVdsEnable;    /* fem.gpio.28v_vds_enable */
	double  gpioTxRfEnable;      /* fem.gpio.tx_rf_enable */
	double  gpioRxRfEnable;      /* fem.gpio.rx_rf_enable */
	double  gpioPaVdsEnable;     /* fem.gpio.pa_vds_enable */
	double  gpioRfPalEnable;     /* fem.gpio.rf_pal_enable */
	double  gpioPsuPgood;        /* fem.gpio.psu_pgood */
	double  snapshotTs;          /* fem.snapshot.timestamp */
} FemUnitMetrics;

typedef struct {
	int controllerOk;            /* controller.ok */
	FemUnitMetrics fem[2];       /* fem_unit 1 -> [0], fem_unit 2 -> [1] */
} FemMetrics;

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
		usys_log_error("femd_agent: curl_easy_init failed");
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
		usys_log_error("femd_agent: curl failed url=%s err=%s",
		               url, curl_easy_strerror(cres));
		goto done;
	}

	*outRoot = json_loads(buf.data ? buf.data : "{}", 0, &jerr);
	if (!*outRoot) {
		usys_log_error("femd_agent: json parse failed url=%s line=%d: %s",
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

static json_t *j_get_obj(json_t *root, const char *key) {
	json_t *v = NULL;
	if (!root || !key) return NULL;
	v = json_object_get(root, key);
	if (!v || !json_is_object(v)) return NULL;
	return v;
}

static json_t *j_get_arr(json_t *root, const char *key) {
	json_t *v = NULL;
	if (!root || !key) return NULL;
	v = json_object_get(root, key);
	if (!v || !json_is_array(v)) return NULL;
	return v;
}

static int femd_read_metrics(const char *url, FemMetrics *out) {
	int ret = RETURN_NOTOK;
	json_t *root = NULL;
	json_t *controller = NULL;
	json_t *fems = NULL;

	if (!out) return RETURN_NOTOK;
	memset(out, 0, sizeof(*out));

	if (http_get_json(url, &root) != RETURN_OK) {
		return RETURN_NOTOK;
	}

	controller = j_get_obj(root, "controller");
	out->controllerOk = controller ? j_get_bool(controller, "ok", 0) : 0;

	fems = j_get_arr(root, "fems");
	if (fems) {
		size_t i = 0;
		size_t n = json_array_size(fems);

		for (i = 0; i < n; i++) {
			json_t *fem = json_array_get(fems, i);
			json_t *gpio = NULL;
			json_t *temp = NULL;
			json_t *adc  = NULL;
			json_t *snap = NULL;
			int femUnit = 0;
			int idx = -1;

			if (!fem || !json_is_object(fem)) continue;

			femUnit = (int)j_get_num(fem, "fem_unit");
			if (femUnit == 1) idx = 0;
			else if (femUnit == 2) idx = 1;
			else continue;

			out->fem[idx].present = 1;
			out->fem[idx].ok = j_get_bool(fem, "ok", 0);

			/* temperature.temperature */
			temp = j_get_obj(fem, "temperature");
			if (temp) out->fem[idx].temperatureC = j_get_num(temp, "temperature");

			/* adc fields */
			adc = j_get_obj(fem, "adc");
			if (adc) {
				out->fem[idx].forwardPower = j_get_num(adc, "forward_power");
				out->fem[idx].reversePower = j_get_num(adc, "reverse_power");
				out->fem[idx].paCurrent    = j_get_num(adc, "pa_current");
			}

			/* gpio fields */
			gpio = j_get_obj(fem, "gpio");
			if (gpio) {
				out->fem[idx].gpio28vVdsEnable = (double)j_get_bool(gpio, "28v_vds_enable", 0);
				out->fem[idx].gpioTxRfEnable   = (double)j_get_bool(gpio, "tx_rf_enable", 0);
				out->fem[idx].gpioRxRfEnable   = (double)j_get_bool(gpio, "rx_rf_enable", 0);
				out->fem[idx].gpioPaVdsEnable  = (double)j_get_bool(gpio, "pa_vds_enable", 0);
				out->fem[idx].gpioRfPalEnable  = (double)j_get_bool(gpio, "rf_pal_enable", 0);
				out->fem[idx].gpioPsuPgood     = (double)j_get_bool(gpio, "psu_pgood", 0);
			}

			/* snapshot.timestamp (if present) */
			snap = j_get_obj(fem, "snapshot");
			if (snap) out->fem[idx].snapshotTs = j_get_num(snap, "timestamp");
		}
	}

	ret = RETURN_OK;

	if (root) json_decref(root);
	return ret;
}

static FemUnitMetrics *pick_fem_by_kpi(KPIConfig *kpi, FemMetrics *m) {
	int isFem1 = 0;
	int isFem2 = 0;

	if (!kpi || !kpi->fqname || !m) return NULL;

	isFem1 = str_icontains(kpi->fqname, "fem1");
	isFem2 = str_icontains(kpi->fqname, "fem2");

	if (isFem1) return &m->fem[0];
	if (isFem2) return &m->fem[1];

	return NULL;
}

static int femd_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                          FemMetrics *m,
                                          metricAddFunc addFunc) {
	for (int idx = 0; idx < cfgStat->kpiCount; idx++) {
		KPIConfig *kpi = &(cfgStat->kpi[idx]);
		double val = 0;
		FemUnitMetrics *fu = NULL;

		if (!kpi || !kpi->fqname) continue;

		/* Controller metric */
		if (str_icontains(kpi->fqname, "controller_ok")) {
			val = (double)(m->controllerOk ? 1 : 0);
			addFunc(kpi, &val);
			continue;
		}

		/* Everything else is per FEM (fem1/fem2) */
		fu = pick_fem_by_kpi(kpi, m);
		if (!fu || !fu->present) {
			/* If FEM is missing in payload, skip silently */
			continue;
		}

		if (str_icontains(kpi->fqname, "temperature_c")) {
			val = fu->temperatureC;

		} else if (str_icontains(kpi->fqname, "forward_power")) {
			val = fu->forwardPower;

		} else if (str_icontains(kpi->fqname, "reverse_power")) {
			val = fu->reversePower;

		} else if (str_icontains(kpi->fqname, "pa_current")) {
			val = fu->paCurrent;

		} else if (str_icontains(kpi->fqname, "gpio_28v_vds_enable")) {
			val = fu->gpio28vVdsEnable;

		} else if (str_icontains(kpi->fqname, "gpio_tx_rf_enable")) {
			val = fu->gpioTxRfEnable;

		} else if (str_icontains(kpi->fqname, "gpio_rx_rf_enable")) {
			val = fu->gpioRxRfEnable;

		} else if (str_icontains(kpi->fqname, "gpio_pa_vds_enable")) {
			val = fu->gpioPaVdsEnable;

		} else if (str_icontains(kpi->fqname, "gpio_rf_pal_enable")) {
			val = fu->gpioRfPalEnable;

		} else if (str_icontains(kpi->fqname, "gpio_psu_pgood")) {
			val = fu->gpioPsuPgood;

		} else if (str_icontains(kpi->fqname, "snapshot_timestamp")) {
			val = fu->snapshotTs;

		} else if (str_icontains(kpi->fqname, "ok")) {
			/* keep this last so it doesn't steal matches like "controller_ok" */
			val = (double)(fu->ok ? 1 : 0);

		} else {
			continue;
		}

		addFunc(kpi, &val);
	}

	return RETURN_OK;
}

int femd_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {
	int port = 0;
	char urlBuf[256] = {0};
	const char *path = NULL;

	FemMetrics *m = calloc(1, sizeof(FemMetrics));
	if (!m) {
		usys_log_error("femd_agent: oom allocating metrics");
		return RETURN_NOTOK;
	}

	/* Resolve fem.d port from /etc/services */
	port = usys_find_service_port(SERVICE_FEM);
	if (port <= 0) {
		usys_log_error("femd_agent: could not resolve service port for '%s'",
		               SERVICE_FEM);
		free(m);
		return RETURN_NOTOK;
	}

	/* cfgStat->url is treated as endpoint path, default /v1/status */
	path = (cfgStat && cfgStat->url && cfgStat->url[0]) ? cfgStat->url : "/v1/status";

	/* Build full URL */
	snprintf(urlBuf, sizeof(urlBuf), "http://127.0.0.1:%d%s", port, path);

	if (femd_read_metrics(urlBuf, m) != RETURN_OK) {
		usys_log_error("femd_agent: failed to read %s", urlBuf);
		free(m);
		return RETURN_NOTOK;
	}

	if (femd_push_stat_to_metric_server(cfgStat, m, addFunc) != RETURN_OK) {
		usys_log_error("femd_agent: failed to push metrics to server");
		free(m);
		return RETURN_NOTOK;
	}

	free(m);
	return RETURN_OK;
}

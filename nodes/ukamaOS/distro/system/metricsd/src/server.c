/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "microhttpd.h"
#include "prom.h"
#include "promhttp.h"

#include "usys_log.h"

static struct MHD_Daemon *mhdDaemon = NULL;
static char gMetricNodePrefix[64] = {'\0'};

static int append_buf(char **buf, size_t *len, const char *text) {

    char *next = NULL;
    size_t add = 0;

    if (text == NULL) {
        return RETURN_OK;
    }

    add = strlen(text);
    next = realloc(*buf, *len + add + 1);
    if (next == NULL) {
        return RETURN_NOTOK;
    }

    *buf = next;
    memcpy(*buf + *len, text, add);
    *len += add;
    (*buf)[*len] = '\0';

    return RETURN_OK;
}

static int append_nbuf(char **buf, size_t *len, const char *text, size_t n) {

    char *next = NULL;

    if (text == NULL || n == 0) {
        return RETURN_OK;
    }

    next = realloc(*buf, *len + n + 1);
    if (next == NULL) {
        return RETURN_NOTOK;
    }

    *buf = next;
    memcpy(*buf + *len, text, n);
    *len += n;
    (*buf)[*len] = '\0';

    return RETURN_OK;
}

static void metric_server_capture_node_prefix(const char *fqname) {

    const char *sep = NULL;
    size_t len = 0;

    if (gMetricNodePrefix[0] != '\0' ||
        fqname == NULL ||
        fqname[0] == '\0') {
        return;
    }

    sep = strchr(fqname, '_');
    if (sep == NULL) {
        return;
    }

    len = (size_t)(sep - fqname);
    if (len == 0 || len >= sizeof(gMetricNodePrefix)) {
        return;
    }

    memcpy(gMetricNodePrefix, fqname, len);
    gMetricNodePrefix[len] = '\0';
}

static int metric_name_has_node_prefix(const char *name) {

    size_t plen = 0;

    if (name == NULL || gMetricNodePrefix[0] == '\0') {
        return 0;
    }

    plen = strlen(gMetricNodePrefix);

    if (strncmp(name, gMetricNodePrefix, plen) != 0) {
        return 0;
    }

    return name[plen] == '_';
}

static int append_prefixed_metric_line(char **buf, size_t *len, const char *line) {

    const char *name = NULL;
    const char *rest = NULL;
    char metricName[256] = {'\0'};
    size_t nameLen = 0;

    if (line == NULL) {
        return RETURN_OK;
    }

    if (gMetricNodePrefix[0] == '\0') {
        if (append_buf(buf, len, line) != RETURN_OK) return RETURN_NOTOK;
        if (append_buf(buf, len, "\n") != RETURN_OK) return RETURN_NOTOK;
        return RETURN_OK;
    }

    if (strncmp(line, "# HELP ", 7) == 0 || strncmp(line, "# TYPE ", 7) == 0) {
        name = line + 7;
        rest = strchr(name, ' ');
        if (rest == NULL) {
            if (append_buf(buf, len, line) != RETURN_OK) return RETURN_NOTOK;
            if (append_buf(buf, len, "\n") != RETURN_OK) return RETURN_NOTOK;
            return RETURN_OK;
        }

        nameLen = (size_t)(rest - name);
        if (nameLen >= sizeof(metricName)) {
            nameLen = sizeof(metricName) - 1;
        }

        memcpy(metricName, name, nameLen);
        metricName[nameLen] = '\0';

        if (metric_name_has_node_prefix(metricName)) {
            if (append_buf(buf, len, line) != RETURN_OK) return RETURN_NOTOK;
            if (append_buf(buf, len, "\n") != RETURN_OK) return RETURN_NOTOK;
            return RETURN_OK;
        }

        if (append_nbuf(buf, len, line, 7) != RETURN_OK)          return RETURN_NOTOK;
        if (append_buf(buf, len, gMetricNodePrefix) != RETURN_OK) return RETURN_NOTOK;
        if (append_buf(buf, len, "_") != RETURN_OK)               return RETURN_NOTOK;
        if (append_nbuf(buf, len, name, nameLen) != RETURN_OK)    return RETURN_NOTOK;
        if (append_buf(buf, len, rest) != RETURN_OK)              return RETURN_NOTOK;
        if (append_buf(buf, len, "\n") != RETURN_OK)              return RETURN_NOTOK;
        return RETURN_OK;
    }

    if (line[0] == '#'
        || line[0] == '\0') {
        if (append_buf(buf, len, line) != RETURN_OK) return RETURN_NOTOK;
        if (append_buf(buf, len, "\n") != RETURN_OK) return RETURN_NOTOK;
        return RETURN_OK;
    }

    name = line;
    rest = name;

    while (*rest != '\0' && *rest != ' ' && *rest != '\t' && *rest != '{') {
        rest++;
    }

    nameLen = (size_t)(rest - name);
    if (nameLen == 0) {
        if (append_buf(buf, len, line) != RETURN_OK) return RETURN_NOTOK;
        if (append_buf(buf, len, "\n") != RETURN_OK) return RETURN_NOTOK;
        return RETURN_OK;
    }

    if (nameLen >= sizeof(metricName)) {
        nameLen = sizeof(metricName) - 1;
    }

    memcpy(metricName, name, nameLen);
    metricName[nameLen] = '\0';

    if (metric_name_has_node_prefix(metricName)) {
        if (append_buf(buf, len, line) != RETURN_OK) return RETURN_NOTOK;
        if (append_buf(buf, len, "\n") != RETURN_OK) return RETURN_NOTOK;
        return RETURN_OK;
    }

    if (append_buf(buf, len, gMetricNodePrefix) != RETURN_OK) return RETURN_NOTOK;
    if (append_buf(buf, len, "_") != RETURN_OK)               return RETURN_NOTOK;
    if (append_nbuf(buf, len, name, nameLen) != RETURN_OK)    return RETURN_NOTOK;
    if (append_buf(buf, len, rest) != RETURN_OK)              return RETURN_NOTOK;
    if (append_buf(buf, len, "\n") != RETURN_OK)              return RETURN_NOTOK;

    return RETURN_OK;
}

static char *metric_server_build_prefixed_metrics(void) {

    const char *raw = NULL;
    const char *cur = NULL;
    char *buf = NULL;
    size_t len = 0;

    raw = prom_collector_registry_bridge(PROM_COLLECTOR_REGISTRY_DEFAULT);
    if (raw == NULL) {
        return NULL;
    }

    cur = raw;
    while (*cur != '\0') {
        const char *end = strchr(cur, '\n');
        char line[4096] = {'\0'};
        size_t lineLen = 0;

        if (end == NULL) {
            lineLen = strlen(cur);
        } else {
            lineLen = (size_t)(end - cur);
        }

        if (lineLen >= sizeof(line)) {
            lineLen = sizeof(line) - 1;
        }

        memcpy(line, cur, lineLen);
        line[lineLen] = '\0';

        if (append_prefixed_metric_line(&buf, &len, line) != RETURN_OK) {
            free(buf);
            return NULL;
        }

        if (end == NULL) {
            break;
        }

        cur = end + 1;
    }

    return buf;
}

static enum MHD_Result metric_server_metrics_handler(void *cls,
                                                     struct MHD_Connection *connection,
                                                     const char *url,
                                                     const char *method,
                                                     const char *version,
                                                     const char *upload_data,
                                                     size_t *upload_data_size,
                                                     void **con_cls) {

    struct MHD_Response *response = NULL;
    char *body = NULL;
    int status = MHD_HTTP_OK;
    enum MHD_Result ret = MHD_NO;

    (void)cls;
    (void)version;
    (void)upload_data;
    (void)upload_data_size;
    (void)con_cls;

    if (strcmp(method, "GET") != 0) {
        return MHD_NO;
    }

    if (strcmp(url, "/") != 0 && strcmp(url, "/metrics") != 0) {
        return MHD_NO;
    }

    body = metric_server_build_prefixed_metrics();
    if (body == NULL) {
        body = strdup("# failed to build metrics\n");
        if (body == NULL) {
            return MHD_NO;
        }
        status = MHD_HTTP_INTERNAL_SERVER_ERROR;
    }

    response = MHD_create_response_from_buffer(strlen(body), (void *)body,
                                               MHD_RESPMEM_MUST_COPY);
    if (response == NULL) {
        free(body);
        return MHD_NO;
    }

    MHD_add_response_header(response, "Content-Type",
                            "text/plain; version=0.0.4; charset=utf-8");
    ret = MHD_queue_response(connection, status, response);
    MHD_destroy_response(response);
    free(body);

    return ret;
}

void metric_server_free_kpi(KPIConfig *kpi) {

    if ((kpi == NULL) || (kpi->registry == NULL)) {
        return;
    }

    /*
     * Registered Prometheus metrics are owned by the collector registry.
     * They must be destroyed only by prom_collector_registry_destroy().
     *
     * We only free our wrapper object here.
     */
    free(kpi->registry);
    kpi->registry = NULL;
}


int metric_server_register_kpi(KPIConfig *kpi) {

    usys_log_trace("registering prometheus metric for %s", kpi->fqname);

    metric_server_capture_node_prefix(kpi->fqname);

    kpi->registry = calloc(1, sizeof(PromRegistry));
    if (kpi->registry == NULL) {
        return RETURN_NOTOK;
    }

    switch (kpi->type) {
    case METRICTYPE_COUNTER:
        kpi->registry->counter =
            prom_collector_registry_must_register_metric(
                prom_counter_new(kpi->fqname, kpi->desc,
                                 kpi->numLabels,
                                 (const char **)kpi->labels));
        if (kpi->registry->counter == NULL) {
            free(kpi->registry);
            kpi->registry = NULL;
            return RETURN_NOTOK;
        }
        break;

    case METRICTYPE_GAUGE:
        kpi->registry->gauge =
            prom_collector_registry_must_register_metric(
                prom_gauge_new(kpi->fqname, kpi->desc, kpi->numLabels,
                               (const char **)kpi->labels));
        if (kpi->registry->gauge == NULL) {
            free(kpi->registry);
            kpi->registry = NULL;
            return RETURN_NOTOK;
        }
        break;

    case METRICTYPE_HISTOGRAM:
        kpi->registry->histogram =
            prom_collector_registry_must_register_metric(
                prom_histogram_new(kpi->fqname, kpi->desc, kpi->buckets,
                                   kpi->numLabels,
                                   (const char **)kpi->labels));
        if (kpi->registry->histogram == NULL) {
            free(kpi->registry);
            kpi->registry = NULL;
            return RETURN_NOTOK;
        }
        break;

    default:
        usys_log_error("invalid kpi type %d for %s", kpi->type, kpi->name);
        free(kpi->registry);
        kpi->registry = NULL;
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

int metric_server_add_kpi_data(KPIConfig *kpi, void *value) {

    int ret = RETURN_OK;

    usys_log_trace("adding kpi %s type %d", kpi->name, kpi->type);

    switch (kpi->type) {
    case METRICTYPE_COUNTER:
        ret = prom_counter_inc(kpi->registry->counter,
                               (const char **)kpi->labels);
        break;

    case METRICTYPE_GAUGE:
        ret = prom_gauge_set(kpi->registry->gauge, *(double *)value,
                             (const char **)kpi->labels);
        break;

    case METRICTYPE_HISTOGRAM:
        ret = prom_histogram_observe(kpi->registry->histogram,
                                     *(double *)value,
                                     (const char **)kpi->labels);
        break;

    default:
        ret = RETURN_NOTOK;
        usys_log_error("invalid kpi type %d for %s", kpi->type, kpi->name);
        break;
    }

    return ret;
}

void metric_server_registry_init(void) {

    prom_collector_registry_default_init();
}

void metric_server_set_active_registry(void) {

    /*
     * We expose metrics through our own MHD handler now, not promhttp.
     * Keep this as a no-op so the rest of the code path stays unchanged.
     */
}

void metric_server_registry_destroy(void) {

    prom_collector_registry_destroy(PROM_COLLECTOR_REGISTRY_DEFAULT);
}

int metric_server_start(int port) {

    mhdDaemon = MHD_start_daemon(MHD_USE_SELECT_INTERNALLY,
                                 (unsigned short)port,
                                 NULL, NULL,
                                 &metric_server_metrics_handler, NULL,
                                 MHD_OPTION_END);
    if (mhdDaemon == NULL) {
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

void metric_server_stop(void) {

    if (mhdDaemon != NULL) {
        MHD_stop_daemon(mhdDaemon);
        mhdDaemon = NULL;
    }
}

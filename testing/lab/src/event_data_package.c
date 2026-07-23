#define _GNU_SOURCE

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "event.h"
#include "log.h"
#include "selector.h"
#include "util.h"

typedef struct {
    bff_client_t client;
    ue_t ue;
    package_t *pkg;
    subscriber_t *subscriber;
    double amount;
    char currency[ULAB_MAX_REF];
    bff_payment_t payment;
    ulab_error_t err;
    int rc;
} purchase_thread_t;

static package_t *event_package_for_ue(event_ctx_t *ctx,
                                       const char *ref,
                                       const ue_t *ue,
                                       ulab_error_t *err) {
    package_t *pkg;

    pkg = world_package_for_network(ctx->world, ref, ue->network_ref);
    if (pkg == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "unknown package %.128s for UE %.128s", ref, ue->ref);
    }
    return pkg;
}

static int event_allocate_sim(event_ctx_t *ctx,
                              const event_spec_t *event,
                              ulab_error_t *err) {
    selector_result_t ues;
    size_t i;

    memset(&ues, 0, sizeof(ues));
    if (selector_resolve_ues(ctx->world, &event->ues, &ues, err)) {
        return ULAB_ERR;
    }
    for (i = 0; i < ues.count; i++) {
        ue_t *ue;
        subscriber_t *subscriber;
        network_t *network;
        package_t *pkg;

        ue = &ctx->world->ues[ues.idx[i]];
        subscriber = world_subscriber_by_ref(ctx->world,
                                             ue->subscriber_ref);
        network = world_network_by_ref(ctx->world, ue->network_ref);
        pkg = event_package_for_ue(ctx, event->package_ref, ue, err);
        if (subscriber == NULL || network == NULL || pkg == NULL) {
            selector_result_free(&ues);
            if (err->msg[0] == '\0') {
                snprintf(err->msg, sizeof(err->msg),
                         "allocate_sim cannot resolve scenario resources");
            }
            return ULAB_ERR;
        }
        if (bff_allocate_sim(ctx->bff, ue, subscriber, network, pkg,
                             ctx->sim_type, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
    }
    selector_result_free(&ues);
    return ULAB_OK;
}

static int event_create_invalid_package(event_ctx_t *ctx,
                                        const event_spec_t *event,
                                        ulab_error_t *err) {
    package_t *source;
    network_t *network;
    char created_id[ULAB_MAX_ID];
    int rc;

    source = world_package_by_base_ref(ctx->world, event->package_ref);
    if (source == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "unknown package %.128s", event->package_ref);
        return ULAB_ERR;
    }
    network = source->network_ref[0] != '\0' ?
        world_network_by_ref(ctx->world, source->network_ref) :
        (ctx->world->network_count > 0 ? &ctx->world->networks[0] : NULL);
    memset(created_id, 0, sizeof(created_id));
    rc = bff_add_invalid_package(ctx->bff, source, network, event->variant,
                                 created_id, sizeof(created_id), err);
    if (rc == ULAB_OK && created_id[0] != '\0') {
        package_t *grown;
        package_t *created;
        package_t source_copy;

        source_copy = *source;
        grown = realloc(ctx->world->packages,
                        (ctx->world->package_count + 1) * sizeof(package_t));
        if (grown == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "accepted invalid package could not be tracked");
            return ULAB_ERR;
        }
        ctx->world->packages = grown;
        created = &ctx->world->packages[ctx->world->package_count++];
        *created = source_copy;
        snprintf(created->ref, sizeof(created->ref),
                 "%.80s-invalid-%.32s", source_copy.ref, event->variant);
        snprintf(created->base_ref, sizeof(created->base_ref),
                 "%.80s-invalid-%.32s", source_copy.base_ref,
                 event->variant);
        snprintf(created->name, sizeof(created->name),
                 "%.180s invalid %.48s", source_copy.name, event->variant);
        ulab_copy(created->bff_id, sizeof(created->bff_id), created_id);
    }
    return rc;
}

static int parse_iso_utc(const char *value, time_t *out) {
    struct tm parsed;
    int year;
    int month;
    int day;
    int hour;
    int minute;
    int second;

    if (value == NULL || out == NULL) {
        return ULAB_ERR;
    }
    year = 0;
    month = 0;
    day = 0;
    hour = 0;
    minute = 0;
    second = 0;
    if (sscanf(value, "%d-%d-%dT%d:%d:%d", &year, &month, &day,
               &hour, &minute, &second) != 6) {
        return ULAB_ERR;
    }
    memset(&parsed, 0, sizeof(parsed));
    parsed.tm_year = year - 1900;
    parsed.tm_mon = month - 1;
    parsed.tm_mday = day;
    parsed.tm_hour = hour;
    parsed.tm_min = minute;
    parsed.tm_sec = second;
    parsed.tm_isdst = 0;
    *out = timegm(&parsed);
    return *out == (time_t)-1 ? ULAB_ERR : ULAB_OK;
}

static int event_wait_package_boundary(event_ctx_t *ctx,
                                       const event_spec_t *event,
                                       ulab_error_t *err) {
    selector_result_t ues;
    size_t i;

    memset(&ues, 0, sizeof(ues));
    if (selector_resolve_ues(ctx->world, &event->ues, &ues, err)) {
        return ULAB_ERR;
    }
    for (i = 0; i < ues.count; i++) {
        ue_t *ue;
        package_t *pkg;
        bff_sim_package_t assignments[ULAB_MAX_BFF_SIM_PACKAGES];
        size_t assignment_count;
        size_t j;
        const char *end_date;
        time_t boundary;
        time_t target;
        time_t now;

        ue = &ctx->world->ues[ues.idx[i]];
        pkg = event_package_for_ue(ctx, event->package_ref, ue, err);
        if (pkg == NULL) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        if (bff_get_sim_packages(ctx->bff, ue, assignments,
                                 ULAB_MAX_BFF_SIM_PACKAGES,
                                 &assignment_count, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        end_date = NULL;
        for (j = 0; j < assignment_count; j++) {
            if (ulab_streq(assignments[j].package_id, pkg->bff_id)) {
                end_date = assignments[j].end_date;
                break;
            }
        }
        if (end_date == NULL || parse_iso_utc(end_date, &boundary)) {
            snprintf(err->msg, sizeof(err->msg),
                     "package %.128s has no parseable end date", pkg->ref);
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        target = boundary + (time_t)event->offset_seconds;
        now = time(NULL);
        if (event->offset_seconds < 0 && target <= now) {
            snprintf(err->msg, sizeof(err->msg),
                     "package %.128s pre-boundary target already passed",
                     pkg->ref);
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        while (now < target) {
            time_t remaining;
            unsigned int chunk;

            remaining = target - now;
            chunk = remaining > 60 ? 60u : (unsigned int)remaining;
            ulab_status("WAIT", "package=%s boundary=%s remaining=%lds",
                        pkg->ref, end_date, (long)remaining);
            sleep(chunk);
            now = time(NULL);
        }
    }
    selector_result_free(&ues);
    return ULAB_OK;
}

static void *purchase_thread(void *arg) {
    purchase_thread_t *thread;

    thread = arg;
    memset(&thread->payment, 0, sizeof(thread->payment));
    memset(&thread->err, 0, sizeof(thread->err));
    thread->rc = bff_record_cash_package_sale(&thread->client,
                                              &thread->ue,
                                              thread->pkg,
                                              thread->subscriber,
                                              thread->amount,
                                              thread->currency,
                                              &thread->payment,
                                              &thread->err);
    return NULL;
}

static int event_purchase_parallel(event_ctx_t *ctx,
                                   const event_spec_t *event,
                                   ulab_error_t *err) {
    selector_result_t ues;
    size_t i;

    memset(&ues, 0, sizeof(ues));
    if (selector_resolve_ues(ctx->world, &event->ues, &ues, err)) {
        return ULAB_ERR;
    }
    for (i = 0; i < ues.count; i++) {
        ue_t *ue;
        subscriber_t *subscriber;
        package_t *packages[2];
        purchase_thread_t work[2];
        pthread_t threads[2];
        size_t j;
        size_t started;

        ue = &ctx->world->ues[ues.idx[i]];
        subscriber = world_subscriber_by_ref(ctx->world,
                                             ue->subscriber_ref);
        packages[0] = event_package_for_ue(ctx, event->package_ref, ue, err);
        packages[1] = event_package_for_ue(ctx, event->other_package_ref,
                                           ue, err);
        if (subscriber == NULL || packages[0] == NULL || packages[1] == NULL) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        memset(work, 0, sizeof(work));
        started = 0;
        for (j = 0; j < 2; j++) {
            work[j].client = *ctx->bff;
            work[j].client.logf = NULL;
            work[j].ue = *ue;
            work[j].pkg = packages[j];
            work[j].subscriber = subscriber;
            work[j].amount = 0;
            if (pthread_create(&threads[j], NULL, purchase_thread,
                               &work[j]) != 0) {
                snprintf(err->msg, sizeof(err->msg),
                         "failed to start parallel purchase thread");
                break;
            }
            started++;
        }
        for (j = 0; j < started; j++) {
            pthread_join(threads[j], NULL);
        }
        if (started != 2) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        for (j = 0; j < 2; j++) {
            if (work[j].rc != ULAB_OK) {
                *err = work[j].err;
                selector_result_free(&ues);
                return ULAB_ERR;
            }
            ulab_status("PAYMENT", "parallel id=%s sim=%s package=%s",
                        work[j].payment.id, ue->ref, packages[j]->ref);
        }
    }
    selector_result_free(&ues);
    return ULAB_OK;
}

int event_data_package(event_ctx_t *ctx, const event_spec_t *event,
                       ulab_error_t *err) {
    switch (event->type) {
    case EVT_ALLOCATE_SIM:
        return event_allocate_sim(ctx, event, err);
    case EVT_CREATE_INVALID_PACKAGE:
        return event_create_invalid_package(ctx, event, err);
    case EVT_WAIT_PACKAGE_BOUNDARY:
        return event_wait_package_boundary(ctx, event, err);
    case EVT_PURCHASE_PACKAGES_PARALLEL:
        return event_purchase_parallel(ctx, event, err);
    default:
        snprintf(err->msg, sizeof(err->msg),
                 "not a data-package event");
        return ULAB_ERR;
    }
}

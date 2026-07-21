/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_BFF_H_
#define ULAB_BFF_H_

#include <stdint.h>
#include <stdio.h>

#include "world.h"
#include "ulab.h"

typedef struct {
    char url[ULAB_MAX_URL];
    char pauth_url[ULAB_MAX_URL];
    char bff_base_url[ULAB_MAX_URL];
    char session_token[4096];
    char token[4096];
    int  authenticated;

    char access_id[ULAB_MAX_ID];
    char backhaul_id[ULAB_MAX_ID];
    char power_id[ULAB_MAX_ID];
    char spectrum_id[ULAB_MAX_ID];
    char switch_id[ULAB_MAX_ID];

    FILE *logf;
} bff_client_t;

typedef struct {
    char id[ULAB_MAX_ID];
    char connectivity[ULAB_MAX_REF];
    char state[ULAB_MAX_REF];
} bff_node_status_t;

typedef struct {
    char name[ULAB_MAX_NAME];
    char version[ULAB_MAX_REF];
    char tag[ULAB_MAX_REF];
    char status[ULAB_MAX_REF];
} bff_app_state_t;

typedef struct {
    char name[ULAB_MAX_NAME];
    char type[ULAB_MAX_REF];
    char version[ULAB_MAX_REF];
    char uploaded_at[ULAB_MAX_REF];
    int  available;
    int  chunked;
    int  desired;
} bff_release_t;

#define ULAB_MAX_BFF_PAYMENTS 64

typedef struct {
    char id[ULAB_MAX_ID];
    char item_id[ULAB_MAX_ID];
    char item_type[ULAB_MAX_REF];
    char amount[ULAB_MAX_REF];
    char currency[ULAB_MAX_REF];
    char payment_method[ULAB_MAX_REF];
    char status[ULAB_MAX_REF];
    char paid_at[ULAB_MAX_REF];
    char payer_email[ULAB_MAX_NAME];
    char payer_phone[ULAB_MAX_REF];
    char metadata[ULAB_MAX_LINE];
} bff_payment_t;

#define ULAB_MAX_BFF_SIM_PACKAGES 32

typedef struct {
    char id[ULAB_MAX_ID];
    char package_id[ULAB_MAX_ID];
    char start_date[ULAB_MAX_REF];
    char end_date[ULAB_MAX_REF];
    int  active;
} bff_sim_package_t;

#define ULAB_MAX_BFF_KPI_SCOPES 8

typedef struct {
    char key[ULAB_MAX_REF];
    char value[ULAB_MAX_ID];
} bff_scope_entry_t;

typedef struct {
    char direction[ULAB_MAX_REF];
    double change_pct;
    double change_abs;
    double previous_value;
    int has_previous;
} bff_kpi_trend_t;

typedef struct {
    char kpi[ULAB_MAX_REF];
    double value;
    char span[ULAB_MAX_REF];
    char op[ULAB_MAX_REF];
    char unit[ULAB_MAX_REF];
    char symbol[ULAB_MAX_REF];
    char from[ULAB_MAX_REF];
    char to[ULAB_MAX_REF];
    char computed_at[ULAB_MAX_REF];
    int is_partial;
    bff_scope_entry_t scope[ULAB_MAX_BFF_KPI_SCOPES];
    size_t scope_count;
    bff_kpi_trend_t trend;
} bff_kpi_value_t;

int bff_init(bff_client_t *c,
             const char *url,
             const char *run_dir);

int bff_login(bff_client_t *c,
              const char *identifier,
              const char *password,
              ulab_error_t *err);

void bff_close(bff_client_t *c);

int bff_add_network(bff_client_t *c,
                    network_t *n,
                    ulab_error_t *err);

int bff_wait_site_anchor_online(bff_client_t *c,
                                site_t *site,
                                ulab_error_t *err);

int bff_add_site(bff_client_t *c,
                 site_t *s,
                 const network_t *n,
                 ulab_error_t *err);

int bff_add_package(bff_client_t *c,
                    package_t *p,
                    const network_t *net,
                    ulab_error_t *err);

int bff_set_package_active(bff_client_t *c,
                           package_t *p,
                           int active,
                           ulab_error_t *err);

int bff_add_subscriber(bff_client_t *c,
                       subscriber_t *sub,
                       const network_t *net,
                       ulab_error_t *err);

int bff_upload_sims_from_csv(bff_client_t *c,
                             const char *csv_path,
                             const char *sim_type,
                             ulab_error_t *err);

int bff_get_sims_from_pool(bff_client_t *c,
                           const char *sim_type,
                           char iccids[][ULAB_MAX_ID],
                           char pool_sim_ids[][ULAB_MAX_ID],
                           size_t max_iccids,
                           size_t *iccid_count,
                           ulab_error_t *err);

int bff_allocate_sim_from_pool(bff_client_t *c,
                               ue_t *ue,
                               const subscriber_t *sub,
                               const network_t *net,
                               const package_t *pkg,
                               const char *sim_type,
                               ulab_error_t *err);

int bff_allocate_sim(bff_client_t *c,
                     ue_t *ue,
                     const subscriber_t *sub,
                     const network_t *net,
                     const package_t *pkg,
                     const char *sim_type,
                     ulab_error_t *err);

int bff_add_package_to_sim(bff_client_t *c,
                           ue_t *ue,
                           const package_t *pkg,
                           ulab_error_t *err);

int bff_record_cash_package_sale(bff_client_t *c,
                                 ue_t *ue,
                                 const package_t *pkg,
                                 const subscriber_t *subscriber,
                                 double amount,
                                 const char *currency,
                                 bff_payment_t *payment,
                                 ulab_error_t *err);

int bff_get_package_payments(bff_client_t *c,
                             const package_t *pkg,
                             bff_payment_t payments[],
                             size_t max_payments,
                             size_t *payment_count,
                             ulab_error_t *err);

int bff_get_kpi_value(bff_client_t *c,
                      const char *key,
                      const char *span,
                      const char *op,
                      const char *network_id,
                      const char *scope_key,
                      const char *scope_value,
                      bff_kpi_value_t *value,
                      int *found,
                      ulab_error_t *err);

int bff_get_performance_report_cell(bff_client_t *c,
                                    const char *report,
                                    const char *span,
                                    const char *network_id,
                                    const char *entity_id,
                                    const char *column,
                                    double *value,
                                    char *unit,
                                    size_t unit_len,
                                    int *found,
                                    ulab_error_t *err);

int bff_clear_sim_packages(bff_client_t *c,
                           const ue_t *ue,
                           ulab_error_t *err);

int bff_toggle_sim_status(bff_client_t *c,
                          const ue_t *ue,
                          const char *status,
                          ulab_error_t *err);

int bff_get_sim_usage(bff_client_t *c,
                      const ue_t *ue,
                      uint64_t *used_mb,
                      ulab_error_t *err);

int bff_get_packages_for_sim(bff_client_t *c,
                             const ue_t *ue,
                             const char *package_id,
                             int *active,
                             ulab_error_t *err);

int bff_get_sim_packages(bff_client_t *c,
                         const ue_t *ue,
                         bff_sim_package_t packages[],
                         size_t max_packages,
                         size_t *package_count,
                         ulab_error_t *err);

int bff_get_node_status(bff_client_t *c,
                        const node_t *node,
                        bff_node_status_t *status,
                        ulab_error_t *err);

int bff_restart_node(bff_client_t *c,
                     const node_t *node,
                     ulab_error_t *err);

int bff_get_release(bff_client_t *c,
                    const char *name,
                    const char *type,
                    const char *version,
                    bff_release_t *release,
                    int *found,
                    ulab_error_t *err);

int bff_promote_release(bff_client_t *c,
                        const char *name,
                        const char *type,
                        const char *version,
                        ulab_error_t *err);

int bff_update_software(bff_client_t *c,
                        const node_t *node,
                        const char *app,
                        const char *tag,
                        ulab_error_t *err);

int bff_get_node_app(bff_client_t *c,
                     const node_t *node,
                     const char *app,
                     bff_app_state_t *state,
                     ulab_error_t *err);

int bff_toggle_site_service(bff_client_t *c,
                            const site_t *site,
                            int enabled,
                            ulab_error_t *err);

int bff_toggle_site_radio(bff_client_t *c,
                          const site_t *site,
                          int enabled,
                          ulab_error_t *err);

int bff_network_overview_loads(bff_client_t *c,
                               const network_t *net,
                               ulab_error_t *err);

int bff_site_view_loads(bff_client_t *c,
                        const site_t *site,
                        ulab_error_t *err);

int bff_backend_count(bff_client_t *c,
                      const char *target,
                      const world_t *w,
                      size_t *count,
                      ulab_error_t *err);

int bff_backend_contains(bff_client_t *c,
                         const char *view,
                         const char *id,
                         const world_t *w,
                         int *found,
                         ulab_error_t *err);

int bff_cleanup_world(bff_client_t *c,
                      const world_t *w,
                      ulab_error_t *err);

#endif /* ULAB_BFF_H_ */

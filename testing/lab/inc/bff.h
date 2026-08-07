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
    char id[ULAB_MAX_ID];
    char release_date[ULAB_MAX_REF];
    char node_id[ULAB_MAX_ID];
    char status[ULAB_MAX_REF];
    char current_version[ULAB_MAX_REF];
    char desired_version[ULAB_MAX_REF];
    char name[ULAB_MAX_NAME];
    char created_at[ULAB_MAX_REF];
    char updated_at[ULAB_MAX_REF];
} bff_software_t;

typedef struct {
    char id[ULAB_MAX_ID];
    char type[ULAB_MAX_REF];
    char status[ULAB_MAX_REF];
    char requested_by[ULAB_MAX_NAME];
    char started_at[ULAB_MAX_REF];
    char lease_expires_at[ULAB_MAX_REF];
} bff_operation_t;

typedef struct {
    char node_id[ULAB_MAX_ID];
    char node_type[ULAB_MAX_REF];
    int  busy;
    int  has_operation;
    bff_operation_t operation;
} bff_node_operation_status_t;

typedef struct {
    int  available;
    char reason[ULAB_MAX_ERR];
} bff_action_availability_t;

typedef struct {
    char site_id[ULAB_MAX_ID];
    int  busy;
    int  degraded;
    size_t node_count;
    bff_node_operation_status_t nodes[ULAB_MAX_LIST];
    bff_action_availability_t restart_site;
    bff_action_availability_t rf;
    bff_action_availability_t service;
} bff_site_operation_status_t;

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

typedef struct {
    char uuid[ULAB_MAX_ID];
    char name[ULAB_MAX_NAME];
    uint64_t data_volume;
    uint32_t duration_minutes;
    double amount;
    char data_unit[ULAB_MAX_REF];
    char currency[ULAB_MAX_REF];
    char country[ULAB_MAX_REF];
    char network_id[ULAB_MAX_ID];
    int active;
} bff_package_t;

typedef struct {
    char package_id[ULAB_MAX_ID];
    double revenue;
    uint32_t attach_count;
    int has_attach_count;
} bff_package_metrics_t;

typedef struct {
    double total_paid;
    double total_pending;
    double month_paid;
    double previous_month_paid;
    double month_over_month_percent;
} bff_revenue_summary_t;

typedef struct {
    double mrr;
    double arpu;
    int    has_mrr;
    int    has_arpu;
} bff_package_kpis_t;

typedef struct {
    uint32_t subscribers_total;
    uint32_t subscribers_active;
    uint32_t subscribers_inactive;
    uint32_t sites_total;
    uint32_t nodes_total;
    uint32_t nodes_online;
    uint32_t nodes_offline;
} bff_network_summary_t;

typedef struct {
    uint32_t component_total;
    uint32_t component_category_total;
    uint32_t sim_total;
    uint32_t sim_available;
    uint32_t sim_consumed;
    uint32_t sim_pool_total;
    uint32_t sim_pool_available;
    uint32_t sim_pool_consumed;
} bff_inventory_summary_t;

typedef struct {
    uint32_t payment_count;
    uint32_t settled_count;
    double   settled_amount;
} bff_subscriber_billing_t;

typedef struct {
    char status[ULAB_MAX_REF];
    /* Value of the row's "active" attribute ("true"/"false"). Distinct from
     * status, which is a sales-performance label, not the plan flag. */
    char active[ULAB_MAX_REF];
    uint32_t row_index;
    uint32_t row_count;
    int has_name;
    int has_price;
    int has_validity;
    int has_active;
} bff_performance_row_t;

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

int bff_update_package_name(bff_client_t *c,
                            package_t *p,
                            const char *name,
                            ulab_error_t *err);

int bff_package_name_available(bff_client_t *c,
                               const char *name,
                               int *available,
                               ulab_error_t *err);

int bff_get_package(bff_client_t *c,
                    const package_t *pkg,
                    bff_package_t *actual,
                    ulab_error_t *err);

int bff_package_visible_for_network(bff_client_t *c,
                                    const package_t *pkg,
                                    const network_t *network,
                                    int *visible,
                                    ulab_error_t *err);

int bff_invalid_package_name_available(bff_client_t *c,
                                       const package_t *pkg,
                                       const char *variant,
                                       int *available,
                                       ulab_error_t *err);

int bff_add_invalid_package(bff_client_t *c,
                            const package_t *pkg,
                            const network_t *network,
                            const char *variant,
                            char *created_id,
                            size_t created_id_len,
                            ulab_error_t *err);

int bff_get_package_metrics(bff_client_t *c,
                            const package_t *pkg,
                            const network_t *network,
                            bff_package_metrics_t *metrics,
                            int *found,
                            ulab_error_t *err);

int bff_get_revenue_summary(bff_client_t *c,
                            const network_t *network,
                            bff_revenue_summary_t *summary,
                            ulab_error_t *err);

int bff_get_package_kpis(bff_client_t *c,
                         const network_t *network,
                         bff_package_kpis_t *kpis,
                         ulab_error_t *err);

int bff_get_network_summary(bff_client_t *c,
                            const network_t *network,
                            bff_network_summary_t *summary,
                            ulab_error_t *err);

int bff_get_nodes_count(bff_client_t *c,
                        const network_t *network,
                        uint32_t *count,
                        ulab_error_t *err);

int bff_get_inventory_summary(bff_client_t *c,
                              const char *sim_type,
                              bff_inventory_summary_t *summary,
                              ulab_error_t *err);

int bff_get_subscriber_payment_summary(
    bff_client_t *c,
    const subscriber_t *subscriber,
    bff_subscriber_billing_t *billing,
    ulab_error_t *err);

int bff_sim_is_unallocated(bff_client_t *c,
                           const ue_t *ue,
                           const char *sim_type,
                           int *unallocated,
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

int bff_get_performance_report_row(bff_client_t *c,
                                   const char *report,
                                   const char *span,
                                   const char *network_id,
                                   const char *entity_id,
                                   bff_performance_row_t *row,
                                   int *found,
                                   ulab_error_t *err);

int bff_clear_sim_packages(bff_client_t *c,
                           const ue_t *ue,
                           ulab_error_t *err);

int bff_toggle_sim_status(bff_client_t *c,
                          const ue_t *ue,
                          const char *status,
                          ulab_error_t *err);

int bff_get_sim_status(bff_client_t *c,
                       const ue_t *ue,
                       char *status,
                       size_t status_len,
                       ulab_error_t *err);

int bff_get_sim_usage(bff_client_t *c,
                      const ue_t *ue,
                      const network_t *network,
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

int bff_restart_site(bff_client_t *c,
                     const site_t *site,
                     const network_t *network,
                     ulab_error_t *err);

int bff_toggle_internet_switch(bff_client_t *c,
                               const site_t *site,
                               uint32_t port,
                               int enabled,
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

int bff_console_network_loads(bff_client_t *c,
                              const network_t *net,
                              ulab_error_t *err);

int bff_backend_count(bff_client_t *c,
                      const char *target,
                      const world_t *w,
                      size_t *count,
                      ulab_error_t *err);

/*
 * Count entries in a backend list. When world is non-NULL the returned
 * count is RUN-SCOPED: only entries whose id/uuid belongs to this run's
 * world are counted, so residue from earlier runs cannot inflate it.
 * backend_total (optional) receives the raw list size for diagnostics.
 */
int bff_get_list_count(bff_client_t *c,
                       const char *target,
                       const network_t *network,
                       const world_t *world,
                       size_t *count,
                       size_t *backend_total,
                       ulab_error_t *err);

int bff_get_node_list_count(bff_client_t *c,
                            const network_t *network,
                            const char *node_type,
                            const char *view,
                            size_t *count,
                            ulab_error_t *err);

int bff_get_site_node_count(bff_client_t *c,
                            const site_t *site,
                            const char *node_type,
                            size_t *count,
                            ulab_error_t *err);

int bff_entity_fields_match_world(bff_client_t *c,
                                  const char *entity,
                                  const char *ref,
                                  const world_t *world,
                                  const char *view,
                                  int *matched,
                                  char *detail,
                                  size_t detail_len,
                                  ulab_error_t *err);

int bff_entity_list_detail_reconciles(bff_client_t *c,
                                      const char *entity,
                                      const char *ref,
                                      const world_t *world,
                                      const char *view,
                                      int *matched,
                                      char *detail,
                                      size_t detail_len,
                                      ulab_error_t *err);

int bff_get_software(bff_client_t *c,
                     const node_t *node,
                     const char *app,
                     const char *view,
                     bff_software_t *software,
                     int *found,
                     ulab_error_t *err);

int bff_get_software_count(bff_client_t *c,
                           const node_t *node,
                           size_t *count,
                           ulab_error_t *err);

int bff_get_software_list(bff_client_t *c,
                          const node_t *node,
                          const char *view,
                          bff_software_t software[],
                          size_t max_software,
                          size_t *software_count,
                          ulab_error_t *err);

int bff_get_node_status_for_view(bff_client_t *c,
                                 const network_t *network,
                                 const node_t *node,
                                 const char *view,
                                 bff_node_status_t *status,
                                 ulab_error_t *err);

int bff_get_node_operation_status(
    bff_client_t *c,
    const node_t *node,
    bff_node_operation_status_t *status,
    ulab_error_t *err);

int bff_get_site_operation_status(
    bff_client_t *c,
    const site_t *site,
    bff_site_operation_status_t *status,
    ulab_error_t *err);

int bff_get_kpi_timeseries(bff_client_t *c,
                           const char *key,
                           const char *span,
                           const char *op,
                           const char *from,
                           const char *to,
                           const char *network_id,
                           const char *site_id,
                           bff_kpi_value_t values[],
                           size_t max_values,
                           size_t *value_count,
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

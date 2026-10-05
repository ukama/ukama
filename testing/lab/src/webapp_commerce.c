/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include "webapp.h"
#include "util.h"
#include <string.h>

const char *webapp_commerce_kind(const event_spec_t *event) {
    if (event->type != EVT_WEB_COMMERCE) return NULL;
    if (!strcmp(event->target, "create_plan")) return "package";
    if (!strcmp(event->target, "create_customer")) return "subscriber";
    if ((!strcmp(event->target, "allocate_sim") || !strcmp(event->target, "allocate_auto"))) return "sim";
    if (!strcmp(event->target, "top_up") || !strcmp(event->target,"top_up_rapid") || !strcmp(event->target,"failed_top_up")) return "payment";
    return NULL;
}
static json_t *plan_input(const package_t *p, const char *unit) {
    return json_pack("{s:s,s:s,s:s,s:I,s:i,s:f,s:s,s:s,s:s,s:b}",
        "ref", p->ref, "id", p->bff_id, "name", p->web_name[0] ? p->web_name : p->name, "data_mb", (json_int_t)p->data_mb,
        "duration_minutes", (int)(p->duration_minutes ? p->duration_minutes : p->duration_days * 1440),
        "amount", p->amount, "currency", p->currency, "country", p->country,
        "unit", unit, "organization", !p->network_ref[0]);
}
int webapp_commerce_inputs(const event_spec_t *e, world_t *w, json_t *inputs, ulab_error_t *err) {
    network_t *n = world_network_by_ref(w, e->networks.value);
    package_t *p = e->package_ref[0] ? world_package_by_ref(w, e->package_ref) : NULL;
    ue_t *ue = e->ues.kind == SEL_REF ? world_ue_by_ref(w, e->ues.value) : NULL;
    subscriber_t *sub = ue ? world_subscriber_by_ref(w, ue->subscriber_ref) : NULL;
    const char *kind = webapp_commerce_kind(e);
    const char *ref = NULL;
    const char *name = NULL;
    const char *id = NULL;
    if (!n || !n->bff_id[0] || (p && p->network_ref[0] && strcmp(p->network_ref, n->ref)) ||
        (ue && (!sub || strcmp(ue->network_ref, n->ref))))
        return webapp_error(err, "commerce references do not match the selected network");
    json_object_set_new(inputs, "network_id", json_string(n->bff_id));
    json_object_set_new(inputs, "view", json_string(e->view));
    json_object_set_new(inputs, "action", json_string(e->target));
    if (p) {
        if (strcmp(e->target, "create_plan") && strcmp(e->target, "name_pending") && strcmp(e->target, "name_failure") && !p->bff_id[0]) return webapp_error(err, "commerce plan has not been created through the UI");
        json_object_set_new(inputs, "plan", plan_input(p, e->variant[0] ? e->variant : "MB"));
    }
    if (!strcmp(e->target, "rename_plan")) {
        char name[ULAB_MAX_NAME];
        if (!p || strlen(p->web_name[0] ? p->web_name : p->name) + strlen("-renamed") >= sizeof(name)) return webapp_error(err,"renamed plan name too long");
        snprintf(name,sizeof(name),"%.247s-renamed",p->web_name[0] ? p->web_name : p->name);
        json_object_set_new(inputs,"new_name",json_string(name));
    }
    if (ue && sub) {
        if (strcmp(e->target, "create_customer") && !sub->bff_id[0]) return webapp_error(err, "customer has not been created through the UI");
        if ((!strcmp(e->target, "top_up") || !strcmp(e->target,"top_up_rapid") || !strcmp(e->target,"failed_top_up") || !strcmp(e->target,"download_receipt") || !strcmp(e->target, "cancel_top_up") || !strcmp(e->target, "open_receipt") ||
             !strcmp(e->target, "activate_sim") || !strcmp(e->target, "deactivate_sim")) && !ue->bff_id[0])
            return webapp_error(err, "SIM has not been allocated through the UI");
        json_object_set_new(inputs, "customer", json_pack("{s:s,s:s,s:s,s:s,s:s,s:s,s:s}",
            "payment_id",ue->last_payment_id,"ref", sub->ref, "id", sub->bff_id, "name", sub->name, "email", sub->email, "iccid", ue->iccid, "sim_id", ue->bff_id));
    }
    if (kind) {
        if (!strcmp(kind, "package")) { ref = p->ref; name = p->name; id = p->bff_id; }
        else if (!strcmp(kind, "subscriber")) { ref = sub->ref; name = sub->name; id = sub->bff_id; }
        else { ref = ue->ref; name = ue->iccid; id = !strcmp(kind, "sim") ? ue->bff_id : ue->last_payment_id; }
        if (id[0]) return webapp_error(err, "commerce mutation already has an identity; refusing duplicate submission");
        json_object_set_new(inputs, "creation", json_pack("{s:s,s:s,s:s}", "kind", kind, "ref", ref, "name", name));
    }
    return ULAB_OK;
}
int webapp_commerce_check(const check_spec_t *check, check_spec_t *resolved, world_t *w, json_t **inputs, ulab_error_t *err) {
    package_t *p = check->package_ref[0] ? world_package_by_ref(w, check->package_ref) : NULL;
    ue_t *ue = check->ues.kind == SEL_REF ? world_ue_by_ref(w, check->ues.value) : NULL;
    subscriber_t *sub = ue ? world_subscriber_by_ref(w, ue->subscriber_ref) : NULL;
    if ((check->package_ref[0] && (!p || (!p->bff_id[0] && strcmp(check->label,"Commerce fault")))) || (check->ues.kind != SEL_NONE && (!ue || !sub)))
        return webapp_error(err, "commerce check identity is unresolved");
    if (check->key[0]) {
        const char *value = !strcmp(check->key, "plan_name") ? (p ? (p->web_name[0] ? p->web_name : p->name) : "") :
            !strcmp(check->key, "iccid") ? (ue ? ue->iccid : "") : ue ? ue->last_payment_id : "";
        if (!*value || ulab_copy(resolved->expected, sizeof(resolved->expected), value))
            return webapp_error(err, "commerce expected identity is missing or too long");
    }
    *inputs = webapp_check_inputs(resolved);
    if (!*inputs) return webapp_error(err, "cannot encode commerce expectation");
    if (p) {
        json_object_set_new(*inputs,"plan_id",json_string(p->bff_id));
        json_object_set_new(*inputs, "plan_name", json_string(p->web_name[0] ? p->web_name : p->name));
    }
    if (ue) {
        json_object_set_new(*inputs, "customer_name", json_string(sub->name));
        json_object_set_new(*inputs, "iccid", json_string(ue->iccid));
        json_object_set_new(*inputs,"sim_id",json_string(ue->bff_id));
        json_object_set_new(*inputs,"payment_id",json_string(ue->last_payment_id));
    }
    return ULAB_OK;
}

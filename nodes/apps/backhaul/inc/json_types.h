/*
 * JSON serialization routines for backhaul.d
 * Keep all jansson usage in this file pair only.
 */

#ifndef JSON_SERDES_H
#define JSON_SERDES_H

#include "jansson.h"
#include "metrics_store.h"

/* All exported JSON routines must be prefixed json_ */
json_t* json_backhaul_status(MetricsStore *store);
json_t* json_backhaul_metrics(MetricsStore *store);

#endif /* JSON_TYPES_H_ */

/*
 * Small HTTP/CORS + URL param helpers for FEMD web API (Ulfius)
 * C89-compatible
 */

#include <stdlib.h>
#include "api_http.h"

void api_set_cors_allow(UResponse *resp, const char *allow_str) {
    if (allow_str != NULL) {
        u_map_put(resp->map_header, "Allow", allow_str);
        u_map_put(resp->map_header, "Access-Control-Allow-Methods", allow_str);
    }
    u_map_put(resp->map_header, "Access-Control-Allow-Headers", "Content-Type");
    u_map_put(resp->map_header, "Access-Control-Allow-Origin", "*");
}

int api_parse_fem_id(const URequest *req, FemUnit *out_unit) {
    const char *fem_str;
    long v;
    if (out_unit == NULL) return 0;

    fem_str = u_map_get(req->map_url, "femId");
    if (fem_str == NULL || *fem_str == '\0') return 0;

    v = strtol(fem_str, NULL, 10);
    if (v != 1 && v != 2) return 0;

    *out_unit = (v == 1) ? FEM_UNIT_1 : FEM_UNIT_2;
    return 1;
}

int api_parse_channel_id(const URequest *req, int *out_ch) {
    const char *ch_str;
    long v;
    if (out_ch == NULL) return 0;

    ch_str = u_map_get(req->map_url, "channel");
    if (ch_str == NULL || *ch_str == '\0') return 0;

    v = strtol(ch_str, NULL, 10);
    if (v < 0 || v >= (long)ADC_CHANNEL_MAX) return 0;

    *out_ch = (int)v;
    return 1;
}


/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */


#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>
#include <termios.h>
#include <math.h>
#include <sys/select.h>
#include <pthread.h>

#include "drv_victron.h"
#include "time_util.h"
#include "usys_log.h"

static int configure_serial(int fd, int baud) {
    struct termios tty;
    speed_t speed;

    if (tcgetattr(fd, &tty) != 0) {
        usys_log_error("victron: tcgetattr failed: %s", strerror(errno));
        return -1;
    }

    switch (baud) {
    case 19200:  speed = B19200;  break;
    case 9600:   speed = B9600;   break;
    case 115200: speed = B115200; break;
    default:
        usys_log_error("victron: unsupported baud rate %d", baud);
        return -1;
    }
    cfsetispeed(&tty, speed);
    cfsetospeed(&tty, speed);

    tty.c_cflag &= ~PARENB;
    tty.c_cflag &= ~CSTOPB;
    tty.c_cflag &= ~CSIZE;
    tty.c_cflag |= CS8;
    tty.c_cflag &= ~CRTSCTS;
    tty.c_cflag |= CREAD | CLOCAL;

    tty.c_lflag &= ~(ICANON | ECHO | ECHOE | ISIG);
    tty.c_iflag &= ~(IXON | IXOFF | IXANY);
    tty.c_iflag &= ~(INLCR | ICRNL | IGNCR);
    tty.c_oflag &= ~OPOST;

    tty.c_cc[VMIN]  = 0;
    tty.c_cc[VTIME] = 1;

    if (tcsetattr(fd, TCSANOW, &tty) != 0) {
        usys_log_error("victron: tcsetattr failed: %s", strerror(errno));
        return -1;
    }

    tcflush(fd, TCIOFLUSH);
    return 0;
}

static int parse_field(const char *line, VeDirectField *field) {
    const char *tab = strchr(line, '\t');
    size_t label_len, value_len;
    const char *value;

    if (!tab) return -1;

    label_len = tab - line;
    if (label_len >= VICTRON_MAX_LABEL_LEN) return -1;

    strncpy(field->label, line, label_len);
    field->label[label_len] = '\0';

    value = tab + 1;
    value_len = strlen(value);

    while (value_len > 0 && (value[value_len-1] == '\r' || value[value_len-1] == '\n')) {
        value_len--;
    }

    if (value_len >= VICTRON_MAX_VALUE_LEN) return -1;

    strncpy(field->value, value, value_len);
    field->value[value_len] = '\0';

    return 0;
}

static uint8_t calculate_checksum(const char *block, size_t len) {
    uint8_t sum = 0;
    for (size_t i = 0; i < len; i++) {
        sum += (uint8_t)block[i];
    }
    return sum;
}

static int parse_frame(VictronCtx *ctx, const char *frame, size_t len) {
    char *buf, *line, *saveptr;

    if (calculate_checksum(frame, len) != 0) {
        ctx->checksum_errors++;
        usys_log_debug("victron: checksum error");
        return -1;
    }

    buf = strndup(frame, len);
    if (!buf) return -1;

    ctx->field_count = 0;

    line = strtok_r(buf, "\r\n", &saveptr);
    while (line && ctx->field_count < VICTRON_MAX_FIELDS) {
        if (strlen(line) > 0 && strchr(line, '\t')) {
            if (parse_field(line, &ctx->fields[ctx->field_count]) == 0) {
                ctx->field_count++;
            }
        }
        line = strtok_r(NULL, "\r\n", &saveptr);
    }

    free(buf);

    if (ctx->field_count == 0) {
        ctx->parse_errors++;
        return -1;
    }

    ctx->frame_valid  = true;
    ctx->frames_received++;
    ctx->last_frame_ts = time_now_ms();

    return 0;
}

static const char *find_field(VictronCtx *ctx, const char *label) {
    for (int i = 0; i < ctx->field_count; i++) {
        if (strcmp(ctx->fields[i].label, label) == 0) {
            return ctx->fields[i].value;
        }
    }
    return NULL;
}

/* Map Victron CS enum values to driver-agnostic ChargeState */
static ChargeState map_charge_state(int cs) {
    switch (cs) {
    case VCS_OFF:                   return CHARGE_STATE_OFF;
    case VCS_LOW_POWER:             return CHARGE_STATE_OFF;
    case VCS_FAULT:                 return CHARGE_STATE_FAULT;
    case VCS_BULK:                  return CHARGE_STATE_BULK;
    case VCS_ABSORPTION:            return CHARGE_STATE_ABSORPTION;
    case VCS_REPEATED_ABSORPTION:   return CHARGE_STATE_ABSORPTION;
    case VCS_FLOAT:                 return CHARGE_STATE_FLOAT;
    case VCS_BATTERY_SAFE:          return CHARGE_STATE_FLOAT;
    case VCS_STORAGE:               return CHARGE_STATE_STORAGE;
    case VCS_EQUALIZE:              return CHARGE_STATE_EQUALIZE;
    case VCS_AUTO_EQUALIZE:         return CHARGE_STATE_EQUALIZE;
    case VCS_STARTING:              return CHARGE_STATE_BULK;
    default:                        return CHARGE_STATE_UNKNOWN;
    }
}

static void fields_to_data(VictronCtx *ctx, ControllerData *data) {
    const char *val;

    memset(data, 0, sizeof(*data));
    data->timestamp_ms  = time_now_ms();
    data->batt_soc_pct  = -1;
    data->temperature_c = NAN;

    val = find_field(ctx, VE_LABEL_VOLTAGE);
    if (val) data->batt_voltage_v = atof(val) / 1000.0;

    val = find_field(ctx, VE_LABEL_CURRENT);
    if (val) data->batt_current_a = atof(val) / 1000.0;

    val = find_field(ctx, VE_LABEL_PV_VOLTAGE);
    if (val) data->pv_voltage_v = atof(val) / 1000.0;

    val = find_field(ctx, VE_LABEL_PV_POWER);
    if (val) data->pv_power_w = atof(val);

    if (data->pv_voltage_v > 0 && data->pv_power_w > 0) {
        data->pv_current_a = data->pv_power_w / data->pv_voltage_v;
    }

    val = find_field(ctx, VE_LABEL_CHARGE_STATE);
    if (val) data->charge_state = map_charge_state(atoi(val));

    val = find_field(ctx, VE_LABEL_ERROR);
    if (val) data->error_code = (uint32_t)atoi(val);

    /* H19 and H20 are in units of 0.01 kWh */
    val = find_field(ctx, VE_LABEL_YIELD_TOTAL);
    if (val) data->yield_total_kwh = atof(val) / 100.0;

    val = find_field(ctx, VE_LABEL_YIELD_TODAY);
    if (val) data->yield_today_kwh = atof(val) / 100.0;

    val = find_field(ctx, VE_LABEL_FIRMWARE);
    if (val) snprintf(data->firmware, sizeof(data->firmware), "%s", val);

    val = find_field(ctx, VE_LABEL_SERIAL);
    if (val) snprintf(data->serial, sizeof(data->serial), "%s", val);

    val = find_field(ctx, VE_LABEL_PRODUCT_ID);
    if (val) snprintf(data->product_id, sizeof(data->product_id), "%s", val);

    val = find_field(ctx, VE_LABEL_LOAD);
    if (val) {
        data->load_output_available = true;
        data->load_output_state     = (strcmp(val, "ON") == 0);
    }

    val = find_field(ctx, VE_LABEL_LOAD_CURRENT);
    if (val) data->load_current_a = atof(val) / 1000.0;

    val = find_field(ctx, VE_LABEL_RELAY);
    if (val) {
        data->relay_available = true;
        data->relay_state     = (strcmp(val, "ON") == 0);
    }

    /* Battery temperature sent only when a Smart Battery Sense is paired */
    val = find_field(ctx, VE_LABEL_TEMP);
    if (val) data->temperature_c = atof(val);

    data->comm_ok = true;
}

static int read_frame(VictronCtx *ctx, int timeout_ms) {
    char buf[256];
    ssize_t n;
    uint64_t start = time_now_ms();

    while ((time_now_ms() - start) < (uint64_t)timeout_ms) {
        fd_set fds;
        struct timeval tv;
        FD_ZERO(&fds);
        FD_SET(ctx->fd, &fds);
        tv.tv_sec  = 0;
        tv.tv_usec = 100000;

        int ret = select(ctx->fd + 1, &fds, NULL, NULL, &tv);
        if (ret < 0) {
            if (errno == EINTR) continue;
            usys_log_error("victron: select failed: %s", strerror(errno));
            return -1;
        }
        if (ret == 0) continue;

        n = read(ctx->fd, buf, sizeof(buf) - 1);
        if (n < 0) {
            if (errno == EAGAIN || errno == EINTR) continue;
            usys_log_error("victron: read failed: %s", strerror(errno));
            return -1;
        }
        if (n == 0) continue;

        if (ctx->rx_len + n >= (int)sizeof(ctx->rx_buf)) {
            ctx->rx_len  = 0;
        }

        memcpy(ctx->rx_buf + ctx->rx_len, buf, n);
        ctx->rx_len += n;
        ctx->rx_buf[ctx->rx_len] = '\0';

        /* "Checksum\t<byte>" marks the end of a frame */
        char *checksum_pos = strstr(ctx->rx_buf, "Checksum\t");
        if (checksum_pos) {
            char *frame_end = checksum_pos + 9 + 1;  /* past the checksum byte */
            if (frame_end <= ctx->rx_buf + ctx->rx_len) {
                size_t frame_len = frame_end - ctx->rx_buf;

                /* Skip the trailing \r\n that follows the checksum byte.
                 * Without this, \r\n stays in rx_buf and corrupts the
                 * checksum of the next frame (+13+10 = +23 offset). */
                char *next_frame = frame_end;
                if (next_frame + 1 < ctx->rx_buf + ctx->rx_len &&
                    next_frame[0] == '\r' && next_frame[1] == '\n') {
                    next_frame += 2;
                }

                if (parse_frame(ctx, ctx->rx_buf, frame_len) == 0) {
                    int remaining = ctx->rx_len - (int)(next_frame - ctx->rx_buf);
                    if (remaining > 0) {
                        memmove(ctx->rx_buf, next_frame, remaining);
                    }
                    ctx->rx_len = remaining;
                    return 0;
                }

                int remaining = ctx->rx_len - (int)(next_frame - ctx->rx_buf);
                if (remaining > 0) {
                    memmove(ctx->rx_buf, next_frame, remaining);
                }
                ctx->rx_len = remaining;
            }
        }

        if (ctx->rx_len > 400) {
            memmove(ctx->rx_buf, ctx->rx_buf + 200, ctx->rx_len - 200);
            ctx->rx_len -= 200;
        }
    }

    return -1;
}

/*
 * Public driver interface
 */

int victron_open(void *vctx, const char *port, int baud) {
    VictronCtx *ctx = (VictronCtx *)vctx;

    if (!ctx || !port) return -1;

    memset(ctx, 0, sizeof(*ctx));
    ctx->fd = -1;

    if (baud <= 0) baud = VICTRON_BAUD_RATE;

    ctx->fd = open(port, O_RDWR | O_NOCTTY | O_NONBLOCK);
    if (ctx->fd < 0) {
        usys_log_error("victron: failed to open %s: %s", port, strerror(errno));
        return -1;
    }

    if (configure_serial(ctx->fd, baud) != 0) {
        close(ctx->fd);
        ctx->fd = -1;
        return -1;
    }

    strncpy(ctx->port, port, sizeof(ctx->port) - 1);
    ctx->baud = baud;

    if (pthread_mutex_init(&ctx->lock, NULL) != 0) {
        close(ctx->fd);
        ctx->fd = -1;
        return -1;
    }

    usys_log_info("victron: opened %s at %d baud", port, baud);
    return 0;
}

void victron_close(void *vctx) {
    VictronCtx *ctx = (VictronCtx *)vctx;

    if (!ctx) return;

    if (ctx->fd >= 0) {
        close(ctx->fd);
        ctx->fd = -1;
    }

    pthread_mutex_destroy(&ctx->lock);

    usys_log_info("victron: closed %s (frames=%u, checksum_err=%u, parse_err=%u)",
                  ctx->port, ctx->frames_received,
                  ctx->checksum_errors, ctx->parse_errors);
}

int victron_read_data(void *vctx, ControllerData *out) {
    VictronCtx *ctx = (VictronCtx *)vctx;

    if (!ctx || !out || ctx->fd < 0) return -1;

    pthread_mutex_lock(&ctx->lock);

    if (read_frame(ctx, VICTRON_FRAME_TIMEOUT_MS) != 0) {
        ctx->cached_data.comm_ok = false;
        ctx->cached_data.comm_errors++;
        memcpy(out, &ctx->cached_data, sizeof(*out));
        pthread_mutex_unlock(&ctx->lock);
        return -1;
    }

    fields_to_data(ctx, &ctx->cached_data);
    memcpy(out, &ctx->cached_data, sizeof(*out));

    pthread_mutex_unlock(&ctx->lock);
    return 0;
}

/*
 * Control via VE.Direct HEX protocol (optional feature).
 * Sending a HEX command pauses TEXT frames for a few seconds.
 * Not all MPPT models support all registers via HEX.
 */

const ControllerDriver victron_driver = {
    .name        = "victron",
    .description = "Victron VE.Direct (MPPT SmartSolar/BlueSolar)",
    .open        = victron_open,
    .close       = victron_close,
    .read_data   = victron_read_data,
    .set_absorption_voltage = NULL,
    .set_float_voltage      = NULL,
    .set_charge_mode        = NULL,
    .set_relay              = NULL,
    .set_load_output        = NULL,
    .ctx_size    = sizeof(VictronCtx)
};


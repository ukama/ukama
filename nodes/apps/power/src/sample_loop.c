/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/time.h>

#include "sample_loop.h"
#include "drv_lm25066.h"
#include "drv_ads1015.h"
#include "drv_lm75.h"
#include "usys_log.h"

static uint64_t now_unix_ms(void) {

	struct timeval tv;
	gettimeofday(&tv, NULL);
	return (uint64_t)tv.tv_sec * 1000ULL + (uint64_t)(tv.tv_usec / 1000);
}

static void set_err(PowerMetrics *m, const char *msg) {

	m->ok = 0;
	strncpy(m->err, msg ? msg : "error", sizeof(m->err)-1);
}

static void *sample_thread(void *arg) {

	SampleLoop *l = (SampleLoop *)arg;

	Lm25066 lm25066;
	Lm75    lm75;
	Ads1015 ads1015;

	int haveLm25066 = 0;
	int haveLm75 = 0;
	int haveAds = 0;

	memset(&lm25066, 0, sizeof(lm25066));
	memset(&lm75,    0, sizeof(lm75));
	memset(&ads1015, 0, sizeof(ads1015));

	if (l->config->lm25066Dev && l->config->lm25066Addr) {
		if (drv_lm25066_open(&lm25066, l->config->lm25066Dev,
                             l->config->lm25066Addr,
		                     l->config->lm25066ClHigh,
                             l->config->lm25066RsMohm) == 0) {
			haveLm25066 = 1;
			usys_log_info("lm25066 enabled: dev=%s addr=0x%02x rs=%dmohm clHigh=%d",
			              l->config->lm25066Dev,
                          l->config->lm25066Addr,
                          l->config->lm25066RsMohm,
                          l->config->lm25066ClHigh);
		} else {
			usys_log_error("lm25066 open failed (disabled)");
		}
	}

	if (l->config->lm75Dev && l->config->lm75Addr) {
		if (drv_lm75_open(&lm75, l->config->lm75Dev, l->config->lm75Addr) == 0) {
			haveLm75 = 1;
			usys_log_info("lm75 enabled: dev=%s addr=0x%02x",
                          l->config->lm75Dev,
                          l->config->lm75Addr);
		} else {
			usys_log_error("lm75 open failed (disabled)");
		}
	}

	if (l->config->ads1015Dev && l->config->ads1015Addr) {
		if (drv_ads1015_open(&ads1015, l->config->ads1015Dev, l->config->ads1015Addr) == 0) {
			haveAds = 1;
			usys_log_info("ads1015 enabled: dev=%s addr=0x%02x chmap(vin=%d,vpa=%d,aux=%d)",
			              l->config->ads1015Dev,
                          l->config->ads1015Addr,
			              l->config->adsChVin,
                          l->config->adsChVpa,
                          l->config->adsChAux);
		} else {
			usys_log_error("ads1015 open failed (disabled)");
		}
	}

	while (!l->stop) {
		PowerMetrics m;
		memset(&m, 0, sizeof(m));

		m.sampleUnixMs = now_unix_ms();
		strncpy(m.board, l->config->boardName ?
                l->config->boardName : "unknown", sizeof(m.board)-1);
		m.ok = 1;
		m.err[0] = 0;

		m.haveLm25066 = haveLm25066;
		m.haveLm75 = haveLm75;
		m.haveAds1015 = haveAds;

		if (haveLm25066) {
			Lm25066Sample s;
			if (drv_lm25066_read_sample(&lm25066, &s) == 0) {
				m.inVolts        = s.vinV;
				m.outVolts       = s.voutV;
				m.inAmps         = s.iinA;
				m.inWatts        = s.pinW;
				m.hsTempC        = s.tempC;
				m.statusWord     = s.statusWord;
				m.diagnosticWord = s.diagnosticWord;
			} else {
				set_err(&m, "lm25066 read failed");
			}
		}

		if (haveLm75) {
			double t;
			if (drv_lm75_read_temp_c(&lm75, &t) == 0) {
				m.boardTempC = t;
			} else {
				set_err(&m, "lm75 read failed");
			}
		}

		if (haveAds) {
			double v;
			if (l->config->adsChVin >= 0) {
				if (drv_ads1015_read_single_ended(&ads1015,
                                                  l->config->adsChVin,
                                                  &v) == 0) m.adcVin = v;
				else set_err(&m, "ads1015 vin read failed");
			}

			if (l->config->adsChVpa >= 0) {
				if (drv_ads1015_read_single_ended(&ads1015,
                                                  l->config->adsChVpa,
                                                  &v) == 0) m.adcVpa = v;
				else set_err(&m, "ads1015 vpa read failed");
			}

			if (l->config->adsChAux >= 0) {
				if (drv_ads1015_read_single_ended(&ads1015,
                                                  l->config->adsChAux,
                                                  &v) == 0) m.adcAux = v;
				else set_err(&m, "ads1015 aux read failed");
			}
		}

		metrics_store_set(l->store, &m);

		usleep(l->config->sampleMs * 1000);
	}

	if (haveLm25066) drv_lm25066_close(&lm25066);
	if (haveLm75)    drv_lm75_close(&lm75);
	if (haveAds)     drv_ads1015_close(&ads1015);

	return NULL;
}

int sample_loop_start(SampleLoop *l, const Config *config, MetricsStore *store) {

	memset(l, 0, sizeof(*l));
	l->config = config;
	l->store  = store;
	l->stop   = 0;

	if (pthread_create(&l->thread, NULL, sample_thread, l) != 0) {
		usys_log_error("sample_loop: pthread_create failed");
		return -1;
	}
	return 0;
}

void sample_loop_stop(SampleLoop *l) {

	if (!l) return;

	l->stop = 1;
	pthread_join(l->thread, NULL);
	memset(l, 0, sizeof(*l));
}

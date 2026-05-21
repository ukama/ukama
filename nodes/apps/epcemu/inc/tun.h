/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef TUN_H_
#define TUN_H_

int tun_create(const char *name);
int tun_configure(const char *name, const char *addr);
void tun_close(int fd);

#endif /* TUN_H_ */

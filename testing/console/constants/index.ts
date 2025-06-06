/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import dotenv from 'dotenv';
dotenv.config();

export const CONSOLE_ROOT_URL =
  process.env.CONSOLE_BASE_URL || 'http://localhost:3000';
export const LOGIN_URL = `${process.env.AUTH_BASE_URL || 'http://localhost:4455'}/auth/login`;
export const MANAGE_URL = `${CONSOLE_ROOT_URL}/manage`;
export const TEST_USER_EMAIL = process.env.TEST_USER_EMAIL || 'admin@ukama.com';
export const TEST_USER_PASSWORD =
  process.env.TEST_USER_PASSWORD || '@Pass2025.';
export const TMP_SIMS_PATH = 'test-temp';
export const LIGHTHOUSE_SCORE_THRESHOLD =
  process.env.LIGHTHOUSE_SCORE_THRESHOLD || 0.2;

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Page } from '@playwright/test';
import { TEST_USER_EMAIL, TEST_USER_PASSWORD } from '../constants';
export async function login(
  page: Page,
  email: string = TEST_USER_EMAIL,
  password: string = TEST_USER_PASSWORD,
) {
  await page.goto('http://localhost:4455/auth/login');
  await page.getByRole('textbox', { name: 'EMAIL' }).click();
  await page.getByRole('textbox', { name: 'EMAIL' }).fill(email);
  await page.getByRole('textbox', { name: 'PASSWORD' }).click();
  await page.getByRole('textbox', { name: 'PASSWORD' }).fill(password);
  await page.getByRole('button', { name: 'LOG IN' }).click();
}

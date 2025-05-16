/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { test as base } from '@playwright/test';
import { TEST_USER_EMAIL, TEST_USER_PASSWORD } from '../constants';
import { ConsoleTests } from '../types';

export const authTest = base.extend<ConsoleTests>({
  login: async ({ page }, use) => {
    await use(async () => {
      await page.goto('http://localhost:4455/auth/login');
      await page.getByRole('textbox', { name: 'EMAIL' }).click();
      await page.getByRole('textbox', { name: 'EMAIL' }).fill(TEST_USER_EMAIL);
      await page.getByRole('textbox', { name: 'PASSWORD' }).click();
      await page
        .getByRole('textbox', { name: 'PASSWORD' })
        .fill(TEST_USER_PASSWORD);
      await page.getByRole('button', { name: 'LOG IN' }).click();
    });
  },
});

export { expect } from '@playwright/test';

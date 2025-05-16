/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { faker } from '@faker-js/faker';
import { test as base } from '@playwright/test';
import { CONSOLE_ROOT_URL, SIM_POOL_CSV_PATH } from '../constants';
import { ConsoleTests } from '../types';

export const onboardingTest = base.extend<ConsoleTests>({
  onboarding: async ({ page }, use) => {
    await use(async () => {
      await page.waitForTimeout(1000);
      await page.goto(`${CONSOLE_ROOT_URL}/welcome`);
      await page.getByRole('button', { name: 'Next' }).click();
      await page
        .getByRole('textbox', { name: 'Network name' })
        .fill(`${faker.lorem.word(4)}-network`);
      await page.getByRole('button', { name: 'NAME NETWORK' }).click();
      await page.goto(`${CONSOLE_ROOT_URL}/configure?step=2&flow=onb`);
      await page
        .getByRole('checkbox', { name: 'I have installed my site' })
        .check();
      await page.getByRole('button', { name: 'Next' }).click();
      await page.waitForTimeout(5000);
      while (!(await page.getByRole('button', { name: 'Next' }).isVisible())) {
        await page.getByRole('button', { name: 'Retry' }).click();
        await page.waitForTimeout(5000);
      }
      await page.getByRole('button', { name: 'Next' }).click();
      await page.getByRole('textbox', { name: 'Site name' }).click();
      await page
        .getByRole('textbox', { name: 'Site name' })
        .fill(`${faker.lorem.word(4)}-site`);
      await page.getByRole('button', { name: 'Next' }).click();
      await page.getByRole('combobox', { name: 'SWITCH' }).click();
      await page.locator('li:nth-child(1)').click();
      await page.getByRole('combobox', { name: 'POWER' }).click();
      await page.locator('li:nth-child(1)').click();
      await page.getByRole('combobox', { name: 'BACKHAUL' }).click();
      await page.locator('li:nth-child(1)').click();
      await page.getByRole('button', { name: 'Configure site' }).click();
      await page
        .locator('input[type="file"]')
        .setInputFiles(
          `${process.cwd()}/${SIM_POOL_CSV_PATH}/100Sims_part_1.csv`,
        );

      await page.getByRole('button', { name: 'Upload sims' }).click();
      await page.getByRole('button', { name: 'Continue to console' }).click();
    });
  },
});

export { expect } from '@playwright/test';

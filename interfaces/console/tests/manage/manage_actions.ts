/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { faker } from '@faker-js/faker';
import { test as base } from '@playwright/test';
import { SIM_POOL_CSV_PATH } from '../constants';
import { ConsoleTests } from '../types';

export const manageTest = base.extend<ConsoleTests>({
  createMonthlyPlan: async ({ page }, use) => {
    await use(async () => {
      await page.locator('#manage-btn').click();
      await page.getByRole('link', { name: 'Data plans' }).click();
      await page.getByRole('button', { name: 'CREATE DATA PLAN' }).click();
      await page.getByRole('textbox', { name: 'DATA PLAN NAME' }).click();
      await page
        .getByRole('textbox', { name: 'DATA PLAN NAME' })
        .fill(faker.lorem.word(6));
      await page.getByRole('textbox', { name: 'PRICE' }).click();
      await page
        .getByRole('textbox', { name: 'PRICE' })
        .fill(faker.number.int({ min: 1000, max: 2000 }).toString());
      await page.getByRole('textbox', { name: 'DATA LIMIT' }).click();
      await page
        .getByRole('textbox', { name: 'DATA LIMIT' })
        .fill(faker.number.int({ min: 10, max: 20 }).toString());
      await page.getByRole('combobox', { name: 'DURATION' }).click();
      await page.getByRole('option', { name: 'Month' }).click();
      await page.getByRole('button', { name: 'Save Data Plan' }).click();
    });
  },
  createDailyPlan: async ({ page }, use) => {
    await use(async () => {
      await page.getByRole('button', { name: 'CREATE DATA PLAN' }).click();
      await page.getByRole('textbox', { name: 'DATA PLAN NAME' }).click();
      await page
        .getByRole('textbox', { name: 'DATA PLAN NAME' })
        .fill(faker.lorem.word(6));
      await page.getByRole('textbox', { name: 'PRICE' }).click();
      await page
        .getByRole('textbox', { name: 'PRICE' })
        .fill(faker.number.int({ min: 700, max: 1200 }).toString());
      await page.getByRole('textbox', { name: 'DATA LIMIT' }).click();
      await page
        .getByRole('textbox', { name: 'DATA LIMIT' })
        .fill(faker.number.int({ min: 6, max: 12 }).toString());
      await page.getByRole('combobox', { name: 'DURATION' }).click();
      await page.getByRole('option', { name: 'Day' }).click();
      await page.getByRole('button', { name: 'Save Data Plan' }).click();
    });
  },
  editPlan: async ({ page }, use) => {
    await use(async () => {
      await page.locator('#data-table-action-popover').first().click();
      await page.getByRole('menuitem', { name: 'Edit' }).click();
      await page.getByRole('textbox', { name: 'DATA PLAN NAME' }).click();
      await page
        .getByRole('textbox', { name: 'DATA PLAN NAME' })
        .press('ControlOrMeta+a');
      await page
        .getByRole('textbox', { name: 'DATA PLAN NAME' })
        .fill(faker.lorem.word(6));
      await page.getByRole('button', { name: 'Update Data Plan' }).click();
    });
  },
  uploadSims: async ({ page }, use) => {
    await use(async () => {
      await page.locator('#manage-btn').click();
      await page.getByRole('link', { name: 'SIM pool' }).click();
      await page.getByRole('button', { name: 'CLAIM SIMS' }).click();
      await page
        .locator('input[type="file"]')
        .setInputFiles(
          `${process.cwd()}/${SIM_POOL_CSV_PATH}/100Sims_part_${faker.number.int({ min: 2, max: 50 })}.csv`,
        );
      await page.getByRole('button', { name: 'Claim' }).click();
    });
  },
});

export { expect } from '@playwright/test';

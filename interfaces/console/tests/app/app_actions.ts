import { faker } from '@faker-js/faker';
/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { test as base } from '@playwright/test';
import { log } from 'node:console';
import { ConsoleTests } from '../types';
import { selectRandomOption } from '../utils';

export const appTest = base.extend<ConsoleTests>({
  createNetwork: async ({ page }, use) => {
    await use(async () => {
      await page.getByTestId('create-network-select').click();
      const addButton = page.getByTestId('create-network-add-button');
      if (await addButton.isEnabled()) {
        await addButton.click();
      } else {
        log('Network creation limit reached');
        return;
      }
      await page.getByRole('textbox', { name: 'Network name' }).click();
      await page
        .getByRole('textbox', { name: 'Network name' })
        .fill(faker.lorem.word(5));
      await page.getByRole('button', { name: 'Submit' }).click();
    });
  },
  createSubscriber: async ({ page }, use) => {
    await use(async () => {
      await page.getByRole('link', { name: 'Subscribers' }).click();
      await page.getByRole('button', { name: 'Add Subscriber' }).click();
      await page.getByRole('textbox', { name: 'Name' }).click();
      await page
        .getByRole('textbox', { name: 'Name' })
        .fill(faker.person.fullName());
      await page.getByRole('textbox', { name: 'Email' }).click();
      await page
        .getByRole('textbox', { name: 'Email' })
        .fill(faker.internet.email());
      await page.getByRole('button', { name: 'Next' }).click();
      await page.getByRole('combobox', { name: 'SIM ICCID*' }).click();
      await selectRandomOption(
        page,
        page.getByRole('combobox', { name: 'SIM ICCID*' }),
      );
      await page.waitForSelector('button:not([disabled])', {
        state: 'visible',
      });
      await page.getByRole('button', { name: 'Next' }).click();
      await page.getByLabel('', { exact: true }).click();
      await page.locator(`li:nth-child(1) > .MuiTypography-root`).click();
      await page.waitForSelector('button:not([disabled])', {
        state: 'visible',
      });
      await page.getByRole('button', { name: 'ADD SUBSCRIBER' }).click();
      await page.getByRole('button', { name: 'Close', exact: true }).click();
    });
  },
  topupSubscriberData: async ({ page }, use) => {
    await use(async () => {
      await page.getByRole('link', { name: 'Subscribers' }).click();
      await page
        .locator('table tbody tr:first-child')
        .locator('#data-table-action-popover')
        .click();
      await page.getByText('Top up data').click();
      await page.getByLabel('', { exact: true }).click();
      await page.waitForSelector('li[role="option"]', { state: 'visible' });
      await page.locator('li[role="option"]').first().click();
      await page.waitForSelector('button:not([disabled])', {
        state: 'visible',
      });
      await page.getByRole('button', { name: 'TOP UP' }).click();
    });
  },
  restartNode: async ({ page }, use) => {
    await use(async () => {
      await page.getByRole('link', { name: 'Nodes' }).click();
      await page.locator('table tbody tr:first-child td:first-child a').click();
      await page.getByRole('button', { name: 'Restart' }).click();
      await page.getByRole('button', { name: 'Confirm' }).click();
    });
  },
  nodeRFOff: async ({ page }, use) => {
    await use(async () => {
      await page.getByRole('link', { name: 'Nodes' }).click();
      await page.locator('table tbody tr:first-child td:first-child a').click();
      await page.getByRole('button', { name: 'select merge strategy' }).click();
      await page.getByRole('menuitem', { name: 'Turn RF Off' }).click();
      await page.getByRole('button', { name: 'Turn RF Off' }).click();
      await page.getByRole('button', { name: 'Confirm' }).click();
    });
  },
  nodeRFOn: async ({ page }, use) => {
    await use(async () => {
      await page.getByRole('link', { name: 'Nodes' }).click();
      await page.locator('table tbody tr:first-child td:first-child a').click();
      await page.getByRole('button', { name: 'select merge strategy' }).click();
      await page.getByRole('menuitem', { name: 'Turn RF On' }).click();
      await page.getByRole('button', { name: 'Turn RF On' }).click();
      await page.getByRole('button', { name: 'Confirm' }).click();
    });
  },
});

export { expect } from '@playwright/test';

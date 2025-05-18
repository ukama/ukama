/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { faker } from '@faker-js/faker';
import { test as base } from '@playwright/test';
import { log } from 'node:console';
import { ConsoleTests } from '../../types';
import { selectRandomOption } from '../../utils';

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
  renameNode: async ({ page }, use) => {
    await use(async () => {
      await page.getByRole('link', { name: 'Nodes' }).click();
      await page.locator('table tbody tr:first-child td:first-child a').click();
      await page.getByRole('button', { name: 'Edit NODE' }).click();
      await page.getByRole('textbox', { name: 'NODE NAME' }).click();
      await page
        .getByRole('textbox', { name: 'NODE NAME' })
        .press('ControlOrMeta+a');
      await page
        .getByRole('textbox', { name: 'NODE NAME' })
        .fill(`${faker.lorem.word(5)}-node`);
      await page.getByRole('button', { name: 'Save' }).click();
    });
  },
  renameSite: async ({ page }, use) => {
    await use(async () => {
      await page.getByRole('link', { name: 'Sites' }).click();
      await page.getByRole('button', { name: 'menu' }).click();
      await page.getByRole('menuitem', { name: 'Edit Name' }).click();
      await page.getByRole('textbox', { name: 'Site Name' }).click();
      await page
        .getByRole('textbox', { name: 'Site Name' })
        .press('ControlOrMeta+a');
      await page
        .getByRole('textbox', { name: 'Site Name' })
        .fill(`${faker.lorem.word(4)}-site`);
      await page.getByRole('button', { name: 'Save' }).click();
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
  logout: async ({ page }, use) => {
    await use(async () => {
      await page.locator('#account-settings-btn').click();
      await page.getByRole('link', { name: 'Logout of account' }).click();
    });
  },
});

export { expect } from '@playwright/test';

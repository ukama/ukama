/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import path from 'path';
import { applyPatch } from '../common';

const applyCreateSubscriberPatch = async () => {
  const version = path.basename(__dirname);
  const customReplacements = [
    {
      regex:
        /await page\.getByRole\('link', { name: 'Subscribers' }\)\.click\(\);/g,
      replacement: `await page.waitForURL('**/console/home');\nawait page.getByRole('link', { name: 'Subscribers' }).click();\n await page.waitForURL('**/console/subscribers');`,
    },
    {
      regex:
        /await page\.getByRole\('button', { name: 'Add Subscriber' }\)\.click\(\);/g,
      replacement: `await page.getByRole('button', { name: 'Add Subscriber' }).click();
        await page.waitForSelector('[role="dialog"]', {
          state: 'visible',
          timeout: 30000,
        });`,
    },
    {
      regex:
        /await page\.getByRole\('textbox', { name: 'Name' }\)\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'Name' }).fill(\`\${faker.person.fullName()}\`);`,
    },
    {
      regex:
        /await page\.getByRole\('textbox', { name: 'Email' }\)\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'Email' }).fill(\`\${faker.internet.email()}\`);`,
    },
    {
      regex:
        /await page\.getByRole\('combobox', { name: 'SIM ICCID\*' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '\d+' }\)\.click\(\);/g,
      replacement: `await page.getByRole('combobox', { name: 'SIM ICCID*' }).click();
        await selectRandomOption(
          page,
          page.getByRole('combobox', { name: 'SIM ICCID*' }),
        );`,
    },
    {
      regex:
        /await page\s*\.getByLabel\('', { exact: true }\)\.click\(\);\s*await page\s*\.getByText\('textor - .*'\)\.click\(\);/g,
      replacement: `const dropdown = page.getByRole('combobox');
        await dropdown.waitFor({ state: 'visible' });
        await dropdown.click();
        const options = page.locator('[role="option"]');
        await options.first().waitFor({ state: 'visible', timeout: 30000 });
        await options.first().click();
        await page.waitForSelector('button:not([disabled])', {
          state: 'visible',
        });`,
    },
  ];

  await applyPatch(
    'create-subscriber',
    version,
    'subscriber',
    customReplacements,
  );
};

export default applyCreateSubscriberPatch;

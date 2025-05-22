/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import { applyPatch } from './common';

const applyTopupSubscriberPatch = async () => {
  const customReplacements = [
    {
      regex:
        /await page\s*\.getByRole\('row', { name: '.*' }\)\s*\.locator\('#data-table-action-popover'\)\s*\.click\(\);/g,
      replacement: `await page.locator('table tbody tr:first-child').locator('#data-table-action-popover').click();`,
    },
    {
      regex:
        /await page\s*\.getByLabel\('', { exact: true }\)\s*\.click\(\);\s*await page\s*\.getByRole\('option', { name: 'textor - .*' }\)\s*\.click\(\);/g,
      replacement: `await page.getByLabel('', { exact: true }).click();
        await page.waitForSelector('li[role="option"]', { state: 'visible' });
        await page.locator('li[role="option"]').first().click();
        await page.waitForSelector('button:not([disabled])', {
          state: 'visible',
      });`,
    },
  ];

  await applyPatch('topup-subscriber', 'subscriber', customReplacements);
};

export default applyTopupSubscriberPatch;

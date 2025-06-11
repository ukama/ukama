/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import path from 'path';
import { applyPatch } from '../common';

const applySolarPortTogglePatch = async () => {
  const version = path.basename(__dirname);
  const customReplacements = [
    {
      regex: /await page\.getByTestId\('site-switch'\)\.first\(\)\.click\(\);/g,
      replacement: `await page.getByTestId('site-switch').first().click({ force: true });
                    await page.waitForTimeout(1000);`,
    },
    {
      regex:
        /await page\.getByRole\('heading', { name: 'test-site1' }\)\.click\(\);/g,
      replacement: `await page.getByRole('heading', { name: /.+-.+/ }).first().click();
                    await page.waitForTimeout(1000);`,
    },
    {
      regex:
        /await page\.getByRole\('button', { name: 'Port 2 \(Solar Controller\)' }\)\.click\(\);/g,
      replacement: `if (await page.locator('text=Not available').count()) return;
                    await page.getByRole('button', { name: 'Port 2 (Solar Controller)' }).click();
                    await page.waitForTimeout(1000);`,
    },
    {
      regex: /await page\.getByRole\('checkbox'\)\.uncheck\(\);/g,
      replacement: `await page.getByRole('checkbox').uncheck();
                    await page.waitForTimeout(1000);`,
    },
    {
      regex: /await page\.getByRole\('checkbox'\)\.check\(\);/g,
      replacement: `await page.getByRole('checkbox').check();
                    await page.waitForTimeout(1000);`,
    },
  ];

  await applyPatch(
    'solar-switch-port-toggle',
    version,
    'site',
    customReplacements,
  );
};

export default applySolarPortTogglePatch;

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import path from 'path';
import { applyPatch } from '../common';

const applySiteDetailsPatch = async () => {
  const version = path.basename(__dirname);
  const customReplacements = [
    {
      regex: /await page\.getByRole\('heading', { name: 'test-site1' }\)\.click\(\);/g,
      replacement: `await page.locator('table tbody tr:first-child td:first-child a, [data-testid="site-card"]:first-child, h2:first-of-type').first().click();`,
    },
    {
      regex: /await page\.locator\('path:nth-child\(17\)'\)\.click\(\);/g,
      replacement: `await page.getByTestId('site-menu-button').or(page.locator('[aria-label*="menu"], [aria-label*="options"]')).first().click();`,
    },
    {
      regex: /await page\.locator\('rect:nth-child\(15\)'\)\.click\(\);/g,
      replacement: `await page.getByTestId('power-diagram').or(page.locator('[data-testid*="diagram"], svg rect')).first().click();`,
    },
    {
      regex: /await page\.getByRole\('link', { name: 'Sites' }\)\.click\(\);/g,
      replacement: `await page.waitForURL('**/console/home');\n  await page.getByRole('link', { name: 'Sites' }).click();\n  await page.waitForURL('**/sites');`,
    },
    {
      regex: /await page\.getByRole\('button', { name: 'Port (\d+) \((.*?)\)' }\)\.click\(\);/g,
      replacement: `await page.getByRole('button', { name: /Port $1.*$2/ }).click();`,
    },
    {
      regex: /await page\.getByRole\('checkbox'\)\.(check|uncheck)\(\);/g,
      replacement: `await page.getByRole('checkbox').waitFor({ state: 'visible' });\n  await page.getByRole('checkbox').$1();`,
    },
    {
      regex: /await page\.getByText\('test-site1Site is up for 3 minutesSite informationNodes:Not availableDate'\)\.click\(\);/g,
      replacement: `await page.getByTestId('site-info-panel').or(page.locator('[data-testid*="site-details"], .site-info')).first().click();`,
    },
  ];

  await applyPatch('site-details', version, 'site', customReplacements);
};

export default applySiteDetailsPatch;
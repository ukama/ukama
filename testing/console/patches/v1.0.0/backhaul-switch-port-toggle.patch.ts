/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import path from 'path';
import { applyPatch } from '../common';

const applyBackhaulPortTogglePatch = async () => {
  const version = path.basename(__dirname);
  const customReplacements = [
    {
      regex:
        /await page\.getByText\('HomeSitesNodesSubscribers'\)\.click\(\);/g,
      replacement: `// Navigation menu click removed as it's not needed`,
    },
    {
      regex:
        /await page\.getByRole\('heading', { name: 'cinis-site' }\)\.click\(\);/g,
      replacement: `await page.getByRole('heading', { name: /.*-.*/ }).first().click();`,
    },
    {
      regex: /await page\.locator\('rect:nth-child\(15\)'\)\.click\(\);/g,
      replacement: `await page.locator('rect:nth-child(15)').click();`,
    },
    {
      regex:
        /await page\.getByRole\('button', { name: 'Port 3 \(Backhaul\)' }\)\.click\(\);/g,
      replacement: `await page.getByRole('button', { name: /Port.*Backhaul/i }).click();`,
    },
  ];

  await applyPatch(
    'backhaul-switch-port-toggle',
    version,
    'site',
    customReplacements,
  );
};

export default applyBackhaulPortTogglePatch;

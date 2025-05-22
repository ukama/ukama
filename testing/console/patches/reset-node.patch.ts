/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import { applyPatch } from './common';

const applyResetNodePatch = async () => {
  const customReplacements = [
    {
      regex: /await page\.getByRole\('link', { name: 'uk-.*' }\)\.click\(\);/g,
      replacement: `await page.locator('table tbody tr:first-child td:first-child a').click();`,
    },
  ];

  await applyPatch('reset-node', 'node', customReplacements);
};

export default applyResetNodePatch;

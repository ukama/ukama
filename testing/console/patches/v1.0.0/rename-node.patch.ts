/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import path from 'path';
import { applyPatch } from '../common';

const applyRenameNodePatch = async () => {
  const version = path.basename(__dirname);
  const customReplacements = [
    {
      regex: /await page\.getByRole\('link', { name: 'uk-.*' }\)\.click\(\);/g,
      replacement: `await page.locator('table tbody tr:first-child td:first-child a').click();`,
    },
    {
      regex:
        /await page\.getByRole\('textbox', { name: 'NODE NAME' }\)\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'NODE NAME' }).fill(\`\${faker.lorem.word(5)}-node\`);`,
    },
  ];

  await applyPatch('rename-node', version, 'node', customReplacements);
};

export default applyRenameNodePatch;

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import path from 'path';
import { applyPatch } from '../common';

const applyRenameSitePatch = async () => {
  const version = path.basename(__dirname);
  const customReplacements = [
    {
      regex:
        /await page\s*\.getByRole\('textbox', { name: 'Site Name' }\)\s*\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'Site Name' }).fill(\`\${faker.lorem.word(5)}-site\`);`,
    },
    {
      regex:
        /await page\.getByRole\('textbox', { name: 'NODE NAME' }\)\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'NODE NAME' }).fill(\`\${faker.lorem.word(5)}-node\`);`,
    },
  ];

  await applyPatch('rename-site', version, 'site', customReplacements);
};

export default applyRenameSitePatch;

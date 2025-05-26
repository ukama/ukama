/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import path from 'path';
import { applyPatch } from '../common';

const applyCreateNetworkPatch = async () => {
  const version = path.basename(__dirname);
  const customReplacements = [
    {
      regex:
        /await page\.getByRole\('combobox', { name: '.*-network' }\)\.click\(\);/g,
      replacement: `await page.waitForURL('**/console/home');\nawait page.getByTestId("create-network-select").click();`,
    },
    {
      regex: /'test-network'/g,
      replacement: 'faker.lorem.word(5)',
    },
    {
      regex:
        /await page\.getByTestId\('create-network-add-button'\)\.click\(\);/,
      replacement: `const addButton = page.getByTestId('create-network-add-button');
          if (await addButton.isEnabled()) {
            await addButton.click();
          } else {
            console.log('Network creation limit reached');
            return;
          }`,
    },
  ];

  await applyPatch(`create-network`, version, 'network', customReplacements);
};

export default applyCreateNetworkPatch;

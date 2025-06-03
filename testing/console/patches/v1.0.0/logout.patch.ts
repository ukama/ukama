/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import path from 'path';
import { applyPatch } from '../common';

const applyLogoutPatch = async () => {
  const version = path.basename(__dirname);
  const customReplacements = [
    {
      regex:
        /await page\.getByTestId\('account-settings-btn'\)\.click\(\);\s*await page\.getByTestId\('logout-link'\)\.click\(\);/g,
      replacement: `await page.waitForURL('**/console/home');
          await page.getByTestId('account-settings-btn').click();
          await page.waitForSelector('[data-testid="logout-link"]', {
            state: 'visible',
            timeout: 30000,
          });
          await page.getByTestId('logout-link').click();`,
    },
  ];
  await applyPatch('logout', version, 'auth', customReplacements);
};

export default applyLogoutPatch;

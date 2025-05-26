/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import { TMP_SIMS_PATH } from '@/constants';
import { createFakeSimCSV } from '@/helpers';
import { faker } from '@faker-js/faker';
import path from 'path';
import { applyPatch } from '../common';

const applyClaimSimsPatch = async () => {
  const version = path.basename(__dirname);
  const simpoolFileName = `manage-sims-${faker.number.int({ min: 50, max: 100 })}`;
  await Promise.all([
    createFakeSimCSV(
      5,
      `${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFileName}.csv`,
    ),
  ]);

  const customReplacements = [
    {
      regex:
        /await page\.getByTestId\('manage-btn'\)\.click\(\);\s*await page\.getByTestId\('manage-sim'\)\.click\(\);/g,
      replacement: `await page.waitForURL('**/console/home');\n  await page.getByTestId('manage-btn').click();\n  await page.waitForURL('**/manage');\n  await page.getByTestId('manage-sim').click();\n  await page.waitForURL('**/manage/sims');`,
    },
    {
      regex:
        /await page\s*\.locator\('#csv-file-input'\)\s*\.setInputFiles\('.*'\);/g,
      replacement: `await page.locator('#csv-file-input input[type="file"]').setInputFiles('${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFileName}.csv');`,
    },
  ];

  await applyPatch('claim-sims', version, 'manage', customReplacements);
};

export default applyClaimSimsPatch;

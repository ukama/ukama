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

const applyDataPlanCreationPatch = async () => {
  const version = path.basename(__dirname);
  const simpoolFileName = `manage-sims-${faker.number.int({ min: 50, max: 100 })}`;
  await Promise.all([
    createFakeSimCSV(
      2,
      `${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFileName}.csv`,
    ),
  ]);

  const customReplacements = [
    {
      regex:
        /await page\.getByTestId\('manage-btn'\)\.click\(\);\s*await page\.getByTestId\('manage-data-plan'\)\.click\(\);/g,
      replacement: `await page.waitForURL('**/console/home');\n  await page.getByTestId('manage-btn').click();\n  await page.waitForURL('**/manage');\n  await page.getByTestId('manage-data-plan').click();\n  await page.waitForURL('**/manage/data-plans');`,
    },
    {
      regex:
        /await page\.getByRole\('textbox', { name: 'DATA PLAN NAME' }\)\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'DATA PLAN NAME' }).fill(\`\${faker.lorem.word(6)}-plan\`);`,
    },
    {
      regex:
        /await page\.getByRole\('textbox', { name: 'PRICE' }\)\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'PRICE' }).fill(\`\${faker.number.int({ min: 1000, max: 2000 })}\`);`,
    },
    {
      regex:
        /await page\.getByRole\('textbox', { name: 'DATA LIMIT' }\)\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'DATA LIMIT' }).fill(\`\${faker.number.int({ min: 10, max: 20 })}\`);`,
    },
    {
      regex:
        /await page\.getByRole\('combobox', { name: 'DURATION' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '.*' }\)\.click\(\);/g,
      replacement: `await page.getByRole('combobox', { name: 'DURATION' }).click();
        await page.getByRole('option', { name: '${['Month', 'Day'][Math.floor(Math.random() * 2)]}' }).click();`,
    },
  ];

  await applyPatch('dataplan-creation', version, 'manage', customReplacements);
};

export default applyDataPlanCreationPatch;

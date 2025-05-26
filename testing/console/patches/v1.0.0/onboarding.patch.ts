/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import { CONSOLE_ROOT_URL, TMP_SIMS_PATH } from '@/constants';
import { createFakeSimCSV } from '@/helpers';
import { faker } from '@faker-js/faker';
import path from 'path';
import { applyPatch } from '../common';

const applyOnboardingPatch = async () => {
  const version = path.basename(__dirname);
  const simpoolFileName = `onboarding-sims-${faker.number.int({ min: 50, max: 100 })}`;
  await Promise.all([
    createFakeSimCSV(
      5,
      `${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFileName}.csv`,
    ),
  ]);
  const customReplacements = [
    {
      regex: /'http:\/\/localhost:3000\/welcome'/g,
      replacement: '`${CONSOLE_ROOT_URL}/welcome`',
    },
    {
      regex:
        /await page\s*\.getByRole\('textbox', { name: 'Network name' }\)\s*\.fill\s*\(\s*faker\.lorem\.word\(\d+\)\s*\)\s*;/g,
      replacement: `await page.getByRole('textbox', { name: 'Network name' }).fill(\`\${faker.lorem.word(4)}-network\`);`,
    },
    {
      regex:
        /http:\/\/localhost:3000\/configure\?step=2&step=1&networkid=[a-f0-9-]+&flow=onb/,
      replacement: `${CONSOLE_ROOT_URL}/configure?step=2&flow=onb`,
    },
    {
      regex:
        /await page\.getByRole\('button', { name: 'Retry' }\)\.click\(\);/g,
      replacement: `await page.waitForTimeout(5000);
        while (!(await page.getByRole('button', { name: 'Next' }).isVisible())) {
          await page.getByRole('button', { name: 'Retry' }).click();
          await page.waitForTimeout(5000);
        }`,
    },
    {
      regex:
        /await page\.getByRole\('textbox', { name: 'Site name' }\)\.fill\('.*'\);/g,
      replacement: `await page.getByRole('textbox', { name: 'Site name' }).fill(\`\${faker.lorem.word(5)}-site\`);`,
    },
    {
      regex:
        /await page\.getByRole\('combobox', { name: 'SWITCH' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '.*' }\)\.click\(\);/g,
      replacement: `await page.getByRole('combobox', { name: 'SWITCH' }).click();
        await page.locator('li:nth-child(1)').click();`,
    },
    {
      regex:
        /await page\.getByRole\('combobox', { name: 'POWER' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '.*' }\)\.click\(\);/g,
      replacement: `await page.getByRole('combobox', { name: 'POWER' }).click();
        await page.locator('li:nth-child(1)').click();`,
    },
    {
      regex:
        /await page\.getByRole\('combobox', { name: 'BACKHAUL' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '.*' }\)\.click\(\);/g,
      replacement: `await page.getByRole('combobox', { name: 'BACKHAUL' }).click();
        await page.locator('li:nth-child(1)').click();`,
    },
    {
      regex:
        /await page\s*\.locator\('#csv-file-input-onboarding'\)\s*\.setInputFiles\('.*'\);/g,
      replacement: `await page.locator('#csv-file-input-onboarding input[type="file"]').setInputFiles('${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFileName}.csv');`,
    },
  ];
  await applyPatch('onboarding', version, 'onboarding', customReplacements);
};

export default applyOnboardingPatch;

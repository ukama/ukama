/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import { faker } from '@faker-js/faker';
import { test } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';
import { CONSOLE_ROOT_URL, TMP_SIMS_PATH } from '../constants';
import { createFakeSimCSV } from '../helpers';

const PATTERNS = {
  DRAG_DROP: 'Drag & Drop Or Choose file to upload',
  NETWORK_SELECT: 'create-network-select',
  FILE_UPLOAD: 'input[type="file"]',
} as const;

test('apply patch to generated test', async () => {
  try {
    const generatedTestPath = path.join(
      __dirname,
      '../tests/generated_test.ts',
    );
    const patchedTestPath = path.join(
      __dirname,
      '../tests/patched_test.spec.ts',
    );

    const simpoolFile1 = `Sims_${faker.number.int({ min: 1, max: 10 })}`;
    const simpoolFile2 = `Sims_${faker.number.int({ min: 10, max: 50 })}`;
    const simpoolFile3 = `Sims_${faker.number.int({ min: 50, max: 100 })}`;

    await Promise.all([
      createFakeSimCSV(
        2,
        `${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFile1}.csv`,
      ),
      createFakeSimCSV(
        2,
        `${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFile2}.csv`,
      ),
      createFakeSimCSV(
        2,
        `${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFile3}.csv`,
      ),
    ]);

    // Read and process the original content
    const originalContent = fs.readFileSync(generatedTestPath, 'utf8');
    const importStatement =
      "import { TEST_USER_EMAIL, TEST_USER_PASSWORD, LOGIN_URL, CONSOLE_ROOT_URL } from '../constants';\nimport { faker } from '@faker-js/faker';\nimport { selectRandomOption } from '../utils';\n";

    const lines = originalContent.split('\n');
    const insertIndex = lines.findIndex((line) => line.includes('import')) + 1;
    lines.splice(insertIndex, 0, importStatement);
    const contentWithImport = lines.join('\n');

    const patchedContent = contentWithImport
      .replace(/'http:\/\/localhost:4455\/auth\/login'/g, 'LOGIN_URL')
      .replace(/'admin@ukama\.com'/g, 'TEST_USER_EMAIL')
      .replace(/'@Pass2025\.'/g, 'TEST_USER_PASSWORD')
      .replace(
        /await page\.getByRole\('combobox', { name: '.*-network' }\)\.click\(\);/g,
        `await page.getByTestId("${PATTERNS.NETWORK_SELECT}").click();`,
      )
      .replace(/'test-network'/g, 'faker.lorem.word(5)')
      .replace(
        /await page\.getByText\('Drag & Drop Or Choose file to'\).click\(\);/g,
        `await page.locator('${PATTERNS.FILE_UPLOAD}').setInputFiles('${process.cwd()}/${TMP_SIMS_PATH}/${simpoolFile1}.csv');`,
      )
      .replace(
        /await page\.getByRole\('textbox', { name: 'DATA PLAN NAME' }\)\.fill\('.*'\);/g,
        `await page.getByRole('textbox', { name: 'DATA PLAN NAME' }).fill(\`\${faker.lorem.word(6)}-plan\`);`,
      )
      .replace(
        /await page\.getByRole\('textbox', { name: 'PRICE' }\)\.fill\('.*'\);/g,
        `await page.getByRole('textbox', { name: 'PRICE' }).fill(\`\${faker.number.int({ min: 1000, max: 2000 })}\`);`,
      )
      .replace(
        /await page\.getByRole\('textbox', { name: 'DATA LIMIT' }\)\.fill\('.*'\);/g,
        `await page.getByRole('textbox', { name: 'DATA LIMIT' }).fill(\`\${faker.number.int({ min: 10, max: 20 })}\`);`,
      )
      .replace(
        /await page\.getByRole\('combobox', { name: 'DURATION' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '.*' }\)\.click\(\);/g,
        `await page.getByRole('combobox', { name: 'DURATION' }).click();
        await page.getByRole('option', { name: '${['Month', 'Day'][Math.floor(Math.random() * 2)]}' }).click();`,
      )
      .replace(
        /await page\.getByRole\('link', { name: 'uk-.*' }\)\.click\(\);/g,
        `await page.locator('table tbody tr:first-child td:first-child a').click();`,
      )
      .replace(
        /await page\s*\.getByRole\('textbox', { name: 'NODE NAME' }\)\s*\.fill\('.*'\);/g,
        `await page.getByRole('textbox', { name: 'NODE NAME' }).fill(\`\${faker.lorem.word(5)}-node\`);`,
      )
      .replace(
        /await page\s*\.getByRole\('textbox', { name: 'Site Name' }\)\s*\.fill\('.*'\);/g,
        `await page.getByRole('textbox', { name: 'Site Name' }).fill(\`\${faker.lorem.word(5)}-site\`);`,
      )
      .replace(
        /await page\.getByRole\('textbox', { name: 'Name' }\)\.fill\('.*'\);/g,
        `await page.getByRole('textbox', { name: 'Name' }).fill(\`\${faker.person.fullName()}\`);`,
      )
      .replace(
        /await page\.getByRole\('textbox', { name: 'Email' }\)\.fill\('.*'\);/g,
        `await page.getByRole('textbox', { name: 'Email' }).fill(\`\${faker.internet.email()}\`);`,
      )
      .replace(
        /await page\.getByRole\('combobox', { name: 'SIM ICCID\*' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '\d+' }\)\.click\(\);/g,
        `await page.getByRole('combobox', { name: 'SIM ICCID*' }).click();
      await selectRandomOption(
        page,
        page.getByRole('combobox', { name: 'SIM ICCID*' }),
      );`,
      )
      .replace(
        /await page\s*\.getByLabel\('', { exact: true }\)\.click\(\);\s*await page\s*\.getByText\('textor - .*'\)\.click\(\);/g,
        `await page.getByLabel('', { exact: true }).click();
        await page.waitForSelector('li[role="option"]', { state: 'visible' });
        await page.locator('li[role="option"]').first().click();
        await page.waitForSelector('button:not([disabled])', {
          state: 'visible',
        });`,
      )
      .replace(
        /await page\s*\.getByRole\('row', { name: '.*' }\)\s*\.locator\('#data-table-action-popover'\)\s*\.click\(\);/g,
        `await page.locator('table tbody tr:first-child').locator('#data-table-action-popover').click();`,
      )
      .replace(
        /await page\s*\.getByLabel\('', { exact: true }\)\s*\.click\(\);\s*await page\s*\.getByRole\('option', { name: 'textor - .*' }\)\s*\.click\(\);/g,
        `await page.getByLabel('', { exact: true }).click();
        await page.waitForSelector('li[role="option"]', { state: 'visible' });
        await page.locator('li[role="option"]').first().click();
        await page.waitForSelector('button:not([disabled])', {
          state: 'visible',
        });`,
      )
      .replace(
        /'http:\/\/localhost:3000\/welcome'/g,
        '`${CONSOLE_ROOT_URL}/welcome`',
      )
      .replace(
        /await page\s*\.getByRole\('textbox', { name: 'Network name' }\)\s*\.fill\s*\(\s*faker\.lorem\.word\(\d+\)\s*\)\s*;/g,
        `await page.getByRole('textbox', { name: 'Network name' }).fill(\`\${faker.lorem.word(4)}-network\`);`,
      )
      .replace(
        /http:\/\/localhost:3000\/configure\?step=2&step=1&networkid=[a-f0-9-]+&flow=onb/,
        `${CONSOLE_ROOT_URL}/configure?step=2&flow=onb`,
      )
      .replace(
        /await page\.getByRole\('button', { name: 'Retry' }\)\.click\(\);/g,
        `await page.waitForTimeout(5000);
          while (!(await page.getByRole('button', { name: 'Next' }).isVisible())) {
            await page.getByRole('button', { name: 'Retry' }).click();
            await page.waitForTimeout(5000);
          }`,
      )
      .replace(
        /await page\.getByRole\('textbox', { name: 'Site name' }\)\.fill\('.*'\);/g,
        `await page.getByRole('textbox', { name: 'Site name' }).fill(\`\${faker.lorem.word(5)}-site\`);`,
      )
      .replace(
        /await page\.getByRole\('combobox', { name: 'SWITCH' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '.*' }\)\.click\(\);/g,
        `await page.getByRole('combobox', { name: 'SWITCH' }).click();
      await page.locator('li:nth-child(1)').click();`,
      )
      .replace(
        /await page\.getByRole\('combobox', { name: 'POWER' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '.*' }\)\.click\(\);/g,
        `await page.getByRole('combobox', { name: 'POWER' }).click();
      await page.locator('li:nth-child(1)').click();`,
      )
      .replace(
        /await page\.getByRole\('combobox', { name: 'BACKHAUL' }\)\.click\(\);\s*await page\.getByRole\('option', { name: '.*' }\)\.click\(\);/g,
        `await page.getByRole('combobox', { name: 'BACKHAUL' }).click();
      await page.locator('li:nth-child(1)').click();`,
      );

    const finalContent = patchedContent
      .split('\n')
      .filter((line) => !line.includes(PATTERNS.DRAG_DROP))
      .join('\n');

    fs.writeFileSync(patchedTestPath, finalContent);
  } catch (error) {
    console.error('Error while patching test file:', error);
    throw error;
  }
});

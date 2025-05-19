/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
import { Locator, Page } from '@playwright/test';

export async function selectRandomOption(page: Page, locator: Locator) {
  await page.waitForSelector('li[role="option"].MuiAutocomplete-option', {
    state: 'visible',
  });
  const options = await page
    .locator('li[role="option"].MuiAutocomplete-option')
    .all();
  if (options.length > 0) {
    const randomIndex = Math.floor(Math.random() * options.length);
    await options[randomIndex].click();
    await page.waitForSelector('button:not([disabled])', {
      state: 'visible',
    });
  } else {
    console.error('No options found in the dropdown');
  }
}

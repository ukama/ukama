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
  } else {
    console.error('No options found in the dropdown');
  }
}

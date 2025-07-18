// -----------------------------------------------------------------------------
// :warning: This file was automatically generated by Playwright AI Codegen.
// DO NOT MODIFY THIS FILE DIRECTLY.
// Any changes will be overwritten the next time the code is regenerated.
// -----------------------------------------------------------------------------
import { test } from '@playwright/test';

test('Node RF On Test', async ({ page }) => {
  await page.goto('http://localhost:4455/auth/login');
  await page.getByRole('textbox', { name: 'EMAIL' }).click();
  await page.getByRole('textbox', { name: 'EMAIL' }).fill('admin@ukama.com');
  await page.getByRole('textbox', { name: 'EMAIL' }).press('Tab');
  await page.getByRole('textbox', { name: 'PASSWORD' }).fill('@Pass2025.');
  await page.getByRole('button', { name: 'LOG IN' }).click();
  await page.getByRole('link', { name: 'Nodes' }).click();
  await page.getByRole('link', { name: 'uk-sa2450-tnode-v0-4e86' }).click();
  await page.getByRole('button', { name: 'select merge strategy' }).click();
  await page.getByRole('menuitem', { name: 'Turn RF On' }).click();
  await page.getByRole('button', { name: 'Turn RF On' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
});

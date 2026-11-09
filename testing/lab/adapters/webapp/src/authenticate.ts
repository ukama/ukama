/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { mkdir, open, rename, rm } from 'node:fs/promises';
import { dirname } from 'node:path';
import { randomUUID } from 'node:crypto';
import { chromium, firefox, webkit } from 'playwright';
import { assertSession } from './console-app.js';
import { baseURL, Budget, WorkerError } from './contract.js';

export type BrowserName = 'chromium' | 'firefox' | 'webkit';

async function privateJSON(path: string, value: unknown): Promise<void> {
  await mkdir(dirname(path), { recursive: true, mode: 0o700 });
  const temporary = `${path}.${randomUUID()}.tmp`;
  try {
    const file = await open(temporary, 'wx', 0o600);
    try { await file.writeFile(JSON.stringify(value)); } finally { await file.close(); }
    await rename(temporary, path);
  } finally { await rm(temporary, { force: true }); }
}

export function allowedLoginOrigin(consoleOrigin: string, loginOrigin: string, configured?: string): boolean {
  if (configured) return loginOrigin === baseURL(configured);
  const consoleURL = new URL(consoleOrigin), loginURL = new URL(loginOrigin);
  if (loginURL.protocol !== consoleURL.protocol) return false;
  if (loginURL.hostname === consoleURL.hostname) return true;
  const labels = consoleURL.hostname.split('.');
  // The Console and password app may be sibling hosts in the deployment.
  return labels.length >= 4 && loginURL.hostname.endsWith(`.${labels.slice(1).join('.')}`);
}

export async function automaticAuth(origin: string, destination: string,
                                    engine: BrowserName, metadata?: string): Promise<void> {
  const identifier = process.env.UKAMA_IDENTIFIER, password = process.env.UKAMA_PASSWORD;
  if (!identifier || !password)
    throw new WorkerError('AUTH_CREDENTIALS', 'Export UKAMA_IDENTIFIER and UKAMA_PASSWORD before running the batch');
  const executablePath = process.env[`ULAB_WEBAPP_${engine.toUpperCase()}_EXECUTABLE_PATH`] || process.env.ULAB_WEBAPP_EXECUTABLE_PATH;
  const browser = await ({ chromium, firefox, webkit })[engine].launch({
    headless: true, timeout: 30000, ...(executablePath ? { executablePath } : {}),
  });
  try {
    const context = await browser.newContext();
    const page = await context.newPage();
    page.setDefaultTimeout(30000);
    const response = await page.goto(origin, { waitUntil: 'domcontentloaded', timeout: 60000 });
    if (response && response.status() >= 400)
      throw new WorkerError('AUTH_UNAVAILABLE', 'Console sign-in returned an HTTP error');
    const passwordInput = page.locator('input[type="password"]').filter({ visible: true });
    await passwordInput.waitFor({ state: 'visible' });
    const authOrigin = new URL(page.url()).origin;
    if (!allowedLoginOrigin(origin, authOrigin, process.env.ULAB_WEBAPP_AUTH_ORIGIN))
      throw new WorkerError('AUTH_ORIGIN', 'Unexpected sign-in host; configure ULAB_WEBAPP_AUTH_ORIGIN for the intended password app');
    const form = page.locator('form').filter({ has: passwordInput });
    if (await form.count() !== 1)
      throw new WorkerError('AUTH_FORM', 'Expected one visible password sign-in form');
    const identity = form.locator('input[name="password_identifier"], input[name="identifier"], input[name="email"], input[name="username"], input[type="email"], input[autocomplete="username"]').filter({ visible: true });
    if (await identity.count() !== 1)
      throw new WorkerError('AUTH_FORM', 'Expected one visible login or email field');
    let submit = form.locator('button[type="submit"][value="password"]').filter({ visible: true });
    if (!await submit.count()) submit = form.getByRole('button', { name: /^(sign in|log in|login)$/i }).filter({ visible: true });
    if (!await submit.count()) submit = form.locator('button[type="submit"], input[type="submit"]').filter({ visible: true });
    if (await submit.count() !== 1)
      throw new WorkerError('AUTH_FORM', 'Expected one password sign-in submit control');
    // Login has no trace, screenshots or request logging containing credentials.
    await identity.fill(identifier);
    await passwordInput.fill(password);
    await submit.click();
    await page.locator('header.topbar .viewseg').waitFor({ state: 'visible', timeout: 60000 });
    await assertSession(page, origin, new Budget(10000));
    await privateJSON(destination, await context.storageState({ indexedDB: true }));
    if (metadata) await privateJSON(metadata, { auth_origin: authOrigin });
    process.stderr.write('WEBAPP authentication ready\n');
  } catch (error) {
    if (error instanceof WorkerError) throw error;
    // Playwright call logs can include filled values. Never forward them.
    throw new WorkerError('AUTH_FAILED', 'Automatic login failed; check the credentials, login form and Console availability');
  } finally { await browser.close(); }
}

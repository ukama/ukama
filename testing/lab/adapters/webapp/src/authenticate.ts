/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { mkdir, open, rename, rm } from 'node:fs/promises';
import { dirname } from 'node:path';
import { randomUUID } from 'node:crypto';
import { chromium, errors, firefox, webkit, type Browser, type Page } from 'playwright';
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

// Report only known route names. Redirect URLs can contain flow IDs and tokens.
function loginLocation(page: Page | undefined, consoleOrigin: string, loginOrigin: string): string {
  if (!page) return 'unopened';
  if (page.isClosed()) return 'closed';
  const url = new URL(page.url());
  if (url.protocol === 'about:') return 'blank';
  const host = url.origin === consoleOrigin ? 'console' : url.origin === loginOrigin ? 'sign-in' : 'other';
  const route = url.pathname.match(/^\/(auth\/login|login|sign-in|signin|welcome|configure|unauthorized|business|network|customer|error)(?:\/|$)/);
  return host + '/' + (route?.[1] ?? (url.pathname === '/' ? '' : 'other'));
}

export async function automaticAuth(origin: string, destination: string,
                                    engine: BrowserName, metadata?: string): Promise<void> {
  const identifier = process.env.UKAMA_IDENTIFIER, password = process.env.UKAMA_PASSWORD;
  if (!identifier || !password)
    throw new WorkerError('AUTH_CREDENTIALS', 'Export UKAMA_IDENTIFIER and UKAMA_PASSWORD before running the batch');
  const executablePath = process.env[`ULAB_WEBAPP_${engine.toUpperCase()}_EXECUTABLE_PATH`] || process.env.ULAB_WEBAPP_EXECUTABLE_PATH;
  let browser: Browser | undefined, activePage: Page | undefined;
  let step = 'start-browser', authOrigin = '', loginHttp: number | undefined;
  try {
    browser = await ({ chromium, firefox, webkit })[engine].launch({
      headless: process.env.ULAB_WEBAPP_HEADLESS !== 'false', timeout: 30000,
      ...(executablePath ? { executablePath } : {}),
    });
    step = 'create-page';
    const context = await browser.newContext();
    const page = await context.newPage();
    activePage = page;
    page.setDefaultTimeout(30000);
    step = 'open-console';
    const response = await page.goto(origin, { waitUntil: 'domcontentloaded', timeout: 60000 });
    if (response && response.status() >= 400)
      throw new WorkerError('AUTH_UNAVAILABLE', `Console sign-in returned HTTP ${response.status()}`);
    step = 'find-password';
    const passwordInput = page.locator('input[type="password"]').filter({ visible: true });
    await passwordInput.waitFor({ state: 'visible' });
    authOrigin = new URL(page.url()).origin;
    if (!allowedLoginOrigin(origin, authOrigin, process.env.ULAB_WEBAPP_AUTH_ORIGIN))
      throw new WorkerError('AUTH_ORIGIN', 'Unexpected sign-in host; configure ULAB_WEBAPP_AUTH_ORIGIN for the intended password app');
    step = 'find-login-controls';
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
    step = 'fill-identifier';
    await identity.fill(identifier);
    step = 'fill-password';
    await passwordInput.fill(password);
    step = 'submit-login';
    page.on('response', reply => {
      if (reply.request().method() !== 'POST') return;
      const url = new URL(reply.url());
      if (url.protocol !== 'https:' && url.protocol !== 'http:') return;
      if ((url.origin === authOrigin || allowedLoginOrigin(origin, url.origin)) &&
          /\/(login|sign-in|signin)\/?$/.test(url.pathname)) loginHttp = reply.status();
    });
    await submit.click();
    step = 'wait-console-landing';
    const header = page.locator('header.topbar .viewseg');
    const welcome = page.getByRole('heading', { name: 'Welcome to Ukama!', exact: true });
    const landing = await new Budget(60000).poll(async () => {
      const url = new URL(page.url());
      if (url.origin !== origin) return null;
      if (url.pathname === '/welcome') return await welcome.isVisible() ? 'welcome' : null;
      return await header.isVisible() ? 'dashboard' : null;
    }, Boolean, 'Waiting for the Console welcome page or dashboard');
    if (landing === 'welcome') {
      step = 'confirm-welcome';
      await page.getByRole('button', { name: 'Continue', exact: true }).click();
      step = 'wait-dashboard-after-welcome';
      const welcomeError = page.locator('main.welcome-root .welcome-error');
      await header.or(welcomeError).filter({ visible: true }).first().waitFor({ timeout: 60000 });
      if (await welcomeError.isVisible())
        throw new WorkerError('AUTH_WELCOME', 'Console could not save the welcome confirmation');
    }
    step = 'verify-dashboard';
    await assertSession(page, origin, new Budget(10000));
    step = 'save-session';
    await privateJSON(destination, await context.storageState({ indexedDB: true }));
    if (metadata) await privateJSON(metadata, { auth_origin: authOrigin });
    process.stderr.write('WEBAPP authentication ready\n');
  } catch (error) {
    // Playwright exceptions may contain filled values. Keep only their timeout category.
    const code = error instanceof WorkerError ? error.code : 'AUTH_FAILED';
    const reason = error instanceof WorkerError ? error.message :
      error instanceof errors.TimeoutError ? 'Timed out' : 'Operation failed';
    const location = loginLocation(activePage, origin, authOrigin);
    throw new WorkerError(code, `${reason} (step=${step} page=${location} login_post_http=${loginHttp ?? 'unobserved'})`);
  } finally { await browser?.close(); }
}

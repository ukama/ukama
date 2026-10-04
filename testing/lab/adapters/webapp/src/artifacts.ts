/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { appendFile, chmod, mkdir, writeFile } from 'node:fs/promises';
import { join, resolve } from 'node:path';
import type { BrowserContext, Page } from 'playwright';
import { safeURL, type Result } from './contract.js';

export class Artifacts {
  readonly paths: string[] = [];
  readonly directory: string;
  private diagnostics: unknown[] = [];
  private tracing = false;
  private captured = false;
  constructor(root: string, run: string) { this.directory = resolve(root, run); }
  async create(): Promise<void> {
    await mkdir(resolve(this.directory, '..'), { recursive: true, mode: 0o700 });
    // A new run may never overwrite another run's evidence or follow its dir.
    await mkdir(this.directory, { mode: 0o700 });
    this.paths.push(join(this.directory, 'results.jsonl'));
    await writeFile(this.paths[0]!, '', { mode: 0o600, flag: 'wx' });
  }
  private note(value: unknown): void {
    if (this.diagnostics.length === 200) this.diagnostics.shift();
    this.diagnostics.push({ at: new Date().toISOString(), detail: value });
  }
  attach(page: Page): void {
    // Do not record console arguments, headers, bodies or URLs with queries:
    // they can contain session tokens and customer data. Trace is separate,
    // private raw evidence, not a sanitized/shareable report.
    page.on('console', message => this.note({ type: 'console', level: message.type(), source: safeURL(message.location().url) }));
    page.on('pageerror', () => this.note({ type: 'pageerror', message: 'Uncaught browser exception; inspect private trace' }));
    page.on('requestfailed', request => this.note({ type: 'requestfailed', method: request.method(), url: safeURL(request.url()) }));
    page.on('response', response => {
      if (response.status() >= 400) this.note({ type: 'http_error', status: response.status(), url: safeURL(response.url()) });
    });
  }
  async start(context: BrowserContext): Promise<void> {
    await context.tracing.start({ screenshots: true, snapshots: true, sources: false });
    this.tracing = true;
  }
  async record(result: Result): Promise<void> {
    await appendFile(join(this.directory, 'results.jsonl'), `${JSON.stringify(result)}\n`, { mode: 0o600 });
  }
  async summary(status: 'passed' | 'failed', reason: string): Promise<void> {
    const path = join(this.directory, 'worker-summary.json');
    await writeFile(path, JSON.stringify({ status, reason, finished_at: new Date().toISOString(),
      scope: 'worker commands only; not scenario or product coverage', artifacts: this.paths }, null, 2) + '\n', { mode: 0o600 });
    if (!this.paths.includes(path)) this.paths.push(path);
  }
  async capture(page?: Page, context?: BrowserContext): Promise<void> {
    if (this.captured) return;
    this.captured = true;
    if (page && !page.isClosed()) {
      const path = join(this.directory, 'failure.png');
      try {
        await page.screenshot({ path, fullPage: true, timeout: 2500 });
        await chmod(path, 0o600); this.paths.push(path);
      } catch { this.note({ type: 'artifact_error', artifact: 'failure.png' }); }
    }
    if (context && this.tracing) {
      const path = join(this.directory, 'trace.zip');
      try {
        await context.tracing.stop({ path });
        await chmod(path, 0o600); this.paths.push(path);
      } catch { this.note({ type: 'artifact_error', artifact: 'trace.zip' }); }
      this.tracing = false;
    }
    const path = join(this.directory, 'diagnostics.json');
    await writeFile(path, JSON.stringify({ page_url: page ? safeURL(page.url()) : null, events: this.diagnostics }, null, 2) + '\n', { mode: 0o600 });
    this.paths.push(path);
  }
  async discardTrace(context?: BrowserContext): Promise<void> {
    if (context && this.tracing) { this.tracing = false; await context.tracing.stop(); }
  }
}

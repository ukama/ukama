/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { readFile, lstat } from 'node:fs/promises';
import { performance } from 'node:perf_hooks';
import { chromium, firefox, webkit, type Browser, type BrowserContext, type Page } from 'playwright';
import { Artifacts } from './artifacts.js';
import { ConsoleApp, assertSession } from './console-app.js';
import { Budget, WorkerError, canonical, command, keys, profile, safeURL, str, type Command, type Profile, type Result } from './contract.js';

export const BROWSERS = { chromium, firefox, webkit };
export async function loadState(path: string): Promise<NonNullable<Parameters<Browser['newContext']>[0]>['storageState']> {
  try {
    const stat = await lstat(path);
    if (!stat.isFile() || stat.size > 4 * 1024 * 1024) throw new Error();
    const state = JSON.parse(await readFile(path, 'utf8'));
    if (!state || !Array.isArray(state.cookies) || !Array.isArray(state.origins)) throw new Error();
    return state;
  } catch { throw new WorkerError('AUTH_STATE_INVALID', 'Cannot read a valid authentication state file'); }
}
export class Worker {
  failed = false;
  closed = false;
  private run?: string;
  private config?: Profile;
  private browser?: Browser;
  private context?: BrowserContext;
  private page?: Page;
  private app?: ConsoleApp;
  private evidence?: Artifacts;
  private end = Infinity;
  private watchdog?: NodeJS.Timeout;
  private highest = 0;
  private responses = new Map<number, { fingerprint: string; result: Result }>();
  private cleanup?: Promise<void>;
  private busy = false;
  private failureReason?: string;
  onScenarioTimeout?: () => void;

  private base(c?: Command): Result {
    return { protocol: 1, run_id: c?.run_id ?? null, command_id: c?.command_id ?? null, action: c?.action ?? null,
      status: 'ok', run_status: this.failed ? 'failed' : 'running', duration_ms: 0,
      expected: null, actual: null, bindings: [], artifacts: [...(this.evidence?.paths ?? [])] };
  }
  async handle(raw: unknown): Promise<Result> {
    const started = performance.now();
    let c: Command | undefined;
    let fingerprint = '';
    let result = this.base();
    let owned = false;
    let activeBudget: Budget | undefined;
    try {
      c = command(raw); result = this.base(c); fingerprint = canonical(c);
      if (this.busy) throw new WorkerError('CONCURRENT_COMMAND', 'Only one command may execute at a time');
      this.busy = true; owned = true;
      if (this.run && c.run_id !== this.run) throw new WorkerError('RUN_MISMATCH', 'Worker belongs to a different run');
      const previous = this.responses.get(c.command_id);
      if (previous) {
        if (previous.fingerprint !== fingerprint) throw new WorkerError('COMMAND_CONFLICT', 'Command ID was already used with different content');
        return structuredClone(previous.result);
      }
      if (c.command_id <= this.highest) throw new WorkerError('COMMAND_ORDER', 'New command IDs must increase');
      this.highest = c.command_id;
      if (this.closed && c.action !== 'close') throw new WorkerError('WORKER_CLOSED', 'Worker has already closed');
      if (this.failed && c.action !== 'close') throw new WorkerError('RUN_FAILED', 'Run already failed; only close is accepted');
      const remaining = c.deadline_ms - Date.now();
      if (remaining <= 0 || remaining > 900000) throw new WorkerError('DEADLINE_EXCEEDED', 'Command deadline must be in the next 900 seconds');
      if (c.action === 'close') {
        keys(c.inputs, ['failed', 'reason']);
        if (c.inputs.failed !== undefined && typeof c.inputs.failed !== 'boolean')
          throw new WorkerError('INVALID_INPUT', 'close.failed must be a boolean');
        if (c.inputs.reason !== undefined && !/^[A-Z_]{1,64}$/.test(str(c.inputs.reason, 'close.reason')))
          throw new WorkerError('INVALID_INPUT', 'close.reason must be a diagnostic code');
        if (!this.run) throw new WorkerError('NOT_INITIALIZED', 'Cannot finish a run that was never initialized');
        await this.shutdown(c.inputs.failed === true, c.inputs.reason as string | undefined);
        result.actual = { closed: true }; result.run_status = this.failed ? 'failed' : 'passed';
      } else {
        const milliseconds = Math.min(remaining, this.end - performance.now());
        const budget = new Budget(milliseconds); activeBudget = budget;
        // The timer also aborts Playwright calls that have no timeout option.
        let timer: NodeJS.Timeout | undefined;
        const timeout = new Promise<never>((_, reject) => {
          timer = setTimeout(() => reject(new WorkerError('DEADLINE_EXCEEDED', 'Command deadline exceeded')), Math.max(1, milliseconds));
        });
        try {
          const observed = await Promise.race([this.execute(c, budget), timeout]);
          Object.assign(result, observed);
        } finally { clearTimeout(timer); }
      }
    } catch (error) {
      this.failed = true;
      const e = error instanceof WorkerError ? error : new WorkerError('BROWSER_ERROR', 'Browser operation failed; inspect private failure artifacts');
      result.status = 'error'; result.run_status = 'failed';
      result.error = { code: e.code, message: e.message };
      this.failureReason ??= e.code;
      result.expected = c?.inputs.expected ?? c?.inputs.expected_count ?? c?.inputs.available ?? null;
      result.actual = e.actual ?? activeBudget?.actual ?? null;
      // Failure is terminal for browser work. Close the context to stop any
      // command still settling after its deadline; no actions can be retried.
      await this.shutdown(true);
    } finally { if (owned) this.busy = false; }
    result.duration_ms = Math.round(performance.now() - started);
    result.run_status = this.failed ? 'failed' : this.closed ? 'passed' : 'running';
    result.artifacts = [...(this.evidence?.paths ?? [])];
    if (this.evidence) {
      try { await this.evidence.record(result); }
      catch {
        this.failed = true; result.status = 'error'; result.run_status = 'failed';
        result.error = { code: 'ARTIFACT_WRITE_FAILED', message: 'Could not persist the command result' };
        await this.shutdown(true);
      }
    }
    if (c && !this.responses.has(c.command_id)) this.responses.set(c.command_id, { fingerprint, result: structuredClone(result) });
    return result;
  }
  private async execute(c: Command, budget: Budget): Promise<Partial<Result>> {
    if (c.action === 'init') {
      if (this.run) throw new WorkerError('ALREADY_INITIALIZED', 'Worker can only initialize once');
      keys(c.inputs, ['profile', 'artifacts_dir']);
      const config = profile(c.inputs.profile);
      this.run = c.run_id; this.config = config;
      this.end = performance.now() + config.scenario_timeout_seconds * 1000;
      this.watchdog = setTimeout(() => {
        void this.shutdown(true, 'SCENARIO_DEADLINE').finally(() => this.onScenarioTimeout?.());
      }, config.scenario_timeout_seconds * 1000);
      const evidence = new Artifacts(str(c.inputs.artifacts_dir, 'artifacts_dir'), this.run);
      try { await evidence.create(); } catch { throw new WorkerError('ARTIFACT_DIRECTORY', 'Artifact run directory must be new and writable'); }
      this.evidence = evidence;
      const state = await loadState(config.auth_state);
      try {
        this.browser = await BROWSERS[config.browser].launch({ headless: config.headless, timeout: budget.remaining(),
          ...(process.env.ULAB_WEBAPP_EXECUTABLE_PATH ? { executablePath: process.env.ULAB_WEBAPP_EXECUTABLE_PATH } : {}) });
      } catch {
        throw new WorkerError('BROWSER_START_FAILED', 'Cannot start browser; install the Playwright browser/dependencies and check the display for headed mode');
      }
      if (this.closed) { await this.browser.close(); throw new WorkerError('CANCELLED', 'Worker stopped during initialization'); }
      this.context = await this.browser.newContext({ storageState: state, viewport: { width: 1440, height: 1000 }, locale: 'en-US', timezoneId: 'UTC', acceptDownloads: false });
      await evidence.start(this.context);
      this.page = await this.context.newPage(); evidence.attach(this.page);
      this.app = new ConsoleApp(this.page, config.base_url);
      const response = await this.page.goto(config.base_url, { waitUntil: 'domcontentloaded', timeout: budget.remaining() });
      if (response && response.status() >= 400) throw new WorkerError('APP_UNAVAILABLE', 'Console initial navigation returned an HTTP error', response.status());
      await assertSession(this.page, config.base_url, budget);
      return { actual: { authenticated: true, browser: config.browser, browser_version: this.browser.version(), page_url: safeURL(this.page.url()) } };
    }
    if (!this.app || !this.config) throw new WorkerError('NOT_INITIALIZED', 'init must precede browser commands');
    if (c.action === 'web_open') return { bindings: await this.app.open(c.inputs, budget), actual: { page_url: safeURL(this.page!.url()) } };
    if (c.action === 'web_select_network') {
      keys(c.inputs, ['network_name']); await this.app.selectNetwork(str(c.inputs.network_name, 'network_name'), budget);
      return { actual: { network_name: c.inputs.network_name } };
    }
    if (c.action === 'web_reload') { await this.app.reload(c.inputs, budget); return { actual: { page_url: safeURL(this.page!.url()) } }; }
    return this.app.check(c.action, c.inputs, budget);
  }
  async shutdown(failure: boolean, reason = 'RUN_ABORTED'): Promise<void> {
    if (failure) this.failed = true;
    if (failure) this.failureReason ??= reason;
    if (this.cleanup) {
      await this.cleanup;
      if (failure) {
        try { await this.evidence?.summary('failed', this.failureReason ?? reason); }
        catch { this.failed = true; }
      }
      return;
    }
    this.closed = true; clearTimeout(this.watchdog);
    this.cleanup = (async () => {
      // Separate bounded evidence/cleanup budget, never charged as acceptance.
      let timer: NodeJS.Timeout | undefined;
      const capture = (async () => {
        if (this.failed) await this.evidence?.capture(this.page, this.context);
        else await this.evidence?.discardTrace(this.context);
      })();
      try { await Promise.race([capture, new Promise<void>(r => { timer = setTimeout(r, 5000); })]); }
      catch { this.failed = true; }
      finally { clearTimeout(timer); }
      const bounded = async (operation: () => Promise<unknown>) => {
        let timeout: NodeJS.Timeout | undefined;
        try { await Promise.race([operation(), new Promise<never>((_, reject) => {
          timeout = setTimeout(() => reject(new Error('cleanup timeout')), 2000);
        })]); }
        catch { this.failed = true; this.failureReason ??= 'CLEANUP_FAILED'; }
        finally { clearTimeout(timeout); }
      };
      await bounded(async () => this.context?.close());
      await bounded(async () => this.browser?.close());
      try { await this.evidence?.summary(this.failed ? 'failed' : 'passed', this.failureReason ?? 'EXPLICIT_CLOSE'); }
      catch { this.failed = true; }
    })();
    return this.cleanup;
  }
}

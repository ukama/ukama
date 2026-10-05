/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { mkdir, open, readFile, rename, rm } from 'node:fs/promises';
import { dirname, resolve } from 'node:path';
import { createInterface } from 'node:readline/promises';
import { randomUUID } from 'node:crypto';
import { parseArgs } from 'node:util';
import { chromium } from 'playwright';
import { assertSession } from './console-app.js';
import { Budget, WorkerError, baseURL, profile, type Command } from './contract.js';
import { Worker } from './worker.js';

process.umask(0o077);
const usage = `Ukama-lab local web-app worker (Node 22+)
  node adapters/webapp/dist/cli.js worker
  node adapters/webapp/dist/cli.js auth --base-url http://localhost:3000 --out .auth/owner.json
  node adapters/webapp/dist/cli.js smoke --profile adapters/webapp/profile.example.json \\
    --network lab-network --view network_home --label 'Sites online' --expected '2/2'

worker: JSONL stdin/stdout, one process per scenario. EOF without close is failure.
auth: manual sign-in, save private state after visible landing verification.
  --landing dashboard|welcome|unauthorized (default dashboard)
smoke: one visible KPI/field/action/table check; no provisioning or coverage credit.
  --kind kpi|field|action|table (default kpi)
  --artifacts runs/webapp-worker (default)
All relative paths (including auth_state inside profile) resolve from the working directory.
`;

async function serve(): Promise<void> {
  const worker = new Worker();
  let explicitClose = false;
  let stopping = false;
  const terminate = async (code: number, reason: string) => {
    if (stopping) return; stopping = true;
    const kill = setTimeout(() => process.exit(code), 10000);
    await worker.shutdown(true, reason);
    clearTimeout(kill); process.exit(code);
  };
  worker.onScenarioTimeout = () => { void terminate(1, 'SCENARIO_DEADLINE'); };
  process.once('SIGTERM', () => { void terminate(143, 'SIGTERM'); });
  process.once('SIGINT', () => { void terminate(130, 'SIGINT'); });
  process.stdout.on('error', () => { void terminate(1, 'PIPE_CLOSED'); });
  let buffered = '';
  const respond = async (line: string) => {
    let raw: unknown;
    try { raw = JSON.parse(line); } catch { raw = null; }
    const result = await worker.handle(raw);
    await new Promise<void>((done, reject) => process.stdout.write(`${JSON.stringify(result)}\n`, error => error ? reject(error) : done()));
    if (result.action === 'close' && result.status === 'ok') explicitClose = true;
  };
  try {
    process.stdin.setEncoding('utf8');
    for await (const chunk of process.stdin) {
      buffered += chunk;
      let newline: number;
      while ((newline = buffered.indexOf('\n')) >= 0) {
        const line = buffered.slice(0, newline); buffered = buffered.slice(newline + 1);
        if (Buffer.byteLength(line) > 1024 * 1024) throw new WorkerError('PROTOCOL_SIZE', 'Command exceeds 1 MiB');
        if (line.trim()) await respond(line);
      }
      if (Buffer.byteLength(buffered) > 1024 * 1024) throw new WorkerError('PROTOCOL_SIZE', 'Command exceeds 1 MiB');
    }
    if (buffered.trim()) await respond(buffered);
  } catch {
    process.stderr.write('WEBAPP worker input/output failed\n');
    await worker.shutdown(true);
  } finally {
    if (!explicitClose) {
      process.stderr.write('WEBAPP input ended without an explicit close command\n');
      await worker.shutdown(true, 'INPUT_EOF');
    }
  }
  process.exitCode = worker.failed ? 1 : 0;
}

async function authenticate(args: string[]): Promise<void> {
  const { values } = parseArgs({ args, options: { 'base-url': { type: 'string' }, out: { type: 'string' }, landing: { type: 'string', default: 'dashboard' } }, strict: true });
  const origin = baseURL(values['base-url']);
  if (!['dashboard','welcome','unauthorized'].includes(values.landing!)) throw new WorkerError('INVALID_INPUT', '--landing must be dashboard, welcome or unauthorized');
  if (!values.out) throw new WorkerError('INVALID_INPUT', '--out is required');
  if (!process.stdin.isTTY) throw new WorkerError('INVALID_INPUT', 'Authentication requires an interactive terminal and display');
  const destination = resolve(values.out);
  const browser = await chromium.launch({ headless: false,
    ...(process.env.ULAB_WEBAPP_EXECUTABLE_PATH ? { executablePath: process.env.ULAB_WEBAPP_EXECUTABLE_PATH } : {}) });
  const context = await browser.newContext({ viewport: { width: 1440, height: 1000 } });
  const prompt = createInterface({ input: process.stdin, output: process.stderr });
  const temporary = `${destination}.${randomUUID()}.tmp`;
  try {
    const page = await context.newPage();
    await page.goto(origin, { waitUntil: 'domcontentloaded', timeout: 30000 });
    await prompt.question(`Sign in through the browser. Leave the ${values.landing} screen visible and press Enter here. `);
    if (values.landing === 'dashboard') await assertSession(page, origin, new Budget(10000));
    else {
      const path = values.landing === 'welcome' ? '/welcome' : '/unauthorized';
      if (new URL(page.url()).origin !== origin || new URL(page.url()).pathname !== path) throw new WorkerError('AUTH_REQUIRED', 'Expected console landing is not visible');
      const visible = values.landing === 'welcome' ? page.locator('main.welcome-root').getByRole('heading', {name:'Welcome to Ukama!',exact:true}) : page.getByText("Your account isn't set up for this console", {exact:true});
      await visible.waitFor({state:'visible',timeout:10000});
    }
    const state = await context.storageState({ indexedDB: true });
    await mkdir(dirname(destination), { recursive: true, mode: 0o700 });
    const file = await open(temporary, 'wx', 0o600);
    try { await file.writeFile(JSON.stringify(state)); } finally { await file.close(); }
    await rename(temporary, destination);
    process.stderr.write(`WEBAPP authentication saved to ${destination}\n`);
  } finally {
    prompt.close(); await rm(temporary, { force: true }); await context.close(); await browser.close();
  }
}

async function smoke(args: string[]): Promise<void> {
  const { values } = parseArgs({ args, strict: true, options: {
    profile: { type: 'string' }, network: { type: 'string' }, view: { type: 'string' },
    label: { type: 'string' }, expected: { type: 'string' }, kind: { type: 'string', default: 'kpi' },
    artifacts: { type: 'string', default: 'runs/webapp-worker' },
  } });
  if (!values.profile || !values.view || !values.label || values.expected === undefined)
    throw new WorkerError('INVALID_INPUT', '--profile, --view, --label and --expected are required');
  const config = profile(JSON.parse(await readFile(values.profile, 'utf8')));
  const actions: Record<string, Command['action']> = { kpi: 'web_kpi_equals', field: 'web_field_equals', action: 'web_action_available', table: 'web_table_count_equals' };
  const action = actions[values.kind!];
  if (!action) throw new WorkerError('INVALID_INPUT', 'Unknown check kind');
  if (values.kind === 'action' && !['true', 'false'].includes(values.expected)) throw new WorkerError('INVALID_INPUT', 'Action expected value must be true or false');
  if (values.kind === 'table' && !/^\d+$/.test(values.expected)) throw new WorkerError('INVALID_INPUT', 'Table expected value must be an integer');
  const worker = new Worker(); const run = `smoke-${Date.now()}-${randomUUID().slice(0, 8)}`;
  let id = 0;
  const send = async (type: Command['action'], inputs: Command['inputs']) => {
    const result = await worker.handle({ protocol: 1, run_id: run, command_id: ++id, action: type,
      deadline_ms: Date.now() + 1000 * (type.includes('equals') || type === 'web_action_available' ? config.check_timeout_seconds : config.action_timeout_seconds), inputs });
    process.stdout.write(JSON.stringify(result) + '\n');
    process.stderr.write(`WEBAPP ${type} ${result.status === 'ok' ? 'OK' : 'FAIL'}${result.error ? ` ${result.error.code}: ${result.error.message}` : ''}\n`);
    return result.status === 'ok';
  };
  const stop = async () => { await worker.shutdown(true); process.exit(130); };
  process.once('SIGINT', () => { void stop(); }); process.once('SIGTERM', () => { void stop(); });
  try {
    if (await send('init', { profile: config, artifacts_dir: values.artifacts }) &&
        await send('web_open', { view: values.view, ...(values.network ? { network_name: values.network } : {}) })) {
      const expectation = values.kind === 'table' ? { expected_count: Number(values.expected) } :
        values.kind === 'action' ? { available: values.expected === 'true' } : { expected: values.expected };
      await send(action, { view: values.view, label: values.label, requirement: 'WEB-SMOKE-LOCAL', ...expectation });
    }
    await send('close', {});
  } finally { await worker.shutdown(worker.failed); }
  process.stderr.write(`WEBAPP ${worker.failed ? 'FAIL' : 'PASS'} local smoke; no scenario coverage credit\n`);
  process.exitCode = worker.failed ? 1 : 0;
}

async function main(): Promise<void> {
  const [mode, ...args] = process.argv.slice(2);
  if (!mode || mode === '--help') { process.stdout.write(usage); return; }
  if (mode === 'worker' && args.length === 0) return serve();
  if (mode === 'auth') return authenticate(args);
  if (mode === 'smoke') return smoke(args);
  throw new WorkerError('INVALID_INPUT', 'Unknown command; use --help');
}
main().catch(error => {
  const message = error instanceof WorkerError ? `${error.code}: ${error.message}` : 'Setup failed; check profile, browser installation, endpoint and output paths';
  process.stderr.write(`WEBAPP FAIL ${message}\n`); process.exitCode = 1;
});

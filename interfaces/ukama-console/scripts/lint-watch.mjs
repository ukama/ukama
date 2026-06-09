/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Watches src/ and re-runs ESLint on change, writing results to .logs/lint.log
 * so they can be read without terminal access. Zero-dependency.
 */
import { spawn } from 'node:child_process';
import { watch, mkdirSync, writeFileSync, createWriteStream } from 'node:fs';
import { join } from 'node:path';

const ROOT = process.cwd();
const LOG_DIR = join(ROOT, '.logs');
const LOG = join(LOG_DIR, 'lint.log');
mkdirSync(LOG_DIR, { recursive: true });

let running = false;
let queued = false;
/** @type {NodeJS.Timeout | undefined} */
let debounce;

function runLint() {
  if (running) {
    queued = true;
    return;
  }
  running = true;
  const out = createWriteStream(LOG);
  out.write(`# eslint run ${new Date().toISOString()}\n`);
  const child = spawn('pnpm', ['exec', 'eslint', '.'], { cwd: ROOT });
  child.stdout.pipe(out, { end: false });
  child.stderr.pipe(out, { end: false });
  child.on('close', (code) => {
    out.end(`\n# exit ${code} (${code === 0 ? 'clean' : 'issues found'})\n`);
    console.log(`[lint-watch] eslint exit ${code} → .logs/lint.log`);
    running = false;
    if (queued) {
      queued = false;
      runLint();
    }
  });
  child.on('error', (err) => {
    writeFileSync(LOG, `# lint-watch failed to spawn eslint: ${err.message}\n`);
    running = false;
  });
}

console.log('[lint-watch] watching src/ — results → .logs/lint.log');
runLint();

watch(join(ROOT, 'src'), { recursive: true }, () => {
  clearTimeout(debounce);
  debounce = setTimeout(runLint, 1200);
});

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Verifies every source file starts with the MPL-2.0 license header.
 * Zero-dependency; run via `pnpm check:headers`.
 */
import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, extname } from 'node:path';

const ROOT = process.cwd();
const EXTS = new Set(['.ts', '.tsx', '.js', '.jsx', '.mjs']);
const SKIP_DIRS = new Set(['node_modules', '.next', '.logs', '.git', 'out']);
const MARKER = 'Mozilla Public';

/** @param {string} dir @returns {string[]} */
function walk(dir) {
  /** @type {string[]} */
  const out = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const st = statSync(full);
    if (st.isDirectory()) {
      if (!SKIP_DIRS.has(entry)) out.push(...walk(full));
    } else if (EXTS.has(extname(entry))) {
      out.push(full);
    }
  }
  return out;
}

const missing = walk(ROOT).filter((file) => {
  const head = readFileSync(file, 'utf8').slice(0, 400);
  return !head.includes(MARKER);
});

if (missing.length > 0) {
  console.error('Missing MPL license header in:');
  for (const f of missing) console.error('  ' + f.replace(ROOT + '/', ''));
  process.exit(1);
}
console.log('License headers OK');

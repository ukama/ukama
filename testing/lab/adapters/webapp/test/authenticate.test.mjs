// SPDX-License-Identifier: MPL-2.0
import assert from 'node:assert/strict';
import { spawn } from 'node:child_process';
import { once } from 'node:events';
import http from 'node:http';
import { mkdtemp, readFile, rm, stat } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { allowedLoginOrigin } from '../dist/authenticate.js';

const cli = new URL('../dist/cli.js', import.meta.url);
const identifier = 'fixture-user@example.test', password = 'fixture-password-do-not-log';

function run(args, environment = {}) {
  const process = spawn(globalThis.process.execPath, [cli.pathname, ...args], {
    env: {...globalThis.process.env, ULAB_WEBAPP_AUTH_ORIGIN: '', UKAMA_IDENTIFIER: identifier, UKAMA_PASSWORD: password, ...environment},
    stdio: ['ignore', 'pipe', 'pipe'],
  });
  let output = '';
  process.stdout.on('data', data => { output += data; });
  process.stderr.on('data', data => { output += data; });
  return once(process, 'close').then(([code]) => ({code, output}));
}

test('login origins stay within the deployment or an explicitly configured origin', () => {
  assert.equal(allowedLoginOrigin('https://app.udev.ukama.com', 'https://auth.udev.ukama.com'), true);
  assert.equal(allowedLoginOrigin('https://app.udev.ukama.com', 'https://auth.attacker.test'), false);
  assert.equal(allowedLoginOrigin('https://app.udev.ukama.com', 'http://auth.udev.ukama.com'), false);
  assert.equal(allowedLoginOrigin('https://app.udev.ukama.com', 'https://auth.other.test', 'https://auth.other.test'), true);
});

test('automatic auth needs credentials and does not open an interactive prompt', async () => {
  const result = await run(['auth', '--automatic', '--base-url', 'http://localhost:3000', '--out', '/unused/state.json'], {UKAMA_PASSWORD: ''});
  assert.equal(result.code, 1);
  assert.match(result.output, /AUTH_CREDENTIALS/);
  assert.doesNotMatch(result.output, /press Enter|interactive terminal/);
});

test('password login saves a private real browser session after dashboard verification', async t => {
  const directory = await mkdtemp(path.join(tmpdir(), 'ulab-auth-'));
  t.after(() => rm(directory, {recursive: true, force: true}));
  let submissions = 0;
  const server = http.createServer(async (request, response) => {
    if (request.url === '/') { response.writeHead(302, {Location: '/auth/login'}); response.end(); return; }
    if (request.url === '/auth/login') {
      response.setHeader('Content-Type', 'text/html');
      response.end('<form method="post" action="/login"><label>Email<input name="identifier" type="email"></label><label>Password<input name="password" type="password"></label><button type="submit" name="method" value="password">Sign in</button></form>');
      return;
    }
    if (request.url === '/login' && request.method === 'POST') {
      let body = ''; for await (const chunk of request) body += chunk;
      const fields = new URLSearchParams(body);
      assert.equal(fields.get('identifier'), identifier);
      assert.equal(fields.get('password'), password);
      submissions++;
      response.writeHead(302, {Location: '/business', 'Set-Cookie': 'ukama_session=fixture-session; Path=/; HttpOnly; SameSite=Lax'});
      response.end(); return;
    }
    response.setHeader('Content-Type', 'text/html');
    response.end('<header class="topbar"><div class="viewseg">Business</div></header><aside class="sidebar">Navigation</aside><main class="main">Dashboard</main>');
  });
  server.listen(0, '127.0.0.1'); await once(server, 'listening');
  t.after(() => { server.closeAllConnections(); server.close(); });
  const origin = `http://127.0.0.1:${server.address().port}`;
  const destination = path.join(directory, 'auth/owner.json'), metadata = path.join(directory, 'profile.json');
  const result = await run(['auth', '--automatic', '--base-url', origin, '--out', destination, '--metadata', metadata]);
  assert.equal(result.code, 0, result.output);
  assert.equal(submissions, 1);
  assert.equal((await stat(destination)).mode & 0o777, 0o600);
  assert.equal((await stat(path.dirname(destination))).mode & 0o777, 0o700);
  const state = JSON.parse(await readFile(destination));
  assert.equal(state.cookies.find(c => c.name === 'ukama_session')?.value, 'fixture-session');
  assert.equal(JSON.parse(await readFile(metadata)).auth_origin, origin);
  assert.ok(!result.output.includes(identifier) && !result.output.includes(password));
  assert.doesNotMatch(result.output, /press Enter/);
});

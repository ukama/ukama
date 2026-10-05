/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { appendFile, chmod, mkdir, rename, rm, writeFile } from 'node:fs/promises';
import { randomUUID } from 'node:crypto';
import { join, resolve } from 'node:path';
import type { BrowserContext, Page, Request } from 'playwright';
import { safeURL, type Result, type Command } from './contract.js';

export class Artifacts {
  readonly paths: string[] = [];
  readonly directory: string;
  private diagnostics: unknown[] = [];
  private step?: {command_id:number; action:string};
  private dropped = 0;
  private problems = new WeakMap<Page,number>();
  private requestSteps = new WeakMap<Request,unknown>();
  begin(command:Command):void { this.step={command_id:command.command_id,action:command.action}; }
  end():void { this.step=undefined; }
  health(page:Page):string { return (this.problems.get(page)??0)===0?'clean':'errors'; }
  private tracing = new Map<BrowserContext, string>();
  private captured = false;
  constructor(root: string, run: string) { this.directory = resolve(root, run); }
  async create(): Promise<void> {
    await mkdir(resolve(this.directory, '..'), { recursive: true, mode: 0o700 });
    // A new run may never overwrite another run's evidence or follow its dir.
    await mkdir(this.directory, { mode: 0o700 });
    this.paths.push(join(this.directory, 'results.jsonl'));
    await writeFile(this.paths[0]!, '', { mode: 0o600, flag: 'wx' });
  }
  private note(value: unknown, initiated?:unknown): void {
    if (this.diagnostics.length === 200) { this.diagnostics.shift(); this.dropped++; }
    this.diagnostics.push({ at: new Date().toISOString(), command: this.step??null, initiated_by: initiated??null, detail: value });
  }
  attach(page: Page): void {
    // Do not record console arguments, headers, bodies or URLs with queries:
    // they can contain session tokens and customer data. Trace is separate,
    // private raw evidence, not a sanitized/shareable report.
    page.on('console', message => this.note({ type: 'console', level: message.type(), source: safeURL(message.location().url) }));
    const problem=()=>this.problems.set(page,(this.problems.get(page)??0)+1);
    const critical=(r:Request)=>['document','fetch','xhr'].includes(r.resourceType());
    page.on('request',request=>this.requestSteps.set(request,this.step??null));
    page.on('pageerror', () => {problem();this.note({ type: 'pageerror', message: 'Uncaught browser exception; inspect private trace' });});
    page.on('requestfailed', request => {if(critical(request))problem();this.note({ type: 'requestfailed', critical:critical(request), method: request.method(), url: safeURL(request.url()) },this.requestSteps.get(request));});
    page.on('response', response => {
      if (response.status() >= 400) {if(critical(response.request()))problem();this.note({ type: 'http_error', critical:critical(response.request()), status: response.status(), url: safeURL(response.url()) },this.requestSteps.get(response.request()));}
    });
  }
  async start(context: BrowserContext): Promise<void> {
    await context.tracing.start({ screenshots: true, snapshots: true, sources: false });
    this.tracing.set(context, this.tracing.size ? 'trace-peer.zip' : 'trace.zip');
  }
  async record(result: Result): Promise<void> {
    // Each command references its own snapshot; asynchronous failures retain
    // both their request's initiating command and the command observing them.
    if(this.diagnostics.length||this.dropped){
      const path=join(this.directory,`diagnostics-${result.command_id??'protocol'}-${randomUUID()}.json`);
      await writeFile(path,JSON.stringify({command:this.step??null,dropped_events:this.dropped,events:this.diagnostics},null,2)+'\n',{mode:0o600});
      this.paths.push(path);result.artifacts.push(path);
      this.diagnostics=[];this.dropped=0;
    }
    await appendFile(join(this.directory, 'results.jsonl'), `${JSON.stringify(result)}\n`, { mode: 0o600 });
  }
  async summary(status: 'passed' | 'failed', reason: string): Promise<void> {
    const path = join(this.directory, 'worker-summary.json');
    // Signal termination can overlap a command's failure shutdown. Never
    // truncate the last complete summary while another caller is exiting.
    const temporary = `${path}.${randomUUID()}.tmp`;
    try {
      await writeFile(temporary, JSON.stringify({ status, reason, finished_at: new Date().toISOString(),
        scope: 'worker commands only; not scenario or product coverage', artifacts: this.paths }, null, 2) + '\n', { mode: 0o600, flag: 'wx' });
      await rename(temporary, path);
    } finally { await rm(temporary, { force: true }); }
    if (!this.paths.includes(path)) this.paths.push(path);
  }
  async capture(page?: Page, _context?: BrowserContext): Promise<void> {
    if (this.captured) return;
    this.captured = true;
    if (page && !page.isClosed()) {
      const path = join(this.directory, 'failure.png');
      try {
        await page.screenshot({ path, fullPage: true, timeout: 2500 });
        await chmod(path, 0o600); this.paths.push(path);
      } catch { this.note({ type: 'artifact_error', artifact: 'failure.png' }); }
    }
    for (const [traced, filename] of this.tracing) {
      const path = join(this.directory, filename);
      try {
        await traced.tracing.stop({ path });
        await chmod(path, 0o600); this.paths.push(path);
      } catch { this.note({ type: 'artifact_error', artifact: filename }); }
    }
    this.tracing.clear();
    const path = join(this.directory, 'diagnostics.json');
    await writeFile(path, JSON.stringify({ command:this.step??null,dropped_events:this.dropped,page_url: page ? safeURL(page.url()) : null, events: this.diagnostics }, null, 2) + '\n', { mode: 0o600 });
    this.paths.push(path);
  }
  async discardTrace(_context?: BrowserContext): Promise<void> {
    for (const traced of this.tracing.keys()) await traced.tracing.stop();
    this.tracing.clear();
  }
}

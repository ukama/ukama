/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { closeSync, fsyncSync, openSync, renameSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import type { Page, Response } from 'playwright';
import { ConsoleApp } from './console-app.js';
import { Budget, WorkerError, keys, object, str, type ObjectValue } from './contract.js';

// A durable receipt establishes ownership, never acceptance. It is updated
// before a submit and on passive UI response metadata, including late failures.
export class Creation {
  private receipt: ObjectValue;
  private pending = new Set<Promise<void>>();
  private error?: Error;
  private submitted = false;
  private requestSeen = false;
  private binding?: ObjectValue;
  private path: string;
  constructor(private page: Page, private inputs: ObjectValue, private kind: 'network' | 'site',
    directory: string, command: number) {
    keys(inputs, kind === 'network' ? ['ref', 'name'] : ['ref', 'name', 'network_name', 'network_id', 'tower_id', 'components']);
    for (const key of kind === 'network' ? ['ref', 'name'] : ['ref', 'name', 'network_name', 'network_id', 'tower_id']) str(inputs[key], key);
    if (!/^[a-z0-9-]{3,40}$/.test(String(inputs.name))) throw new WorkerError('INVALID_INPUT', 'Invalid planned resource name');
    this.path = join(directory, `creation-${command}.json`);
    this.receipt = { command_id: command, kind, ref: inputs.ref, name: inputs.name, state: 'prepared', bindings: [] };
    this.save();
  }
  private save(): void {
    const temp = `${this.path}.tmp`;
    const fd = openSync(temp, 'w', 0o600);
    try { writeFileSync(fd, JSON.stringify(this.receipt)); fsyncSync(fd); } finally { closeSync(fd); }
    renameSync(temp, this.path);
    const dir = openSync(join(this.path, '..'), 'r'); try { fsyncSync(dir); } finally { closeSync(dir); }
  }
  private identify(value: unknown): void {
    if (!value || typeof value !== 'object') return;
    const row = value as ObjectValue;
    if (row.name !== this.inputs.name || typeof row.id !== 'string' || !/^[A-Za-z0-9_-]{1,127}$/.test(row.id)) return;
    if (this.binding && this.binding.id !== row.id) throw new WorkerError('OWNERSHIP_CONFLICT', 'Multiple IDs observed for this creation');
    this.binding = { kind: this.kind, ref: this.inputs.ref, name: this.inputs.name, id: row.id, observed_via: 'ui_result' };
    this.receipt.state = 'identified'; this.receipt.bindings = [this.binding]; this.save();
  }
  private response = (response: Response): void => {
    const work = (async () => {
      if (!this.submitted) return;
      let request: any;
      try { request = response.request().postDataJSON(); } catch { return; }
      const data = request?.variables?.data;
      const mutation = this.kind === 'network' ? 'addNetwork' : 'addSite';
      const matches = typeof request?.query === 'string' && new RegExp(`\\b${mutation}\\s*\\(`).test(request.query) &&
        data?.name === this.inputs.name && (this.kind === 'network' || data?.network_id === this.inputs.network_id);
      // Site creation's UI can confirm an uncertain transport result through
      // its own getSites polling. Scope that metadata to this network/name.
      const poll = this.kind === 'site' && this.requestSeen &&
        /\bgetSites\s*\(/.test(request?.query ?? '') && data?.networkId === this.inputs.network_id;
      if (!matches && !poll) return;
      let body: any; try { body = await response.json(); } catch { return; }
      if (matches) this.identify(body?.data?.[mutation]);
      if (poll) for (const site of body?.data?.getSites?.sites ?? []) this.identify(site);
    })().catch(error => { this.error = error; });
    this.pending.add(work); void work.finally(() => this.pending.delete(work));
  };
  private request = (request: import('playwright').Request): void => {
    if (!this.submitted) return;
    let body: any; try { body = request.postDataJSON(); } catch { return; }
    const data = body?.variables?.data;
    const mutation = this.kind === 'network' ? 'addNetwork' : 'addSite';
    if (new RegExp(`\\b${mutation}\\s*\\(`).test(body?.query ?? '') && data?.name === this.inputs.name &&
        (this.kind === 'network' || data?.network_id === this.inputs.network_id)) this.requestSeen = true;
  };
  private submit(): void { this.submitted = true; this.receipt.state = 'submitted'; this.save(); }
  async run(app: ConsoleApp, budget: Budget): Promise<ObjectValue[]> {
    this.page.on('request', this.request); this.page.on('response', this.response);
    try {
      if (this.kind === 'network') await this.network(app, budget); else await this.site(app, budget);
      await budget.poll(async () => { if (this.error) throw this.error; return this.binding; }, Boolean, 'Created resource ID was not observed');
      await Promise.all(this.pending); if (this.error) throw this.error;
      return [this.binding!];
    } finally {
      await Promise.all(this.pending);
      this.page.off('request', this.request); this.page.off('response', this.response);
    }
  }
  private async network(app: ConsoleApp, budget: Budget): Promise<void> {
    // Leave Manage through the existing visible navigation before opening the switcher.
    await app.networkHome(budget);
    const switcher = this.page.locator('header.topbar button.netswitch');
    await switcher.click({ timeout: budget.remaining() });
    if (await this.page.getByRole('menuitem').filter({ has: this.page.getByText(String(this.inputs.name), { exact: true }) }).count())
      throw new WorkerError('NAME_EXISTS', 'Planned network already exists; refusing to claim it');
    await this.page.getByRole('menuitem', { name: 'Add network', exact: true }).click({ timeout: budget.remaining() });
    const dialog = this.page.getByRole('dialog');
    await dialog.getByPlaceholder('network-name', { exact: true }).fill(String(this.inputs.name), { timeout: budget.remaining() });
    this.submit();
    await dialog.getByRole('button', { name: 'Create network', exact: true }).click({ timeout: budget.remaining() });
    await dialog.waitFor({ state: 'hidden', timeout: budget.remaining() });
    await budget.poll(() => switcher.locator('.nm').innerText({ timeout: budget.remaining() }), value => value === this.inputs.name, 'Created network not visible');
  }
  private async site(app: ConsoleApp, budget: Budget): Promise<void> {
    const name = String(this.inputs.name), network = String(this.inputs.network_name), tower = String(this.inputs.tower_id);
    await app.open({ view: 'network_sites', network_name: network }, budget);
    await budget.poll(async () => !(await this.page.locator('main.main .MuiSkeleton-root:visible').count()), Boolean, 'Sites are still loading');
    if (await this.page.locator('main.main .ecard').filter({ has: this.page.getByText(name, { exact: true }) }).count())
      throw new WorkerError('NAME_EXISTS', 'Planned site already exists; refusing to claim it');
    await app.open({ view: 'network_node_pool' }, budget);
    const row = this.page.getByRole('row').filter({ has: this.page.getByText(tower, { exact: true }) });
    await row.getByRole('button', { name: 'Configure', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('heading', { name: 'Select a network', exact: true }).waitFor({ timeout: budget.remaining() });
    if (new URL(this.page.url()).searchParams.get('nid') !== tower) throw new WorkerError('WRONG_TOWER', 'Configure lost the selected tower; console defect, no mutation submitted');
    await this.page.getByRole('radio').filter({ has: this.page.getByText(network, { exact: true }) }).click({ timeout: budget.remaining() });
    await this.page.getByRole('button', { name: 'Continue', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('checkbox', { name: "I've installed and powered on all my units", exact: true }).check({ timeout: budget.remaining() });
    await this.page.getByRole('button', { name: 'Next', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('heading', { name: 'Name your site', exact: true }).waitFor({ timeout: budget.remaining() });
    await this.page.getByPlaceholder('site-name', { exact: true }).fill(name, { timeout: budget.remaining() });
    await this.page.getByRole('button', { name: 'Name site', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('heading', { name: 'Configure site settings', exact: true }).waitFor({ timeout: budget.remaining() });
    const params = new URL(this.page.url()).searchParams;
    if (params.get('nid') !== tower || params.get('networkid') !== this.inputs.network_id ||
        (await this.page.locator('.cfg-readonly').innerText({ timeout: budget.remaining() })).trim() !== tower)
      throw new WorkerError('WRONG_TOWER', 'Site wizard identity differs from the planned tower/network');
    const components = object(this.inputs.components ?? {}, 'components'); keys(components, ['switch', 'backhaul', 'power']);
    for (const field of ['switch', 'backhaul', 'power']) {
      const select = this.page.locator(`select[name="${field}Id"]`);
      if (components[field]) await select.selectOption({ label: str(components[field], field) }, { timeout: budget.remaining() });
      await budget.poll(() => select.inputValue({ timeout: budget.remaining() }), Boolean, `Select an unambiguous ${field} component or set its visible label in webapp`);
    }
    this.submit();
    await this.page.getByRole('button', { name: 'Create site', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('heading', { name: 'Upload SIMs', exact: true }).waitFor({ timeout: budget.remaining() });
    await this.page.getByRole('button', { name: 'Finish setup', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('button', { name: 'Go to Console', exact: true }).click({ timeout: budget.remaining() });
    await app.open({ view: 'network_sites', network_name: network }, budget);
    const card = this.page.locator('main.main .ecard[role="button"]').filter({ has: this.page.getByText(name, { exact: true }) });
    await card.waitFor({ state: 'visible', timeout: budget.remaining() });
    await budget.poll(async () => { if (this.error) throw this.error; return this.binding; }, Boolean, 'Site receipt missing');
    await card.click({ timeout: budget.remaining() });
    await budget.poll(async () => new URL(this.page.url()).pathname, p => p === `/network/sites/${this.binding!.id}`, 'Created site detail identity differs');
    await this.page.locator('main.main').getByRole('heading', { name, exact: true }).waitFor({ timeout: budget.remaining() });
  }
}

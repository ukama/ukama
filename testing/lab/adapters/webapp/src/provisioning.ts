/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { closeSync, fsyncSync, openSync, renameSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { errors, type Page, type Response } from 'playwright';
import { ConsoleApp } from './console-app.js';
import { Budget, WorkerError, keys, normalize, object, str, type ObjectValue } from './contract.js';

// The name step copies the current address into the settings URL on click.
// The button is already enabled while reverse geocoding is still pending.
export async function advanceSiteName(page: Page, budget: Budget): Promise<void> {
  const field = page.locator('label.ff').filter({ has: page.getByText('Site location', { exact: true }) });
  const resolved = field.locator('.cfg-readonly').filter({
    hasText: /\S/, hasNotText: /^\s*Resolving location(?:…|\.{3})\s*$/i,
  });
  try {
    await resolved.waitFor({ state: 'visible', timeout: Math.max(1, Math.min(60_000, budget.remaining() - 250)) });
  } catch (error) {
    if (!(error instanceof errors.TimeoutError)) throw error;
    throw new WorkerError('SITE_LOCATION_UNAVAILABLE', 'The console did not display a site address or coordinates before naming the site');
  }
  const displayed = normalize(await resolved.innerText({ timeout: budget.remaining() }));
  await page.getByRole('button', { name: 'Name site', exact: true }).click({ timeout: budget.remaining() });
  await page.getByRole('heading', { name: 'Configure site settings', exact: true }).waitFor({ timeout: budget.remaining() });
  const carried = new URL(page.url()).searchParams.get('location');
  if (!carried || normalize(carried) !== displayed)
    throw new WorkerError('SITE_LOCATION_MISMATCH', 'The site settings did not carry the location displayed by the console');
}

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
  begin(): void { this.page.on('request',this.request); this.page.on('response',this.response); }
  markSubmitted(): void { this.submit(); }
  markIntercepted(): void {
    if (this.binding) throw new WorkerError('OWNERSHIP_CONFLICT','Cannot mark an identified creation as intercepted');
    this.receipt.state = 'intercepted'; this.save();
  }
  async finish(budget: Budget): Promise<ObjectValue[]> {
    await budget.poll(async()=>{if(this.error)throw this.error;return this.binding;},Boolean,'Created resource identity was not observed');
    await Promise.all(this.pending); if(this.error)throw this.error;
    this.page.off('request',this.request); this.page.off('response',this.response);
    return [this.binding!];
  }
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
    await app.selectNetwork(String(this.inputs.name), budget);
  }
  private async configureTower(app: ConsoleApp, tower: string, budget: Budget): Promise<void> {
    const main = this.page.locator('main.main');
    const row = main.getByRole('row').filter({ has: this.page.getByRole('cell', { name: tower, exact: true }) });
    const configure = row.getByRole('button', { name: 'Configure', exact: true });
    const end = performance.now() + Math.min(180_000, Math.max(1, budget.remaining() - 500));
    let nextRefresh = performance.now() + 2_000;
    let state = 'tower row missing';

    for (;;) {
      const rows = await row.count();
      if (rows > 1) throw new WorkerError('AMBIGUOUS_LOCATOR', `Node pool has multiple rows for ${tower}`);
      if (await configure.isVisible() && await configure.isEnabled()) break;
      if (await main.locator('.MuiSkeleton-root:visible').count()) state = 'loading node inventory';
      else if (await main.getByText("Couldn't load node pool", { exact: true }).isVisible()) state = 'inventory read failed';
      else if (!rows) state = 'tower row missing';
      else if (await row.getByRole('button', { name: 'View detail', exact: true }).isVisible()) state = 'row shows View detail';
      else state = await configure.isVisible() ? 'Configure disabled' : 'row has no Configure button';
      budget.actual = { tower, node_pool: state };
      if (performance.now() >= end) throw new WorkerError('NODE_POOL_NOT_READY', `Node pool: ${tower}: ${state}`, budget.actual);
      // NodePool uses a cache-first read. Refresh the visible inventory while
      // the newly booted tower registers, before dispatching any Configure click.
      if (performance.now() >= nextRefresh && state !== 'loading node inventory') {
        await app.reload({}, budget);
        nextRefresh = performance.now() + 15_000;
      } else {
        await new Promise(resolve => setTimeout(resolve, Math.min(500, Math.max(1, end - performance.now()))));
      }
    }
    await configure.click({ timeout: Math.min(10_000, budget.remaining()) });
  }
  private async selectSiteComponent(field: string, requested: unknown, budget: Budget): Promise<void> {
    const select = this.page.locator(`select[name="${field}Id"]`);
    await select.waitFor({ state: 'visible', timeout: budget.remaining() });
    // The Console renders these native selects after its inventory query settles.
    const options = await select.locator('option:enabled').evaluateAll(elements =>
      (elements as HTMLOptionElement[]).filter(o => o.value.trim()).map(o => ({
        value: o.value, label: o.label, selected: o.selected,
      })));
    if (!options.length)
      throw new WorkerError('SITE_COMPONENT_UNAVAILABLE', `No selectable ${field} component in the Console inventory`);

    let choice = options.find(o => o.selected) ?? options[0]!;
    if (requested === 'first-named' || requested === undefined || requested === '') {
      const named = options.find(o => o.label.trim() && !/^others?\b/i.test(o.label.trim()));
      if (!named && requested === 'first-named')
        throw new WorkerError('SITE_COMPONENT_UNAVAILABLE',
          `No named ${field} component in the Console dropdown`,
          { field, available_labels: options.map(o => o.label) });
      if (named) choice = named;
    } else {
      const label = str(requested, `${field}_component`);
      const matches = options.filter(o => o.label === label);
      if (matches.length !== 1)
        throw new WorkerError('SITE_COMPONENT_UNAVAILABLE',
          `Expected one ${field} component labeled "${label}", found ${matches.length}`,
          { field, available_labels: options.slice(0, 12).map(o => o.label) });
      choice = matches[0]!;
    }
    if (!choice.selected)
      await select.selectOption({ value: choice.value }, { timeout: budget.remaining() });
    const actual = await select.inputValue({ timeout: budget.remaining() });
    if (actual !== choice.value)
      throw new WorkerError('SITE_COMPONENT_UNAVAILABLE', `Console did not retain the selected ${field} component`, actual);
  }
  private async site(app: ConsoleApp, budget: Budget): Promise<void> {
    const name = String(this.inputs.name), network = String(this.inputs.network_name), tower = String(this.inputs.tower_id);
    await app.open({ view: 'network_sites', network_name: network }, budget);
    await budget.poll(async () => !(await this.page.locator('main.main .MuiSkeleton-root:visible').count()), Boolean, 'Sites are still loading');
    if (await this.page.locator('main.main .ecard').filter({ has: this.page.getByText(name, { exact: true }) }).count())
      throw new WorkerError('NAME_EXISTS', 'Planned site already exists; refusing to claim it');
    await app.open({ view: 'network_node_pool' }, budget);
    await this.configureTower(app, tower, budget);
    await this.page.getByRole('heading', { name: 'Select a network', exact: true }).waitFor({ timeout: budget.remaining() });
    // Some Console versions auto-detect the trio after network selection. In
    // either flow the settings screen must identify our exact tower before submit.
    const selectedTower = new URL(this.page.url()).searchParams.get('nid');
    if (selectedTower !== null && selectedTower !== tower)
      throw new WorkerError('WRONG_TOWER', `Configure lost the selected tower: expected ${tower}, observed ${selectedTower}`);
    await this.page.getByRole('radio').filter({ has: this.page.getByText(network, { exact: true }) }).click({ timeout: budget.remaining() });
    await this.page.getByRole('button', { name: 'Continue', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('checkbox', { name: "I've installed and powered on all my units", exact: true }).check({ timeout: budget.remaining() });
    await this.page.getByRole('button', { name: 'Next', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('heading', { name: 'Name your site', exact: true }).waitFor({ timeout: budget.remaining() });
    await this.page.getByPlaceholder('site-name', { exact: true }).fill(name, { timeout: budget.remaining() });
    await advanceSiteName(this.page, budget);
    const params = new URL(this.page.url()).searchParams;
    const urlTower = params.get('nid'), urlNetwork = params.get('networkid');
    const shown = (await this.page.locator('.cfg-readonly:visible').allInnerTexts()).map(value => value.trim());
    if (urlTower !== tower || urlNetwork !== this.inputs.network_id || shown.length !== 1 || shown[0] !== tower)
      throw new WorkerError('WRONG_TOWER',
        `Site wizard identity differs: expected tower=${tower} network=${this.inputs.network_id}; ` +
        `observed nid=${urlTower ?? '<missing>'} displayed=${shown.join(',') || '<missing>'} networkid=${urlNetwork ?? '<missing>'}`,
        { url_tower: urlTower, displayed: shown, url_network: urlNetwork });
    const components = object(this.inputs.components ?? {}, 'components'); keys(components, ['switch', 'backhaul', 'power']);
    for (const field of ['switch', 'backhaul', 'power'])
      await this.selectSiteComponent(field, components[field], budget);
    this.submit();
    await this.page.getByRole('button', { name: 'Create site', exact: true }).click({ timeout: budget.remaining() });
    const upload = this.page.getByRole('heading', { name: 'Upload SIMs', exact: true });
    const formError = this.page.locator('.cfg-fields .cfg-error:visible, .cfg-fields .ff-err:visible').first();
    await upload.or(formError).first().waitFor({ state: 'visible', timeout: budget.remaining() });
    if (await formError.isVisible()) {
      const message = await formError.innerText({ timeout: budget.remaining() });
      // These Console validation errors return before useCreateSite runs.
      const rejectedByForm = await formError.evaluate(el => el.classList.contains('ff-err')) ||
        message === 'Your node was not found in your inventory (no access component). Please contact support.' ||
        message === 'No spectrum component found in your inventory. Please contact support.';
      if (rejectedByForm && !this.requestSeen && !this.binding) {
        this.receipt.state = 'prepared';
        this.save();
      }
      throw new WorkerError('SITE_CREATE_FAILED', message);
    }
    await this.page.getByRole('button', { name: 'Finish setup', exact: true }).click({ timeout: budget.remaining() });
    await this.page.getByRole('button', { name: 'Go to Console', exact: true }).click({ timeout: budget.remaining() });
    await app.open({ view: 'network_sites', network_name: network }, budget);
    const card = this.page.locator('main.main .ecard[role="button"]').filter({ has: this.page.getByText(name, { exact: true }) });
    await card.waitFor({ state: 'visible', timeout: budget.remaining() });
    await budget.poll(async () => { if (this.error) throw this.error; return this.binding; }, Boolean, 'Site receipt missing');
    await card.click({ timeout: budget.remaining() });
    await budget.poll(async () => new URL(this.page.url()).pathname, p => p === `/network/sites/${this.binding!.id}`, 'Created site detail identity differs');
    await budget.poll(() => this.page.locator('main.main .pagehead .pagetitle').innerText({timeout:budget.remaining()}), value => value.trim() === name, 'Created site title differs from the owned site');
  }
}

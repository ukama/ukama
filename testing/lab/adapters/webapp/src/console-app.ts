/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Commerce } from './commerce.js';
import { Operations } from './operations.js';
import type { Locator, Page } from 'playwright';
import { Budget, WorkerError, bool, integer, keys, normalize, object, str, type ObjectValue } from './contract.js';

type View = { path: string; lens: string; label: string; network: boolean; detail?: 'node' | 'site' };
const view = (path: string, label: string, network = true): View => ({ path, label, network, lens: path.split('/')[1]! });
export const VIEWS: Record<string, View> = {
  business_home: view('/business', 'Home'),
  business_revenue: view('/business/revenue', 'Revenue'),
  business_customers: view('/business/customers', 'Customers'),
  business_packages: view('/business/packages', 'Packages'),
  business_data_plans: view('/business/manage/data-plans', 'Data plans', false),
  business_members: view('/business/manage/members', 'Members', false),
  business_sim_pool: view('/business/manage/sim-pool', 'SIM pool', false),
  business_support: view('/business/support', 'Support'),
  business_settings: view('/business/settings', 'Settings', false),
  network_home: view('/network', 'Home'),
  network_sites: view('/network/sites', 'Sites'),
  network_site_detail: { ...view('/network/sites', 'Sites'), detail: 'site' },
  network_nodes: view('/network/nodes', 'Nodes'),
  network_node_detail: { ...view('/network/nodes', 'Nodes'), detail: 'node' },
  network_customers: view('/network/customers', 'Customers'),
  network_node_pool: view('/network/manage/node-pool', 'Node pool', false),
  network_sim_pool: view('/network/manage/sim-pool', 'SIM pool', false),
  network_support: view('/network/support', 'Support'),
  network_settings: view('/network/settings', 'Settings', false),
  customer_customers: view('/customer/customers', 'Customers'),
  customer_data_plans: view('/customer/data-plans', 'Data plans'),
  customer_settings: view('/customer/settings', 'Settings', false),
};
export function getView(name: unknown): View {
  const result = VIEWS[str(name, 'view')];
  if (!result) throw new WorkerError('UNSUPPORTED_VIEW', 'This view has no Patch 2 navigation handler');
  return result;
}
async function visible(locator: Locator): Promise<Locator | undefined> {
  const candidates = locator.filter({ visible: true });
  const count = await candidates.count();
  if (count > 1) throw new WorkerError('AMBIGUOUS_LOCATOR', 'More than one visible element matches the semantic locator');
  return count === 1 ? candidates : undefined;
}
async function text(locator: Locator, allowEmpty = false): Promise<string | null> {
  const target = await visible(locator);
  if (target) return normalize(await target.innerText({ timeout: 500 }));
  // An empty inline field has no box, so isVisible() is false. It is still a
  // valid empty value when its row is visible and no CSS/hidden ancestor hides
  // it. Missing nodes and hidden non-empty values must never become "".
  if (allowEmpty && await locator.count() === 1 && await locator.innerText({ timeout: 500 }) === '') {
    const rendered = await locator.evaluate(element => {
      for (let e: Element | null = element; e; e = e.parentElement) {
        const style = getComputedStyle(e);
        if (e.hasAttribute('hidden') || style.display === 'none' || style.visibility === 'hidden' || style.visibility === 'collapse') return false;
      }
      return true;
    });
    if (rendered) return '';
  }
  return null;
}
export async function assertSession(page: Page, origin: string, budget: Budget): Promise<void> {
  await budget.poll(async () => {
    const url = new URL(page.url());
    if (url.origin !== origin || /\/(login|sign-in|signin|unauthorized|welcome)(\/|$)/.test(url.pathname))
      throw new WorkerError('AUTH_REQUIRED', 'Saved session did not reach the dashboard; capture authentication again');
    return (await page.locator('header.topbar .viewseg').isVisible()) &&
      (await page.locator('aside.sidebar:visible').count()) === 1 && (await page.locator('main.main').isVisible());
  }, Boolean, 'Authenticated dashboard was not visible');
}

export class ConsoleApp {
  private current?: { name: string; path: string; network?: string };
  constructor(private page: Page, private origin: string) {}
  private main(): Locator { return this.page.locator('main.main'); }
  private async path(path: string, budget: Budget): Promise<void> {
    await budget.poll(async () => {
      const url = new URL(this.page.url());
      if (url.origin !== this.origin) throw new WorkerError('AUTH_REQUIRED', 'Navigation left the console origin');
      return url.pathname;
    }, actual => actual === path, 'Console navigation did not reach the requested view');
    await this.main().waitFor({ state: 'visible', timeout: budget.remaining() });
  }
  private async lens(name: string, budget: Budget): Promise<void> {
    const current = new URL(this.page.url()).pathname.split('/')[1];
    if (current === name) return;
    await this.page.locator('header.topbar .viewseg').getByRole('button', {
      name: name[0]!.toUpperCase() + name.slice(1), exact: true,
    }).click({ timeout: budget.remaining() });
    await this.path(name === 'customer' ? '/customer/customers' : `/${name}`, budget);
  }
  private async sidebar(route: View, budget: Budget): Promise<void> {
    if (new URL(this.page.url()).pathname !== route.path) {
      await this.page.locator('aside.sidebar:visible').getByRole('link', { name: route.label, exact: true }).click({ timeout: budget.remaining() });
    }
    await this.path(route.path, budget);
  }
  async networkHome(budget: Budget): Promise<void> {
    this.current = undefined;
    await assertSession(this.page, this.origin, budget);
    await this.lens('network', budget);
    await this.sidebar(VIEWS.network_home!, budget);
  }
  async selectNetwork(name: string, budget: Budget): Promise<void> {
    this.current = undefined;
    // Manage pages deliberately have no network switcher. Leave them through
    // the visible Home/Customers link before selecting the requested network.
    if (/\/manage(\/|$)/.test(new URL(this.page.url()).pathname)) {
      const lens = new URL(this.page.url()).pathname.split('/')[1]!;
      await this.sidebar(view(`/${lens}`, 'Home'), budget);
    }
    const switcher = this.page.locator('header.topbar button.netswitch');
    await switcher.waitFor({ state: 'visible', timeout: budget.remaining() });
    if (await text(switcher.locator('.nm')) !== name) {
      await switcher.click({ timeout: budget.remaining() });
      // Menu accessible names include "Default network"; match the name text
      // inside the item, rather than a substring of another network's name.
      const option = this.page.getByRole('menuitem').filter({ has: this.page.getByText(name, { exact: true }) });
      await option.click({ timeout: budget.remaining() });
    }
    await budget.poll(() => text(switcher.locator('.nm')), v => v === name, 'Network switcher did not show the requested network');
  }
  async open(inputs: ObjectValue, budget: Budget): Promise<unknown[]> {
    keys(inputs, ['view', 'network_name', 'entity']);
    const name = str(inputs.view, 'view');
    const route = getView(name);
    const network = inputs.network_name === undefined ? undefined : str(inputs.network_name, 'network_name');
    if (route.network && !network) throw new WorkerError('INVALID_INPUT', 'This view requires a resolved network_name');
    let entity: { ref: string; id: string; text: string } | undefined;
    if (route.detail) {
      const raw = object(inputs.entity, 'entity'); keys(raw, ['ref', 'id', 'text']);
      entity = { ref: str(raw.ref, 'entity.ref'), id: str(raw.id, 'entity.id'), text: str(raw.text, 'entity.text') };
      if (!/^[A-Za-z0-9_-]+$/.test(entity.id)) throw new WorkerError('INVALID_INPUT', 'Invalid entity ID');
    } else if (inputs.entity !== undefined) throw new WorkerError('INVALID_INPUT', 'entity is only valid for detail views');
    this.current = undefined;
    await assertSession(this.page, this.origin, budget);
    await this.lens(route.lens, budget);
    if (network) await this.selectNetwork(network, budget);
    await this.sidebar(route, budget);
    let path = route.path;
    const bindings: unknown[] = [];
    if (entity) {
      const escaped = entity.text.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      const cards = this.main().locator('.ecard[role="button"]');
      const card = route.detail === 'node'
        ? cards.filter({ has: this.page.locator('.tnum').filter({ hasText: new RegExp(`(?:^|[\\s·])${escaped}$`) }) })
        : cards.filter({ has: this.page.getByText(entity.text, { exact: true }) });
      await card.click({ timeout: budget.remaining() });
      path += `/${encodeURIComponent(entity.id)}`;
      await this.path(path, budget);
      if (route.detail === 'node') {
        // A client-side URL can change before the list is replaced. Wait for
        // the bound detail's visible serial, not the still-visible list dots.
        const serial = this.main().locator('.kv-row').filter({ has: this.page.getByText('Serial #', { exact: true }) }).locator(':scope > span > .tnum');
        await budget.poll(() => text(serial), value => value === entity!.id, 'Node detail has not rendered the bound serial');
      }
      bindings.push({ ref: entity.ref, kind: route.detail, id: entity.id, observed_via: 'url' });
    }
    this.current = { name, path, network: route.network ? network : undefined };
    return bindings;
  }
  async reload(inputs: ObjectValue, budget: Budget): Promise<void> {
    keys(inputs, []);
    await this.page.reload({ waitUntil: 'domcontentloaded', timeout: budget.remaining() });
    await assertSession(this.page, this.origin, budget);
    if (this.current) await this.assertView(this.current.name);
  }
  private async assertView(name: string): Promise<void> {
    getView(name);
    if (!this.current || this.current.name !== name) throw new WorkerError('WRONG_VIEW', 'Check requires a preceding web_open for this view');
    const url = new URL(this.page.url());
    if (url.origin !== this.origin || url.pathname !== this.current.path)
      throw new WorkerError('WRONG_VIEW', 'Browser is no longer on the expected view');
    if (this.current.network && await text(this.page.locator('header.topbar .netswitch .nm')) !== this.current.network)
      throw new WorkerError('WRONG_NETWORK', 'Visible network selection changed before the check');
    if (!await this.main().isVisible()) throw new WorkerError('WRONG_VIEW', 'Dashboard content is not visible');
  }
  async assertCommerce(inputs: ObjectValue): Promise<void> {
    const name = str(inputs.view, 'view');
    if (!['business_data_plans', 'customer_customers'].includes(name)) throw new WorkerError('INVALID_INPUT', 'Unsupported commerce view');
    await this.assertView(name);
    if (name === 'customer_customers' && this.current!.network !== inputs.network_name) throw new WorkerError('WRONG_NETWORK', 'Commerce network differs from the opened network');
  }
  async operation(inputs: ObjectValue, budget: Budget): Promise<void> {
    keys(inputs, ['view', 'network_name', 'entity', 'action', 'value', 'app', 'tag']);
    const name = str(inputs.view, 'view'), route = getView(name);
    if (!route.detail) throw new WorkerError('INVALID_INPUT', 'Operations require a detail view');
    const entity = object(inputs.entity, 'entity'); keys(entity, ['ref', 'id', 'text']);
    str(entity.ref, 'entity.ref'); str(entity.text, 'entity.text');
    const id = str(entity.id, 'entity.id');
    await this.assertView(name);
    if (this.current!.path !== `${route.path}/${encodeURIComponent(id)}` ||
        this.current!.network !== str(inputs.network_name, 'network_name'))
      throw new WorkerError('WRONG_ENTITY', 'Operation does not match the opened entity and network');
    await new Operations(this.page, route.detail === 'site').run(inputs, budget);
  }
  async check(action: string, inputs: ObjectValue, budget: Budget): Promise<{ expected: unknown; actual: unknown }> {
    const field = action === 'web_table_count_equals' ? 'expected_count' : action === 'web_action_available' ? 'available' : 'expected';
    keys(inputs, action === 'web_commerce_equals' ? ['view', 'label', 'requirement', 'expected', 'plan_name', 'customer_name', 'iccid'] : ['view', 'label', 'requirement', field, 'node_id', 'app', 'match']);
    const name = str(inputs.view, 'view');
    const label = str(inputs.label, 'label');
    if (!/^WEB-[A-Z0-9][A-Z0-9-]*$/.test(str(inputs.requirement, 'requirement')))
      throw new WorkerError('INVALID_INPUT', 'Invalid requirement identifier');
    const expected = field === 'expected_count' ? integer(inputs[field], field, 0, 1000000) :
      field === 'available' ? bool(inputs[field], field) : normalize(str(inputs[field], field, true));
    const nodeId = inputs.node_id === undefined ? undefined : str(inputs.node_id, 'node_id');
    if (nodeId && (name !== 'network_nodes' || action !== 'web_field_equals' || !/^[A-Za-z0-9_-]+$/.test(nodeId)))
      throw new WorkerError('INVALID_INPUT', 'node_id requires a node-list field check');
    const main = this.main();
    const ops = getView(name).detail ? new Operations(this.page, getView(name).detail === 'site') : undefined;
    if (inputs.app !== undefined && (!ops || name !== 'network_node_detail' ||
        !['Software status', 'Current version', 'Target version', 'Update Now', 'Retry update'].includes(label)))
      throw new WorkerError('INVALID_INPUT', 'app only scopes node software checks');
    const contains = inputs.match === 'contains';
    if (inputs.match !== undefined && (!contains || action !== 'web_field_equals' || !label.endsWith('reason') || !expected))
      throw new WorkerError('INVALID_INPUT', 'contains is limited to nonempty reason checks');
    if (action === 'web_table_count_equals' && ['network_nodes', 'network_sites'].includes(name))
      throw new WorkerError('UNSUPPORTED_LOCATOR', 'This view renders cards; use its visible Nodes count or Sites count field');
    const actual = await budget.poll(async () => {
      await this.assertView(name);
      if (action === 'web_commerce_equals') return new Commerce(this.page).observe(inputs);
      if (action === 'web_kpi_equals') {
        const card = main.locator('.MuiCard-root').filter({ has: this.page.getByText(label, { exact: true }) });
        const target = await visible(card);
        return target ? text(target.locator(':scope > div').nth(1)) : null;
      }
      if (action === 'web_field_equals') {
        if (ops?.supportsField(label)) return ops.field(label, inputs.app);
        if (nodeId) {
          const card = await visible(main.locator('.ecard[role="button"]').filter({ has: this.page.locator('.tnum').filter({ hasText: new RegExp(` · ${nodeId}$`) }) }));
          if (!card) return null;
          if (label === 'Site') return text(card.locator(':scope > div:last-child > span'));
          if (label === 'Serial #' || label === 'Model type') {
            const identity = await text(card.locator('.tnum'));
            if (identity === null) return null;
            const parts = identity.split(' · ');
            return parts.length === 2 ? parts[label === 'Serial #' ? 1 : 0] : null;
          }
          if (label === 'Connectivity') {
            const dot = await visible(card.locator('[title^="Connectivity: "]'));
            return dot ? (await dot.getAttribute('title'))!.slice('Connectivity: '.length) : null;
          }
          throw new WorkerError('UNSUPPORTED_LOCATOR', 'Unknown node card field');
        }
        if (name === 'network_node_detail' && label === 'Connectivity') {
          const dot = await visible(main.locator('[title^="Connectivity: "]'));
          return dot ? (await dot.getAttribute('title'))!.slice('Connectivity: '.length) : null;
        }
        if ((name === 'network_nodes' && label === 'Nodes count') || (name === 'network_sites' && label === 'Sites count')) {
          if (await main.locator('.MuiSkeleton-root:visible').count()) return null;
          // PageHeader hides a numeric zero; only an explicit empty state is
          // proof of zero. Missing/hidden count during loading is not zero.
          const count = await text(main.locator('.pagehead .pagetitle .cnt'));
          if (count !== null) return count;
          const empty = name === 'network_nodes' ? 'No nodes yet' : 'No sites yet';
          return await main.getByText(empty, { exact: true }).isVisible() ? '0' : null;
        }
        const row = main.locator('.kv-row').filter({ has: this.page.getByText(label, { exact: true }) });
        const target = await visible(row);
        return target ? text(target.locator(':scope > span > .tnum'), true) : null;
      }
      if (action === 'web_action_available') {
        if (ops?.control(label, inputs.app)) return ops.available(label, inputs.app);
        const button = await visible(main.getByRole('button', { name: label, exact: true }));
        return button ? await button.isEnabled() : null;
      }
      // The supplied DataTable has no accessible table name. Only accept the
      // current view's exact list label and a single rendered table. No table
      // (including an error/skeleton/empty-state replacement) is not count 0.
      if (getView(name).label !== label) throw new WorkerError('UNSUPPORTED_LOCATOR', 'Table label must match the current list view');
      const table = await visible(main.getByRole('table'));
      if (!table || await table.locator('.MuiSkeleton-root').count()) return null;
      return table.locator('tbody > tr:visible').count();
    }, value => value !== null && (contains ? typeof value === 'string' && value.includes(String(expected)) : value === expected), 'Visible value did not match the scenario expectation');
    return { expected, actual };
  }
}

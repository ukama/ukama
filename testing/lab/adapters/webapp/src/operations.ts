/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 * Only rendered controls and text: no GraphQL, app-store or response assertions.
 */
import type { Locator, Page } from 'playwright';
import { Budget, WorkerError, normalize, str, type ObjectValue } from './contract.js';

async function one(target: Locator): Promise<Locator | undefined> {
  const visible = target.filter({ visible: true });
  const count = await visible.count();
  if (count > 1) throw new WorkerError('AMBIGUOUS_LOCATOR', 'Operation locator matched multiple visible controls');
  return count ? visible : undefined;
}
async function rendered(target: Locator): Promise<string | null> {
  const locator = await one(target);
  return locator ? normalize(await locator.innerText({ timeout: 500 })) : null;
}
export class Operations {
  constructor(private page: Page, private site: boolean) {}
  private main() { return this.page.locator('main.main'); }
  private dialog() {
    return this.page.getByRole('dialog').filter({ has: this.page.getByText(this.site ? 'Restart site' : 'Restart node', { exact: true }) });
  }
  private menu(label: string) {
    return this.page.getByRole('menuitem').filter({ has: this.page.locator('.MuiListItemText-primary').getByText(label, { exact: true }) });
  }
  private software(app: unknown) { return this.main().getByTestId(`software-${str(app, 'app')}`); }
  control(label: string, app?: unknown): Locator | undefined {
    if (label === 'Confirm restart') return this.dialog().getByRole('button').filter({ hasText: /^(Restart node|Restart|Restarting…|Node is busy|Site is busy)$/ });
    if (label === 'Cancel') return this.dialog().getByRole('button', { name: 'Cancel', exact: true });
    if (this.site) {
      if (label === 'Site actions') return this.main().getByTestId('site-actions');
      if (label === 'Restart site') return this.menu(label);
      if (label === 'Radio' || label === 'Cellular') return this.menu(label).getByRole('checkbox');
    } else {
      if (label === 'Restart node') return this.main().getByTestId('node-restart');
      if (label === 'Update Now' || label === 'Retry update') return this.software(app).getByRole('button', { name: label, exact: true });
    }
    return undefined;
  }
  supportsField(label: string): boolean {
    return ['Dialog', 'Confirmation', 'Restart progress', 'Restart reason', 'Radio state', 'Cellular state',
      'Radio reason', 'Cellular reason', 'Software status', 'Current version', 'Target version', 'Lifecycle', 'Notification'].includes(label);
  }
  async field(label: string, app?: unknown): Promise<string | null> {
    if (label === 'Dialog') return await one(this.dialog()) ? 'open' : 'closed';
    if (label === 'Confirmation') {
      const input = await one(this.dialog().getByRole('textbox'));
      return input ? input.inputValue({ timeout: 500 }) : null;
    }
    if (label === 'Notification') return rendered(this.page.getByRole('alert'));
    if (label === 'Lifecycle') return rendered(this.main().getByTestId('node-lifecycle'));
    if (label === 'Restart progress') {
      const content = await rendered(this.control(this.site ? 'Site actions' : 'Restart node')!);
      return content?.replace(/ \([^()]*\)$/, '') ?? null; // Ignore only the displayed elapsed timer.
    }
    if (label === 'Restart reason' && !this.site) {
      const target = await one(this.control('Restart node')!);
      return target ? normalize(await target.getAttribute('title') ?? '') : null;
    }
    if (label.endsWith('reason')) return rendered(this.menu(label === 'Restart reason' ? 'Restart site' : label.split(' ')[0]!).locator('.MuiListItemText-secondary'));
    if (label === 'Radio state' || label === 'Cellular state') return rendered(this.menu(label.split(' ')[0]!).getByTestId('toggle-state'));
    if (label === 'Current version') return rendered(this.software(app).getByTestId('software-version'));
    if (label === 'Target version') return rendered(this.software(app).getByTestId('software-target'));
    if (label === 'Software status') return rendered(this.software(app).getByTestId('software-status'));
    throw new WorkerError('UNSUPPORTED_LOCATOR', 'Unknown operation field');
  }
  async available(label: string, app?: unknown): Promise<boolean | null> {
    const control = this.control(label, app);
    if (!control) throw new WorkerError('UNSUPPORTED_LOCATOR', 'Unknown operation control');
    const target = await one(control);
    return target ? target.isEnabled() : null;
  }
  private async click(target: Locator, budget: Budget): Promise<void> {
    // One click only. Playwright waits for actionability before dispatching it;
    // no repeat mutation after a transport failure or assertion timeout.
    await target.click({ timeout: budget.remaining() });
  }
  async run(inputs: ObjectValue, budget: Budget): Promise<void> {
    const action = str(inputs.action, 'action');
    const siteOnly = ['open_site_actions', 'close_site_actions', 'fill_confirmation', 'set_radio', 'set_service'];
    const nodeOnly = ['open_software', 'update_software', 'retry_update'];
    if ((siteOnly.includes(action) && !this.site) || (nodeOnly.includes(action) && this.site))
      throw new WorkerError('INVALID_INPUT', 'Action does not belong to this detail view');
    const fill = action === 'fill_confirmation', toggle = action === 'set_radio' || action === 'set_service';
    const update = action === 'update_software' || action === 'retry_update';
    if ((!fill && !toggle && inputs.value !== undefined) || (!update && (inputs.app !== undefined || inputs.tag !== undefined)))
      throw new WorkerError('INVALID_INPUT', 'Action has unrelated parameters');
    if (action === 'open_restart') return this.click(this.control(this.site ? 'Restart site' : 'Restart node')!, budget);
    if (action === 'confirm_restart') return this.click(this.control('Confirm restart')!, budget);
    if (action === 'cancel_dialog') return this.click(this.control('Cancel')!, budget);
    if (action === 'open_site_actions') return this.click(this.control('Site actions')!, budget);
    if (action === 'close_site_actions') {
      await this.page.getByRole('menu').waitFor({ state: 'visible', timeout: budget.remaining() });
      await this.page.keyboard.press('Escape'); return;
    }
    if (fill) {
      await this.dialog().getByRole('textbox').fill(str(inputs.value, 'value', true), { timeout: budget.remaining() }); return;
    }
    if (toggle) {
      const value = str(inputs.value, 'value');
      if (!['on', 'off'].includes(value)) throw new WorkerError('INVALID_INPUT', 'Toggle value must be on or off');
      const label = action === 'set_radio' ? 'Radio' : 'Cellular';
      await budget.poll(() => this.field(`${label} state`), v => v === 'On' || v === 'Off', 'Reported toggle state is unavailable');
      if (await this.field(`${label} state`) === (value === 'on' ? 'On' : 'Off'))
        throw new WorkerError('NO_STATE_CHANGE', 'Requested toggle is already in that state');
      await this.click(this.menu(label), budget); return; // UI confirms checked state only after a health report.
    }
    if (action === 'open_software') return this.click(this.main().getByRole('tab', { name: 'Software', exact: true }), budget);
    if (update) {
      const tag = normalize(str(inputs.tag, 'tag'));
      await budget.poll(() => this.field('Target version', inputs.app), v => v !== null, 'Software target is not visible');
      const target = await this.field('Target version', inputs.app);
      if (target !== tag) throw new WorkerError('TARGET_MISMATCH', 'Visible software target differs from requested tag', target);
      const expectedState = action === 'retry_update' ? 'Update failed' : 'Update available';
      if (await this.field('Software status', inputs.app) !== expectedState)
        throw new WorkerError('WRONG_UPDATE_STATE', `Expected ${expectedState} before submitting update`);
      return this.click(this.control(action === 'retry_update' ? 'Retry update' : 'Update Now', inputs.app)!, budget);
    }
    throw new WorkerError('UNSUPPORTED_ACTION', 'Unknown semantic operation');
  }
}

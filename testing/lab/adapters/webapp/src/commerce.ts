/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { closeSync, fsyncSync, openSync, renameSync, writeFileSync } from 'node:fs';
import { lstat } from 'node:fs/promises';
import { join, isAbsolute } from 'node:path';
import { paymentWindow, paymentDate, downloadReceipt, receiptPdf, rejectPayment, paymentRejection } from './payments.js';
import { commerceFault, commerceFaultState } from './commerce-faults.js';
import type { Locator, Page, Response, Request } from 'playwright';
import { Budget, WorkerError, integer, keys, normalize, object, str, bool, type ObjectValue } from './contract.js';

const escape = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
const field = (root: Locator, page: Page, label: string) => root.locator('.ff').filter({ has: page.locator('.ff-label').filter({ hasText: new RegExp(`^${escape(label)}(?:\\s*\\*)?$`) }) });
async function visibleText(target: Locator): Promise<string | null> {
  const live = target.filter({ visible: true });
  if (await live.count() > 1) throw new WorkerError('AMBIGUOUS_LOCATOR', 'Commerce field matches multiple visible elements');
  return await live.count() === 1 ? normalize(await live.innerText({ timeout: 500 })) : null;
}
export function planValues(raw: unknown): ObjectValue & {days: number; volume: number} {
  const p = object(raw, 'plan');
  keys(p, ['ref', 'id', 'name', 'data_mb', 'duration_minutes', 'amount', 'currency', 'country', 'unit', 'organization']);
  for (const k of ['ref', 'name', 'currency', 'country']) str(p[k], k);
  str(p.id, 'id', true);
  const minutes = integer(p.duration_minutes, 'duration_minutes', 1, 43200);
  if (![1440, 10080, 43200].includes(minutes)) throw new WorkerError('INVALID_INPUT', 'Plan validity must be exactly 1, 7 or 30 days');
  const mb = integer(p.data_mb, 'data_mb', 1, 1048576);
  if (!['MB', 'GB'].includes(String(p.unit)) || (p.unit === 'GB' && mb % 1024)) throw new WorkerError('INVALID_INPUT', 'Invalid binary data unit or fractional GB');
  if (typeof p.amount !== 'number' || !Number.isFinite(p.amount) || p.amount <= 0 || p.amount > 1000000) throw new WorkerError('INVALID_INPUT', 'Invalid price');
  bool(p.organization, 'organization');
  return { ...p, days: minutes / 1440, volume: p.unit === 'GB' ? mb / 1024 : mb };
}
// Receipts are private ownership/correlation records. Only the DOM supplies
// assertion values. The listener never sends requests or changes application state.
class MutationReceipt {
  private document: ObjectValue;
  private path: string;
  private submitted = false;
  private started = 0;
  private attempted = new Set<Request>();
  private identified = new Set<Request>();
  private lastAttempt = 0;
  private binding?: ObjectValue;
  private failure?: Error;
  private pending = new Set<Promise<void>>();
  constructor(private page: Page, private inputs: ObjectValue, directory: string, command: number) {
    const creation = object(inputs.creation, 'creation'); keys(creation, ['kind', 'ref', 'name']);
    for (const k of ['kind', 'ref', 'name']) str(creation[k], k);
    if (!['package', 'subscriber', 'sim', 'payment'].includes(String(creation.kind))) throw new WorkerError('INVALID_INPUT', 'Unsupported commerce creation kind');
    this.path = join(directory, `creation-${command}.json`);
    this.document = { ...creation, command_id: command, state: 'prepared', bindings: [] };
    this.save(); this.page.on('response', this.response); this.page.on('request', this.request);
  }
  private request = (request: Request) => {
    if (!this.submitted || this.inputs.action !== 'top_up_rapid') return;
    let b: any; try { b = request.postDataJSON(); } catch { return; }
    const c = this.inputs.customer as ObjectValue, p = this.inputs.plan as ObjectValue;
    if (!/\baddPayment\s*\(/.test(b?.query ?? '') || b?.variables?.data?.sim !== c.sim_id || b?.variables?.data?.itemId !== p.id || b?.variables?.data?.payerEmail !== c.email) return;
    this.attempted.add(request); this.lastAttempt = performance.now();
    this.document.attempt_count = this.attempted.size;
    if ((this.document.bindings as unknown[]).length < 2) this.document.state = 'submitted';
    this.save();
  };
  private save() {
    const fd = openSync(`${this.path}.tmp`, 'w', 0o600);
    try { writeFileSync(fd, JSON.stringify(this.document)); fsyncSync(fd); } finally { closeSync(fd); }
    renameSync(`${this.path}.tmp`, this.path);
    const dir = openSync(join(this.path, '..'), 'r'); try { fsyncSync(dir); } finally { closeSync(dir); }
  }
  private response = (response: Response) => {
    const task = (async () => {
      if (!this.submitted) return;
      let request: any; try { request = response.request().postDataJSON(); } catch { return; }
      const kind = String(this.document.kind), plan = this.inputs.plan as ObjectValue | undefined, customer = this.inputs.customer as ObjectValue | undefined;
      const mutation = ({ package: 'addPackage', subscriber: 'addSubscriber', sim: 'allocateSim', payment: 'addPayment' } as const)[kind as 'package'];
      if (!new RegExp(`\\b${mutation}\\s*\\(`).test(request?.query ?? '')) return;
      const data = request?.variables?.data;
      const matches = kind === 'package' ? data?.name === plan?.name && data?.networkId === (plan?.organization ? '' : this.inputs.network_id) :
        kind === 'subscriber' ? data?.name === customer?.name && data?.email === customer?.email && data?.network_id === this.inputs.network_id :
        kind === 'sim' ? (this.inputs.action === 'allocate_auto' ? !data?.iccid : data?.iccid === customer?.iccid) && data?.subscriber_id === customer?.id && data?.network_id === this.inputs.network_id && data?.package_id === plan?.id :
        data?.sim === customer?.sim_id && data?.itemId === plan?.id && data?.payerEmail === customer?.email;
      if (!matches) return;
      let body: any; try { body = await response.json(); } catch { return; }
      const row = body?.data?.[mutation];
      const id = row?.[kind === 'package' || kind === 'subscriber' ? 'uuid' : 'id'];
      if (typeof id !== 'string' || !/^[A-Za-z0-9_-]{1,127}$/.test(id)) return;
      if ((kind === 'package' || kind === 'subscriber') && row.name !== this.document.name) return;
      if (kind === 'sim' && (row.iccid !== customer?.iccid || row.subscriber_id !== customer?.id || row.network_id !== this.inputs.network_id)) return;
      if (this.binding && this.binding.id !== id) {
        const bindings = this.document.bindings as ObjectValue[];
        if (!bindings.some(b => b.id === id)) this.document.bindings = [...bindings, {kind,ref:this.document.ref,name:this.document.name,id,observed_via:'ui_result'}];
        this.document.state = 'ambiguous';
        this.save();
        throw new WorkerError('OWNERSHIP_CONFLICT', 'Multiple IDs observed; retained all private receipts for reconciliation');
      }
      if ((this.document.bindings as unknown[]).length > 1) return;
      this.binding = { kind, ref: this.document.ref, name: this.document.name, id, observed_via: 'ui_result' };
      this.identified.add(response.request());
      this.document.state = this.inputs.action === 'top_up_rapid' && this.identified.size !== this.attempted.size ? 'submitted' : 'identified'; this.document.bindings = [this.binding]; this.save();
    })().catch(error => { this.failure = error; });
    this.pending.add(task); void task.finally(() => this.pending.delete(task));
  };
  submit() { this.started = Date.now(); this.submitted = true; this.document.state = 'submitted'; this.save(); }
  intercepted() { this.document.state = 'intercepted'; this.save(); }
  async result(budget: Budget) {
    await budget.poll(async () => {
      if (this.failure) throw this.failure;
      if (this.inputs.action === 'top_up_rapid' && (!this.attempted.size || this.attempted.size !== this.identified.size || performance.now() - this.lastAttempt < 250)) return undefined;
      return this.binding;
    }, Boolean, 'Mutation ownership was not identified');
    if (this.document.kind === 'payment') paymentWindow(this.page,String(this.binding!.id),this.started,Date.now());
    return [this.binding!];
  }
  async close() { await Promise.all(this.pending); this.page.off('response', this.response); this.page.off('request', this.request); if (this.failure) throw this.failure; }
}
export class Commerce {
  constructor(private page: Page) {}
  private main() { return this.page.locator('main.main'); }
  private dialog() { return this.page.getByRole('dialog').filter({ visible: true }); }
  private drawer() { return this.page.locator('.MuiDrawer-paper:visible'); }
  private card(name: string) { return this.main().locator('.card').filter({ has: this.page.getByText(name, { exact: true }) }); }
  private packages(name: string) { return this.drawer().locator('.card.card-pad').filter({ has: this.page.getByRole('button', { name: 'Package options', exact: true }) }).filter({ has: this.page.getByText(name, { exact: true }) }); }
  private detail(label: string) { return this.drawer().getByText(label, { exact: true }).locator('..').locator(':scope > span.tnum'); }
  private async customer(customer: ObjectValue, budget: Budget, requireSim = true) {
    await budget.poll(() => visibleText(this.drawer().getByText(String(customer.name), { exact: true })), v => v === customer.name, 'Wrong customer drawer');
    if (requireSim) await budget.poll(() => visibleText(this.detail('ICCID')), v => v === customer.iccid, 'Wrong customer SIM');
  }
  async run(inputs: ObjectValue, directory: string, command: number, budget: Budget): Promise<unknown[]> {
    keys(inputs, ['view', 'network_name', 'network_id', 'action', 'plan', 'customer', 'creation', 'new_name']);
    const action = str(inputs.action, 'action');
    const p = inputs.plan === undefined ? undefined : planValues(inputs.plan);
    const c = inputs.customer === undefined ? undefined : object(inputs.customer, 'customer');
    if (c) { keys(c, ['ref', 'id', 'name', 'email', 'iccid', 'sim_id', 'payment_id']); for (const k of ['ref', 'name', 'email', 'iccid']) str(c[k], k); }
    const requiredKind = ({create_plan:'package',create_customer:'subscriber',allocate_sim:'sim',allocate_auto:'sim',top_up:'payment',top_up_rapid:'payment',failed_top_up:'payment'} as Record<string,string>)[action];
    if (requiredKind ? object(inputs.creation).kind !== requiredKind : inputs.creation !== undefined) throw new WorkerError('INVALID_INPUT', 'Mutation intent does not match action');
    const receipt = requiredKind ? new MutationReceipt(this.page, inputs, directory, command) : undefined;
    const click = async (root: Locator, name: string) => root.getByRole('button', { name, exact: true }).click({ timeout: budget.remaining() });
    const submit = async (name: string) => { receipt?.submit(); await click(this.dialog(), name); await this.dialog().waitFor({ state: 'hidden', timeout: budget.remaining() }); };
    try {
      if (action === 'download_receipt' && c && p) {
        await this.customer(c,budget);
        await downloadReceipt(this.page,String(p.name),str(c.payment_id,'payment_id'),directory,command,budget);
      } else if (['name_pending','name_failure','pool_failure','clear_fault'].includes(action)) {
        await commerceFault(this.page, action, p ? String(p.name) : undefined, budget);
      } else if (action === 'rename_plan' && p) {
        const next = str(inputs.new_name, 'new_name');
        if (next !== String(p.name) + '-renamed') throw new WorkerError('INVALID_INPUT','Rename must use the planned suffix');
        if (await this.card(next).count()) throw new WorkerError('NAME_EXISTS','Renamed plan already exists');
        await click(this.card(String(p.name)), 'Plan actions');
        await this.page.getByRole('menuitem',{name:'Edit plan',exact:true}).click({timeout:budget.remaining()});
        await field(this.dialog(), this.page, 'Data plan name').locator('input').fill(next,{timeout:budget.remaining()});
        await this.dialog().getByText('✓ Name is available',{exact:true}).waitFor({timeout:budget.remaining()});
        await submit('Save changes'); await this.card(next).waitFor({timeout:budget.remaining()});
      } else if (action === 'create_plan' && p) {
        if (await this.card(String(p.name)).count()) throw new WorkerError('NAME_EXISTS', 'Planned data plan already exists');
        await click(this.main().locator('.pagehead'), 'Create plan');
        const d = this.dialog();
        await field(d, this.page, 'Data plan name').locator('input').fill(String(p.name), { timeout: budget.remaining() });
        await field(d, this.page, 'Price').locator('input').fill(String(p.amount), { timeout: budget.remaining() });
        await field(d, this.page, 'Data volume').locator('input').fill(String(p.volume), { timeout: budget.remaining() });
        await field(d, this.page, 'Unit').locator('select').selectOption({ label: String(p.unit) }, { timeout: budget.remaining() });
        await field(d, this.page, 'Validity').locator('select').selectOption(String(p.days), { timeout: budget.remaining() });
        await d.getByRole('checkbox', { name: 'Available to all networks', exact: true }).setChecked(Boolean(p.organization), { timeout: budget.remaining() });
        if (!p.organization) await field(d, this.page, 'Network').locator('select').selectOption({ label: str(inputs.network_name, 'network_name') }, { timeout: budget.remaining() });
        await d.getByText('✓ Name is available', { exact: true }).waitFor({ timeout: budget.remaining() });
        await submit('Create plan'); await this.card(String(p.name)).waitFor({ timeout: budget.remaining() });
      } else if (action === 'edit_plan' && p) {
        await click(this.card(String(p.name)), 'Plan actions'); await this.page.getByRole('menuitem', { name: 'Edit plan', exact: true }).click({ timeout: budget.remaining() });
      } else if (action === 'create_customer' && c) {
        if (await this.main().locator('tbody > tr').filter({ has: this.page.getByText(String(c.name), { exact: true }) }).count()) throw new WorkerError('NAME_EXISTS', 'Planned customer already exists');
        await click(this.main().locator('.pagehead'), 'Add customer');
        const [first, ...last] = String(c.name).split(' '), d = this.dialog();
        await field(d, this.page, 'First name').locator('input').fill(first!, { timeout: budget.remaining() });
        await field(d, this.page, 'Last name').locator('input').fill(last.join(' '), { timeout: budget.remaining() });
        await field(d, this.page, 'Email').locator('input').fill(String(c.email), { timeout: budget.remaining() });
        // Leave plan empty: creation and SIM allocation have separate ownership receipts.
        if (await field(d, this.page, 'Data plan').locator('select').inputValue() !== '') throw new WorkerError('UNEXPECTED_DEFAULT', 'Customer form preselected a plan');
        await submit('Add customer');
        await this.main().locator('tbody > tr').filter({ has: this.page.getByText(String(c.name), { exact: true }) }).waitFor({ timeout: budget.remaining() });
      } else if (action === 'open_customer' && c) {
        await this.main().locator('tbody > tr').filter({ has: this.page.getByText(String(c.name), { exact: true }) }).click({ timeout: budget.remaining() });
        await this.customer(c, budget, Boolean(c.sim_id));
      } else if (action === 'close_customer' && c) {
        await this.customer(c, budget, Boolean(c.sim_id)); await click(this.drawer(), 'Close'); await this.drawer().waitFor({ state: 'hidden', timeout: budget.remaining() });
      } else if (['allocate_sim', 'allocate_auto', 'top_up', 'top_up_rapid', 'failed_top_up', 'cancel_top_up'].includes(action) && c && p) {
        await this.customer(c, budget, !action.startsWith('allocate_'));
        await click(this.drawer(), action.startsWith('allocate_') ? 'Allocate a SIM' : 'Top up');
        const select = field(this.dialog(), this.page, 'Data plan').locator('select');
        const option = select.locator('option').filter({ hasText: new RegExp(`^${escape(String(p.name))} · `) });
        await option.waitFor({ state: 'attached', timeout: budget.remaining() });
        if (await option.count() !== 1) throw new WorkerError('AMBIGUOUS_LOCATOR', 'Plan option is not unique');
        const value = await option.getAttribute('value');
        if (value !== p.id) throw new WorkerError('WRONG_ENTITY', 'Visible plan option has another identity');
        await select.selectOption(value!, { timeout: budget.remaining() });
        if (action === 'allocate_sim') await field(this.dialog(), this.page, 'SIM').locator('select').selectOption({ label: String(c.iccid) }, { timeout: budget.remaining() });
        if (action === 'allocate_auto') {
          const select = field(this.dialog(),this.page,'SIM').locator('select');
          const options = await select.locator('option').evaluateAll(es => es.map(e => ({value:(e as HTMLOptionElement).value,text:e.textContent?.trim()})));
          if (options.length !== 2 || options.filter(o => o.value === '' && o.text === 'Auto-assign from pool').length !== 1 || options.filter(o => o.value === c.iccid && o.text === c.iccid).length !== 1)
            throw new WorkerError('UNSAFE_AUTO_ASSIGN','Auto-assignment requires the sole available SIM to be owned by this run');
          await select.selectOption('',{timeout:budget.remaining()});
        }
        if (action === 'failed_top_up') {
          await rejectPayment(this.page,inputs,async()=>{receipt!.intercepted();await click(this.dialog(),'Top up');},budget);
        } else if (action === 'top_up_rapid') {
          receipt!.submit();
          await this.dialog().getByRole('button',{name:'Top up',exact:true}).dblclick({timeout:budget.remaining(),delay:0});
          await this.dialog().waitFor({state:'hidden',timeout:budget.remaining()});
        } else if (action === 'cancel_top_up') { await click(this.dialog(), 'Cancel'); await this.dialog().waitFor({ state: 'hidden', timeout: budget.remaining() }); }
        else await submit(action.startsWith('allocate_') ? 'Allocate SIM' : 'Top up');
      } else if (action === 'open_receipt' && c && p) {
        await this.customer(c, budget); await click(this.packages(String(p.name)), 'Package options');
        await this.page.getByRole('menuitem', { name: 'View receipt', exact: true }).click({ timeout: budget.remaining() });
        await this.dialog().getByText('Payment receipt', { exact: true }).waitFor({ timeout: budget.remaining() });
      } else if (action === 'close_dialog') {
        const d = this.dialog();
        if (await d.getByRole('button', { name: 'Cancel', exact: true }).count()) await click(d, 'Cancel'); else await click(d, 'Close');
        await d.waitFor({ state: 'hidden', timeout: budget.remaining() });
      } else if (['activate_sim', 'deactivate_sim'].includes(action) && c) {
        await this.customer(c, budget); const name = action === 'activate_sim' ? 'Activate SIM' : 'Deactivate SIM';
        await click(this.drawer(), name); await submit(name);
      } else throw new WorkerError('INVALID_INPUT', 'Unknown commerce action or missing inputs');
      return receipt && action !== 'failed_top_up' ? await receipt.result(budget) : [];
    } finally { await receipt?.close(); }
  }
  async importSims(inputs: ObjectValue, budget: Budget) {
    keys(inputs, ['csv_path', 'iccids']);
    const path = str(inputs.csv_path, 'csv_path');
    if (!isAbsolute(path) || !path.endsWith('.csv')) throw new WorkerError('INVALID_INPUT', 'Factory CSV must be an absolute .csv path');
    const stat = await lstat(path);
    if (!stat.isFile() || stat.size === 0 || stat.size > 10 * 1024 * 1024) throw new WorkerError('INVALID_INPUT', 'Invalid SIM CSV');
    if (!Array.isArray(inputs.iccids) || !inputs.iccids.length || new Set(inputs.iccids).size !== inputs.iccids.length || inputs.iccids.some(x => typeof x !== 'string' || !/^\d{18,22}$/.test(x))) throw new WorkerError('INVALID_INPUT', 'Invalid expected SIM identities');
    await this.main().locator('.pagehead').getByRole('button', { name: 'Upload SIMs', exact: true }).click({ timeout: budget.remaining() });
    await this.dialog().locator('input[type="file"]').setInputFiles(path, { timeout: budget.remaining() });
    await this.dialog().getByRole('button', { name: 'Upload', exact: true }).click({ timeout: budget.remaining() });
    await this.dialog().waitFor({ state: 'hidden', timeout: budget.remaining() });
    // This console renders a capped list without a search control. Missing
    // imported rows fail here; do not fetch or insert them into the page.
    for (const iccid of inputs.iccids) {
      const row = this.main().getByRole('row').filter({ has: this.page.getByText(iccid, { exact: true }) });
      await budget.poll(() => visibleText(row.locator('.MuiChip-label')), value => value === 'Available', 'Imported SIM was not visibly available');
    }
  }
  async observe(inputs: ObjectValue): Promise<string | null> {
    const label = str(inputs.label, 'label'), name = inputs.plan_name === undefined ? undefined : str(inputs.plan_name, 'plan_name');
    if (label === 'Commerce fault') return commerceFaultState(this.page);
    if (label === 'SIM option present' && inputs.view === 'customer_customers' && inputs.iccid) {
      const select = field(this.dialog(),this.page,'SIM').locator('select');
      if (!await this.dialog().getByText('Allocate a SIM',{exact:true}).isVisible() || !await select.isVisible() || await select.isDisabled()) return null;
      const matches = select.locator('option').filter({hasText:new RegExp(`^${escape(String(inputs.iccid))}$`)});
      if (await matches.count() > 1) throw new WorkerError('AMBIGUOUS_LOCATOR','Duplicate ICCID options');
      return String(await matches.count() === 1 && await matches.getAttribute('value') === inputs.iccid);
    }
    if (label === 'Customer plan' && String(inputs.view).endsWith('_customers') && inputs.customer_name) {
      const row=this.main().locator('tbody > tr').filter({has:this.page.getByText(String(inputs.customer_name),{exact:true})});
      return visibleText(row.locator('td').nth(1));
    }
    if (inputs.view === 'business_packages' && name) {
      const columns: Record<string, number> = { 'Performance price': 1, 'Performance sold': 2, 'Performance revenue': 3, 'Performance share': 4 };
      if (columns[label] === undefined) throw new WorkerError('UNSUPPORTED_LOCATOR', 'Unknown package performance column');
      const row = this.main().getByRole('row').filter({ has: this.page.getByRole('cell', { name, exact: true }) });
      return visibleText(row.getByRole('cell').nth(columns[label]!));
    }
    if (inputs.view === 'business_sim_pool') {
      if (label === 'Pool reconciliation') {
        if (await this.main().locator('.MuiSkeleton-root:visible').count() || await this.main().getByText("Couldn't load SIMs",{exact:true}).isVisible()) return null;
        const rows=this.main().locator('tbody > tr:visible'), count=await rows.count();
        if (!count && !await this.main().getByText('No SIMs',{exact:true}).isVisible()) return null;
        const footer=await visibleText(this.main().locator('.tbl-foot .tnum'));
        if (footer !== `Showing ${count} of ${count.toLocaleString('en-US')}`) return 'incomplete inventory';
        const totals: Record<string,number>={Available:0,Assigned:0,Faulty:0}, ids=new Set<string>();
        for (let n=0;n<count;n++) {
          const row=rows.nth(n), id=await visibleText(row.locator('td').first()), status=await visibleText(row.locator('.MuiChip-label'));
          if (!id || ids.has(id) || !status || !(status in totals)) return 'invalid inventory';
          ids.add(id); totals[status]!++;
        }
        for (const [status,total] of Object.entries(totals)) {
          const tile=this.main().locator('.MuiCard-root').filter({has:this.page.getByText(status,{exact:true})});
          const value=await visibleText(tile.locator(':scope > div').nth(1));
          if (value !== total.toLocaleString('en-US')) return `${status}: ${value} != ${total}`;
        }
        return 'matched';
      }
      if (label !== 'Pool status' || !inputs.iccid) throw new WorkerError('INVALID_INPUT', 'Pool status needs a SIM identity');
      return visibleText(this.main().getByRole('row').filter({ has: this.page.getByText(String(inputs.iccid), { exact: true }) }).locator('.MuiChip-label'));
    }
    if (inputs.view === 'business_data_plans' && name) {
      if (label === 'Plan terms') return visibleText(this.card(name).locator('.card-pad > div').nth(2));
      if (label === 'Plan price') return visibleText(this.card(name).locator('.card-pad > div').nth(1));
      if (label === 'Plan scope') return visibleText(this.card(name).locator('span[title]'));
      if (['Validity', 'Price', 'Data volume', 'Unit'].includes(label)) {
        const input = field(this.dialog(), this.page, 'Data plan name').locator('input');
        if (!await input.isVisible() || await input.inputValue() !== name) throw new WorkerError('WRONG_ENTITY', 'Edit dialog is not for the selected plan');
        const f = field(this.dialog(), this.page, label);
        if (await f.locator('input,select,textarea,[contenteditable="true"]').count()) return 'editable';
        return visibleText(f.locator('.ff-readonly'));
      }
    }
    if (inputs.view === 'customer_customers' && inputs.customer_name) {
      if (await visibleText(this.drawer().getByText(String(inputs.customer_name), { exact: true })) !== inputs.customer_name) throw new WorkerError('WRONG_ENTITY', 'Check is not on the requested customer drawer');
      if (label === 'ICCID') return visibleText(this.detail(label));
      if (await visibleText(this.detail('ICCID')) !== inputs.iccid) throw new WorkerError('WRONG_ENTITY', 'Drawer ICCID differs from the selected SIM');
      if (['SIM status', 'Total usage', 'Phone'].includes(label)) return visibleText(this.detail(label));
      if (label === 'Cycle usage') return visibleText(this.drawer().locator('.card.card-pad').filter({ hasNot: this.page.getByRole('button', { name: 'Package options', exact: true }) }).locator('.tnum'));
      if (label === 'Active plan') return visibleText(this.drawer().locator('.card.card-pad').filter({ hasNot: this.page.getByRole('button', { name: 'Package options', exact: true }) }).locator(':scope > div > span').first());
      if (label === 'Package count' && name) {
        if (!await this.drawer().getByText('Packages', {exact: true}).isVisible() || await this.drawer().locator('.MuiSkeleton-root:visible').count()) return null;
        return String(await this.packages(name).count());
      }
      if (label === 'Package status' && name) return visibleText(this.packages(name).locator('.MuiChip-label'));
      if (label === 'Package dates' && name) return visibleText(this.packages(name).locator('.tnum'));
      if (label === 'Package days' && name) {
        const dates = await visibleText(this.packages(name).locator('.tnum'));
        if (!dates) return null;
        const parts = dates.split(' – ');
        const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
        const parse = (date: string) => {
          const m = /^(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (\d{1,2}), (\d{4})$/.exec(date);
          return m ? Date.UTC(Number(m[3]), months.indexOf(m[1]!), Number(m[2])) : NaN;
        };
        const days = parts.length === 2 ? (parse(parts[1]!) - parse(parts[0]!)) / 86400000 : NaN;
        return Number.isInteger(days) && days >= 0 ? String(days) : null;
      }
      if (label === 'Payment rejection' && inputs.sim_id && inputs.plan_id) return paymentRejection(this.page,String(inputs.sim_id),String(inputs.plan_id));
      const d = this.dialog();
      if (label.startsWith('Receipt ')) {
        if (!await d.getByText('Payment receipt', { exact: true }).isVisible()) return null;
        if (label === 'Receipt PDF' && name && inputs.payment_id) return receiptPdf(this.page,name,String(inputs.payment_id));
        if (label === 'Receipt date window' && inputs.payment_id) {
          const date = await visibleText(d.getByText('Paid on',{exact:true}).locator('..').locator(':scope > div').nth(1));
          return date ? paymentDate(this.page,String(inputs.payment_id),date) : null;
        }
        if (label === 'Receipt payer') return visibleText(d.getByText('Billed to',{exact:true}).locator('..').locator(':scope > div').nth(1));
        if (label === 'Receipt total') return visibleText(d.getByText('Total paid', { exact: true }).locator('..').locator(':scope > span').nth(1));
        if (label === 'Receipt payment ID') return visibleText(d.getByText('Payment ID', { exact: true }).locator('..').locator(':scope > div').nth(1));
        if (label === 'Receipt method') return visibleText(d.getByText('Method', { exact: true }).locator('..').locator(':scope > div').nth(1));
        if (label === 'Receipt status') return visibleText(d.getByText('Completed', { exact: true }));
        if (label === 'Receipt plan' && name) return visibleText(d.getByText(name, { exact: true }));
        if (label === 'Receipt empty') return visibleText(d.getByText('No receipt found for this package. It may have been allocated without a recorded payment.', { exact: true }));
      }
    }
    throw new WorkerError('UNSUPPORTED_LOCATOR', 'Unknown commerce field or missing resource scope');
  }
}

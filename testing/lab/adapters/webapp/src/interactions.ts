/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 * UI-only interactions. Selectors are derived from the unchanged console.
 */
import type { Locator, Page } from 'playwright';
import { Budget, WorkerError, keys, normalize, str, type ObjectValue } from './contract.js';
import {TeamSupport, TEAM_ACTIONS, TEAM_LABELS} from './team-support.js';
import { getView, type ConsoleApp } from './console-app.js';
import {Analytics, ANALYTICS_ACTIONS, ANALYTICS_LABELS, analyticsView} from './analytics.js';
const escape = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
async function text(target: Locator): Promise<string | null> {
  const visible = target.filter({ visible: true }), count = await visible.count();
  if (count > 1) throw new WorkerError('AMBIGUOUS_LOCATOR', 'UI field matches multiple visible elements');
  return count === 1 ? normalize(await visible.innerText({ timeout: 500 })) : null;
}
export class Interactions {
  constructor(private page: Page, private origin: string, private app?: ConsoleApp) {}
  private main() { return this.page.locator('main.main'); }
  private dialog() { return this.page.getByRole('dialog').filter({ visible: true }); }
  private field(label: string) { return this.dialog().locator('.ff').filter({ has: this.page.locator('.ff-label').filter({ hasText: new RegExp(`^${escape(label)}(?:\\s*\\*)?$`) }) }); }
  private async scope() { return await this.dialog().count() ? this.dialog() : this.main(); }
  private async route(view: string) {
    const url = new URL(this.page.url());
    if (view === 'session') return;
    if (url.origin !== this.origin || url.pathname !== getView(view).path || !await this.main().isVisible())
      throw new WorkerError('WRONG_VIEW', 'UI command is not on its declared console view');
  }
  async run(i: ObjectValue, budget: Budget) {
    keys(i, ['view', 'action', 'label', 'value', 'network_name', 'network_id']);
    const view = str(i.view, 'view'), action = str(i.action, 'action');
    const label = str(i.label ?? '', 'label', true), value = str(i.value ?? '', 'value', true);
    if (TEAM_ACTIONS.includes(action)) {
      if (!this.app) throw new WorkerError('WRONG_VIEW','Team/support requires an opened view');
      await this.app.assertView(view,i.network_name as string | undefined);
      await new TeamSupport(this.page).run(view,action,label,value,budget,i.network_id as string | undefined); return;
    }
    if (ANALYTICS_ACTIONS.includes(action)) {
      if (!analyticsView(view) || !this.app) throw new WorkerError('WRONG_VIEW','Analytics requires an opened analytics view');
      await this.app.assertView(view);
      await new Analytics(this.page).run(action,label,value,budget); return;
    }
    await this.route(view);
    if (view !== 'session' && i.network_name && !view.includes('_data_plans') && !view.includes('_members') && !view.includes('_sim_pool') &&
        await text(this.page.locator('header .netswitch .nm')) !== i.network_name)
      throw new WorkerError('WRONG_NETWORK', 'UI action is not on its selected network');
    const click = (target: Locator) => target.click({ timeout: budget.remaining() });
    if (action === 'search') {
      const placeholders: Record<string, string> = {customer_customers: 'Search name or phone', business_customers: 'Search name or phone', network_customers: 'Search name or phone'};
      if (!placeholders[view]) throw new WorkerError('UNSUPPORTED_LOCATOR', 'Search not mapped for this view');
      await this.main().getByPlaceholder(placeholders[view]!, {exact: true}).fill(value, {timeout: budget.remaining()});
    } else if (action === 'sort') {
      if (!['Customer', 'Data usage', 'Last seen'].includes(label) || !view.endsWith('_customers')) throw new WorkerError('INVALID_INPUT', 'Unsupported sortable column');
      await click(this.main().getByRole('columnheader', {name: label, exact: true}));
    } else if (action === 'filter') {
      if (view === 'business_sim_pool' && label === 'Status' && ['All statuses','Available','Assigned','Faulty'].includes(value)) {
        await click(this.main().getByRole('button',{name:value,exact:true})); return;
      }
      if (!['Active plan', 'SIM'].includes(label) || !view.endsWith('_customers')) throw new WorkerError('INVALID_INPUT', 'Unsupported filter column');
      await click(this.main().getByRole('columnheader').getByRole('button', {name: label, exact: true}));
      await click(this.page.getByRole('menuitem', {name: value, exact: true}));
    } else if (action === 'open_form') {
      if (!['Create plan', 'Add customer', 'Invite member'].includes(label)) throw new WorkerError('INVALID_INPUT', 'Unsupported form');
      const opener=this.main().locator('.pagehead').getByRole('button', {name: label, exact: true});
      if (value === 'keyboard') { await opener.focus({timeout:budget.remaining()}); await opener.press('Enter',{timeout:budget.remaining()}); }
      else if (!value) await click(opener);
      else throw new WorkerError('INVALID_INPUT','Form activation is mouse (empty) or keyboard');
      // A missing prerequisite may produce a visible toast instead of a dialog.
    } else if (action === 'fill' || action === 'select') {
      if (!['Data plan name', 'Price', 'Data volume', 'Unit', 'Validity', 'Network', 'First name', 'Last name', 'Email', 'Role', 'Data plan', 'SIM'].includes(label)) throw new WorkerError('INVALID_INPUT', 'Unsupported form field');
      if (action === 'fill') await this.field(label).locator('input').fill(value, {timeout: budget.remaining()});
      else await this.field(label).locator('select').selectOption({label: value}, {timeout: budget.remaining()});
    } else if (action === 'cancel') {
      await click(this.dialog().getByRole('button', {name: 'Cancel', exact: true}));
      await this.dialog().waitFor({state: 'hidden', timeout: budget.remaining()});
    } else if (action === 'press') {
      if (!['Tab', 'Shift+Tab', 'Escape'].includes(value)) throw new WorkerError('INVALID_INPUT', 'Unsupported UI key');
      if (label) await (await this.scope()).getByRole('button',{name:label,exact:true}).focus({timeout:budget.remaining()});
      await this.page.keyboard.press(value);
    } else if (action === 'viewport') {
      if (!['desktop', 'narrow'].includes(value)) throw new WorkerError('INVALID_INPUT', 'Viewport is desktop or narrow');
      await this.page.setViewportSize(value === 'desktop' ? {width:1440,height:1000} : {width:390,height:844});
    } else if (action === 'date_range') {
      if (!['Last 24h','Last 7 days','Last 30 days'].includes(value)) throw new WorkerError('INVALID_INPUT', 'Unknown date range');
      await click(this.main().locator('.pagehead').getByRole('button').filter({hasText:/^Last (24h|7 days|30 days)$/}));
      await click(this.page.getByRole('menuitem',{name:value,exact:true}));
    } else if (action === 'go_back') {
      await this.page.goBack({waitUntil:'domcontentloaded',timeout:budget.remaining()});
    } else if (action === 'palette') {
      await this.page.keyboard.press('Control+k');
      await this.dialog().getByPlaceholder('Jump to a page…',{exact:true}).waitFor({timeout:budget.remaining()});
    } else if (action === 'palette_choose') {
      await this.dialog().getByPlaceholder('Jump to a page…',{exact:true}).fill(value,{timeout:budget.remaining()});
      const result=this.dialog().getByRole('button',{name:value,exact:true});
      await result.focus({timeout:budget.remaining()}); await result.press('Enter',{timeout:budget.remaining()});
      await this.dialog().waitFor({state:'hidden',timeout:budget.remaining()});
    } else if (action === 'mobile_open') {
      await click(this.page.getByRole('button',{name:'Open navigation',exact:true}));
      await this.page.locator('.mobile-nav').waitFor({timeout:budget.remaining()});
    } else if (action === 'mobile_link') {
      await click(this.page.locator('.mobile-nav').getByRole('link',{name:value,exact:true}));
      await this.page.locator('.mobile-nav').waitFor({state:'hidden',timeout:budget.remaining()});
    } else if (action === 'open_allocate' || action === 'open_topup') {
      await click(this.page.locator('.MuiDrawer-paper:visible').getByRole('button',{name:action==='open_allocate'?'Allocate a SIM':'Top up',exact:true}));
    } else if (action === 'clear_session' && view === 'session') {
      await this.page.context().clearCookies();
      await this.page.goto(this.origin,{waitUntil:'domcontentloaded',timeout:budget.remaining()});
    } else throw new WorkerError('INVALID_INPUT', 'Unsupported UI action');
  }
  async check(i: ObjectValue,budget: Budget) {
    keys(i,['view','label','subject','expected','requirement','customer_name','plan_name','network_name']);
    const view=str(i.view,'view'), label=str(i.label,'label'), subject=str(i.subject??'','subject',true), expected=str(i.expected,'expected',true);
    const observe=async ():Promise<string|null>=>{
      if (TEAM_LABELS.includes(label)) {
        if (!this.app) throw new WorkerError('WRONG_VIEW','Team/support requires an opened view');
        await this.app.assertView(view,i.network_name as string | undefined);
        return new TeamSupport(this.page).observe(view,label,subject);
      }
      if (ANALYTICS_LABELS.includes(label)) {
        if (!analyticsView(view) || !this.app) throw new WorkerError('WRONG_VIEW','Analytics requires an opened analytics view');
        await this.app.assertView(view);
        return new Analytics(this.page).observe(label,subject);
      }
      if(view !== 'session' && (new URL(this.page.url()).origin !== this.origin || new URL(this.page.url()).pathname !== getView(view).path || !await this.main().isVisible()))return null;
      if (view === 'session') {
        if (label === 'Auth origin') return new URL(this.page.url()).origin;
        if (label === 'Dashboard visible') return String(await this.main().isVisible());
        throw new WorkerError('INVALID_INPUT','Unknown session assertion');
      }
      if (label==='Path') return new URL(this.page.url()).pathname;
      if (label==='Selected network') return text(this.page.locator('header .netswitch .nm'));
      if (label==='Dialog open') return String(await this.dialog().count()===1);
      if (label==='Mobile navigation open') return String(await this.page.locator('.mobile-nav').isVisible());
      if (label==='Dialog title') return text(this.dialog().getByRole('heading'));
      if (label==='Field error') return text(this.field(subject).locator('.ff-err'));
      if (label==='Field options') {
        const select=this.field(subject).locator('select');
        if(!await select.isVisible()) return null;
        return (await select.locator('option').allTextContents()).map(normalize).join('|');
      }
      if (label==='Field value') { const f=this.field(subject).locator('input,select');return await f.isVisible()?f.inputValue():null; }
      if (label==='Field readonly') return text(this.field(subject).locator('.ff-readonly'));
      if (label==='Text visible') return String(await this.page.getByText(subject,{exact:true}).filter({visible:true}).count()>0);
      if (label==='Button enabled' || label==='Button visible' || label==='Focus on button') {
        const b=(await this.scope()).getByRole('button',{name:subject,exact:true}).filter({visible:true});
        if (await b.count()>1) throw new WorkerError('AMBIGUOUS_LOCATOR','Button is ambiguous');
        if (label==='Button visible') return String(await b.count()===1);
        if (await b.count()!==1)return null;
        return String(label==='Button enabled'?await b.isEnabled():await b.evaluate(e=>e===document.activeElement));
      }
      if (label==='Date range') return text(this.main().locator('.pagehead').getByRole('button').filter({hasText:/^Last (24h|7 days|30 days)$/}));
      if (label==='Plan count') { if(await this.main().locator('.MuiSkeleton-root:visible').count())return null; return String(await this.main().locator('.card').filter({has:this.page.getByText(str(i.plan_name,'plan_name'),{exact:true})}).filter({visible:true}).count()); }
      if (label==='Plan option present') {
        const select=this.field('Data plan').locator('select');
        if(!await select.isVisible()||await select.isDisabled())return null;
        return String(await select.locator('option').filter({hasText:new RegExp(`^${escape(str(i.plan_name,'plan_name'))} · `)}).count()===1);
      }
      if (label==='Customer present' || label==='Customer order') {
        if (await this.main().locator('.MuiSkeleton-root:visible').count())return null;
        const table=this.main().getByRole('table');
        if (!await table.isVisible()) {
          if (label==='Customer present' && (await this.main().getByText('No customers yet',{exact:true}).isVisible() || await this.main().getByText('No customers match',{exact:true}).isVisible()))return 'false';
          return null;
        }
        const rows=table.locator('tbody > tr:visible');
        if(label==='Customer present') {
          const matches=rows.filter({has:this.page.getByText(str(i.customer_name,'customer_name'),{exact:true})});
          if(await matches.count()>1)throw new WorkerError('AMBIGUOUS_LOCATOR','Duplicate customer rows');
          return String(await matches.count()===1);
        }
        const names=await rows.locator('td:first-child > div > div > div:first-child').allInnerTexts();
        return names.length===await rows.count()?names.map(normalize).join('|'):null;
      }
      if(label==='Content fits viewport')return String(await this.page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth));
      throw new WorkerError('UNSUPPORTED_LOCATOR','Unknown UI assertion');
    };
    const actual=await budget.poll(observe,v=>v!==null&&v===expected,'Visible UI state differs from scenario expectation');
    return {expected,actual};
  }
}

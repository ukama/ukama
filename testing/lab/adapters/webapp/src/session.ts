/* SPDX-License-Identifier: MPL-2.0
 * Auth acceptance uses visible state and real saved test sessions. No login
 * tokens are minted, decoded for expectations, or returned in command results.
 */
import { performance } from 'node:perf_hooks';
import type { Locator, Page } from 'playwright';
import { Budget, WorkerError, keys, normalize, str, type ObjectValue } from './contract.js';

export const SESSION_PATHS = ['/', '/business', '/business/settings', '/network/settings', '/customer/settings', '/business/manage/members', '/business/manage/data-plans', '/business/manage/sim-pool', '/customer/customers', '/welcome', '/unauthorized', '/business/manage/billing'];
const settingsFields = ['Full name', 'Email', 'Role', 'Email verified', 'Organization name', 'Country', 'Currency'];
const welcomeFields = ['Network operating country', 'Organization name', 'Role'];
const controls = ['Invite member', 'Create plan', 'Upload SIMs', 'Add customer', 'Continue'];
const navLabels = ['Home', 'Revenue', 'Customers', 'Packages', 'Data plans', 'Members', 'SIM pool', 'Node pool', 'Sites', 'Nodes', 'Settings', 'Billing', 'Support'];
const unauthorized = "Your account isn't set up for this console";
const escape = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
async function oneText(locator: Locator): Promise<string | null> {
  const visible = locator.filter({visible: true}), n = await visible.count();
  if (n > 1) throw new WorkerError('AMBIGUOUS_LOCATOR', 'Session field matches multiple visible elements');
  return n === 1 ? normalize(await visible.innerText({timeout: 500})) : null;
}
export class Session {
  private lastNavigation = performance.now();
  private navigationCount = 0;
  private refreshSeen = false;
  private logoutSeen = false;
  private document?: {url: string; status: number};
  constructor(private page: Page, private origin: string, private authOrigin: string) {
    page.on('request', request => {
      if (!request.isNavigationRequest() || request.frame() !== page.mainFrame()) return;
      this.lastNavigation = performance.now(); this.navigationCount++;
      const u = new URL(request.url());
      if (u.origin === origin && u.pathname === '/api/auth/refresh') this.refreshSeen = true;
      if (u.origin === authOrigin && u.pathname === '/user/logout') this.logoutSeen = true;
    });
    page.on('framenavigated', frame => { if (frame === page.mainFrame()) this.lastNavigation = performance.now(); });
    page.on('response', response => {
      const request = response.request();
      if (request.isNavigationRequest() && request.frame() === page.mainFrame()) {
        const u = new URL(response.url()); this.document = {url: u.origin + u.pathname, status: response.status()};
      }
    });
  }
  private main() { return this.page.locator('main.main'); }
  private stable() {
    if (this.navigationCount > 12) throw new WorkerError('REDIRECT_LOOP', 'Too many document navigations in a session step');
    return this.navigationCount > 0 && performance.now() - this.lastNavigation >= 500;
  }
  private async dashboard() {
    return await this.page.locator('header.topbar').isVisible() && await this.main().isVisible() && await this.page.locator('aside.sidebar').isVisible();
  }
  private async surface(): Promise<string | null> {
    if (!this.stable()) return null;
    const u = new URL(this.page.url());
    if (u.origin === this.authOrigin && u.pathname === '/auth/login' && this.document?.url === u.origin + u.pathname && this.document.status < 400 && await this.page.locator('body').isVisible() && (await this.page.locator('body').innerText()).trim()) return 'auth';
    if (u.origin !== this.origin) {
      if (u.origin !== this.authOrigin && u.protocol !== 'about:') throw new WorkerError('WRONG_ORIGIN', 'Session navigated outside its configured console/auth origins');
      return null;
    }
    if (u.pathname === '/unauthorized' && await this.page.getByText(unauthorized, {exact:true}).isVisible()) return 'unauthorized';
    if (u.pathname === '/welcome' && await this.page.locator('main.welcome-root').getByRole('heading', {name:'Welcome to Ukama!',exact:true}).isVisible()) return 'welcome';
    if (this.document?.url === u.origin + u.pathname && this.document.status === 404 && await this.page.getByRole('heading', {name:'404', exact:true}).isVisible()) return 'not_found';
    return await this.dashboard() ? 'dashboard' : null;
  }
  private async blocked() { const surface = await this.surface(); return surface === 'auth' || surface === 'unauthorized'; }
  private async requireConsole() {
    if (new URL(this.page.url()).origin !== this.origin) throw new WorkerError('WRONG_ORIGIN', 'This action requires the configured console');
  }
  async run(i: ObjectValue, budget: Budget) {
    keys(i, ['view','action','value']);
    if (i.view !== 'session') throw new WorkerError('INVALID_INPUT','Session action requires session view');
    const action = str(i.action,'action'), value = str(i.value ?? '', 'value',true);
    if (!['navigate','settings_tab'].includes(action) && value) throw new WorkerError('INVALID_INPUT','This session action takes no value');
    const click = (locator: Locator) => locator.click({timeout: budget.remaining()});
    if (['drop_token','invalidate_token','reject_token','expire_token'].includes(action)) {
      // Change only this isolated context's gateway cookie. Keep the legitimate
      // session, cookie scope and flags; never derive expected UI claims from it.
      const cookies = await this.page.context().cookies(this.origin);
      if (cookies.filter(c=>c.name === 'ukama_session').length !== 1) throw new WorkerError('AUTH_PRECONDITION','Token recovery requires one existing session cookie');
      const tokens = cookies.filter(c=>c.name === 'token');
      if (tokens.length !== 1) throw new WorkerError('AUTH_PRECONDITION','Token fault requires exactly one existing gateway token');
      const token = tokens[0]!;
      if (action === 'drop_token') await this.page.context().clearCookies({name:token.name,domain:token.domain,path:token.path});
      else {
        let changed = 'ukama-lab-invalid-token';
        if (action === 'reject_token') {
          const dot = token.value.lastIndexOf('.');
          if (dot <= 0) throw new WorkerError('AUTH_PRECONDITION','Signature fault requires a signed gateway token');
          changed = token.value.slice(0,dot) + '.ukama-lab-rejected-signature';
        }
        if (action === 'expire_token') {
          const dot = token.value.lastIndexOf('.');
          const claims = Buffer.from(token.value.slice(0,dot), 'base64').toString('utf8').split(';');
          if (dot <= 0 || claims.length !== 11 || !/^\d+$/.test(claims[10]!)) throw new WorkerError('AUTH_PRECONDITION','Expiry fault requires the source-defined gateway token format');
          claims[10] = '1';
          changed = Buffer.from(claims.join(';')).toString('base64') + '.ukama-lab-expired-signature';
        }
        await this.page.context().addCookies([{...token,value:changed}]);
      }
      return;
    }
    if (action === 'navigate') {
      if (!SESSION_PATHS.includes(value)) throw new WorkerError('INVALID_INPUT','Unsupported session route');
      this.navigationCount = 0; this.document = undefined;
      const response = await this.page.goto(this.origin + value,{waitUntil:'domcontentloaded',timeout:budget.remaining()});
      if (response && response.status() >= 500) throw new WorkerError('APP_UNAVAILABLE','Session navigation returned a server error',response.status());
    } else if (action === 'reload') {
      await this.requireConsole(); this.navigationCount = 0;
      await this.page.reload({waitUntil:'domcontentloaded',timeout:budget.remaining()});

    } else if (action === 'settings_tab') {
      await this.requireConsole();
      if (!['/business/settings','/network/settings','/customer/settings'].includes(new URL(this.page.url()).pathname) || !['My account','Organization','Preferences'].includes(value)) throw new WorkerError('WRONG_VIEW','Unknown settings tab or route');
      await click(this.main().getByRole('tab',{name:value,exact:true}));
    } else if (action === 'open_account') {
      await this.requireConsole(); await click(this.page.locator('header.topbar button.avatar[aria-haspopup="menu"]'));
      await this.page.getByRole('menu').getByRole('menuitem',{name:'Log out',exact:true}).waitFor({timeout:budget.remaining()});
    } else if (action === 'logout') {
      await this.requireConsole(); await click(this.page.getByRole('menu').getByRole('menuitem',{name:'Log out',exact:true}));
    } else if (action === 'ack_welcome') {
      await this.requireConsole();
      if (new URL(this.page.url()).pathname !== '/welcome') throw new WorkerError('WRONG_VIEW','Welcome acknowledgement requires the welcome screen');
      await click(this.page.locator('main.welcome-root').getByRole('button',{name:'Continue',exact:true}));
    } else throw new WorkerError('INVALID_INPUT','Unknown session action');
  }
  async check(i: ObjectValue, budget: Budget) {
    keys(i,['view','label','subject','expected','requirement']);
    if (i.view !== 'session') throw new WorkerError('INVALID_INPUT','Session assertion requires session view');
    const label = str(i.label,'label'), subject = str(i.subject ?? '', 'subject',true), expected = str(i.expected,'expected',true);
    const observe = async (): Promise<string | null> => {
      const surface = await this.surface();
      if (!surface) return null;
      const u = new URL(this.page.url());
      if (label === 'Surface') return surface;
      if (label === 'Path') return u.pathname;
      if (label === 'Auth origin') return u.origin;
      if (label === 'Dashboard visible') return String(await this.dashboard());
      if (label === 'Access blocked') return String(await this.blocked());
      if (label === 'Navigation stable') return String(this.stable());
      if (label === 'Refresh observed') return String(this.refreshSeen);
      if (label === 'Logout handoff') return String(this.logoutSeen);
      if (label === 'Document status') return this.document?.url === u.origin + u.pathname ? String(this.document.status) : null;
      if (label === 'Sensitive text absent') return await this.blocked() ? String(await this.page.getByText(subject,{exact:true}).filter({visible:true}).count() === 0) : null;
      if (label === 'Unauthorized title') return surface === 'unauthorized' ? oneText(this.page.getByText(unauthorized,{exact:true})) : null;
      if (label === 'Unauthorized logout') {
        if (surface !== 'unauthorized') return null;
        const link = this.page.getByRole('link',{name:'Log out',exact:true});
        return await link.isVisible() ? String(await link.getAttribute('href') === this.authOrigin + '/user/logout') : null;
      }
      if (label === 'Unauthorized support') {
        if (surface !== 'unauthorized') return null;
        const link = this.page.getByRole('link',{name:'Contact us',exact:true});
        return await link.isVisible() ? String((await link.getAttribute('href'))?.startsWith('mailto:support@ukama.com?') === true) : null;
      }
      if (label === 'Welcome title') return surface === 'welcome' ? oneText(this.page.locator('.welcome-title')) : null;
      if (label === 'Welcome error') return surface === 'welcome' ? oneText(this.page.locator('.welcome-error')) : null;
      if (label === 'Welcome field') {
        if (!welcomeFields.includes(subject)) throw new WorkerError('INVALID_INPUT','Unknown welcome field');
        return surface === 'welcome' ? oneText(this.page.locator('.welcome-field-label').filter({hasText:new RegExp(`^${escape(subject)}$`)}).locator('..').locator('.welcome-field-value')) : null;
      }
      if (label === 'Settings field') {
        if (!settingsFields.includes(subject)) throw new WorkerError('INVALID_INPUT','Unknown settings field');
        if (surface !== 'dashboard' || !u.pathname.endsWith('/settings')) return null;
        return oneText(this.main().locator('.card.card-pad').filter({has:this.page.locator('label.flabel').filter({hasText:new RegExp(`^${escape(subject)}$`)})}).locator('.ff-readonly'));
      }
      if (label === 'Organization') return surface === 'dashboard' ? oneText(this.page.locator('header.topbar [title^="Organization: "]')) : null;
      if (label === 'Account name' || label === 'Account details') return surface === 'dashboard' ? oneText(this.page.getByRole('menu').locator('.MuiTypography-root').nth(label === 'Account name' ? 0 : 1)) : null;
      if (label === 'Nav visible') {
        if (!navLabels.includes(subject)) throw new WorkerError('INVALID_INPUT','Unknown navigation label');
        return surface === 'dashboard' ? String(await this.page.locator('aside.sidebar').getByRole('link',{name:subject,exact:true}).isVisible()) : null;
      }
      if (label === 'Control state') {
        if (!controls.includes(subject)) throw new WorkerError('INVALID_INPUT','Unknown control');
        if (surface === 'auth' || surface === 'unauthorized' || surface === 'not_found') return 'absent';
        const scope = surface === 'welcome' ? this.page.locator('main.welcome-root') : this.main().locator('.pagehead');
        if (!await scope.isVisible() || await this.main().locator('.MuiSkeleton-root:visible').count()) return null;
        const buttons = scope.getByRole('button',{name:subject,exact:true}).filter({visible:true}), count = await buttons.count();
        if (count > 1) throw new WorkerError('AMBIGUOUS_LOCATOR','Session control is ambiguous');
        return count === 0 ? 'absent' : await buttons.isEnabled() ? 'enabled' : 'disabled';
      }
      throw new WorkerError('UNSUPPORTED_LOCATOR','Unknown session assertion');
    };
    const actual = await budget.poll(observe, value=>value !== null && value === expected, 'Session UI differs from scenario expectation');
    return {expected,actual};
  }
}

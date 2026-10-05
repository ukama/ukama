/* SPDX-License-Identifier: MPL-2.0
 * Controlled read faults; never create inventory or mutate application state.
 */
import type { Page, Request, Route } from 'playwright';
import { Budget, WorkerError } from './contract.js';
type Read = { operation: string; variables: Record<string, unknown> };
type Fault = { path: string; mode: string; applied: number; release: () => void; handler: (r: Route) => Promise<void> };
const endpoints = new WeakMap<Page, Map<string, Set<string>>>();
const poolScopes = new WeakMap<Page, Set<string>>();
const faults = new WeakMap<Page, Fault>();
function read(request: Request): Read | undefined {
  try {
    const b = request.postDataJSON();
    if (request.method() !== 'POST' || !new URL(request.url()).pathname.endsWith('/graphql') ||
        !['isPackageNameAvailable', 'SimPoolOverview'].includes(b?.operationName) ||
        typeof b.query !== 'string' || !new RegExp(`^\\s*query\\s+${b.operationName}\\s*\\(`).test(b.query) || /\bmutation\b/.test(b.query) ||
        !b.variables || typeof b.variables !== 'object' || Array.isArray(b.variables)) return;
    if (b.operationName === 'isPackageNameAvailable' && (Object.keys(b.variables).length !== 1 || typeof b.variables.name !== 'string')) return;
    if (b.operationName === 'SimPoolOverview' && (Object.keys(b.variables).length !== 2 || typeof b.variables.simType !== 'string' || b.variables.limit !== 100)) return;
    return { operation: b.operationName, variables: b.variables };
  } catch { return; }
}
export function trackCommerceReads(page: Page) {
  if (endpoints.has(page)) return;
  const known = new Map<string, Set<string>>(); endpoints.set(page, known); poolScopes.set(page, new Set());
  page.on('close', () => { faults.get(page)?.release(); faults.delete(page); endpoints.delete(page); poolScopes.delete(page); });
  page.on('request', request => {
    const r = read(request); if (!r) return;
    const urls = known.get(r.operation) ?? new Set<string>(); urls.add(request.url()); known.set(r.operation, urls);
    if (r.operation === 'SimPoolOverview') poolScopes.get(page)!.add(String(r.variables.simType));
  });
  page.on('framenavigated', frame => { const f = faults.get(page); if (f && frame === page.mainFrame()) f.applied = 0; });
}
export async function commerceFault(page: Page, mode: string, name: string | undefined, budget: Budget) {
  const previous = faults.get(page);
  if (previous) { faults.delete(page); previous.release(); await page.unroute('**/graphql', previous.handler); }
  if (mode === 'clear_fault') return;
  if (!['name_pending', 'name_failure', 'pool_failure'].includes(mode)) throw new WorkerError('INVALID_INPUT', 'Unsupported commerce read fault');
  const operation = mode === 'pool_failure' ? 'SimPoolOverview' : 'isPackageNameAvailable';
  if (operation === 'isPackageNameAvailable' && !name) throw new WorkerError('INVALID_INPUT', 'Name fault requires its owned planned name');
  const urls = await budget.poll(async () => endpoints.get(page)?.get(operation), x => !!x?.size, 'Read must be observed before arming its fault');
  if (urls!.size !== 1) throw new WorkerError('AMBIGUOUS_ENDPOINT', 'Commerce read used multiple endpoints');
  const endpoint = [...urls!][0];
  const scopes = poolScopes.get(page)!;
  if (operation === 'SimPoolOverview' && scopes.size !== 1) throw new WorkerError('AMBIGUOUS_SCOPE', 'Inventory read must have one observed SIM type');
  const simType = [...scopes][0];
  let release!: () => void;
  const held = new Promise<void>(resolve => { release = resolve; });
  const fault: Fault = { mode, path: new URL(page.url()).pathname, applied: 0, release, handler: async route => {
    const r = read(route.request());
    if (faults.get(page) !== fault || new URL(page.url()).pathname !== fault.path || route.request().url() !== endpoint ||
        r?.operation !== operation || (name && r.variables.name !== name) || (operation === 'SimPoolOverview' && r.variables.simType !== simType)) { await route.fallback(); return; }
    fault.applied++;
    if (mode === 'name_pending') { await held; await route.abort('failed').catch(() => {}); return; }
    const body = mode === 'pool_failure' ? {data: {simPoolView: {simType: r.variables.simType,
      stats: {error: {code: 'LAB_READ_FAILURE', message: 'Lab controlled inventory read failure'}, total: null, available: null, consumed: null, failed: null},
      sims: {error: {code: 'LAB_READ_FAILURE', message: 'Lab controlled inventory read failure'}, sims: []}}}} :
      {errors: [{message: 'Lab controlled plan-name read failure'}]};
    await route.fulfill({status: 200, contentType: 'application/json', headers: {'access-control-allow-origin': new URL(page.url()).origin, 'access-control-allow-credentials': 'true'}, body: JSON.stringify(body)});
  }};
  faults.set(page, fault); await page.route('**/graphql', fault.handler);
}
export function commerceFaultState(page: Page) {
  const f = faults.get(page);
  return !f ? 'none' : new URL(page.url()).pathname !== f.path ? `${f.mode}:wrong_scope` : `${f.mode}:${f.applied ? 'applied' : 'armed'}`;
}
export function commerceEndpoint(page: Page): string {
  const urls = new Set([...(endpoints.get(page)?.values() ?? [])].flatMap(set => [...set]));
  if (urls.size !== 1) throw new WorkerError('AMBIGUOUS_ENDPOINT', 'Payment rejection requires one passively observed commerce endpoint');
  return [...urls][0]!;
}

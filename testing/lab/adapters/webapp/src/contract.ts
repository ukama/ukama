/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { performance } from 'node:perf_hooks';

export type ObjectValue = Record<string, unknown>;
export class WorkerError extends Error {
  constructor(public code: string, message: string, public actual?: unknown) {
    super(message);
  }
}
export function object(value: unknown, label = 'object'): ObjectValue {
  if (!value || typeof value !== 'object' || Array.isArray(value))
    throw new WorkerError('INVALID_INPUT', `${label} must be an object`);
  return value as ObjectValue;
}
export function keys(value: ObjectValue, allowed: string[]): void {
  for (const key of Object.keys(value)) if (!allowed.includes(key))
    throw new WorkerError('INVALID_INPUT', `Unsupported field: ${key}`);
}
export function str(value: unknown, label: string, empty = false): string {
  if (typeof value !== 'string' || value.length > 4096 || (!empty && !value.trim()))
    throw new WorkerError('INVALID_INPUT', `${label} must be a string${empty ? '' : ' with content'}`);
  return value;
}
export function integer(value: unknown, label: string, min: number, max: number): number {
  if (!Number.isSafeInteger(value) || Number(value) < min || Number(value) > max)
    throw new WorkerError('INVALID_INPUT', `${label} is outside ${min}..${max}`);
  return value as number;
}
export function bool(value: unknown, label: string): boolean {
  if (typeof value !== 'boolean') throw new WorkerError('INVALID_INPUT', `${label} must be boolean`);
  return value;
}
export function baseURL(value: unknown): string {
  let url: URL;
  try { url = new URL(str(value, 'base_url')); }
  catch { throw new WorkerError('INVALID_INPUT', 'base_url must be an absolute http(s) URL'); }
  if (!['http:', 'https:'].includes(url.protocol) || url.username || url.password || url.search || url.hash || url.pathname !== '/')
    throw new WorkerError('INVALID_INPUT', 'base_url must be an http(s) origin without credentials, query or path prefix');
  return url.origin;
}
export interface Profile {
  base_url: string; auth_state: string; browser: 'chromium' | 'firefox' | 'webkit';
  headless: boolean; action_timeout_seconds: number; check_timeout_seconds: number;
  scenario_timeout_seconds: number;
  session_mode: 'authenticated' | 'auth_test' | 'onboarding'; auth_origin?: string;
}
export function profile(value: unknown): Profile {
  const p = object(value, 'profile');
  keys(p, ['base_url', 'auth_state', 'browser', 'headless', 'action_timeout_seconds', 'check_timeout_seconds', 'scenario_timeout_seconds', 'session_mode', 'auth_origin']);
  const browser = p.browser ?? 'chromium';
  if (!['chromium', 'firefox', 'webkit'].includes(String(browser)))
    throw new WorkerError('INVALID_INPUT', 'Unknown browser');
  const mode = p.session_mode ?? 'authenticated';
  if (mode !== 'authenticated' && mode !== 'auth_test' && mode !== 'onboarding') throw new WorkerError('INVALID_INPUT','Unknown session_mode');
  const result: Profile = {
    session_mode: mode,
    base_url: baseURL(p.base_url), auth_state: str(p.auth_state, 'auth_state'),
    browser: browser as Profile['browser'], headless: bool(p.headless ?? true, 'headless'),
    action_timeout_seconds: integer(p.action_timeout_seconds ?? 30, 'action timeout', 1, 900),
    check_timeout_seconds: integer(p.check_timeout_seconds ?? 30, 'check timeout', 1, 900),
    scenario_timeout_seconds: integer(p.scenario_timeout_seconds ?? 3600, 'scenario timeout', 1, 86400),
  };
  if (Math.max(result.action_timeout_seconds, result.check_timeout_seconds) > result.scenario_timeout_seconds)
    throw new WorkerError('INVALID_INPUT', 'Step timeout exceeds scenario timeout');
  if (mode === 'auth_test') {
    result.auth_origin = baseURL(p.auth_origin);
    if (result.auth_origin === result.base_url) throw new WorkerError('INVALID_INPUT','auth_origin must differ from console origin');
  } else if (p.auth_origin !== undefined || result.auth_state === 'none') throw new WorkerError('INVALID_INPUT','Auth test fields require auth_test mode');
  return result;
}
export const ACTIONS = ['init', 'web_inventory', 'web_inventory_equals', 'web_onboard', 'web_onboard_equals', 'web_session', 'web_session_equals', 'web_interact', 'web_ui_equals', 'web_commerce', 'web_commerce_equals', 'web_import_sims', 'web_create_network', 'web_create_site', 'web_open', 'web_select_network', 'web_reload', 'web_action', 'web_tab', 'web_kpi_equals', 'web_field_equals', 'web_table_count_equals', 'web_action_available', 'close'] as const;
export type Action = typeof ACTIONS[number];
export interface Command {
  protocol: 1; run_id: string; command_id: number; action: Action;
  deadline_ms: number; inputs: ObjectValue;
}
export function command(value: unknown): Command {
  const c = object(value, 'command');
  keys(c, ['protocol', 'run_id', 'command_id', 'action', 'deadline_ms', 'inputs']);
  if (c.protocol !== 1) throw new WorkerError('PROTOCOL_VERSION', 'Expected protocol 1');
  const run = str(c.run_id, 'run_id');
  if (!/^[A-Za-z0-9][A-Za-z0-9_-]{0,79}$/.test(run)) throw new WorkerError('INVALID_INPUT', 'Invalid run_id');
  if (!ACTIONS.includes(c.action as Action)) throw new WorkerError('UNSUPPORTED_ACTION', 'Unknown worker action');
  return { protocol: 1, run_id: run, command_id: integer(c.command_id, 'command_id', 1, 10000),
    action: c.action as Action, deadline_ms: integer(c.deadline_ms, 'deadline_ms', 1, Number.MAX_SAFE_INTEGER), inputs: object(c.inputs, 'inputs') };
}
export interface Result {
  protocol: 1; run_id: string | null; command_id: number | null; action: string | null;
  status: 'ok' | 'error'; run_status: 'running' | 'passed' | 'failed'; duration_ms: number;
  expected: unknown; actual: unknown; bindings: unknown[]; artifacts: string[];
  error?: { code: string; message: string };
}
export function normalize(value: string): string { return value.replace(/\s+/g, ' ').trim(); }
export function safeURL(value: string): string {
  try { const u = new URL(value); return `${u.origin}${u.pathname}`; } catch { return '<unavailable>'; }
}
// Normalize key order so replay checks do not depend on a serializer's ordering.
export function canonical(value: unknown): string {
  if (Array.isArray(value)) return `[${value.map(canonical).join(',')}]`;
  if (value && typeof value === 'object') return `{${Object.keys(value).sort().map(k => `${JSON.stringify(k)}:${canonical((value as ObjectValue)[k])}`).join(',')}}`;
  return JSON.stringify(value);
}
export class Budget {
  private end: number;
  actual: unknown = null;
  constructor(milliseconds: number) { this.end = performance.now() + milliseconds; }
  remaining(): number {
    const ms = Math.floor(this.end - performance.now());
    if (ms <= 0) throw new WorkerError('DEADLINE_EXCEEDED', 'Command deadline exceeded');
    return ms;
  }
  async poll<T>(observe: () => Promise<T>, matches: (value: T) => boolean, message: string): Promise<T> {
    let actual: T | undefined;
    for (;;) {
      this.remaining();
      actual = await observe(); this.actual = actual;
      if (matches(actual)) return actual;
      const ms = this.end - performance.now();
      if (ms <= 0) throw new WorkerError('ASSERTION_FAILED', message, actual);
      await new Promise(resolve => setTimeout(resolve, Math.min(100, ms)));
    }
  }
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

export const DEFAULT_FANOUT_CONCURRENCY = 10;

/**
 * Map over `items` with `fn`, keeping at most `concurrency` promises in
 * flight. Results preserve input order. Rejections propagate (use
 * `mapWithConcurrencySettled` for per-item results).
 */
export async function mapWithConcurrency<T, R>(
  items: readonly T[],
  fn: (item: T, index: number) => Promise<R>,
  concurrency: number = DEFAULT_FANOUT_CONCURRENCY
): Promise<R[]> {
  if (items.length === 0) return [];
  const limit = Math.max(1, Math.min(concurrency, items.length));
  const results: R[] = new Array(items.length);
  let next = 0;

  const worker = async (): Promise<void> => {
    for (;;) {
      const index = next++;
      if (index >= items.length) return;
      results[index] = await fn(items[index], index);
    }
  };

  await Promise.all(Array.from({ length: limit }, worker));
  return results;
}

/**
 * Like `mapWithConcurrency`, but never rejects: returns a
 * PromiseSettledResult per item, preserving input order.
 */
export async function mapWithConcurrencySettled<T, R>(
  items: readonly T[],
  fn: (item: T, index: number) => Promise<R>,
  concurrency: number = DEFAULT_FANOUT_CONCURRENCY
): Promise<PromiseSettledResult<R>[]> {
  return mapWithConcurrency(
    items,
    async (item, index): Promise<PromiseSettledResult<R>> => {
      try {
        return { status: "fulfilled", value: await fn(item, index) };
      } catch (reason) {
        return { status: "rejected", reason };
      }
    },
    concurrency
  );
}

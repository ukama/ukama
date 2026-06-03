/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Tinted summary row — health summary / recent activity (biz-common.jsx). */
export type FeedTone = 'ok' | 'info' | 'warn' | 'err';

const TONE_COLOR: Record<FeedTone, string> = {
  ok: 'var(--uk-success-bright)',
  info: 'var(--uk-ac)',
  warn: 'var(--uk-orange)',
  err: 'var(--uk-error)',
};

export default function FeedRow({
  tone,
  title,
  detail,
  when,
}: {
  tone: FeedTone;
  title: string;
  detail: string;
  when?: string;
}) {
  return (
    <div style={{ display: 'flex', alignItems: 'flex-start', gap: 11, padding: '11px 0' }}>
      <span
        style={{
          width: 9,
          height: 9,
          borderRadius: '50%',
          flex: 'none',
          marginTop: 5,
          background: TONE_COLOR[tone],
        }}
      />
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{ fontSize: 13.5, fontWeight: 600, color: 'var(--uk-ink)' }}>{title}</div>
        <div style={{ fontSize: 12.5, color: 'var(--uk-ink-2)', marginTop: 1 }}>{detail}</div>
      </div>
      {when && (
        <span style={{ fontSize: 11.5, color: 'var(--uk-ink-3)', whiteSpace: 'nowrap' }}>
          {when}
        </span>
      )}
    </div>
  );
}

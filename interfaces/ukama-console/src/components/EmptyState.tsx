/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Empty / error states with cohesive line-art spot illustrations
 * (table-kit.jsx) — they explain what happened and offer the next action.
 */
import Button from '@mui/material/Button';
import AddRounded from '@mui/icons-material/AddRounded';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';

export type EmptyArtName =
  | 'people'
  | 'sim'
  | 'node'
  | 'site'
  | 'invoice'
  | 'search'
  | 'error';

export function EmptyArt({
  name,
  size = 116,
  tone,
}: {
  name: EmptyArtName;
  size?: number;
  tone?: 'error';
}) {
  const stroke = 'var(--uk-ink-3)';
  const acc = tone === 'error' ? 'var(--uk-orange)' : 'var(--uk-ac)';
  const c = {
    fill: 'none',
    stroke,
    strokeWidth: 3,
    strokeLinecap: 'round' as const,
    strokeLinejoin: 'round' as const,
  };
  const a = { ...c, stroke: acc };

  const art: Record<EmptyArtName, React.ReactNode> = {
    people: (
      <>
        <circle cx="47" cy="49" r="11" {...c} />
        <path d="M29 81c0-12 9-19 18-19s18 7 18 19" {...c} />
        <circle cx="79" cy="43" r="9" {...a} />
        <path d="M65 75c1-10 8-15 16-15 8 0 14 5 16 14" {...a} />
      </>
    ),
    sim: (
      <>
        <path d="M40 31h28l15 15v44a4 4 0 0 1-4 4H40a4 4 0 0 1-4-4V35a4 4 0 0 1 4-4z" {...c} />
        <rect x="50" y="60" width="24" height="23" rx="3" {...a} />
        <path d="M62 60v23M50 71.5h24" {...a} />
      </>
    ),
    node: (
      <>
        <rect x="33" y="63" width="54" height="26" rx="4" {...c} />
        <circle cx="45" cy="76" r="2.6" fill={acc} stroke="none" />
        <path d="M55 53a18 18 0 0 1 30 0M50 45a30 30 0 0 1 40 0" {...a} />
      </>
    ),
    site: (
      <>
        <path d="M60 88s21-17 21-34a21 21 0 1 0-42 0c0 17 21 34 21 34z" {...c} />
        <circle cx="60" cy="52" r="8.5" {...a} />
      </>
    ),
    invoice: (
      <>
        <path d="M43 29h25l13 13v45a4 4 0 0 1-4 4H43a4 4 0 0 1-4-4V33a4 4 0 0 1 4-4z" {...c} />
        <path d="M51 57h28M51 67h28M51 77h18" {...c} />
        <path d="M66 29v13h13" {...a} />
      </>
    ),
    search: (
      <>
        <circle cx="54" cy="52" r="20" {...c} />
        <path d="M69 67l16 16" {...a} />
      </>
    ),
    error: (
      <>
        <path d="M60 31l31 54H29z" {...a} />
        <path d="M60 56v15" {...c} />
        <circle cx="60" cy="79" r="1.8" fill={stroke} stroke="none" />
      </>
    ),
  };

  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 120 120"
      style={{ display: 'block', margin: '0 auto' }}
      aria-hidden="true"
    >
      <circle cx="60" cy="60" r="52" fill="var(--uk-ac-soft)" opacity={tone === 'error' ? 0.55 : 1} />
      {tone === 'error' && <circle cx="60" cy="60" r="52" fill="rgba(226,116,41,.10)" />}
      {art[name]}
    </svg>
  );
}

export function EmptyState({
  art = 'search',
  title,
  sub,
  cta,
  onCta,
}: {
  art?: EmptyArtName;
  title: string;
  sub?: string;
  cta?: string;
  onCta?: () => void;
}) {
  return (
    <div style={{ padding: '46px 24px 50px', textAlign: 'center' }}>
      <EmptyArt name={art} />
      <div
        style={{
          fontFamily: 'var(--font-display)',
          fontSize: 18,
          fontWeight: 500,
          marginTop: 14,
          color: 'var(--uk-ink)',
        }}
      >
        {title}
      </div>
      {sub && (
        <div
          style={{
            fontSize: 13.5,
            color: 'var(--uk-ink-2)',
            margin: '6px auto 0',
            maxWidth: 340,
            textWrap: 'pretty',
          }}
        >
          {sub}
        </div>
      )}
      {cta && (
        <Button
          variant="contained"
          startIcon={<AddRounded />}
          sx={{ mt: 2.25 }}
          onClick={onCta}
        >
          {cta}
        </Button>
      )}
    </div>
  );
}

export function ErrorState({
  title = 'We couldn’t load this',
  sub = 'Something went wrong fetching this data. Check your connection and try again.',
  onRetry,
}: {
  title?: string;
  sub?: string;
  onRetry?: () => void;
}) {
  return (
    <div style={{ padding: '46px 24px 50px', textAlign: 'center' }}>
      <EmptyArt name="error" tone="error" />
      <div
        style={{
          fontFamily: 'var(--font-display)',
          fontSize: 18,
          fontWeight: 500,
          marginTop: 14,
          color: 'var(--uk-ink)',
        }}
      >
        {title}
      </div>
      <div
        style={{
          fontSize: 13.5,
          color: 'var(--uk-ink-2)',
          margin: '6px auto 0',
          maxWidth: 360,
          textWrap: 'pretty',
        }}
      >
        {sub}
      </div>
      <Button
        variant="outlined"
        startIcon={<ChevronRightRounded />}
        sx={{ mt: 2.25 }}
        onClick={onRetry}
      >
        Try again
      </Button>
    </div>
  );
}

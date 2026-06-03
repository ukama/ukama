/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Root error boundary (BUILD-PLAN §13.3) — minimal, no theme dependencies. */
export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <html lang="en">
      <body style={{ fontFamily: 'sans-serif', padding: '3rem' }}>
        <h2>Something went wrong</h2>
        <p style={{ color: '#5a5e66' }}>
          {error.digest ? `Reference: ${error.digest}` : 'Unexpected error.'}
        </p>
        <button type="button" onClick={() => reset()}>
          Try again
        </button>
      </body>
    </html>
  );
}

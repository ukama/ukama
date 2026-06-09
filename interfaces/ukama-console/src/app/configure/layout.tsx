/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Suspense } from 'react';
import ConfigureShell from './_components/ConfigureShell';
import './configure.css';

export const metadata = { title: 'Setup · Ukama Console' };

/** /configure onboarding flow — own split-panel chrome, no dashboard shell. */
export default function ConfigureLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="cfg-page">
      <div className="cfg-gradient" />
      {/* useSearchParams in the shell requires a Suspense boundary at build time. */}
      <Suspense>
        <ConfigureShell>{children}</ConfigureShell>
      </Suspense>
    </div>
  );
}

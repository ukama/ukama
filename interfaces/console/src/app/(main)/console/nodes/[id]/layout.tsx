/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import type { Metadata } from 'next';

interface Props {
  params: { id: string };
  children: React.ReactNode;
}

export function generateMetadata({ params }: Props): Metadata {
  return {
    title: `Node ${params.id} — Ukama Console`,
    description: `Monitor and manage node ${params.id}`,
  };
}

export default function Layout({ children }: { children: React.ReactNode }) {
  return <>{children}</>;
}

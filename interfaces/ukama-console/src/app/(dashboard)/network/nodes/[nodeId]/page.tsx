/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import PageStub from '@/components/PageStub';

export default async function NodeDetailPage({
  params,
}: {
  params: Promise<{ nodeId: string }>;
}) {
  const { nodeId } = await params;
  return (
    <PageStub title={decodeURIComponent(nodeId)} crumb={['Nodes']} phase="5" />
  );
}

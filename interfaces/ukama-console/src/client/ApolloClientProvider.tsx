/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Client-only Apollo provider (plan §15.5 — client + skeletons).
 * Polling-only v1: a single gateway client. The separate metrics client
 * (+ SSE subscriptions) comes with the metrics phase.
 */
import { useState } from 'react';
import { ApolloProvider } from '@apollo/client';
import { makeApolloClient } from './apollo';

export default function ApolloClientProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [client] = useState(makeApolloClient);

  return <ApolloProvider client={client}>{children}</ApolloProvider>;
}

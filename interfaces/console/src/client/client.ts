/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { ApolloClient, HttpLink, InMemoryCache } from '@apollo/client';

export const getMetricsClient = (baseUrl: string) => {
  return new ApolloClient({
    link: new HttpLink({ uri: `${baseUrl}/graphql`, credentials: 'include' }),
    cache: new InMemoryCache({
      typePolicies: {
        Query: {
          fields: {
            getMetricsStat: { merge: false },
            getMetricByTab: { merge: false },
            getMetricBySite: { merge: false },
            getSiteStat: { merge: false },
            getNotifications: { merge: false },
          },
        },
      },
    }),
  });
};

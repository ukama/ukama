/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { metricsClient } from '@/client/ApolloClient';
import {
  Graphs_Type,
  useGetMetricByTabSubSubscription,
} from '@/generated/metrics';
import PubSub from 'pubsub-js';

interface IMetricSubscription {
  from: number;
  type: Graphs_Type;
}

const MetricSubscription = ({ from, type }: IMetricSubscription) => {
  useGetMetricByTabSubSubscription({
    client: metricsClient,
    variables: {
      from: from,
      type: type,
      orgId: 'ukama',
      userId: 'salman',
      nodeId: 'uk-test36-hnode-a1-00ff',
    },
    onData: (data) => {
      PubSub.publish(data.data.data?.getMetricByTabSub.type || '', [
        Math.floor(data.data.data?.getMetricByTabSub.value[0] || 0) * 1000,
        data.data.data?.getMetricByTabSub.value[1],
      ]);
    },
  });
  return <div></div>;
};

export default MetricSubscription;

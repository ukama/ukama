/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Graphs_Type } from '@/client/graphql/generated/subscriptions';

interface IMetricSubscription {
  from: number;
  type: Graphs_Type;
}

const MetricSubscription = ({ from, type }: IMetricSubscription) => {
  console.log(from, type);
  // const { env } = useAppContext();
  // useGetMetricByTabSubSubscription({
  //   // client: getMetricsClient(env.METRIC_URL, env.METRIC_WEBSOCKET_URL),
  //   variables: {
  //     from: from,
  //     type: type,
  //     orgId: 'ukama',
  //     userId: 'salman',
  //     nodeId: 'uk-test36-hnode-a1-00ff',
  //   },
  //   onData: (data) => {
  //     PubSub.publish(data.data.data?.getMetricByTabSub.type ?? '', [
  //       Math.floor(data.data.data?.getMetricByTabSub.value[0] ?? 0) * 1000,
  //       data.data.data?.getMetricByTabSub.value[1],
  //     ]);
  //   },
  // });
  return <div></div>;
};

export default MetricSubscription;

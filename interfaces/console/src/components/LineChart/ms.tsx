// /*
//  * This Source Code Form is subject to the terms of the Mozilla Public
//  * License, v. 2.0. If a copy of the MPL was not distributed with this
//  * file, You can obtain one at https://mozilla.org/MPL/2.0/.
//  *
//  * Copyright (c) 2023-present, Ukama Inc.
//  */

// import {
//   Graphs_Type,
//   useGetMetricByTabSubSubscription,
// } from '@/client/graphql/generated/subscriptions';
// import { useAppContext } from '@/context';

// interface IMetricSubscription {
//   from: number;
//   nodeId: string;
//   type: Graphs_Type;
// }

// const MetricSubscription = ({ from, type, nodeId }: IMetricSubscription) => {
//   const { user, subscriptionClient } = useAppContext();
//   useGetMetricByTabSubSubscription({
//     client: subscriptionClient,
//     variables: {
//       data: {
//         from: from,
//         type: type,
//         nodeId: nodeId,
//         userId: user.id,
//         orgName: user.orgName,
//       },
//     },
//     onData: (data) => {
//       console.log(data);
//       PubSub.publish(data.data.data?.getMetricByTabSub.type ?? '', [
//         Math.floor(data.data.data?.getMetricByTabSub.value[0] ?? 0) * 1000,
//         data.data.data?.getMetricByTabSub.value[1],
//       ]);
//     },
//   });
//   return <div></div>;
// };

// export default MetricSubscription;

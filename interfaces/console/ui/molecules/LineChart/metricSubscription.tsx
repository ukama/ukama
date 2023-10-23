import { metricsClient } from '@/client/ApolloClient';
import { useGetMetricByTabSubSubscription } from '@/generated/metrics';
import { getNodeTabTypeByIndex } from '@/utils';
import PubSub from 'pubsub-js';

interface IMetricSubscription {
  from: number;
}

const MetricSubscription = ({ from }: IMetricSubscription) => {
  useGetMetricByTabSubSubscription({
    client: metricsClient,
    variables: {
      from: from,
      orgId: 'ukama',
      userId: 'salman',
      type: getNodeTabTypeByIndex(0),
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

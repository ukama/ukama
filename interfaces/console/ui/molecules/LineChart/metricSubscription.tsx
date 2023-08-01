import { metricsClient } from '@/client/ApolloClient';
import { useMetricRangeSubscription } from '@/generated/metrics';
import PubSub from 'pubsub-js';

interface IMetricSubscription {
  from: number;
}

const MetricSubscription = ({ from }: IMetricSubscription) => {
  useMetricRangeSubscription({
    client: metricsClient,
    variables: {
      orgId: '123',
      userId: 'salman',
      from: from,
      type: 'memory_trx_used',
      nodeId: 'uk-123456-hnode-77-8888',
    },
    onData: (data) => {
      PubSub.publish(
        'memory_trx_used',
        data.data.data?.getMetricRangeSub.value,
      );
    },
  });
  return <div></div>;
};

export default MetricSubscription;

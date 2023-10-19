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
      type: 'uptime_trx',
      nodeId: 'uk-test36-hnode-a1-00ff',
    },
    onData: (data) => {
      PubSub.publish(
        'uptime_trx',
        data.data.data?.getMetricRangeSub.value,
      );
    },
  });
  return <div></div>;
};

export default MetricSubscription;

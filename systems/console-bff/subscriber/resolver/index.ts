import { NonEmptyArray } from "type-graphql";

import { AddSubscriberResolver } from "./addSubscriber";
import { DeleteSubscriberResolver } from "./deleteSubscriber";
import { GetSubscriberResolver } from "./getSubscriber";
import { GetSubscriberMetricsByNetworkResolver } from "./getSubscriberMetricsByNetwork";
import { GetSubscribersByNetworkResolver } from "./getSubscribersByNetwork";
import { UpdateSubscriberResolver } from "./updateSubscriber";

const resolvers: NonEmptyArray<any> = [
  AddSubscriberResolver,
  DeleteSubscriberResolver,
  GetSubscriberResolver,
  GetSubscriberMetricsByNetworkResolver,
  GetSubscribersByNetworkResolver,
  UpdateSubscriberResolver,
];

export default resolvers;

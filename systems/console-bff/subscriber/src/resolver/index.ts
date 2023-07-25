import { NonEmptyArray } from "type-graphql";

import { AddSubscriberResolver } from "./addSubscriber.resolver";
import { DeleteSubscriberResolver } from "./deleteSubscriber.resolver";
import { GetSubscriberResolver } from "./getSubscriber.resolver";
import { GetSubscriberMetricsByNetworkResolver } from "./getSubscriberMetricsByNetwork.resolver";
import { GetSubscribersByNetworkResolver } from "./getSubscribersByNetwork.resolver";
import { UpdateSubscriberResolver } from "./updateSubscriber.resolver";


const resolvers: NonEmptyArray<Function> = [AddSubscriberResolver,
    DeleteSubscriberResolver,
    GetSubscriberResolver,
    GetSubscriberMetricsByNetworkResolver,
    GetSubscribersByNetworkResolver,
    UpdateSubscriberResolver];

export default resolvers;

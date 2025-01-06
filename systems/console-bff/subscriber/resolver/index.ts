/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { AddSubscriberResolver } from "./addSubscriber";
import { DeleteSubscriberResolver } from "./deleteSubscriber";
import { GetSimsByNetworkResolver } from "./getSimsByNetwork";
import { GetSubscriberResolver } from "./getSubscriber";
import { GetSubscriberMetricsByNetworkResolver } from "./getSubscriberMetricsByNetwork";
import { GetSubscribersByNetworkResolver } from "./getSubscribersByNetwork";
import { UpdateSubscriberResolver } from "./updateSubscriber";

const resolvers: NonEmptyArray<any> = [
  AddSubscriberResolver,
  DeleteSubscriberResolver,
  GetSimsByNetworkResolver,
  GetSubscriberResolver,
  GetSubscriberMetricsByNetworkResolver,
  GetSubscribersByNetworkResolver,
  UpdateSubscriberResolver,
];

export default resolvers;

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { GetCorrespondentsResolver } from "./getCorrespondents";
import { GetPaymentResolver } from "./getPayment";
import { GetPaymentsResolver } from "./getPayments";
import { GetTokenResolver } from "./getToken";
import { ProcessPaymentResolver } from "./processPayment";
import { UpdatePaymentResolver } from "./updatePayment";

const resolvers: NonEmptyArray<any> = [
  GetTokenResolver,
  GetPaymentResolver,
  GetPaymentsResolver,
  UpdatePaymentResolver,
  ProcessPaymentResolver,
  GetCorrespondentsResolver,
];

export default resolvers;

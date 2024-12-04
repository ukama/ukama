/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { CreateCustomerResolver } from "./createCustomer";
import { GetReportResolver } from "./getReport";
import { GetStripeCustomerResolver } from "./getStripeCustomer";
import { RetrivePaymentMethodsResolver } from "./retrivePaymentMethods";

const resolvers: NonEmptyArray<any> = [
  GetReportResolver,
  CreateCustomerResolver,
  GetStripeCustomerResolver,
  RetrivePaymentMethodsResolver,
];

export default resolvers;

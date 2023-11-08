/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NonEmptyArray } from "type-graphql";

import { AttachPaymentWithCustomerResolver } from "./attachPaymentWithCustomer";
import { CreateCustomerResolver } from "./createCustomer";
import { GetBillHistoryResolver } from "./getBillHistory";
import { GetCurrentBillResolver } from "./getCurrentBill";
import { GetStripeCustomerResolver } from "./getStripeCustomer";
import { RetrivePaymentMethodsResolver } from "./retrivePaymentMethods";

const resolvers: NonEmptyArray<any> = [
  AttachPaymentWithCustomerResolver,
  CreateCustomerResolver,
  GetBillHistoryResolver,
  GetCurrentBillResolver,
  GetStripeCustomerResolver,
  RetrivePaymentMethodsResolver,
];

export default resolvers;

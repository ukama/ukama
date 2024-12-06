import { NonEmptyArray } from "type-graphql";

import { AddPaymentResolver } from "./addPayment";
import { GetCorrespondentsResolver } from "./getCorrespondents";
import { GetPaymentResolver } from "./getPayment";
import { GetPaymentsResolver } from "./getPayments";
import { GetTokenResolver } from "./getToken";
import { ProcessPaymentResolver } from "./processPayment";
import { UpdatePaymentResolver } from "./updatePayment";

const resolvers: NonEmptyArray<any> = [
  GetTokenResolver,
  AddPaymentResolver,
  GetPaymentResolver,
  GetPaymentsResolver,
  UpdatePaymentResolver,
  ProcessPaymentResolver,
  GetCorrespondentsResolver,
];

export default resolvers;

import { NonEmptyArray } from "type-graphql";

import { AttachPaymentWithCustomerResolver } from "./attachPaymentWithCustomer";
import { CreateCustomerResolver } from "./createCustomer";
import { GetBillHistoryResolver } from "./getBillHistory";
import { GetCurrentBillResolver } from "./getCurrentBill";
import { GetStripeCustomerResolver } from "./getStripeCustomer";
import { RetrivePaymentMethodsResolver } from "./retrivePaymentMethods";

const resolvers: NonEmptyArray<Function> = [
  AttachPaymentWithCustomerResolver,
  CreateCustomerResolver,
  GetBillHistoryResolver,
  GetCurrentBillResolver,
  GetStripeCustomerResolver,
  RetrivePaymentMethodsResolver,
];

export default resolvers;

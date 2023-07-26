import { NonEmptyArray } from "type-graphql";

import { AttachPaymentWithCustomerResolver } from "./attachPaymentWithCustomer.resolver";
import { CreateCustomerResolver } from "./createCustomer.resolver";
import { GetBillHistoryResolver } from "./getBillHistory.resolver";
import { GetCurrentBillResolver } from "./getCurrentBill.resolver";
import { GetStripeCustomerResolver } from "./getStripeCustomer.resolver";
import { RetrivePaymentMethodsResolver } from "./retrivePaymentMethods.resolver";


const resolvers: NonEmptyArray<Function> = [AttachPaymentWithCustomerResolver,
    CreateCustomerResolver,
    GetBillHistoryResolver,
    GetCurrentBillResolver,
    GetStripeCustomerResolver,
    RetrivePaymentMethodsResolver];

export default resolvers;

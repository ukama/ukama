import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class CurrentBillDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    dataUsed: number;

    @Field()
    rate: number;

    @Field()
    subtotal: number;
}

@ObjectType()
export class CurrentBillResponse {
    @Field()
    status: string;

    @Field(() => [CurrentBillDto])
    data: CurrentBillDto[];
}

@ObjectType()
export class BillResponse {
    @Field(() => [CurrentBillDto])
    bill: CurrentBillDto[];

    @Field()
    total: number;

    @Field()
    billMonth: string;

    @Field()
    dueDate: string;
}

@ObjectType()
export class BillHistoryDto {
    @Field()
    id: string;

    @Field()
    date: string;

    @Field()
    description: string;

    @Field()
    totalUsage: number;

    @Field()
    subtotal: number;
}

@ObjectType()
export class BillHistoryResponse {
    @Field()
    status: string;

    @Field(() => [BillHistoryDto])
    data: BillHistoryDto[];
}
@ObjectType()
export class StripeCustomer {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    email: string;
}

@InputType()
export class CreateCustomerDto {
    @Field()
    name: string;

    @Field()
    email: string;
}
@InputType()
export class AttachPaymentDto {
    @Field()
    customerId: string;

    @Field()
    paymentId: string;
}

@ObjectType()
export class StripePaymentMethods {
    @Field()
    id: string;

    @Field()
    brand: string;

    @Field({ nullable: true })
    cvc_check?: string;

    @Field({ nullable: true })
    country?: string;

    @Field()
    exp_month: number;

    @Field()
    exp_year: number;

    @Field()
    funding: string;

    @Field()
    last4: string;

    @Field()
    type: string;

    @Field()
    created: number;
}

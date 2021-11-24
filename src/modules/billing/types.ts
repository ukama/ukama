import { Field, ObjectType } from "type-graphql";

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

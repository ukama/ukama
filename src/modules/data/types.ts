import { Field, ObjectType } from "type-graphql";

@ObjectType()
export class DataUsageDto {
    @Field()
    id: string;

    @Field()
    dataConsumed: string;

    @Field()
    dataPackage: string;
}

@ObjectType()
export class DataBillDto {
    @Field()
    id: string;

    @Field()
    dataBill: string;

    @Field()
    billDue: string;
}

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
@ObjectType()
export class DataUsageResponse {
    @Field()
    data: DataUsageDto;

    @Field()
    status: string;
}

@ObjectType()
export class DataBillResponse {
    @Field()
    data: DataBillDto;

    @Field()
    status: string;
}

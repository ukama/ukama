import { Field, ObjectType } from "type-graphql";

@ObjectType()
export class DataUsageDto {
    @Field()
    id: string;

    @Field()
    dataConsumed: number;

    @Field()
    dataPackage: string;
}

@ObjectType()
export class DataBillDto {
    @Field()
    id: string;

    @Field()
    dataBill: number;

    @Field()
    billDue: number;
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

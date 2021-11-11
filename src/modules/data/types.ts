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

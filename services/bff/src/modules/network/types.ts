import { Field, ObjectType } from "type-graphql";
import { NETWORK_STATUS } from "../../constants";

@ObjectType()
export class NetworkDto {
    @Field()
    id: string;

    @Field(() => NETWORK_STATUS)
    status: NETWORK_STATUS;

    @Field({ nullable: true })
    description?: string;
}

@ObjectType()
export class NetworkResponse {
    @Field()
    status: string;

    @Field(() => NetworkDto)
    data: NetworkDto;
}

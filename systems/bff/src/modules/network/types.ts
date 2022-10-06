import { Field, ObjectType } from "type-graphql";
import { NETWORK_STATUS } from "../../constants";

@ObjectType()
export class NetworkDto {
    @Field()
    liveNode: number;

    @Field()
    totalNodes: number;

    @Field(() => NETWORK_STATUS)
    status: NETWORK_STATUS;
}

@ObjectType()
export class NetworkResponse {
    @Field()
    status: string;

    @Field(() => NetworkDto)
    data: NetworkDto;
}

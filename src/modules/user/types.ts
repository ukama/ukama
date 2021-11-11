import { Field, ObjectType } from "type-graphql";
import { CONNECTED_USER_TYPE } from "../../constants";

@ObjectType()
export class ConnectedUserDto {
    @Field()
    totalUser: number;

    @Field()
    residentUsers: number;

    @Field()
    guestUsers: number;
}

@ObjectType()
export class UserDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field(() => CONNECTED_USER_TYPE)
    type: CONNECTED_USER_TYPE;

    @Field()
    email: string;
}

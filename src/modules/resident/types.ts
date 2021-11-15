import { Field, ObjectType } from "type-graphql";
import { PaginationResponse } from "../../common/types";

@ObjectType()
export class ResidentDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    usage: string;
}

@ObjectType()
export class ResidentsResponse extends PaginationResponse {
    @Field(() => [ResidentDto])
    residents: ResidentDto[];
}

import { Field, ObjectType } from "type-graphql";

@ObjectType()
export class EsimDto {
    @Field()
    esim: string;

    @Field()
    active: boolean;
}

@ObjectType()
export class EsimResponse {
    @Field()
    status: string;

    @Field(() => [EsimDto])
    data: EsimDto[];
}

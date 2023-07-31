import { Field, InputType, ObjectType } from "type-graphql";

@InputType()
export class DefaultMarkupInputDto {
    @Field({ nullable: false })
    markup: number;
}

@ObjectType()
export class DefaultMarkupResDto {
    @Field()
    markup: number;
}

@ObjectType()
export class DefaultMarkupAPIResDto {
    @Field()
    markup: number;
}

@ObjectType()
export class DefaultMarkupHistoryDto {
    @Field()
    createdAt: string;

    @Field()
    deletedAt: string;

    @Field()
    Markup: number;
}

@ObjectType()
export class DefaultMarkupHistoryAPIResDto {
    @Field(() => [DefaultMarkupHistoryDto], { nullable: true })
    markupRates: DefaultMarkupHistoryDto[];
}

@ObjectType()
export class DefaultMarkupHistoryResDto {
    @Field(() => [DefaultMarkupHistoryDto], { nullable: true })
    markupRates: DefaultMarkupHistoryDto[];
}

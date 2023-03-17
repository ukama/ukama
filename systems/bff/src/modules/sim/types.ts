import { Field, InputType, ObjectType } from "type-graphql";

@InputType()
export class AllocateSimInputDto {
    @Field()
    network_id: string;

    @Field()
    package_id: string;

    @Field()
    sim_token: string;

    @Field()
    sim_type: string;

    @Field()
    subscriber_id: string;
}

@ObjectType()
export class SimAPIDto {
    @Field()
    activationCode: string;

    @Field()
    createdAt: string;

    @Field()
    iccid: string;

    @Field()
    id: string;

    @Field()
    isAllocated: string;

    @Field()
    isPhysical: string;

    @Field()
    msisdn: string;

    @Field()
    qrCode: string;

    @Field()
    simType: string;

    @Field()
    smDpAddress: string;
}

@ObjectType()
export class SimAPIResDto {
    @Field(() => SimAPIDto)
    sim: SimAPIDto;
}

@ObjectType()
export class SimResDto extends SimAPIDto {}

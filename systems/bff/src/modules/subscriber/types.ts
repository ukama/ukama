import { IsEmail, IsEnum, IsPhoneNumber } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";
import { GENDER_TYPE } from "../../constants";

@ObjectType()
export class SubscriberAPIDto {
    @Field()
    subscriber_id: string;

    @Field()
    address: string;

    @Field()
    date_of_birth: string;

    @Field()
    email: string;

    @Field()
    first_name: string;

    @Field()
    last_name: string;

    @Field()
    gender: string;

    @Field()
    id_serial: string;

    @Field()
    network_id: string;

    @Field()
    org_id: string;

    @Field()
    phone_number: string;

    @Field()
    proof_of_identification: string;
}

@ObjectType()
export class SubscriberAPIResDto {
    @Field(() => SubscriberAPIDto)
    Subscriber: SubscriberAPIDto;
}

@InputType()
export class SubscriberInputDto {
    @Field()
    address: string;

    @Field()
    dob: string;

    @Field()
    @IsEmail()
    email: string;

    @Field()
    first_name: string;

    @Field()
    last_name: string;

    @Field()
    @IsEnum(GENDER_TYPE)
    gender: string;

    @Field()
    id_serial: string;

    @Field()
    network_id: string;

    @Field()
    org_id: string;

    @Field()
    @IsPhoneNumber()
    phone: string;

    @Field()
    proof_of_identification: string;
}
@ObjectType()
export class SubscriberDto {
    @Field()
    uuid: string;

    @Field()
    address: string;

    @Field()
    dob: string;

    @Field()
    email: string;

    @Field()
    firstName: string;

    @Field()
    lastName: string;

    @Field()
    gender: string;

    @Field()
    idSerial: string;

    @Field()
    networkId: string;

    @Field()
    orgId: string;

    @Field()
    phone: string;

    @Field()
    proofOfIdentification: string;
}

@InputType()
export class UpdateSubscriberInputDto {
    @Field({ nullable: true })
    address: string;

    @Field({ nullable: true })
    @IsEmail()
    email: string;

    @Field({ nullable: true })
    id_serial: string;

    @Field({ nullable: true })
    @IsPhoneNumber()
    phone: string;

    @Field({ nullable: true })
    proof_of_identification: string;
}

@ObjectType()
export class SubscriberMetricsByNetworkDto {
    @Field()
    total: number;

    @Field()
    active: number;

    @Field()
    inactive: number;

    @Field()
    terminated: number;
}

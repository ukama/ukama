import { Field, InputType, Int, ObjectType } from "type-graphql";

@ObjectType()
export class PackageRateAPIDto {
    @Field()
    sms_mo: string;

    @Field()
    sms_mt: number;

    @Field()
    data: number;

    @Field()
    amount: number;
}
@ObjectType()
export class PackageMarkupAPIDto {
    @Field()
    baserate: string;

    @Field()
    markup: number;
}

@ObjectType()
export class PackageAPIDto {
    @Field()
    uuid: string;

    @Field()
    name: string;

    @Field()
    org_id: string;

    @Field()
    active: boolean;

    @Field()
    duration: string;

    @Field()
    sim_type: string;

    @Field()
    created_at: string;

    @Field()
    deleted_at: string;

    @Field()
    updated_at: string;

    @Field()
    sms_volume: string;

    @Field()
    data_volume: string;

    @Field()
    voice_volume: string;

    @Field()
    ulbr: string;

    @Field()
    dlbr: string;

    @Field()
    type: string;

    @Field()
    data_unit: string;

    @Field()
    voice_unit: string;

    @Field()
    message_unit: string;

    @Field()
    flatrate: boolean;

    @Field()
    currency: string;

    @Field()
    from: string;

    @Field()
    to: string;

    @Field()
    country: string;

    @Field()
    provider: string;

    @Field()
    apn: string;

    @Field()
    owner_id: string;

    @Field()
    amount: number;

    @Field(() => PackageRateAPIDto)
    rate: PackageRateAPIDto;

    @Field(() => PackageMarkupAPIDto)
    markup: PackageMarkupAPIDto;
}
@ObjectType()
export class PackageAPIResDto {
    @Field()
    package: PackageAPIDto;
}

@ObjectType()
export class PackagesAPIResDto {
    @Field(() => [PackageAPIDto])
    packages: PackageAPIDto[];
}

@ObjectType()
export class PackageRateDto {
    @Field()
    smsMo: string;

    @Field()
    smsMt: number;

    @Field()
    data: number;

    @Field()
    amount: number;
}
@ObjectType()
export class PackageMarkupDto {
    @Field()
    baserate: string;

    @Field()
    markup: number;
}

@ObjectType()
export class PackageDto {
    @Field()
    uuid: string;

    @Field()
    name: string;

    @Field()
    orgId: string;

    @Field()
    active: boolean;

    @Field()
    duration: string;

    @Field()
    simType: string;

    @Field()
    createdAt: string;

    @Field()
    deletedAt: string;

    @Field()
    updatedAt: string;

    @Field()
    smsVolume: string;

    @Field()
    dataVolume: string;

    @Field()
    voiceVolume: string;

    @Field()
    ulbr: string;

    @Field()
    dlbr: string;

    @Field()
    type: string;

    @Field()
    dataUnit: string;

    @Field()
    voiceUnit: string;

    @Field()
    messageUnit: string;

    @Field()
    flatrate: boolean;

    @Field()
    currency: string;

    @Field()
    from: string;

    @Field()
    to: string;

    @Field()
    country: string;

    @Field()
    provider: string;

    @Field()
    apn: string;

    @Field()
    ownerId: string;

    @Field()
    amount: number;

    @Field(() => PackageRateAPIDto)
    rate: PackageRateAPIDto;

    @Field(() => PackageMarkupAPIDto)
    markup: PackageMarkupAPIDto;
}

@ObjectType()
export class PackagesResDto {
    @Field(() => [PackageDto])
    packages: PackageDto[];
}

@InputType()
export class AddPackageInputDto {
    @Field()
    name: string;

    @Field()
    active: boolean;

    @Field(() => Int)
    duration: number;

    @Field(() => Int)
    data_volume: number;

    @Field()
    org_id: string;

    @Field(() => Int)
    org_rates_id: number;

    @Field()
    sim_type: string;

    @Field(() => Int)
    sms_volume: number;

    @Field(() => Int)
    voice_volume: number;
}

@InputType()
export class UpdatePackageInputDto {
    @Field({ nullable: true })
    name: string;

    @Field({ nullable: true })
    active: boolean;

    @Field(() => Int, { nullable: true })
    duration: number;

    @Field(() => Int, { nullable: true })
    data_volume: number;

    @Field(() => Int, { nullable: true })
    org_rates_id: number;

    @Field({ nullable: true })
    sim_type: string;

    @Field(() => Int, { nullable: true })
    sms_volume: number;

    @Field(() => Int, { nullable: true })
    voice_volume: number;
}

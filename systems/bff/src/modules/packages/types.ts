import { Field, InputType, ObjectType } from "type-graphql";

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
    org_rates_id: string;
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
    orgRatesId: string;
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

    @Field()
    duration: bigint;

    @Field()
    data_volume: bigint;

    @Field()
    org_id: string;

    @Field()
    org_rates_id: bigint;

    @Field()
    sim_type: string;

    @Field()
    sms_volume: bigint;

    @Field()
    voice_volume: bigint;
}

@InputType()
export class UpdatePackageInputDto {
    @Field({ nullable: true })
    name: string;

    @Field({ nullable: true })
    active: boolean;

    @Field({ nullable: true })
    duration: bigint;

    @Field({ nullable: true })
    data_volume: bigint;

    @Field({ nullable: true })
    org_rates_id: bigint;

    @Field({ nullable: true })
    sim_type: string;

    @Field({ nullable: true })
    sms_volume: bigint;

    @Field({ nullable: true })
    voice_volume: bigint;
}

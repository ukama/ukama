import { Field, InputType, ObjectType } from "type-graphql";

@InputType()
export class AllocateSimInputDto {
    @Field()
    networkId: string;

    @Field()
    packageId: string;

    @Field()
    iccid: string;

    @Field()
    simType: string;

    @Field()
    subscriberId: string;
}
@ObjectType()
export class SimStatusResDto {
    @Field(() => String, { nullable: true })
    simId?: string;
}
@ObjectType()
export class DeleteSimResDto {
    @Field(() => String, { nullable: true })
    simId?: string;
}
@ObjectType()
export class RemovePackageFromSimResDto {
    @Field(() => String, { nullable: true })
    packageId?: string;
}
@ObjectType()
export class AddPackageSimResDto {
    @Field(() => String, { nullable: true })
    packageId?: string;
}
@ObjectType()
export class SetActivePackageForSimResDto {
    @Field(() => String, { nullable: true })
    packageId?: string;
}
@InputType()
export class RemovePackageFormSimInputDto {
    @Field()
    simId: string;

    @Field()
    packageId: string;
}
@ObjectType()
export class GetPackagesForSimResDto {
    @Field(() => [SimPackageDto], { nullable: true })
    Packages?: [SimPackageDto];
}
@InputType()
export class GetPackagesForSimInputDto {
    @Field()
    simId: string;
}
@InputType()
export class ToggleSimStatusInputDto {
    @Field()
    simId: string;

    @Field()
    status: string;
}
@InputType()
export class GetSimInputDto {
    @Field()
    simId: string;
}
@InputType()
export class GetSimBySubscriberIdInputDto {
    @Field()
    subscriberId: string;
}
@InputType()
export class GetSimByNetworkInputDto {
    @Field()
    networkId: string;
}
@InputType()
export class DeleteSimInputDto {
    @Field()
    simId: string;
}
@InputType()
export class AddPackageToSimInputDto {
    @Field()
    simId: string;

    @Field()
    packageId: string;

    @Field()
    startDate: string;
}
@InputType()
export class SetActivePackageForSimInputDto {
    @Field()
    simId: string;

    @Field()
    packageId: string;
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
export class GetSimAPIResDto {
    @Field(() => SimDetailsDto)
    sim: SimDetailsDto;
}
@ObjectType()
export class SimPackageDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    description: string;

    @Field()
    createdAt: string;

    @Field()
    updatedAt: string;
}

@ObjectType()
export class SimDetailsDto {
    @Field()
    id: string;

    @Field()
    subscriberId: string;

    @Field()
    networkId: string;

    @Field()
    orgId: string;

    @Field(() => SimPackageDto)
    Package: SimPackageDto;

    @Field()
    iccid: string;

    @Field()
    msisdn: string;

    @Field()
    imsi: string;

    @Field()
    type: string;

    @Field()
    status: string;

    @Field()
    isPhysical: boolean;

    @Field()
    firstActivatedOn: string;

    @Field()
    lastActivatedOn: string;

    @Field()
    activationsCount: number;

    @Field()
    deactivationsCount: number;

    @Field()
    allocatedAt: string;
}
@ObjectType()
export class SimResDto extends SimAPIDto {}

@ObjectType()
export class SimDataUsage {
    @Field()
    usage: string;
}

@ObjectType()
export class SimPoolStatsDto {
    @Field()
    total: number;

    @Field()
    available: number;

    @Field()
    consumed: number;

    @Field()
    failed: number;

    @Field()
    esim: number;

    @Field()
    physical: number;
}

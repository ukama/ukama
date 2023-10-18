import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

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

  @Field(() => Int)
  duration: number;

  @Field()
  dataUnit: string;

  @Field(() => Float)
  amount: number;

  @Field(() => Int)
  dataVolume: number;
}

@InputType()
export class UpdatePackageInputDto {
  @Field()
  name: string;

  @Field()
  active: boolean;
}

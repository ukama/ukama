import { IsEmail, IsEnum, IsPhoneNumber } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";
import { GENDER_TYPE, GET_STATUS_TYPE } from "../../constants";

@ObjectType()
export class UserAPIObj {
    @Field()
    name: string;

    @Field()
    email: string;

    @Field()
    uuid: string;

    @Field()
    phone: string;

    @Field()
    is_deactivated: boolean;

    @Field()
    registered_since: string;
}
@ObjectType()
export class UserAPIResDto {
    @Field(() => [UserAPIObj])
    user: UserAPIObj;
}

@ObjectType()
export class DeactivateResponse {
    @Field()
    uuid: string;

    @Field()
    name: string;

    @Field()
    email: string;

    @Field()
    isDeactivated: boolean;
}

@InputType()
export class UserInputDto {
    @Field()
    name: string;

    @Field()
    @IsEmail()
    email: string;

    @Field()
    @IsPhoneNumber()
    phone: string;
}

@InputType()
export class UpdateUserInputDto {
    @Field({ nullable: true })
    name: string;

    @Field({ nullable: true })
    @IsEmail()
    email: string;

    @Field({ nullable: true })
    @IsPhoneNumber()
    phone: string;
}

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
export class ConnectedUserDto {
    @Field()
    totalUser: string;
}
@InputType()
export class UserFistVisitInputDto {
    @Field()
    firstVisit: boolean;
}
@ObjectType()
export class UserFistVisitResDto {
    @Field()
    firstVisit: boolean;
}

@ObjectType()
export class GetUserDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    @IsEmail()
    email: string;

    @Field()
    phone: string;
}
@ObjectType()
export class GetAccountDetailsDto {
    @Field()
    email: string;

    @Field()
    isFirstVisit?: boolean;
}
@ObjectType()
export class GetUsersDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    @IsEmail()
    email: string;

    @Field({ nullable: true })
    dataPlan: string;

    @Field({ nullable: true })
    dataUsage: string;
}

@ObjectType()
export class GetUserResponseDto {
    @Field()
    status: string;

    @Field(() => [GetUserDto])
    data: GetUserDto[];

    @Field()
    length: number;
}

@ObjectType()
export class DeleteNodeRes {
    @Field()
    nodeId: string;
}

@ObjectType()
export class OrgUserDto {
    @Field()
    name: string;

    @Field()
    email: string;

    @Field()
    uuid: string;

    @Field()
    isDeactivated: boolean;
}

@ObjectType()
export class UserResDto {
    @Field()
    name: string;

    @Field()
    email: string;

    @Field()
    uuid: string;

    @Field()
    phone: string;

    @Field()
    isDeactivated: boolean;

    @Field()
    registeredSince: string;
}

@ObjectType()
export class OrgUsersResponse {
    @Field()
    org: string;

    @Field(() => [OrgUserDto])
    users: OrgUserDto[];
}

@ObjectType()
export class UserSimServices {
    @Field()
    voice: boolean;
    @Field()
    data: boolean;
    @Field()
    sms: boolean;
}

@ObjectType()
export class UserDataUsageDto {
    @Field({ nullable: true })
    dataUsedBytes: string;

    @Field({ nullable: true })
    dataAllowanceBytes: string;
}

@ObjectType()
export class UserServicesDto {
    @Field(() => GET_STATUS_TYPE)
    status: GET_STATUS_TYPE;

    @Field(() => UserSimServices)
    services: UserSimServices;

    @Field({ nullable: true })
    @Field(() => UserDataUsageDto)
    usage?: UserDataUsageDto;
}

@ObjectType()
export class OrgUserSimDto {
    @Field()
    iccid: string;

    @Field()
    isPhysical: boolean;

    @Field(() => UserServicesDto)
    ukama?: UserServicesDto;

    @Field(() => UserServicesDto)
    carrier?: UserServicesDto;
}

@ObjectType()
export class OrgUserResponse {
    @Field(() => OrgUserSimDto)
    sim: OrgUserSimDto;

    @Field(() => OrgUserDto)
    user: OrgUserDto;
}

@ObjectType()
export class AddUserServiceRes {
    @Field(() => OrgUserDto)
    user: OrgUserDto;

    @Field()
    iccid: string;
}

@InputType()
export class UpdateUserServiceInput {
    @Field()
    simId: string;

    @Field()
    userId: string;

    @Field()
    status: boolean;
}
@InputType()
export class GetESimQRCodeInput {
    @Field()
    simId: string;

    @Field()
    userId: string;
}

@InputType()
export class DataUsageInputDto {
    @Field(() => [String])
    ids: string[];
}

@ObjectType()
export class ESimQRCodeRes {
    @Field()
    qrCode: string;
}

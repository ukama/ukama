import { IsEmail } from "class-validator";
import { GET_STATUS_TYPE } from "../../constants";
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class ConnectedUserDto {
    @Field()
    totalUser: string;
}

@InputType()
export class UserInputDto {
    @Field()
    name: string;

    @Field()
    email: string;

    @Field({ nullable: true })
    status: boolean;
}
@InputType()
export class UserFistVisitInputDto {
    @Field()
    firstVisit: boolean;
    @Field()
    email: string;
}
@ObjectType()
export class UserFistVisitResDto {
    @Field()
    firstVisit: boolean;
}
@ObjectType()
export class ActivateUserResponse {
    @Field()
    success: boolean;
}

@ObjectType()
export class GetUserDto {
    @Field()
    id: string;

    @Field()
    status: boolean;

    @Field()
    name: string;

    @Field()
    eSimNumber?: string;

    @Field()
    iccid: string;

    @Field()
    @IsEmail()
    email: string;

    @Field()
    roaming: boolean;

    @Field()
    dataPlan: string;

    @Field()
    dataUsage: string;
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
    id: string;

    @Field({ nullable: true })
    iccid?: string;
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

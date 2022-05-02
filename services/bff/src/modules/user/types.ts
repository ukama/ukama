import { IsEmail, IsPhoneNumber } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";
import { PaginationResponse } from "../../common/types";
import { GET_STATUS_TYPE } from "../../constants";

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

    @Field()
    phone: string;
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
    @IsPhoneNumber()
    phone: string;

    @Field()
    roaming: boolean;

    @Field()
    dataPlan: number;

    @Field()
    dataUsage: number;
}

@ObjectType()
export class GetUsersDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field({ nullable: true })
    @IsEmail()
    email: string;

    @Field({ nullable: true })
    @IsPhoneNumber()
    phone: string;

    @Field()
    dataPlan: string;

    @Field()
    dataUsage: string;

    @Field()
    isDeactivated: boolean;
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
export class ResidentResponse {
    @Field(() => [GetUserDto])
    residents: GetUserDto[];

    @Field()
    activeResidents: number;

    @Field()
    totalResidents: number;
}

@ObjectType()
export class ResidentsResponse extends PaginationResponse {
    @Field(() => ResidentResponse)
    residents: ResidentResponse;
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
    phone: string;

    @Field()
    isDeactivated: boolean;
}

@ObjectType()
export class OrgUserDto {
    @Field()
    name: string;

    @Field()
    phone: string;

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
    phone: string;

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
export class UserSimUkamaDto {
    @Field(() => GET_STATUS_TYPE)
    status: GET_STATUS_TYPE;

    @Field(() => UserSimServices)
    services: UserSimServices;
}

@ObjectType()
export class OrgUserSimDto {
    @Field()
    iccid: string;

    @Field()
    isPhysical: boolean;

    @Field(() => UserSimUkamaDto)
    ukama?: UserSimUkamaDto;

    @Field(() => UserSimUkamaDto)
    carrier?: UserSimUkamaDto;
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
export class DataUsageInputDto {
    @Field(() => [String])
    ids: string[];
}

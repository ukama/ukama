import { IsEmail, IsPhoneNumber, Length } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";
import { PaginationDto, PaginationResponse } from "../../common/types";
import {
    CONNECTED_USER_TYPE,
    GET_STATUS_TYPE,
    GET_USER_TYPE,
} from "../../constants";

@ObjectType()
export class ConnectedUserDto {
    @Field()
    totalUser: string;
}

@ObjectType()
export class UserDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field(() => CONNECTED_USER_TYPE)
    type: CONNECTED_USER_TYPE;

    @Field()
    email: string;
}

@ObjectType()
export class ConnectedUserResponse {
    @Field()
    data: ConnectedUserDto;

    @Field()
    status: string;
}

@InputType()
export class ActivateUserDto {
    @Field()
    eSimNumber: string;

    @Field()
    iccid: string;

    @Field()
    @Length(3, 255)
    name: string;

    @Field({ nullable: true })
    @IsEmail()
    email?: string;

    @Field({ nullable: true })
    @IsPhoneNumber()
    phone?: string;

    @Field()
    roaming: boolean;

    @Field()
    dataUsage: number;

    @Field()
    dataPlan: number;
}

@InputType()
export class UpdateUserDto {
    @Field()
    id: string;

    @Field({ nullable: true })
    eSimNumber: string;

    @Field({ nullable: true })
    @Length(3, 255)
    firstName: string;

    @Field({ nullable: true })
    @Length(3, 255)
    lastName: string;

    @Field({ nullable: true })
    @IsEmail()
    email?: string;

    @Field({ nullable: true })
    @IsPhoneNumber()
    phone?: string;
}
@ObjectType()
export class UserResponse {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    sim: string;

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
export class ActiveUserResponseDto {
    @Field()
    status: string;

    @Field()
    data: ActivateUserResponse;
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

    @Field({ nullable: true })
    @IsEmail()
    email: string;

    @Field({ nullable: true })
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
    dataPlan: number;

    @Field()
    dataUsage: number;
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

@InputType()
export class GetUserPaginationDto extends PaginationDto {
    @Field(() => GET_USER_TYPE)
    type: GET_USER_TYPE;
}

@ObjectType()
export class GetUserResponse extends PaginationResponse {
    @Field(() => [GetUserDto])
    users: GetUserDto[];
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
    id: string;

    @Field()
    success: boolean;
}

@ObjectType()
export class OrgUserResponseDto {
    @Field()
    orgName: string;

    @Field(() => [GetUserDto])
    users: GetUserDto[];
}

@InputType()
export class AddUserDto {
    @Field()
    @Length(3, 255)
    name: string;

    @Field()
    @IsEmail()
    email: string;
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
}

@ObjectType()
export class AddUserResponse {
    @Field()
    name: string;

    @Field()
    email: string;

    @Field()
    phone: string;

    @Field()
    uuid: string;

    @Field()
    iccid: string;
}

@ObjectType()
export class OrgUsersResponse {
    @Field()
    org: string;

    @Field(() => [OrgUserDto])
    users: OrgUserDto[];
}

@InputType()
export class UserInput {
    @Field()
    orgId: string;

    @Field()
    userId: string;
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
    orgId: string;

    @Field()
    simId: string;

    @Field()
    userId: string;

    @Field()
    status: boolean;
}
@ObjectType()
export class UpdateUserServiceRes {
    @Field()
    success: boolean;
}

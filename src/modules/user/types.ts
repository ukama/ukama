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
    totalUser: number;
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

    @Field(() => GET_STATUS_TYPE)
    status: GET_STATUS_TYPE;

    @Field()
    name: string;

    @Field()
    eSimNumber: string;

    @Field()
    iccid: string;

    @Field({ nullable: true })
    @IsEmail()
    email?: string;

    @Field({ nullable: true })
    @IsPhoneNumber()
    phone?: string;

    @Field()
    roaming: boolean;

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
    imsi: string;

    @Field()
    @Length(3, 255)
    firstName: string;

    @Field()
    @Length(3, 255)
    lastName: string;

    @Field()
    @IsEmail()
    email: string;
}

@ObjectType()
export class AddUserResponse {
    @Field()
    imsi: string;

    @Field()
    firstName: string;

    @Field()
    lastName: string;

    @Field()
    email: string;

    @Field()
    uuid: string;
}

@ObjectType()
export class OrgUserDto {
    @Field()
    firstName: string;

    @Field()
    lastName: string;

    @Field()
    email: string;

    @Field()
    uuid: string;
}

@ObjectType()
export class OrgUserResponse {
    @Field()
    org: string;

    @Field(() => [OrgUserDto])
    users: OrgUserDto[];
}

@ObjectType()
export class ActiveUserMetricsDto {
    @Field({ nullable: true })
    id: string;

    @Field()
    users: number;

    @Field()
    timestamp: number;
}

@ObjectType()
export class ActiveUserMetricsResponse extends PaginationResponse {
    @Field(() => [ActiveUserMetricsDto])
    data: ActiveUserMetricsDto[];
}

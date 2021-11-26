import { IsEmail, IsPhoneNumber, Length } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";
import { PaginationDto, PaginationResponse } from "../../common/types";
import {
    CONNECTED_USER_TYPE,
    DATA_PLAN_TYPE,
    GET_STATUS_TYPE,
    GET_USER_TYPE,
} from "../../constants";

@ObjectType()
export class ConnectedUserDto {
    @Field()
    totalUser: number;

    @Field()
    residentUsers: number;

    @Field()
    guestUsers: number;
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
    @Length(3, 255)
    firstName: string;

    @Field()
    @Length(3, 255)
    lastName: string;

    @Field({ nullable: true })
    @IsEmail()
    email?: string;

    @Field({ nullable: true })
    @IsPhoneNumber()
    phone?: string;
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
    node: string;

    @Field(() => DATA_PLAN_TYPE)
    dataPlan: DATA_PLAN_TYPE;

    @Field()
    dataUsage: number;

    @Field()
    dlActivity: string;

    @Field()
    ulActivity: string;
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
export class ResidentDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    dataUsage: number;
}

@ObjectType()
export class ResidentResponse {
    @Field(() => [ResidentDto])
    residents: ResidentDto[];

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
export class DeleteUserResponse {
    @Field()
    id: string;

    @Field()
    success: boolean;
}

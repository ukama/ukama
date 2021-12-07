import { IsEmail, IsOptional, IsPhoneNumber, Length } from "class-validator";
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
export class DeactivateResponse {
    @Field()
    id: string;

    @Field()
    success: boolean;
}

@ObjectType()
export class AuthenticationMethodDto {
    @Field()
    method: string;

    @Field()
    completed_at: string;
}

@ObjectType()
export class TraitDto {
    @Field()
    name: string;

    @Field()
    @IsEmail()
    email: string;
}

@ObjectType()
export class AddressDto {
    @Field()
    id: string;

    @Field()
    value: string;

    @Field()
    @IsOptional()
    verified?: boolean;

    @Field()
    via: string;

    @Field()
    @IsOptional()
    status?: string;

    @Field()
    created_at: string;

    @Field()
    updated_at: string;
}

@ObjectType()
export class IdentityDto {
    @Field()
    id: string;

    @Field()
    schema_id: string;

    @Field()
    schema_url: string;

    @Field()
    state: boolean;

    @Field()
    state_changed_at: string;

    @Field(() => TraitDto)
    traits: TraitDto;

    @Field(() => [AddressDto])
    verifiable_addresses: AddressDto[];

    @Field(() => [AddressDto])
    recovery_addresses: AddressDto[];

    @Field()
    created_at: string;

    @Field()
    updated_at: string;
}

@ObjectType()
export class OrganisationDto {
    @Field()
    id: string;

    @Field()
    active: boolean;

    @Field()
    expires_at: string;

    @Field()
    authenticated_at: string;

    @Field()
    authenticator_assurance_level: string;

    @Field(() => [AuthenticationMethodDto])
    authentication_methods: AuthenticationMethodDto[];

    @Field()
    issued_at: string;

    @Field(() => IdentityDto)
    identity: IdentityDto;
}

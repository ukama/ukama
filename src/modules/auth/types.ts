import { IsEmail } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";

@InputType()
export class LoginDto {
    @Field()
    @IsEmail()
    email: string;

    @Field()
    password: string;
}

@ObjectType()
export class ServerLoginDto {
    @Field()
    @IsEmail()
    password_identifier: string;

    @Field()
    password: string;

    @Field({ nullable: true })
    method: string;
}

import { IsEmail } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";

@InputType()
export class SendInvitationInputDto {
  @Field()
  name: string;

  @Field()
  @IsEmail()
  email: string;

  @Field()
  role: string;
}

@ObjectType()
export class SendInvitationResDto {
  @Field()
  id: string;

  @Field()
  message: string;
}

@ObjectType()
export class InvitationDto {
  @Field()
  email: string;

  @Field()
  expiresAt: string;

  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  role: string;

  @Field()
  link: string;

  @Field()
  userId: string;

  @Field()
  org: string;

  @Field()
  status: string;
}

@ObjectType()
export class InvitationAPIResDto {
  @Field()
  email: string;

  @Field()
  name: string;

  @Field()
  expires_at: string;

  @Field()
  user_id: string;

  @Field()
  role: string;

  @Field()
  id: string;

  @Field()
  link: string;

  @Field()
  org: string;

  @Field()
  status: string;
}

@ObjectType()
export class GetInvitationByOrgResDto {
  @Field(() => [InvitationDto])
  invitations: InvitationDto[];
}

@InputType()
export class GetInvitationInputDto {
  @Field()
  id: string;
}

@InputType()
export class GetInvitationByOrgInputDto {
  @Field()
  orgName: string;
}

@ObjectType()
export class UpdateInvitationResDto {
  @Field()
  id: string;
}

@ObjectType()
export class DeleteInvitationResDto {
  @Field()
  id: string;
}

@InputType()
export class UpateInvitationInputDto {
  @Field()
  status: string;
}

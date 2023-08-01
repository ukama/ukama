import { Field, InputType, ObjectType } from "type-graphql";

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
export class WhoamiDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  role: string;

  @Field()
  isFirstVisit?: boolean;
}

@InputType()
export class UserFistVisitInputDto {
  @Field()
  userId: string;

  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  firstVisit: boolean;
}

@ObjectType()
export class UserFistVisitResDto {
  @Field()
  firstVisit: boolean;
}

@ObjectType()
export class UserAPIObj {
  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  id: string;

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
export class WhoamiAPIDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  role: string;

  @Field()
  first_visit: boolean;
}

import "reflect-metadata";
import { Field, ID, InputType, ObjectType } from "type-graphql";

@ObjectType()
class Location {
  @Field()
  lat: string;
  @Field()
  lng: string;
  @Field()
  address: string;
}

@ObjectType()
class Site {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  height: number;

  @Field()
  apOption: string;

  @Field()
  solarUptime: number;

  @Field()
  isSetlite: boolean;

  @Field(() => Location, { nullable: false })
  location: Location;
}

@ObjectType()
class Event {
  @Field(() => ID)
  id: string;

  @Field()
  operation: string;

  @Field()
  value: string;

  @Field()
  createdAt: string;
}

@ObjectType()
export class Draft {
  @Field(() => ID)
  id: string;

  @Field()
  name: string;

  @Field()
  userId: string;

  @Field()
  lastSaved: number;

  @Field(() => [Site], { nullable: false })
  site: Site[];

  @Field(() => [Event], { nullable: false })
  events: Event[];
}

@InputType()
export class UpdateSiteInput {
  @Field()
  lastSaved: number;

  @Field()
  siteName: string;

  @Field()
  height: number;

  @Field()
  apOption: string;

  @Field()
  solarUptime: number;

  @Field()
  isSetlite: boolean;

  @Field()
  lat: string;

  @Field()
  lng: string;

  @Field()
  address: string;
}
@InputType()
export class UpdateEventInput {
  @Field()
  operation: string;

  @Field()
  value: string;
}
@InputType()
export class AddDraftInput {
  @Field()
  name: string;

  @Field()
  lastSaved: number;

  @Field()
  userId: string;
}
@InputType()
export class SiteLocationInput {
  @Field()
  lat: string;

  @Field()
  lng: string;

  @Field()
  address: string;
}

@ObjectType()
export class DeleteDraftRes {
  @Field()
  id: string;
}

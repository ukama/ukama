import "reflect-metadata";
import { Field, ID, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class Location {
  @Field()
  id: string;

  @Field()
  lat: string;

  @Field()
  lng: string;

  @Field()
  address: string;
}
@ObjectType()
export class Link {
  @Field()
  id: string;

  @Field()
  data: string;

  @Field()
  linkWith: string;
}

@ObjectType()
export class Site {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  status: string;

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

  @Field(() => [Link], { nullable: false })
  links: Link[];
}

@ObjectType()
export class Event {
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
  sites: Site[];

  @Field(() => [Event], { nullable: false })
  events: Event[];
}

@InputType()
export class SiteInput {
  @Field()
  locationId: string;

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
  userId: string;

  @Field()
  lastSaved: number;
}
@InputType()
export class LocationInput {
  @Field()
  lastSaved: number;

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

@ObjectType()
export class DeleteSiteRes {
  @Field()
  id: string;
}

@InputType()
export class LinkInput {
  @Field()
  lastSaved: number;

  @Field()
  data: string;

  @Field()
  linkWith: string;
}

@ObjectType()
export class DeleteLinkRes {
  @Field()
  id: string;
}

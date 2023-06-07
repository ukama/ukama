import "reflect-metadata";
import { Field, ID, InputType, ObjectType } from "type-graphql";
import { APOptions } from "../../common/enums";

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
  name: string;

  @Field()
  height: number;

  @Field(type => APOptions)
  apOption: APOptions;

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
}

@ObjectType()
export class Draft {
  @Field(() => ID)
  id: string;

  @Field()
  name: string;

  @Field()
  lastSaved: number;

  @Field(() => Site, { nullable: false })
  site: Site;

  @Field(() => [Event], { nullable: false })
  events: Event[];
}

@InputType()
export class AddDraftInput {
  @Field()
  name: string;

  @Field()
  lastSaved: number;

  @Field()
  siteName: string;

  @Field()
  height: number;

  @Field(type => APOptions)
  apOption: APOptions;

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
export class UpdateEvent {
  @Field()
  operation: string;

  @Field()
  value: string;
}

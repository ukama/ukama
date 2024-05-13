import { Field, ObjectType, registerEnumType } from "type-graphql";

export enum ComponentCategory {
  ALL = 0,
  ACCESS = 1,
  BACKHAUL = 2,
  POWER = 3,
  SWITCH = 4,
}

registerEnumType(ComponentCategory, {
  name: "ComponentCategory",
  description: "Categories for components",
});

@ObjectType()
export class ComponentCategoryType {
  @Field(() => ComponentCategory) // Decorator on the property
  id: ComponentCategory;

  @Field()
  name: string;
}

@ObjectType()
export class ComponentAPIDto {
  @Field()
  id: string;

  @Field()
  inventory: string;

  @Field(() => ComponentCategory) // Decorator on the property
  category: ComponentCategory;

  @Field()
  type: string;

  @Field()
  user_id: string;

  @Field()
  description: string;

  @Field()
  datasheet_url: string;

  @Field()
  images_url: string;

  @Field()
  part_number: string;

  @Field()
  manufacturer: string;

  @Field()
  managed: string;

  @Field()
  warranty: number;

  @Field()
  specification: string;
}

@ObjectType()
export class ComponentDto {
  @Field()
  id: string;

  @Field()
  inventory: string;

  @Field()
  type: string;

  @Field() // Removed typo (userId -> userId)
  userId: string;

  @Field()
  description: string;

  @Field() // Renamed for consistency (datasheetURL -> datasheetUrl)
  datasheetUrl: string;

  @Field() // Renamed for consistency (imagesURL -> imagesUrl)
  imageUrl: string;

  @Field()
  partNumber: string;

  @Field()
  manufacturer: string;

  @Field()
  managed: string;

  @Field()
  warranty: number;

  @Field()
  specification: string;
}
@ObjectType()
export class ComponentsResDto {
  @Field(() => [ComponentDto]) // Decorator on the property
  components: ComponentDto[];
}

@ObjectType()
export class ComponentAPIResDto {
  @Field(() => ComponentAPIDto) // Decorator on the property
  component: ComponentAPIDto;
}

@ObjectType()
export class ComponentsAPIResDto {
  @Field(() => [ComponentAPIDto]) // Decorator on the property
  components: ComponentAPIDto[];
}

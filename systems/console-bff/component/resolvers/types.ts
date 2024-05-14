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
  @Field(() => ComponentCategory)
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

  @Field(() => ComponentCategory)
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

  @Field()
  userId: string;

  @Field()
  description: string;

  @Field()
  datasheetUrl: string;

  @Field()
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
  @Field(() => [ComponentDto])
  components: ComponentDto[];
}

@ObjectType()
export class ComponentAPIResDto {
  @Field(() => ComponentAPIDto)
  component: ComponentAPIDto;
}

@ObjectType()
export class ComponentsAPIResDto {
  @Field(() => [ComponentAPIDto])
  components: ComponentAPIDto[];
}

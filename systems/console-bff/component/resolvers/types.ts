import { Field, ObjectType } from "type-graphql";

import { COMPONENT_CATEGORY } from "../../common/enums";

@ObjectType()
export class ComponentCategoryType {
  @Field(() => COMPONENT_CATEGORY)
  id: COMPONENT_CATEGORY;

  @Field()
  name: string;
}

@ObjectType()
export class ComponentAPIDto {
  @Field()
  id: string;

  @Field()
  inventory: string;

  @Field(() => COMPONENT_CATEGORY)
  category: COMPONENT_CATEGORY;

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

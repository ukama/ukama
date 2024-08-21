import { Field, InputType, ObjectType } from "type-graphql";

import { COMPONENT_TYPE } from "../../common/enums";

@ObjectType()
export class ComponentDto {
  @Field()
  id: string;

  @Field()
  inventoryId: string;

  @Field()
  category: string;

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
export class ComponentAPIDto {
  @Field()
  id: string;
  @Field()
  inventory_id: string;

  @Field()
  category: string;

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
export class ComponentsAPIResDto {
  @Field(() => [ComponentAPIDto])
  components: ComponentAPIDto[];
}

@InputType()
export class ComponentTypeInputDto {
  @Field(() => COMPONENT_TYPE)
  category: COMPONENT_TYPE;
}

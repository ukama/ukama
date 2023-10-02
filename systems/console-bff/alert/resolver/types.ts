import { Field, ObjectType } from "type-graphql";

import { ALERT_TYPE } from "../../common/enums";
import { PaginationResponse } from "../../common/types";

@ObjectType()
export class AlertDto {
  @Field({ nullable: true })
  id: string;

  @Field(() => ALERT_TYPE)
  type: ALERT_TYPE;

  @Field({ nullable: true })
  title: string;

  @Field({ nullable: true })
  description: string;

  @Field({ nullable: true })
  alertDate: Date;
}

@ObjectType()
export class AlertsResponse extends PaginationResponse {
  @Field(() => [AlertDto])
  alerts: AlertDto[];
}

@ObjectType()
export class AlertResponse {
  @Field()
  status: string;

  @Field(() => [AlertDto])
  data: AlertDto[];

  @Field()
  length: number;
}

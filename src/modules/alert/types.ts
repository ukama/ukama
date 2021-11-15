import { Field, ObjectType } from "type-graphql";
import { PaginationResponse } from "../../common/types";
import { ALERT_TYPE } from "../../constants";

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

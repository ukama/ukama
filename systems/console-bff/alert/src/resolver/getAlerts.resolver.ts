import {
    Resolver,
    Query,
    Arg,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { AlertsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Resolver()
export class GetAlertsResolver {

    @Query(() => AlertsResponse)
    @UseMiddleware(Authentication)
    async getAlerts(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<AlertsResponse> {
        const { dataSources } = ctx;
        const alerts = dataSources.dataSource.getAlerts(data);
        pubsub.publish("getAlerts", alerts);
        return alerts;
    }
}

import { Arg, Ctx, PubSub, PubSubEngine, Query, Resolver } from "type-graphql";

import { PaginationDto } from "../../common/types";
import { Context } from "../context";
import { AlertsResponse } from "./types";

@Resolver()
export class GetAlertsResolver {
  @Query(() => AlertsResponse)
  async getAlerts(
    @Arg("data") data: PaginationDto,
    @PubSub() pubsub: PubSubEngine,
    @Ctx() context: Context
  ): Promise<AlertsResponse> {
    const { dataSources } = context;
    const alerts = dataSources.dataSource.getAlerts(data);
    pubsub.publish("getAlerts", alerts);
    return alerts;
  }
}

import {
    Resolver,
    Query,
    Arg,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { AlertsResponse } from "../types";
import { AlertService } from "../service";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetAlertsResolver {
    constructor(private readonly alertService: AlertService) {}

    @Query(() => AlertsResponse)
    @UseMiddleware(Authentication)
    async getAlerts(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<AlertsResponse> {
        const alerts = this.alertService.getAlerts(data);
        await pubsub.publish("GETALERTS", alerts);
        return alerts;
    }
}

import {
    Resolver,
    Arg,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { UserService } from "../service";
import { ActiveUserMetricsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetActiveUserMetricsResolver {
    constructor(private readonly activeUserMetricsService: UserService) {}

    @Query(() => ActiveUserMetricsResponse)
    @UseMiddleware(Authentication)
    async getActiveUserMetrics(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<ActiveUserMetricsResponse | null> {
        const activeUserMetrics =
            this.activeUserMetricsService.activeUserMetricsService(data);

        pubsub.publish("activeUserMetrics", activeUserMetrics);
        return activeUserMetrics;
    }
}

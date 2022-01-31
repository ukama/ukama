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
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { UsersAttachedMetricsResponse } from "../types";

@Service()
@Resolver()
export class GetActiveUserMetricsResolver {
    constructor(private readonly activeUserMetricsService: UserService) {}

    @Query(() => UsersAttachedMetricsResponse)
    @UseMiddleware(Authentication)
    async getActiveUserMetrics(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<UsersAttachedMetricsResponse | null> {
        const activeUserMetrics =
            this.activeUserMetricsService.usersAttachedMetricsService(data);

        pubsub.publish("activeUserMetrics", activeUserMetrics);
        return activeUserMetrics;
    }
}

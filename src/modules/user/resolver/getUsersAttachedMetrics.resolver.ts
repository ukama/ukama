import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
    Arg,
} from "type-graphql";
import { Service } from "typedi";
import { UserService } from "../service";
import { UsersAttachedMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { GRAPH_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetUsersAttachedMetricsResolver {
    constructor(private readonly activeUserMetricsService: UserService) {}

    @Query(() => [UsersAttachedMetricsDto])
    @UseMiddleware(Authentication)
    async getUsersAttachedMetrics(
        @Arg("filter", () => GRAPH_FILTER) filter: GRAPH_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<[UsersAttachedMetricsDto] | null> {
        const activeUserMetrics =
            this.activeUserMetricsService.usersAttachedMetricsService(filter);
        pubsub.publish("usersAttachedMetrics", activeUserMetrics);
        return activeUserMetrics;
    }
}

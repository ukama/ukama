import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { UserService } from "../service";
import { UsersAttachedMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetUsersAttachedMetricsResolver {
    constructor(private readonly activeUserMetricsService: UserService) {}

    @Query(() => [UsersAttachedMetricsDto])
    @UseMiddleware(Authentication)
    async getUsersAttachedMetrics(
        @PubSub() pubsub: PubSubEngine
    ): Promise<[UsersAttachedMetricsDto] | null> {
        const activeUserMetrics =
            this.activeUserMetricsService.usersAttachedMetricsService();
        pubsub.publish("usersAttachedMetrics", activeUserMetrics);
        return activeUserMetrics;
    }
}

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
import { UsersAttachedMetricsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetUsersAttachedMetricsResolver {
    constructor(private readonly activeUserMetricsService: UserService) {}

    @Query(() => UsersAttachedMetricsResponse)
    @UseMiddleware(Authentication)
    async getUsersAttachedMetrics(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<UsersAttachedMetricsResponse | null> {
        const activeUserMetrics =
            this.activeUserMetricsService.usersAttachedMetricsService(data);

        pubsub.publish("usersAttachedMetrics", activeUserMetrics);
        return activeUserMetrics;
    }
}

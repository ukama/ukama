import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import {
    UsersAttachedMetricsDto,
    UsersAttachedMetricsResponse,
} from "../../user/types";

@Service()
@Resolver()
export class GetActiveUserMetricsSubscriptionResolver {
    @Subscription(() => UsersAttachedMetricsDto, {
        topics: "activeUserMetrics",
    })
    async getActiveUserMetrics(
        @Root() data: UsersAttachedMetricsResponse
    ): Promise<UsersAttachedMetricsDto> {
        return data.data[0];
    }
}

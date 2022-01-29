import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import {
    ActiveUserMetricsDto,
    ActiveUserMetricsResponse,
} from "../../user/types";

@Service()
@Resolver()
export class GetActiveUserMetricsSubscriptionResolver {
    @Subscription(() => ActiveUserMetricsDto, {
        topics: "activeUserMetrics",
    })
    async getActiveUserMetrics(
        @Root() data: ActiveUserMetricsResponse
    ): Promise<ActiveUserMetricsDto> {
        return data.data[0];
    }
}

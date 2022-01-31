import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import {
    UsersAttachedMetricsDto,
    UsersAttachedMetricsResponse,
} from "../../user/types";

@Service()
@Resolver()
export class GetUsersAttachedMetricsSubscriptionResolver {
    @Subscription(() => UsersAttachedMetricsDto, {
        topics: "usersAttachedMetrics",
    })
    async getUsersAttachedMetrics(
        @Root() data: UsersAttachedMetricsResponse
    ): Promise<UsersAttachedMetricsDto> {
        return data.data[0];
    }
}

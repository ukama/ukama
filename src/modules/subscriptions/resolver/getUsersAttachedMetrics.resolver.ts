import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { UsersAttachedMetricsDto } from "../../user/types";

@Service()
@Resolver()
export class GetUsersAttachedMetricsSubscriptionResolver {
    @Subscription(() => UsersAttachedMetricsDto, {
        topics: "usersAttachedMetrics",
    })
    async getUsersAttachedMetrics(
        @Root() data: [UsersAttachedMetricsDto]
    ): Promise<UsersAttachedMetricsDto> {
        return data[0];
    }
}

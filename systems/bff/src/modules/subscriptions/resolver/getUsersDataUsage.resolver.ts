import { Service } from "typedi";
import { GetUserDto } from "../../user/types";
import { Resolver, Root, Subscription } from "type-graphql";

@Service()
@Resolver()
export class GetUsersDataUsageSubscriptionResolver {
    @Subscription(() => GetUserDto, {
        topics: "getUsersSub",
    })
    async getUsersDataUsage(@Root() user: GetUserDto): Promise<GetUserDto> {
        return user;
    }
}

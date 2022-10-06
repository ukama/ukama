import {
    Arg,
    Ctx,
    Query,
    PubSub,
    Resolver,
    UseMiddleware,
    PubSubEngine,
} from "type-graphql";
import { Service } from "typedi";
import { UserService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { DataUsageInputDto, GetUserDto } from "../types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetUsersDataUsageResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => [GetUserDto])
    @UseMiddleware(Authentication)
    async getUsersDataUsage(
        @Arg("data") data: DataUsageInputDto,
        @PubSub() pubsub: PubSubEngine,
        @Ctx() ctx: Context
    ): Promise<GetUserDto[]> {
        const users: GetUserDto[] = [];
        if (data.ids.length > 0) {
            for (let i = 0; i < data.ids.length; i++) {
                const user = await this.userService.getUser(
                    data.ids[i],
                    parseCookie(ctx)
                );
                pubsub.publish("getUsersSub", user);
                users.push(user);
            }
        }
        return users;
    }
}

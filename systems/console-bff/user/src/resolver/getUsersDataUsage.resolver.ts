import {
    Arg,
    Ctx,
    PubSub,
    PubSubEngine,
    Query,
    Resolver,
    UseMiddleware,
} from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { DataUsageInputDto, GetUserDto } from "../types";

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
        const { dataSources } = ctx;

        const users: GetUserDto[] = [];
        if (data.ids.length > 0) {
            for (let i = 0; i < data.ids.length; i++) {
                const user = await dataSources.dataSource.getUser(
                    data.ids[i],
                    parseHeaders(ctx)
                );
                pubsub.publish("getUsersSub", user);
                // users.push(user);
            }
        }
        return users;
    }
}

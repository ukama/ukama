import { Resolver, Query, Arg, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { ConnectedUserDto } from "../types";
import { UserService } from "../service";
import { TIME_FILTER } from "../../../constants";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetConnectedUsersResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => ConnectedUserDto)
    @UseMiddleware(Authentication)
    async getConnectedUsers(
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER
    ): Promise<ConnectedUserDto> {
        return this.userService.getConnectedUsers(filter);
    }
}

import { Resolver, Query, Arg } from "type-graphql";
import { Service } from "typedi";
import { ConnectedUserDto } from "../types";
import { UserService } from "../service";
import { TIME_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetConnectedUsersResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => ConnectedUserDto)
    async getConnectedUsers(
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER
    ): Promise<ConnectedUserDto> {
        return this.userService.getConnectedUsers(filter);
    }
}

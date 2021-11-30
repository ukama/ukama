import { Resolver, Arg, Query, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { GetUserDto } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class DeleteUserResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => GetUserDto)
    @UseMiddleware(Authentication)
    async getUser(
        @Arg("id")
        id: string
    ): Promise<GetUserDto | null> {
        return this.userService.getUser(id);
    }
}

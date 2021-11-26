import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { UpdateUserDto, UserResponse } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class UpdateUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserResponse)
    @UseMiddleware(Authentication)
    async updateUser(
        @Arg("data")
        req: UpdateUserDto
    ): Promise<UserResponse | null> {
        return this.userService.updateUser(req);
    }
}

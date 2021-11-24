import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { UpdateUserDto, ActivateUserResponse } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class UpdateUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => ActivateUserResponse)
    @UseMiddleware(Authentication)
    async updateUser(
        @Arg("data")
        req: UpdateUserDto
    ): Promise<ActivateUserResponse | null> {
        return await this.userService.updateUser(req);
    }
}

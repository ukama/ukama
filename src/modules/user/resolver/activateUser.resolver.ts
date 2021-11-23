import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { ActivateUserDto, ActivateUserResponse } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class ActivateUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => ActivateUserResponse)
    @UseMiddleware(Authentication)
    async activateUser(
        @Arg("data")
        req: ActivateUserDto
    ): Promise<ActivateUserResponse | null> {
        return await this.userService.activateUser(req);
    }
}

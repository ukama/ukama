import { Resolver, Arg, Mutation } from "type-graphql";
import { Service } from "typedi";
import { ActivateUserDto, ActivateUserResponse } from "../types";
import { UserService } from "../service";

@Service()
@Resolver()
export class ActivateUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => ActivateUserResponse)
    async activateUser(
        @Arg("data")
        req: ActivateUserDto
    ): Promise<ActivateUserResponse | null> {
        return await this.userService.activateUser(req);
    }
}

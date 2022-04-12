import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { DeactivateResponse } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class DeleteUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => DeactivateResponse)
    @UseMiddleware(Authentication)
    async deactivateUser(
        @Arg("id")
        id: string
    ): Promise<DeactivateResponse | null> {
        return this.userService.deactivateUser(id);
    }
}

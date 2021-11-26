import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { DeleteUserResponse } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class DeleteUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => DeleteUserResponse)
    @UseMiddleware(Authentication)
    async deleteUser(
        @Arg("id")
        id: string
    ): Promise<DeleteUserResponse | null> {
        return this.userService.deleteUser(id);
    }
}

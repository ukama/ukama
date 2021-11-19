import { Resolver, Query, Arg, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { UserService } from "../service";
import { GetUserResponse, GetUserPaginationDto } from "../types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetDataBillResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => GetUserResponse)
    @UseMiddleware(Authentication)
    async getUsers(
        @Arg("data") data: GetUserPaginationDto
    ): Promise<GetUserResponse> {
        return this.userService.getUsers(data);
    }
}

import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkService } from "../service";
import { AddNetworkInputDto, NetworkDto } from "../types";

@Service()
@Resolver()
export class AddNetworkResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Mutation(() => NetworkDto)
    @UseMiddleware(Authentication)
    async addNetwork(
        @Arg("data") data: AddNetworkInputDto,
        @Ctx() ctx: Context,
    ): Promise<NetworkDto> {
        return this.networkService.addNetwork(data, parseCookie(ctx));
    }
}

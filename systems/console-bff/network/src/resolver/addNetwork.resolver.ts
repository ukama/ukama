import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddNetworkInputDto, NetworkDto } from "../types";

@Resolver()
export class AddNetworkResolver {

    @Mutation(() => NetworkDto)
    @UseMiddleware(Authentication)
    async addNetwork(
        @Arg("data") data: AddNetworkInputDto,
        @Ctx() ctx: Context
    ): Promise<NetworkDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.addNetwork(data, parseHeaders(ctx));
    }
}

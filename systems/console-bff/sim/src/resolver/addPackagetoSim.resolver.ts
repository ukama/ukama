import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddPackageSimResDto, AddPackageToSimInputDto } from "../types";

@Resolver()
export class AddPackageToSimResolver {

    @Mutation(() => AddPackageSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: AddPackageToSimInputDto,
        @Ctx() ctx: Context
    ): Promise<AddPackageSimResDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.addPackegeToSim(data, parseHeaders(ctx));
    }
}

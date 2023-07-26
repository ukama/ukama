import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import {
    RemovePackageFormSimInputDto,
    RemovePackageFromSimResDto,
} from "../types";

@Resolver()
export class RemovePackageForSimResolver {

    @Mutation(() => RemovePackageFromSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: RemovePackageFormSimInputDto,
        @Ctx() ctx: Context
    ): Promise<RemovePackageFromSimResDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.removePackageFromSim(
            data,
            parseHeaders(ctx)
        );
    }
}

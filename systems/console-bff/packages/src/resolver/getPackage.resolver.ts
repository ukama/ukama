import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { PackageDto } from "../types";

@Resolver()
export class GetPackageResolver {

    @Query(() => PackageDto)
    @UseMiddleware(Authentication)
    async getPackage(
        @Arg("packageId") packageId: string,
        @Ctx() ctx: Context
    ): Promise<PackageDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getPackage(packageId, parseHeaders(ctx));
    }
}

import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { PackageDto, UpdatePackageInputDto } from "../types";

@Resolver()
export class UpdatePackageResolver {
    @Mutation(() => PackageDto)
    @UseMiddleware(Authentication)
    async updatePackage(
        @Arg("packageId") packageId: string,
        @Arg("data") data: UpdatePackageInputDto,
        @Ctx() ctx: Context
    ): Promise<PackageDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.updatePackage(
            packageId,
            data,
            parseHeaders(ctx)
        );
    }
}

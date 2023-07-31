import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddPackageInputDto, PackageDto } from "../types";

@Resolver()
export class AddPackageResolver {

    @Mutation(() => PackageDto)
    @UseMiddleware(Authentication)
    async addPackage(
        @Arg("data") data: AddPackageInputDto,
        @Ctx() ctx: Context
    ): Promise<PackageDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.addPackage(data, parseHeaders(ctx));
    }
}

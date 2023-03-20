import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { PackageDto } from "../types";
import { PackageService } from "./../service";

@Service()
@Resolver()
export class GetPackageResolver {
    constructor(private readonly packageService: PackageService) {}

    @Query(() => PackageDto)
    @UseMiddleware(Authentication)
    async getPackage(
        @Arg("packageId") packageId: string,
        @Ctx() ctx: Context,
    ): Promise<PackageDto> {
        return this.packageService.getPackage(packageId, parseCookie(ctx));
    }
}

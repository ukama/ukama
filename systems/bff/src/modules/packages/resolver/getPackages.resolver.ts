import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { PackagesResDto } from "../types";
import { PackageService } from "./../service";

@Service()
@Resolver()
export class GetPackagesResolver {
    constructor(private readonly packageService: PackageService) {}

    @Query(() => PackagesResDto)
    @UseMiddleware(Authentication)
    async getPackages(@Ctx() ctx: Context): Promise<PackagesResDto> {
        return this.packageService.getPackages(parseCookie(ctx));
    }
}

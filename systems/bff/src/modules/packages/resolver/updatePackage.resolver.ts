import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { PackageService } from "../service";
import { PackageDto, UpdatePackageInputDto } from "../types";

@Service()
@Resolver()
export class UpdatePackageResolver {
    constructor(private readonly packageService: PackageService) {}

    @Mutation(() => PackageDto)
    @UseMiddleware(Authentication)
    async updatePackage(
        @Arg("packageId") packageId: string,
        @Arg("data") data: UpdatePackageInputDto,
        @Ctx() ctx: Context
    ): Promise<PackageDto> {
        return this.packageService.updatePackage(
            packageId,
            data,
            parseCookie(ctx)
        );
    }
}

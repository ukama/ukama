import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context, IdResponse } from "../../../common/types";
import { PackageService } from "../service";

@Service()
@Resolver()
export class DeletePackageResolver {
    constructor(private readonly packageService: PackageService) {}

    @Mutation(() => IdResponse)
    @UseMiddleware(Authentication)
    async deletePackage(
        @Arg("packageId") packageId: string,
        @Ctx() ctx: Context
    ): Promise<IdResponse> {
        return this.packageService.deletePackage(packageId, parseCookie(ctx));
    }
}

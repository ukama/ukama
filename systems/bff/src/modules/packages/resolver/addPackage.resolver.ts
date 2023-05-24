import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddPackageInputDto, PackageDto } from "../types";
import { PackageService } from "./../service";

@Service()
@Resolver()
export class AddPackageResolver {
    constructor(private readonly packageService: PackageService) {}

    @Mutation(() => PackageDto)
    @UseMiddleware(Authentication)
    async addPackage(
        @Arg("data") data: AddPackageInputDto,
        @Ctx() ctx: Context
    ): Promise<PackageDto> {
        return this.packageService.addPackage(data, parseHeaders(ctx));
    }
}

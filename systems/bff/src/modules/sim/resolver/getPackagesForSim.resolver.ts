import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { GetPackagesForSimInputDto, GetPackagesForSimResDto } from "../types";

@Service()
@Resolver()
export class GetPackagesForSimResolver {
    constructor(private readonly simService: SimService) {}

    @Query(() => GetPackagesForSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: GetPackagesForSimInputDto,
        @Ctx() ctx: Context
    ): Promise<GetPackagesForSimResDto> {
        return await this.simService.getPackagesForSim(data, parseHeaders(ctx));
    }
}
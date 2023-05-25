import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { SimPoolStatsDto } from "../types";

@Service()
@Resolver()
export class GetSimPoolStatsResolver {
    constructor(private readonly simService: SimService) {}

    @Query(() => SimPoolStatsDto)
    @UseMiddleware(Authentication)
    async getSimPoolStats(
        @Arg("type") type: string,
        @Ctx() ctx: Context
    ): Promise<SimPoolStatsDto> {
        return await this.simService.getSimPoolStats(type, parseHeaders(ctx));
    }
}

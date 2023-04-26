import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { DefaultMarkupHistoryResDto } from "../types";
import { RateService } from "./../service";

@Service()
@Resolver()
export class GetDefaultMarkupHistoryResolver {
    constructor(private readonly rateService: RateService) {}

    @Query(() => DefaultMarkupHistoryResDto)
    @UseMiddleware(Authentication)
    async getDefaultMarkupHistory(
        @Ctx() ctx: Context,
    ): Promise<DefaultMarkupHistoryResDto> {
        return this.rateService.getDefaultMarkupHistory(parseCookie(ctx));
    }
}

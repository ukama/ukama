import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { DefaultMarkupResDto } from "../types";
import { RateService } from "./../service";

@Service()
@Resolver()
export class GetDefaultMarkupResolver {
    constructor(private readonly rateService: RateService) {}

    @Query(() => DefaultMarkupResDto)
    @UseMiddleware(Authentication)
    async getDefaultMarkup(@Ctx() ctx: Context): Promise<DefaultMarkupResDto> {
        return this.rateService.getDefaultMarkup(parseCookie(ctx));
    }
}

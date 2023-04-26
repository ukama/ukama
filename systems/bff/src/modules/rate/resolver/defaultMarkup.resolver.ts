import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { DefaultMarkupInputDto } from "../types";
import { RateService } from "./../service";

@Service()
@Resolver()
export class DefaultMarkupResolver {
    constructor(private readonly rateService: RateService) {}

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async defaultMarkup(
        @Arg("data") data: DefaultMarkupInputDto,
        @Ctx() ctx: Context,
    ): Promise<BoolResponse> {
        return this.rateService.defaultMarkup(data, parseCookie(ctx));
    }
}

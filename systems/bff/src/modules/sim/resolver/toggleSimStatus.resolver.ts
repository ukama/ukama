import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { SimStatusResDto, ToggleSimStatusInputDto } from "../types";

@Service()
@Resolver()
export class ToggleSimStatusResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => SimStatusResDto)
    @UseMiddleware(Authentication)
    async toggleSimStatus(
        @Arg("data") data: ToggleSimStatusInputDto,
        @Ctx() ctx: Context
    ): Promise<SimStatusResDto> {
        return await this.simService.toggleSimStatus(data, parseHeaders(ctx));
    }
}

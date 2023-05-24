import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import {
    RemovePackageFormSimInputDto,
    RemovePackageFromSimResDto,
} from "../types";

@Service()
@Resolver()
export class RemovePackageForSimResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => RemovePackageFromSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: RemovePackageFormSimInputDto,
        @Ctx() ctx: Context
    ): Promise<RemovePackageFromSimResDto> {
        return await this.simService.removePackageFromSim(
            data,
            parseHeaders(ctx)
        );
    }
}

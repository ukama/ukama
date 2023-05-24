import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import {
    SetActivePackageForSimInputDto,
    SetActivePackageForSimResDto,
} from "../types";

@Service()
@Resolver()
export class SetActivePackageForSimResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => SetActivePackageForSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: SetActivePackageForSimInputDto,
        @Ctx() ctx: Context
    ): Promise<SetActivePackageForSimResDto> {
        return await this.simService.setActivePackageForSim(
            data,
            parseHeaders(ctx)
        );
    }
}

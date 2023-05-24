import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { AddPackageSimResDto, AddPackageToSimInputDto } from "../types";

@Service()
@Resolver()
export class AddPackageToSimResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => AddPackageSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: AddPackageToSimInputDto,
        @Ctx() ctx: Context
    ): Promise<AddPackageSimResDto> {
        return await this.simService.addPackegeToSim(data, parseHeaders(ctx));
    }
}

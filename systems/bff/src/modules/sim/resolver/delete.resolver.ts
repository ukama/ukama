import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { DeleteSimInputDto, DeleteSimResDto } from "../types";

@Service()
@Resolver()
export class DeleteSimResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => DeleteSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: DeleteSimInputDto,
        @Ctx() ctx: Context
    ): Promise<DeleteSimResDto> {
        return await this.simService.deleteSim(data, parseHeaders(ctx));
    }
}

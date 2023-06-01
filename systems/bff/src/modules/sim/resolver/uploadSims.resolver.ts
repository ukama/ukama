import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { UploadSimsInputDto, UploadSimsResDto } from "../types";

@Service()
@Resolver()
export class UploadSimsResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => UploadSimsResDto)
    @UseMiddleware(Authentication)
    async uploadSims(
        @Arg("data") data: UploadSimsInputDto,
        @Ctx() ctx: Context
    ): Promise<UploadSimsResDto> {
        return await this.simService.uploadSims(data, parseHeaders(ctx));
    }
}

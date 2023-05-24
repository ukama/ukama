import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NodeService } from "../service";
import { UpdateNodeDto, UpdateNodeResponse } from "../types";

@Service()
@Resolver()
export class UpdateNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => UpdateNodeResponse)
    @UseMiddleware(Authentication)
    async updateNode(
        @Arg("data")
        req: UpdateNodeDto,
        @Ctx() ctx: Context
    ): Promise<UpdateNodeResponse | null> {
        return this.nodeService.updateNode(req, parseHeaders(ctx));
    }
}

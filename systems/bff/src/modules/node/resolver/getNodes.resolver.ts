import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NodeService } from "../service";
import { GetNodesResDto } from "../types";

@Service()
@Resolver()
export class GetNodesResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => GetNodesResDto)
    @UseMiddleware(Authentication)
    async getNodes(@Ctx() ctx: Context): Promise<GetNodesResDto> {
        return this.nodeService.getNodes(parseHeaders(ctx));
    }
}

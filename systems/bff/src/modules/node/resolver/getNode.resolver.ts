import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NodeService } from "../service";
import { NodeResponse } from "../types";

@Service()
@Resolver()
export class GetNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodeResponse)
    @UseMiddleware(Authentication)
    async getNode(
        @Arg("nodeId")
        nodeId: string,
        @Ctx() ctx: Context
    ): Promise<NodeResponse> {
        return this.nodeService.getNode(nodeId, parseHeaders(ctx));
    }
}

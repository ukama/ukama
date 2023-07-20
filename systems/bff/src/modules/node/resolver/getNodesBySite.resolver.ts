import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NodeService } from "../service";
import { NodesBySiteResDto } from "../types";

@Service()
@Resolver()
export class GetNodesBySiteResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodesBySiteResDto)
    @UseMiddleware(Authentication)
    async getNodesBySite(
        @Arg("siteId")
        siteId: string,
        @Ctx() ctx: Context
    ): Promise<NodesBySiteResDto> {
        return this.nodeService.getNodesBySite(siteId, parseHeaders(ctx));
    }
}

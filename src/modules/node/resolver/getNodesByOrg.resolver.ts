import { Resolver, Query, UseMiddleware, Arg, Ctx } from "type-graphql";
import { Service } from "typedi";
import { OrgNodeResponse } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";

@Service()
@Resolver()
export class GetNodesByOrgResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => OrgNodeResponse)
    @UseMiddleware(Authentication)
    async getNodesByOrg(
        @Arg("orgId") orgId: string,
        @Ctx() ctx: Context
    ): Promise<OrgNodeResponse> {
        return this.nodeService.getNodesByOrg(orgId, ctx.session);
    }
}

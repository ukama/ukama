import { Resolver, Query, UseMiddleware, Ctx } from "type-graphql";
import { Service } from "typedi";
import { OrgNodeResponseDto } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { parseCookie } from "../../../common";

@Service()
@Resolver()
export class GetNodesByOrgResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => OrgNodeResponseDto)
    @UseMiddleware(Authentication)
    async getNodesByOrg(@Ctx() ctx: Context): Promise<OrgNodeResponseDto> {
        return this.nodeService.getNodesByOrg(parseCookie(ctx));
    }
}

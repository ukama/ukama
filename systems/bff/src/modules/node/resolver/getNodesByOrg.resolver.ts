import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NodeService } from "../service";
import { OrgNodeResponseDto } from "../types";

@Service()
@Resolver()
export class GetNodesByOrgResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => OrgNodeResponseDto)
    @UseMiddleware(Authentication)
    async getNodesByOrg(@Ctx() ctx: Context): Promise<OrgNodeResponseDto> {
        return this.nodeService.getNodesByOrg(parseHeaders(ctx));
    }
}

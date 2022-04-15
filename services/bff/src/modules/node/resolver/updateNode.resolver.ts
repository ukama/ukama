import { Service } from "typedi";
import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { OrgNodeDto, UpdateNodeDto } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { parseCookie } from "../../../common";

@Service()
@Resolver()
export class UpdateNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => OrgNodeDto)
    @UseMiddleware(Authentication)
    async updateNode(
        @Arg("data")
        req: UpdateNodeDto,
        @Ctx() ctx: Context
    ): Promise<OrgNodeDto | null> {
        return this.nodeService.updateNode(req, parseCookie(ctx));
    }
}

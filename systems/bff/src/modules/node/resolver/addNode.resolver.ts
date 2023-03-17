import { Service } from "typedi";
import { NodeService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { getNodes, getTowerNode, linkNodes } from "../../../utils";
import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { AddNodeDto, AddNodeResponse, LinkNodes, NodeObj } from "../types";

@Service()
@Resolver()
export class AddNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => AddNodeResponse)
    @UseMiddleware(Authentication)
    async addNode(
        @Arg("data")
        req: AddNodeDto,
        @Ctx() ctx: Context,
    ): Promise<AddNodeResponse | null> {
        const nodes: NodeObj[] = getNodes(req);
        const rootNode: NodeObj = getTowerNode(req);
        const linkedNode: LinkNodes = linkNodes(nodes, rootNode?.nodeId || "");

        if (nodes.length > 0) {
            for (const node of nodes) {
                await this.nodeService.addNode(node, parseCookie(ctx));
            }
        }

        if (!req.associate) {
            await this.nodeService.addNode(rootNode, parseCookie(ctx));
        }

        if (linkedNode.attached && linkedNode.attached.length > 0) {
            await this.nodeService.linkNodes(linkedNode, parseCookie(ctx));
        }

        return { success: true };
    }
}

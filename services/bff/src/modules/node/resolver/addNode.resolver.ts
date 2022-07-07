import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { AddNodeDto, AddNodeResponse, LinkNodes, NodeObj } from "../types";

const isTowerNode = (nodeId: string) => nodeId.includes("tnode");
const getTowerNode = (payload: AddNodeDto): NodeObj => {
    if (isTowerNode(payload.nodeId))
        return { name: payload.name, nodeId: payload.nodeId };

    if (payload.attached)
        for (const node of payload.attached) {
            if (isTowerNode(node.nodeId)) return node;
        }
    return { name: "", nodeId: "" };
};
const getNodes = (payload: AddNodeDto) => {
    const nodes: NodeObj[] = [];
    if (!isTowerNode(payload.nodeId)) {
        nodes.push({ name: payload.name, nodeId: payload.nodeId });
    }
    if (payload.attached)
        for (const node of payload.attached) {
            if (!isTowerNode(node.nodeId))
                nodes.push({ name: node.name, nodeId: node.nodeId });
        }
    return nodes;
};
const linkNodes = (nodes: NodeObj[], rootNodeId: string) => {
    const nodesLinkingObj: LinkNodes = {
        nodeId: rootNodeId,
        attached: [],
    };
    for (const node of nodes) {
        nodesLinkingObj.attached?.push({ nodeId: node.nodeId });
    }
    return nodesLinkingObj;
};

@Service()
@Resolver()
export class AddNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => AddNodeResponse)
    @UseMiddleware(Authentication)
    async addNode(
        @Arg("data")
        req: AddNodeDto,
        @Ctx() ctx: Context
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

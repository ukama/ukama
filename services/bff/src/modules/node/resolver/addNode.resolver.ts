import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { AddNodeDto, AddNodeResponse, LinkNodes } from "../types";

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
        const res = await this.nodeService.addNode(req, parseCookie(ctx));
        if (req.attached && req.attached.length > 0) {
            const nodesLinkingObj: LinkNodes = {
                nodeId: req.nodeId,
                attached: [],
            };
            if (req.associate) {
                nodesLinkingObj.attached?.push({
                    nodeId: req.attached[0].nodeId,
                });
            } else {
                for (let i = 0; i < req.attached.length; i++) {
                    nodesLinkingObj.attached?.push({
                        nodeId: req.attached[i].nodeId,
                    });
                    await this.nodeService.addNode(
                        req.attached[i],
                        parseCookie(ctx)
                    );
                }
            }

            await this.nodeService.linkNodes(nodesLinkingObj, parseCookie(ctx));
        }
        return res;
    }
}

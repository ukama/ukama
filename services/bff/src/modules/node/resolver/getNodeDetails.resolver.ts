import { Resolver, Query, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { NodeDetailDto } from "../types";
import { NodeService } from "../service";

import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodeDetailsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodeDetailDto)
    @UseMiddleware(Authentication)
    async getNodeDetails(): Promise<NodeDetailDto> {
        return this.nodeService.getNodeDetials();
    }
}

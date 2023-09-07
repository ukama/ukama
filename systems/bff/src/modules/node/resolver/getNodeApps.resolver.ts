import { Service } from "typedi";
import { NodeService } from "../service";
import { NodeAppResponse } from "../types";
import { Resolver, Query, UseMiddleware } from "type-graphql";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodeAppsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [NodeAppResponse])
    @UseMiddleware(Authentication)
    async getNodeApps(): Promise<NodeAppResponse[]> {
        return this.nodeService.getNodeApps();
    }
}

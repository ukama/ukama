import { Service } from "typedi";
import { NodeService } from "../service";
import { NodeAppsVersionLogsResponse } from "../types";
import { Resolver, Query, UseMiddleware } from "type-graphql";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodeAppsVersionLogsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [NodeAppsVersionLogsResponse])
    @UseMiddleware(Authentication)
    async getNodeAppsVersionLogs(): Promise<NodeAppsVersionLogsResponse[]> {
        return this.nodeService.getSoftwareLogs();
    }
}

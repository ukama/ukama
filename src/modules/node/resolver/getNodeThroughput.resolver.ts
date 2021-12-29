import { Resolver, Query, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { GraphDto } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodeThroughputResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => GraphDto)
    @UseMiddleware(Authentication)
    async getNodeThroughput(): Promise<GraphDto> {
        return this.nodeService.getNodeThroughput();
    }
}

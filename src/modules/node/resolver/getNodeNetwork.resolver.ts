import { Resolver, Query, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { NetworkDto } from "../../network/types";

@Service()
@Resolver()
export class GetNodeNetworkResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NetworkDto)
    @UseMiddleware(Authentication)
    async getNodeNetwork(): Promise<NetworkDto> {
        return this.nodeService.getNetwork();
    }
}

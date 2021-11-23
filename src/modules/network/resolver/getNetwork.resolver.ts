import { Resolver, Query, Arg, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { NetworkDto } from "../types";
import { NetworkService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { NETWORK_TYPE } from "../../../constants";

@Service()
@Resolver()
export class GetNetworkResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Query(() => NetworkDto)
    @UseMiddleware(Authentication)
    async getNetwork(
        @Arg("filter", () => NETWORK_TYPE) filter: NETWORK_TYPE
    ): Promise<NetworkDto> {
        return this.networkService.getNetwork(filter);
    }
}

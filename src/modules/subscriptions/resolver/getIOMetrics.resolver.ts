import { Service } from "typedi";
import { Resolver, Root, Subscription } from "type-graphql";
import { IOMetricsDto, IOMetricsResponse } from "../../node/types";

@Service()
@Resolver()
export class GetIOMetricsSubscriptionResolver {
    @Subscription(() => IOMetricsDto, {
        topics: "ioMetrics",
    })
    async getIOMetrics(@Root() data: IOMetricsResponse): Promise<IOMetricsDto> {
        return data.data[0];
    }
}

import { Service } from "typedi";
import { IOMetricsDto } from "../../node/types";
import { Resolver, Root, Subscription } from "type-graphql";

@Service()
@Resolver()
export class GetIOMetricsSubscriptionResolver {
    @Subscription(() => IOMetricsDto, {
        topics: "ioMetrics",
    })
    async getIOMetrics(@Root() data: IOMetricsDto): Promise<IOMetricsDto> {
        return data;
    }
}

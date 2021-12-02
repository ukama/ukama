import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { DataUsageDto } from "../../data/types";

@Service()
@Resolver()
export class DataUsageSubscriptionResolver {
    @Subscription(() => DataUsageDto, {
        topics: "DATAUSAGE",
    })
    async DATAUSAGE(@Root() data: DataUsageDto): Promise<DataUsageDto> {
        return data;
    }
}

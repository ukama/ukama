import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { DataBillDto } from "../../data/types";

@Service()
@Resolver()
export class DataBillSubscriptionResolver {
    @Subscription(() => DataBillDto, {
        topics: "DATABILL",
    })
    async DATABILL(@Root() bill: DataBillDto): Promise<DataBillDto> {
        return bill;
    }
}

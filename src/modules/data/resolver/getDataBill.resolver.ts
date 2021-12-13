import {
    Resolver,
    Query,
    Arg,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { DataBillDto } from "../types";
import { DataService } from "../service";
import { DATA_BILL_FILTER } from "../../../constants";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetDataBillResolver {
    constructor(private readonly dataService: DataService) {}

    @Query(() => DataBillDto)
    @UseMiddleware(Authentication)
    async getDataBill(
        @Arg("filter", () => DATA_BILL_FILTER) filter: DATA_BILL_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<DataBillDto> {
        const bill = this.dataService.getDataBill(filter);
        pubsub.publish("getDataBill", bill);
        return bill;
    }
}

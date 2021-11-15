import { Resolver, Query, Arg } from "type-graphql";
import { Service } from "typedi";
import { DataBillDto } from "../types";
import { DataService } from "../service";
import { DATA_BILL_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetDataBillResolver {
    constructor(private readonly dataService: DataService) {}

    @Query(() => DataBillDto)
    async getDataBill(
        @Arg("filter", () => DATA_BILL_FILTER) filter: DATA_BILL_FILTER
    ): Promise<DataBillDto> {
        return this.dataService.getDataBill(filter);
    }
}

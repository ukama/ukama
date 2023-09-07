import { Resolver, Query, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { BillHistoryDto } from "../types";
import { BillService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetBillHistoryResolver {
    constructor(private readonly billService: BillService) {}

    @Query(() => [BillHistoryDto])
    @UseMiddleware(Authentication)
    async getBillHistory(): Promise<BillHistoryDto[]> {
        return this.billService.getBillHistory();
    }
}

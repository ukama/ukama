import { Resolver, Query, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { BillResponse } from "../types";
import { BillService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetCurrentBillResolver {
    constructor(private readonly billService: BillService) {}

    @Query(() => BillResponse)
    @UseMiddleware(Authentication)
    async getCurrentBill(): Promise<BillResponse> {
        return this.billService.getCurrentBill();
    }
}

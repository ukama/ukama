import { Resolver, Query, UseMiddleware } from "type-graphql";
import { BillResponse } from "../types";
import { Authentication } from "../../../common/Authentication";

@Resolver()
export class GetCurrentBillResolver {

    @Query(() => BillResponse)
    @UseMiddleware(Authentication)
    async getCurrentBill(): Promise<BillResponse> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getCurrentBill();
    }
}

import { Resolver, Query, UseMiddleware } from "type-graphql";
import { BillHistoryDto } from "../types";
import { Authentication } from "../../../common/Authentication";

@Resolver()
export class GetBillHistoryResolver {

    @Query(() => [BillHistoryDto])
    @UseMiddleware(Authentication)
    async getBillHistory(): Promise<BillHistoryDto[]> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getBillHistory();
    }
}

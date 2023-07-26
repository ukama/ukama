import { RESTDataSource } from "@apollo/datasource-rest";
import { BillHistoryDto, BillResponse } from "../types";
import BillMapper from "./mapper";
import { SERVER } from "../../../constants/endpoints";

export class BillingApi extends RESTDataSource  {
    public getCurrentBill = async (): Promise<BillResponse> => {
        return this.get(SERVER.GET_CURRENT_BILL).then(res => 
            BillMapper.dtoToDto(res));
    };

    public getBillHistory = async (): Promise<BillHistoryDto[]> => {
        return this.get(SERVER.GET_BILL_HISTORY).then(res => 
            BillMapper.billHistoryDtoToDto(res));
    };
}

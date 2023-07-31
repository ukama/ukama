import { RESTDataSource } from "@apollo/datasource-rest";

import { SERVER } from "../../constants/endpoints";
import { BillHistoryDto, BillResponse } from "../resolvers/types";
import BillMapper from "./mapper";

class BillingAPI extends RESTDataSource {
  public getCurrentBill = async (): Promise<BillResponse> => {
    return this.get(SERVER.GET_CURRENT_BILL).then(res =>
      BillMapper.dtoToDto(res)
    );
  };

  public getBillHistory = async (): Promise<BillHistoryDto[]> => {
    return this.get(SERVER.GET_BILL_HISTORY).then(res =>
      BillMapper.billHistoryDtoToDto(res)
    );
  };
}

export default BillingAPI;

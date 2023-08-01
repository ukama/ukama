import { RESTDataSource } from "@apollo/datasource-rest";

import { BILLING_API_GW } from "../../common/configs";
import { BillHistoryDto, BillResponse } from "../resolvers/types";
import { billHistoryDtoToDto, dtoToDto } from "./mapper";

const version = "/v1/invoices";
class BillingAPI extends RESTDataSource {
  baseURL = BILLING_API_GW + version;
  public getCurrentBill = async (): Promise<BillResponse> => {
    return this.get("/current").then(res => dtoToDto(res));
  };

  public getBillHistory = async (): Promise<BillHistoryDto[]> => {
    return this.get("/history").then(res => billHistoryDtoToDto(res));
  };
}

export default BillingAPI;

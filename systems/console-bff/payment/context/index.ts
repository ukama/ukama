import { THeaders } from "../../common/types";
import PaymentAPI from "../datasource/payment_api";

export interface Context {
  baseURL: string;
  dataSources: {
    dataSource: PaymentAPI;
  };
  headers: THeaders;
}

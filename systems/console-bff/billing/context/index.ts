import BillingAPI from "../datasource/billing_api";

export interface Context {
  dataSources: {
    dataSource: BillingAPI;
  };
}

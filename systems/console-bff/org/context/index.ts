import OrgAPI from "../datasource/org_api";

export interface Context {
  dataSources: {
    dataSource: OrgAPI;
  };
}

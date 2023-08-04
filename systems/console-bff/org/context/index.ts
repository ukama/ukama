import OrgAPI from "../datasource/org_api";
import UserAPI from "../datasource/user_api";

export interface Context {
  dataSources: {
    dataSource: OrgAPI;
    dataSoureceUser: UserAPI;
  };
  req: any;
}

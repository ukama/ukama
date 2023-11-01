import { THeaders } from "../../common/types";
import MemberAPI from "../datasource/member_api";

export interface Context {
  dataSources: {
    dataSource: MemberAPI;
  };
  headers: THeaders;
}

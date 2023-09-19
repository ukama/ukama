import { THeaders } from "../../common/types";
import InvitationAPI from "../datasource/invitation_api";

export interface Context {
  dataSources: {
    dataSource: InvitationAPI;
  };
  headers: THeaders;
}

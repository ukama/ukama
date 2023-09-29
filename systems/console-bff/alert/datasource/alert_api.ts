import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

import { getPaginatedOutput } from "../../common/utils";
import { AlertsResponse } from "../resolver/types";
import { PaginationDto } from "./../../common/types";
import { dtoToDto } from "./mapper";

class AlertApi extends RESTDataSource {
  baseURL = "";
  getAlerts = async (req: PaginationDto): Promise<AlertsResponse> => {
    return this.get(`/alerts`, {
      params: {
        pageNo: `${req.pageNo}`,
        pageSize: `${req.pageSize}`,
      },
    })
      .then(res => {
        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const alerts = dtoToDto(res);
        return {
          alerts,
          meta,
        };
      })
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default AlertApi;

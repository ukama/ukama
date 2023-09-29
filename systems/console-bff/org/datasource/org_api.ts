import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

import { NUCLEUS_API_GW, VERSION } from "../../common/configs";
import { OrgDto, OrgsResDto } from "../resolver/types";
import { dtoToOrgResDto, dtoToOrgsResDto } from "./mapper";

class OrgApi extends RESTDataSource {
  baseURL = NUCLEUS_API_GW;

  getOrgs = async (userId: string): Promise<OrgsResDto> => {
    return this.get(`/${VERSION}/orgs`, {
      params: {
        user_uuid: userId,
      },
    })
      .then(res => dtoToOrgsResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getOrg = async (orgName: string): Promise<OrgDto> => {
    return this.get(`/${VERSION}/orgs/${orgName}`)
      .then(res => dtoToOrgResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default OrgApi;

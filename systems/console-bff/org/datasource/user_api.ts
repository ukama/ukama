import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW, VERSION } from "../../common/configs";
import { UserResDto } from "./../../user/resolver/types";
import { dtoToUserResDto } from "./mapper";

class UserApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW;

  getUser = async (userId: string): Promise<UserResDto> => {
    return this.get(`/${VERSION}/users/${userId}`, {}).then(res =>
      dtoToUserResDto(res)
    );
  };
}
export default UserApi;

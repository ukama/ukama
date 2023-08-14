import { RESTDataSource } from "@apollo/datasource-rest";

import { NUCLEUS_API_GW, VERSION } from "../../common/configs";
import { UserResDto, WhoamiDto } from "../resolver/types";
import { dtoToUserResDto, dtoToWhoamiResDto } from "./mapper";

class UserApi extends RESTDataSource {
  baseURL = NUCLEUS_API_GW;

  getUser = async (userId: string): Promise<UserResDto> => {
    return this.get(`/${VERSION}/users/${userId}`, {}).then(res =>
      dtoToUserResDto(res)
    );
  };

  whoami = async (userId: string): Promise<WhoamiDto> => {
    return this.get(`/${VERSION}/users/whoami/${userId}`).then(res =>
      dtoToWhoamiResDto(res)
    );
  };
}
export default UserApi;

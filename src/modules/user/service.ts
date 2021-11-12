import { Service } from "typedi";
import { ConnectedUserDto } from "./types";
import { IUserService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import UserMapper from "./mapper";
import { TIME_FILTER } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { getUsersMethod } from "./io";

@Service()
export class UserService implements IUserService {
    getConnectedUsers = async (
        filter: TIME_FILTER
    ): Promise<ConnectedUserDto> => {
        const res = await getUsersMethod(
            SERVER.GET_CONNECTED_USERS,
            `${filter}`,
            null
        );

        if (!res) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        const connectedUsers = UserMapper.connectedUsersDtoToDto(res.data.data);

        if (!connectedUsers) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        return connectedUsers;
    };
}

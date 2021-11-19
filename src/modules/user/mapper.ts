import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserResponseDto,
} from "./types";

class UserMapper implements IUserMapper {
    connectedUsersDtoToDto = (res: ConnectedUserResponse): ConnectedUserDto => {
        const connectedUsers = res.data;
        return connectedUsers;
    };
    dtoToDto = (res: GetUserResponseDto): GetUserDto[] => {
        const users = res.data;
        return users;
    };
}
export default <IUserMapper>new UserMapper();

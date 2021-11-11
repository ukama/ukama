import { TIME_FILTER } from "../../constants";
import { ConnectedUserDto } from "./types";

export interface IUserService {
    getConnectedUsers(filter: TIME_FILTER): Promise<ConnectedUserDto>;
}

export interface IUserMapper {
    connectedUsersDtoToDto(res: ConnectedUserDto): ConnectedUserDto;
}

import { TIME_FILTER } from "../../constants";
import {
    ActivateUserDto,
    ActivateUserResponse,
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserPaginationDto,
    GetUserResponse,
    GetUserResponseDto,
} from "./types";

export interface IUserService {
    getConnectedUsers(filter: TIME_FILTER): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    getUsers(req: GetUserPaginationDto): Promise<GetUserResponse>;
}

export interface IUserMapper {
    connectedUsersDtoToDto(res: ConnectedUserResponse): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
}

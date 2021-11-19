import { Service } from "typedi";
import {
    ActivateUserDto,
    ActivateUserResponse,
    ActiveUserResponseDto,
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserPaginationDto,
    GetUserResponse,
    GetUserResponseDto,
} from "./types";
import { IUserService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import UserMapper from "./mapper";
import { API_METHOD_TYPE, TIME_FILTER } from "../../constants";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { getPaginatedOutput } from "../../utils";

@Service()
export class UserService implements IUserService {
    getConnectedUsers = async (
        filter: TIME_FILTER
    ): Promise<ConnectedUserDto> => {
        const res = await catchAsyncIOMethod<ConnectedUserResponse>({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_CONNECTED_USERS,
            params: `${filter}`,
        });

        const connectedUsers = UserMapper.connectedUsersDtoToDto(res);

        if (!connectedUsers) throw new HTTP404Error(Messages.USERS_NOT_FOUND);

        return connectedUsers;
    };
    activateUser = async (
        req: ActivateUserDto
    ): Promise<ActivateUserResponse> => {
        const res = await catchAsyncIOMethod<ActiveUserResponseDto>({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_ACTIVE_USER,
            body: req,
        });
        return res.data;
    };
    getUsers = async (req: GetUserPaginationDto): Promise<GetUserResponse> => {
        const res = await catchAsyncIOMethod<GetUserResponseDto>({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_USERS,
            params: req,
        });
        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const users = UserMapper.dtoToDto(res);
        if (!users) throw new HTTP404Error(Messages.USERS_NOT_FOUND);

        return {
            users,
            meta,
        };
    };
}

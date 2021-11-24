import { Service } from "typedi";
import {
    ActivateUserDto,
    ActivateUserResponse,
    ConnectedUserDto,
    GetUserPaginationDto,
    GetUserResponse,
    ResidentsResponse,
    UpdateUserDto,
} from "./types";
import { IUserService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import UserMapper from "./mapper";
import { API_METHOD_TYPE, TIME_FILTER } from "../../constants";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { getPaginatedOutput } from "../../utils";
import { PaginationDto } from "../../common/types";

@Service()
export class UserService implements IUserService {
    getConnectedUsers = async (
        filter: TIME_FILTER
    ): Promise<ConnectedUserDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_CONNECTED_USERS,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);
        const connectedUsers = UserMapper.connectedUsersDtoToDto(res);

        if (!connectedUsers) throw new HTTP404Error(Messages.USERS_NOT_FOUND);

        return connectedUsers;
    };

    activateUser = async (
        req: ActivateUserDto
    ): Promise<ActivateUserResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_ACTIVE_USER,
            body: req,
        });
        if (checkError(res)) throw new Error(res.message);

        return res.data;
    };

    updateUser = async (req: UpdateUserDto): Promise<ActivateUserResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_UPDATE_USER,
            body: req,
        });
        if (checkError(res)) throw new Error(res.message);

        return res.data;
    };

    getUsers = async (req: GetUserPaginationDto): Promise<GetUserResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_USERS,
            params: req,
        });
        if (checkError(res)) throw new Error(res.message);

        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const users = UserMapper.dtoToDto(res);
        if (!users) throw new HTTP404Error(Messages.USERS_NOT_FOUND);

        return {
            users,
            meta,
        };
    };

    getResidents = async (req: PaginationDto): Promise<ResidentsResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_USERS,
            params: req,
        });
        if (checkError(res)) throw new Error(res.message);

        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const residents = UserMapper.residentDtoToDto(res);
        if (!residents) throw new HTTP404Error(Messages.RESIDENTS_NOT_FOUND);

        return {
            residents,
            meta,
        };
    };
}

import { PaginationDto } from "../../common/types";
import { TIME_FILTER } from "../../constants";
import {
    ActivateUserDto,
    ActivateUserResponse,
    ConnectedUserDto,
    ConnectedUserResponse,
    DeleteResponse,
    GetUserDto,
    GetUserPaginationDto,
    GetUserResponse,
    GetUserResponseDto,
    ResidentResponse,
    ResidentsResponse,
    UpdateUserDto,
    UserResponse,
} from "./types";

export interface IUserService {
    getConnectedUsers(filter: TIME_FILTER): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    updateUser(req: UpdateUserDto): Promise<UserResponse>;
    deleteUser(id: string): Promise<DeleteResponse>;
    getUsers(req: GetUserPaginationDto): Promise<GetUserResponse>;
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
}

export interface IUserMapper {
    connectedUsersDtoToDto(res: ConnectedUserResponse): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    residentDtoToDto(res: GetUserResponseDto): ResidentResponse;
}

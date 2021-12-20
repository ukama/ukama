import { HeaderType, PaginationDto } from "../../common/types";
import { TIME_FILTER } from "../../constants";
import {
    ActivateUserDto,
    ActivateUserResponse,
    AddUserDto,
    AddUserResponse,
    ConnectedUserDto,
    ConnectedUserResponse,
    DeactivateResponse,
    GetUserDto,
    GetUserPaginationDto,
    GetUserResponse,
    GetUserResponseDto,
    OrgUserResponse,
    OrgUserResponseDto,
    ResidentResponse,
    ResidentsResponse,
    UpdateUserDto,
    UserResponse,
} from "./types";

export interface IUserService {
    getConnectedUsers(filter: TIME_FILTER): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    updateUser(req: UpdateUserDto): Promise<UserResponse>;
    deactivateUser(id: string): Promise<DeactivateResponse>;
    getUser(id: string): Promise<GetUserDto>;
    getUsers(req: GetUserPaginationDto): Promise<GetUserResponse>;
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
    getUsersByOrg(
        orgId: string,
        header: HeaderType
    ): Promise<OrgUserResponseDto>;
    addUser(
        orgId: string,
        req: AddUserDto,
        header: HeaderType
    ): Promise<AddUserResponse>;
    deleteUser(
        orgId: string,
        userId: string,
        header: HeaderType
    ): Promise<ActivateUserResponse>;
}

export interface IUserMapper {
    connectedUsersDtoToDto(res: ConnectedUserResponse): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    residentDtoToDto(res: GetUserResponseDto): ResidentResponse;
    dtoToUsersDto(req: OrgUserResponse): OrgUserResponseDto;
}

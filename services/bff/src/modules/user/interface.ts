import {
    HeaderType,
    PaginationDto,
    MetricServiceRes,
} from "../../common/types";
import {
    ActivateUserDto,
    ActivateUserResponse,
    AddUserDto,
    UserResDto,
    AddUserServiceRes,
    ConnectedUserDto,
    DeactivateResponse,
    GetUserDto,
    GetUserPaginationDto,
    GetUserResponse,
    GetUserResponseDto,
    GetUsersDto,
    OrgUserDto,
    OrgUserResponse,
    OrgUsersResponse,
    ResidentResponse,
    ResidentsResponse,
    UpdateUserDto,
    UserInput,
} from "./types";

export interface IUserService {
    getConnectedUsers(
        orgId: string,
        header: HeaderType
    ): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    updateUser(
        orgId: string,
        req: UpdateUserDto,
        header: HeaderType
    ): Promise<UpdateUserDto>;
    deactivateUser(id: string): Promise<DeactivateResponse>;
    getUser(data: UserInput, header: HeaderType): Promise<GetUserDto>;
    getUsers(req: GetUserPaginationDto): Promise<GetUserResponse>;
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
    getUsersByOrg(orgId: string, header: HeaderType): Promise<GetUsersDto[]>;
    addUser(
        orgId: string,
        req: AddUserDto,
        header: HeaderType
    ): Promise<UserResDto | null>;
    deleteUser(
        orgId: string,
        userId: string,
        header: HeaderType
    ): Promise<ActivateUserResponse>;
}

export interface IUserMapper {
    dtoToAddUserDto(res: AddUserServiceRes): UserResDto | null;
    connectedUsersDtoToDto(res: MetricServiceRes[]): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    residentDtoToDto(res: GetUserResponseDto): ResidentResponse;
    dtoToUsersDto(req: OrgUsersResponse): GetUsersDto[];
    dtoToUserDto(req: OrgUserResponse): GetUserDto;
    dtoToUserResDto(req: OrgUserDto): UserResDto;
}

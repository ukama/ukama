import {
    BoolResponse,
    MetricServiceValueRes,
    ParsedCookie,
} from "../../common/types";
import {
    ConnectedUserDto,
    DeactivateResponse,
    ESimQRCodeRes,
    GetESimQRCodeInput,
    GetUserDto,
    GetUserResponseDto,
    OrgUserSimDto,
    SubscriberAPIResDto,
    SubscriberDto,
    UpdateUserInputDto,
    UpdateUserServiceInput,
    UserAPIResDto,
    UserInputDto,
    UserResDto,
} from "./types";

export interface IUserService {
    getConnectedUsers(cookie: ParsedCookie): Promise<ConnectedUserDto>;
    updateUser(
        userId: string,
        req: UpdateUserInputDto,
        cookie: ParsedCookie,
    ): Promise<UserResDto>;
    deactivateUser(
        uuid: string,
        cookie: ParsedCookie,
    ): Promise<DeactivateResponse>;
    getUser(userId: string, cookie: ParsedCookie): Promise<UserResDto>;
    addUser(req: UserInputDto, cookie: ParsedCookie): Promise<UserResDto>;
    deleteUser(userId: string, cookie: ParsedCookie): Promise<BoolResponse>;
    getEsimQRCode(
        data: GetESimQRCodeInput,
        cookie: ParsedCookie,
    ): Promise<ESimQRCodeRes>;
    updateUserRoaming(
        data: UpdateUserServiceInput,
        cookie: ParsedCookie,
    ): Promise<OrgUserSimDto>;
}

export interface IUserMapper {
    connectedUsersDtoToDto(res: MetricServiceValueRes[]): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    dtoToUserResDto(res: UserAPIResDto): UserResDto;
    dtoToSubscriberResDto(res: SubscriberAPIResDto): SubscriberDto;
}

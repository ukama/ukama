import { NonEmptyArray } from "type-graphql";

import { AddUserResolver } from "./addUser";
import { DeactivateUserResolver } from "./deactivateUser";
import { GetAccountDetailsResolver } from "./getAccountDetails";
import { GetConnectedUsersResolver } from "./getConnectedUsers";
import { GetEsimQRResolver } from "./getEsimQR";
import { GetUserResolver } from "./getUser";
import { GetUsersDataUsageResolver } from "./getUsersDataUsage";
import { updateFirstVisitResolver } from "./updateFirstVisit";
import { UpdateUserResolver } from "./updateUser";
import { UpdateUserRoamingResolver } from "./updateUserRoaming";
import { UpdateUserStatusResolver } from "./updateUserStatus";
import { WhoamiResolver } from "./whoami";


const resolvers: NonEmptyArray<Function> = [AddUserResolver,
    DeactivateUserResolver,
    GetAccountDetailsResolver,
    GetConnectedUsersResolver,
    GetEsimQRResolver,
    GetUserResolver,
    GetUsersDataUsageResolver,
    updateFirstVisitResolver,
    UpdateUserResolver,
    UpdateUserRoamingResolver,
    UpdateUserStatusResolver,
    WhoamiResolver];

export default resolvers;

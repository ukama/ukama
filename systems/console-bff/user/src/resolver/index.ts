import { NonEmptyArray } from "type-graphql";

import { AddUserResolver } from "./addUser.resolver";
import { DeactivateUserResolver } from "./deactivateUser.resolver";
import { GetAccountDetailsResolver } from "./getAccountDetails.resolver";
import { GetConnectedUsersResolver } from "./getConnectedUsers.resolver";
import { GetEsimQRResolver } from "./getEsimQR.resolver";
import { GetUserResolver } from "./getUser.resolver";
import { GetUsersDataUsageResolver } from "./getUsersDataUsage.resolver";
import { updateFirstVisitResolver } from "./updateFirstVisit.resolver";
import { UpdateUserResolver } from "./updateUser.resolver";
import { UpdateUserRoamingResolver } from "./updateUserRoaming.resolver";
import { UpdateUserStatusResolver } from "./updateUserStatus.resolver";
import { WhoamiResolver } from "./whoami.resolver";


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

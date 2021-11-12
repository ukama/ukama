import { Request, Response } from "express";
import { CONNECTED_USER_TYPE, TIME_FILTER } from "../constants";
import { UserDto } from "../modules/user/types";
import casual from "./mockData/casual-extensions";

export const getUser = (req: Request, res: Response): void => {
    let users: UserDto[];
    const filter = req.query[0]?.toString();

    switch (filter) {
        case TIME_FILTER.TODAY:
            users = casual.randomArray<UserDto>(2, 10, casual._user);
            break;
        case TIME_FILTER.WEEK:
            users = casual.randomArray<UserDto>(4, 40, casual._user);
            break;
        case TIME_FILTER.MONTH:
            users = casual.randomArray<UserDto>(14, 100, casual._user);
            break;
        case TIME_FILTER.TOTAL:
            users = casual.randomArray<UserDto>(28, 400, casual._user);
            break;
        default:
            users = [];
            break;
    }
    let residentUsers = 0;
    const totalUser = users.length;
    for (let i = 0; i < totalUser; i++) {
        if (users[i].type === CONNECTED_USER_TYPE.RESIDENTS) {
            residentUsers++;
        }
    }

    res.send({
        status: "success",
        data: {
            totalUser: totalUser,
            residentUsers: residentUsers,
            guestUsers: totalUser - residentUsers,
        },
    });
};

export const getDataUsage = (req: Request, res: Response): void => {
    let data;
    const filter = req.query[0]?.toString();
    if (
        filter !== TIME_FILTER.TODAY ||
        TIME_FILTER.WEEK ||
        TIME_FILTER.MONTH ||
        TIME_FILTER.TOTAL
    )
        data = {};

    data = casual._dataUsage();

    res.send({
        status: "success",
        data,
    });
};

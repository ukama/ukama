import { Request, Response } from "express";
import { CONNECTED_USER_TYPE, TIME_FILTER } from "../constants";
import { HTTP404Error, Messages } from "../errors";
import { DataUsageDto } from "../modules/data/types";
import { UserDto } from "../modules/user/types";
import casual from "./mockData/casual-extensions";

export const getUser = (req: Request, res: Response): void => {
    let users: UserDto[];
    let residentUsers = 0;
    let totalUser = 0;

    const filter = req.query[0]?.toString();

    switch (filter) {
        case TIME_FILTER.TODAY:
            users = casual.randomArray<UserDto>(2, 10, casual._user);
            totalUser = users.length;
            for (let i = 0; i < totalUser; i++) {
                if (users[i].type === CONNECTED_USER_TYPE.RESIDENTS) {
                    residentUsers++;
                }
            }
            break;
        case TIME_FILTER.WEEK:
            users = casual.randomArray<UserDto>(4, 40, casual._user);
            totalUser = users.length;

            for (let i = 0; i < totalUser; i++) {
                if (users[i].type === CONNECTED_USER_TYPE.RESIDENTS) {
                    residentUsers++;
                }
            }
            break;
        case TIME_FILTER.MONTH:
            users = casual.randomArray<UserDto>(14, 100, casual._user);
            totalUser = users.length;
            for (let i = 0; i < totalUser; i++) {
                if (users[i].type === CONNECTED_USER_TYPE.RESIDENTS) {
                    residentUsers++;
                }
            }
            break;
        case TIME_FILTER.TOTAL:
            users = casual.randomArray<UserDto>(28, 400, casual._user);
            totalUser = users.length;
            for (let i = 0; i < totalUser; i++) {
                if (users[i].type === CONNECTED_USER_TYPE.RESIDENTS) {
                    residentUsers++;
                }
            }
            break;
        default:
            throw new HTTP404Error(Messages.INVALID_CONNECTED_USER_FILTER);
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
    let data: DataUsageDto;
    const filter = req.query[0]?.toString();
    switch (filter) {
        case TIME_FILTER.TODAY:
            data = casual._dataUsage();
            break;
        case TIME_FILTER.WEEK:
            data = casual._dataUsage();
            break;
        case TIME_FILTER.MONTH:
            data = casual._dataUsage();
            break;
        case TIME_FILTER.TOTAL:
            data = casual._dataUsage();
            break;
        default:
            throw new HTTP404Error(Messages.INVALID_DATA_USAGE_FILTER);
    }
    res.send({
        status: "success",
        data,
    });
};

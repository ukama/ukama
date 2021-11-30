import { Request, Response } from "express";
import {
    CONNECTED_USER_TYPE,
    DATA_BILL_FILTER,
    GET_USER_TYPE,
    NETWORK_TYPE,
    TIME_FILTER,
} from "../constants";
import { AlertDto } from "../modules/alert/types";
import { BillHistoryDto, CurrentBillDto } from "../modules/billing/types";
import { EsimDto } from "../modules/esim/types";
import { NodeDto } from "../modules/node/types";
import { GetUserDto, UserDto } from "../modules/user/types";
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

export const getDataBill = (req: Request, res: Response): void => {
    let data;
    const filter = req.query[0]?.toString();
    if (
        filter !== DATA_BILL_FILTER.CURRENT ||
        DATA_BILL_FILTER.JANUARY ||
        DATA_BILL_FILTER.FEBRURAY ||
        DATA_BILL_FILTER.MARCH ||
        DATA_BILL_FILTER.APRIL ||
        DATA_BILL_FILTER.MAY ||
        DATA_BILL_FILTER.JUNE ||
        DATA_BILL_FILTER.JULY ||
        DATA_BILL_FILTER.AUGUST ||
        DATA_BILL_FILTER.SEPTEMBER ||
        DATA_BILL_FILTER.OCTOBER ||
        DATA_BILL_FILTER.NOVERMBER ||
        DATA_BILL_FILTER.DECEMBER
    )
        data = {};

    data = casual._dataBill();

    res.send({
        status: "success",
        data,
    });
};

export const getAlerts = (req: Request, res: Response): void => {
    const data = casual.randomArray<AlertDto>(1, 10, casual._alert);

    const pageNo = Number(req.query.pageNo);
    const pageSize = Number(req.query.pageSize);

    let alerts = [];
    if (!pageNo) alerts = data;
    else {
        const index = (pageNo - 1) * pageSize;
        for (let i = index; i < index + pageSize; i++) {
            if (data[i]) alerts.push(data[i]);
        }
    }
    res.send({
        status: "success",
        data: alerts,
        length: data.length,
    });
};

export const getNodes = (req: Request, res: Response): void => {
    const data = casual.randomArray<NodeDto>(1, 10, casual._node);

    const pageNo = Number(req.query.pageNo);
    const pageSize = Number(req.query.pageSize);

    let nodes = [];
    if (!pageNo) nodes = data;
    else {
        const index = (pageNo - 1) * pageSize;
        for (let i = index; i < index + pageSize; i++) {
            if (data[i]) nodes.push(data[i]);
        }
    }
    res.send({
        status: "success",
        data: nodes,
        length: data.length,
    });
};

export const getEsims = (req: Request, res: Response): void => {
    const data = casual.randomArray<EsimDto>(3, 10, casual._esim);
    res.send({
        status: "success",
        data: data,
    });
};

export const activateUser = (req: Request, res: Response): void => {
    const { body } = req;
    const data = {
        success: false,
    };

    if (body.firstName && body.lastName && body.eSimNumber) data.success = true;

    res.send({
        status: "success",
        data: data,
    });
};

export const updateUser = (req: Request, res: Response): void => {
    const { body } = req;
    let data;

    if (
        !(
            body.id &&
            (body.firstName ||
                body.lastName ||
                body.eSimNumber ||
                body.email ||
                body.phone)
        )
    )
        data = {};

    data = casual._updateUser(
        body.id,
        body.firstName,
        body.lastName,
        body.eSimNumber,
        body.email,
        body.phone
    );

    res.send({
        status: "success",
        data: data,
    });
};
export const deactivateUser = (req: Request, res: Response): void => {
    const { body } = req;

    let data;

    if (!body.id) data = {};

    data = casual._deleteRes(body.id.toString());

    res.send({
        status: "success",
        data: data,
    });
};

export const getUsers = (req: Request, res: Response): void => {
    let data;
    const filter = req.query.type?.toString();

    if (
        filter !== GET_USER_TYPE.ALL ||
        GET_USER_TYPE.GUEST ||
        GET_USER_TYPE.HOME ||
        GET_USER_TYPE.RESIDENT ||
        GET_USER_TYPE.VISITOR
    )
        data = {};
    data = casual.randomArray<GetUserDto>(3, 30, casual._getUser);

    const pageNo = Number(req.query.pageNo);
    const pageSize = Number(req.query.pageSize);

    let users = [];
    if (!pageNo) users = data;
    else {
        const index = (pageNo - 1) * pageSize;
        for (let i = index; i < index + pageSize; i++) {
            if (data[i]) users.push(data[i]);
        }
    }
    res.send({
        status: "success",
        data: users,
        length: data.length,
    });
};

export const addNode = (req: Request, res: Response): void => {
    const { body } = req;
    const data = {
        success: false,
    };
    if (body.name && body.serialNo && body.securityCode) data.success = true;
    res.send({
        status: "success",
        data: data,
    });
};
export const updateNode = (req: Request, res: Response): void => {
    const { body } = req;
    let data;
    if (!(body.id && (body.name || body.serialNo || body.securityCode)))
        data = {};
    data = casual._updateNode(body.id, body.name, body.serialNo);

    res.send({
        status: "success",
        data: data,
    });
};
export const deleteNode = (req: Request, res: Response): void => {
    const { body } = req;

    let data;

    if (!body.id) data = {};

    data = casual._deleteRes(body.id.toString());

    res.send({
        status: "success",
        data: data,
    });
};

export const getCurrentBill = (req: Request, res: Response): void => {
    const data = casual.randomArray<CurrentBillDto>(1, 5, casual._currentBill);
    res.send({
        status: "success",
        data: data,
    });
};

export const getBillHistory = (req: Request, res: Response): void => {
    const data = casual.randomArray<BillHistoryDto>(1, 5, casual._billHistory);
    res.send({
        status: "success",
        data: data,
    });
};

export const getNetwork = (req: Request, res: Response): void => {
    let data;
    const filter = req.query[0]?.toString();
    if (filter !== NETWORK_TYPE.PUBLIC || NETWORK_TYPE.PRIVATE) data = {};

    data = casual._network();

    res.send({
        status: "success",
        data,
    });
};

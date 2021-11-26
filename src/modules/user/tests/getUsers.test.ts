import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_USER_QUERY } from "../../../common/graphql";
import { GET_USER_TYPE, HEADER } from "../../../constants";
import { GetUserPaginationDto } from "../types";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "6ca56229-7726-46f9-939b-40fbf8799206",
            status: "ACTIVE",
            name: "Mrs. Baron Harber",
            node: "Intermediate Data Plan",
            dataPlan: "NA",
            dataUsage: 34,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
        {
            id: "e2ac0b99-38fe-457d-92c1-9b9ba38156ac",
            status: "INACTIVE",
            name: "Miss Woodrow Olson",
            node: "Default Data Plan",
            dataPlan: "NA",
            dataUsage: 57,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
        {
            id: "840613cb-57a4-4dc6-931f-ca50cab38b37",
            status: "INACTIVE",
            name: "Dr. Wilhelmine Bailey",
            node: "Intermediate Data Plan",
            dataPlan: "PAID",
            dataUsage: 54,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
        {
            id: "d8334e30-c6a6-4ce4-a2c7-7b3be1c6cd0d",
            status: "INACTIVE",
            name: "Mr. Veda Schimmel",
            node: "Default Data Plan",
            dataPlan: "PAID",
            dataUsage: 52,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
        {
            id: "b6a23033-20a1-457b-acda-ff79399dc4ce",
            status: "ACTIVE",
            name: "Mr. Francisco Renner",
            node: "Intermediate Data Plan",
            dataPlan: "PAID",
            dataUsage: 19,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
    ],
    length: 6,
};

const meta: GetUserPaginationDto = {
    type: GET_USER_TYPE.ALL,
    pageNo: 1,
    pageSize: 5,
};

describe("Get Users", () => {
    beforeEachGetCall(
        "/user/get_users?type=ALL&pageNo=1&pageSize=5",
        nockResponse,
        200
    );
    it("get users", async () => {
        const response = await gCall({
            source: GET_USER_QUERY,
            variableValues: {
                input: meta,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                getUsers: {
                    users: nockResponse.data,
                    meta: {
                        page: 1,
                        count: 6,
                        pages: 2,
                        size: 5,
                    },
                },
            },
        });
    });
});

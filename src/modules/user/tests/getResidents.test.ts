import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_RESIDENTS_QUERY } from "../../../common/graphql";
import { PaginationDto } from "../../../common/types";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "528690b6-bbf1-4d01-b7a3-e4f7a3d8106f",
            status: "ACTIVE",
            name: "Dr. Margarita Fisher",
            node: "Intermediate Data Plan",
            dataPlan: "PAID",
            dataUsage: 189,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
        {
            id: "45cbf15f-fb62-4a14-932d-e7298657970f",
            status: "INACTIVE",
            name: "Mrs. Araceli Schowalter",
            node: "Intermediate Data Plan",
            dataPlan: "PAID",
            dataUsage: 146,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
        {
            id: "d7e9b719-beaf-4ee8-ad9a-2569443af336",
            status: "INACTIVE",
            name: "Mrs. Vena Konopelski",
            node: "Intermediate Data Plan",
            dataPlan: "NA",
            dataUsage: 164,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
        {
            id: "c6504bc0-4a2b-4abd-8e7a-917c4c721d17",
            status: "INACTIVE",
            name: "Mr. Anabelle Mayer",
            node: "Intermediate Data Plan",
            dataPlan: "UNPAID",
            dataUsage: 1,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
        {
            id: "916da5a0-c9e7-4b52-a174-21a2669afbbc",
            status: "INACTIVE",
            name: "Mrs. Emilia Goodwin",
            node: "Intermediate Data Plan",
            dataPlan: "PAID",
            dataUsage: 190,
            dlActivity: "Table cell",
            ulActivity: "Table cell",
        },
    ],
    length: 20,
};
const res = [
    {
        id: "528690b6-bbf1-4d01-b7a3-e4f7a3d8106f",
        name: "Dr. Margarita Fisher",
        dataUsage: 189,
    },
    {
        id: "45cbf15f-fb62-4a14-932d-e7298657970f",
        name: "Mrs. Araceli Schowalter",
        dataUsage: 146,
    },
    {
        id: "d7e9b719-beaf-4ee8-ad9a-2569443af336",
        name: "Mrs. Vena Konopelski",
        dataUsage: 164,
    },
    {
        id: "c6504bc0-4a2b-4abd-8e7a-917c4c721d17",
        name: "Mr. Anabelle Mayer",
        dataUsage: 1,
    },
    {
        id: "916da5a0-c9e7-4b52-a174-21a2669afbbc",
        name: "Mrs. Emilia Goodwin",
        dataUsage: 190,
    },
];

const meta: PaginationDto = {
    pageNo: 1,
    pageSize: 5,
};

describe("Get Residents", () => {
    beforeEachGetCall("/user/get_users?pageNo=1&pageSize=5", nockResponse, 200);
    it("get residents", async () => {
        const response = await gCall({
            source: GET_RESIDENTS_QUERY,
            variableValues: {
                input: meta,
            },
            contextValue: {
                req: {
                    headers: {
                        authorisation: "test",
                    },
                },
            },
        });
        expect(response).toMatchObject({
            data: {
                getResidents: {
                    residents: {
                        residents: res,
                        activeResidents: 1,
                        totalResidents: 20,
                    },
                    meta: {
                        page: 1,
                        count: 20,
                        pages: 4,
                        size: 5,
                    },
                },
            },
        });
    });
});

import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_BILL_HISTORY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "18a02e2f-e95b-4bae-ba04-ea4812252cb4",
            date: "12-08-2021",
            description: "Bill for month",
            totalUsage: 1,
            subtotal: 3,
        },
        {
            id: "143b5c04-c84b-4d04-b81d-37ef6fa5f6d9",
            date: "01-23-2021",
            description: "Bill for month",
            totalUsage: 2,
            subtotal: 6,
        },
        {
            id: "fd4094c2-3ad0-4883-b02d-e026850f1d15",
            date: "07-05-2021",
            description: "Bill for month",
            totalUsage: 1,
            subtotal: 3,
        },
        {
            id: "29b8e3db-b2ed-4951-968b-50fa5ad8a2f2",
            date: "07-28-2021",
            description: "Bill for month",
            totalUsage: 6,
            subtotal: 18,
        },
    ],
};

describe("Get Bill History", () => {
    beforeEachGetCall("/bill/get_bill_history", nockResponse, 200);
    it("get bill history", async () => {
        const response = await gCall({
            source: GET_BILL_HISTORY,
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                getBillHistory: nockResponse.data,
            },
        });
    });
});

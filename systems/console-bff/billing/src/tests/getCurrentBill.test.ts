import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_CURRENT_BILL } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "64e8967c-3877-4f13-a324-02600107c3b3",
            name: "Mrs. Delia Hansen",
            dataUsed: 5,
            rate: 4,
            subtotal: 20,
        },
        {
            id: "0f0e6715-b8fc-4c71-a0ea-5b5dd2a7632f",
            name: "Miss Louisa Ritchie",
            dataUsed: 1,
            rate: 4,
            subtotal: 4,
        },
        {
            id: "d788f704-e311-424b-a187-a430e5867c0d",
            name: "Ms. Max Lueilwitz",
            dataUsed: 9,
            rate: 5,
            subtotal: 45,
        },
    ],
};

describe("Get CurrentBill", () => {
    beforeEachGetCall("/bill/get_current_bill", nockResponse, 200);
    it("get current bill", async () => {
        const response = await gCall({
            source: GET_CURRENT_BILL,
            contextValue: {
                req: HEADER,
            },
        });

        expect(response).toMatchObject({
            data: {
                getCurrentBill: {
                    bill: nockResponse.data,
                    total: 69,
                    dueDate: "10-10-2021",
                    billMonth: "11-10-2021",
                },
            },
        });
    });
});

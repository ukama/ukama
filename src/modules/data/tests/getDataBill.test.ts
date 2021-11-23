import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_DATA_BILL_QUERY } from "../../../common/graphql";
import { DATA_BILL_FILTER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "27f90c88-7c1c-41c4-8186-1741a8f87420",
        dataBill: 50,
        billDue: 10,
    },
};

describe("Get Data Bill", () => {
    beforeEachGetCall("/data/data_bill?0=CURRENT", nockResponse, 200);
    it("Get Data Bill", async () => {
        const response = await gCall({
            source: GET_DATA_BILL_QUERY,
            variableValues: {
                data: DATA_BILL_FILTER.CURRENT,
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
                getDataBill: nockResponse.data,
            },
        });
    });
});

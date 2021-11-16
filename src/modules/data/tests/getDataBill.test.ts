import "reflect-metadata";
import { gCall, beforeEachCall } from "../../../common/utils";
import { GET_DATA_BILL_QUERY } from "../../../common/graphql";
import { DATA_BILL_FILTER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "f0dc1856-e0c0-439b-b1ba-6bee696fb247",
        dataBill: "25$",
        billDue: "25 days",
    },
};

describe("Get Data Bill", () => {
    beforeEachCall("/data/data_bill?0=CURRENT", nockResponse, 200);
    it("Get Data Bill", async () => {
        const response = await gCall({
            source: GET_DATA_BILL_QUERY,
            variableValues: {
                data: DATA_BILL_FILTER.CURRENT,
            },
        });

        expect(response).toMatchObject({
            data: {
                getDataBill: nockResponse.data,
            },
        });
    });
});

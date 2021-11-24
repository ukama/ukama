import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_DATA_USAGE_QUERY } from "../../../common/graphql";
import { HEADER, TIME_FILTER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "7c51cdc5-fc3b-4837-bc7d-94bb16e11d8d",
        dataConsumed: 941,
        dataPackage: "Unlimited",
    },
};

describe("Get Data Usage", () => {
    beforeEachGetCall("/data/data_usage?0=MONTH", nockResponse, 200);
    it("Get Data Usage", async () => {
        const response = await gCall({
            source: GET_DATA_USAGE_QUERY,
            variableValues: {
                data: TIME_FILTER.MONTH,
            },
            contextValue: {
                req: HEADER,
            },
        });

        expect(response).toMatchObject({
            data: {
                getDataUsage: nockResponse.data,
            },
        });
    });
});

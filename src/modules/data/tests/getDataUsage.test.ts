import "reflect-metadata";
import { gCall, beforeEachCall } from "../../../common/utils";
import { GET_DATA_USAGE_QUERY } from "../../../common/graphql";
import { TIME_FILTER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "f6e025fc-b490-462c-acdd-689c4a9d2cb8",
        dataConsumed: "126GBs",
        dataPackage: "Unlimited",
    },
};

describe("Get Data Usage", () => {
    beforeEachCall("/data/data_usage?0=MONTH", nockResponse, 200);
    it("Get Data Usage", async () => {
        const response = await gCall({
            source: GET_DATA_USAGE_QUERY,
            variableValues: {
                data: TIME_FILTER.MONTH,
            },
            contextValue: {
                req: {
                    headers: {
                        authorization: "test",
                    },
                },
            },
        });

        expect(response).toMatchObject({
            data: {
                getDataUsage: nockResponse.data,
            },
        });
    });
});

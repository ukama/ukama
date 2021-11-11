import "reflect-metadata";
import { gCall } from "../../../test/utils";
import { GET_DATA_USAGE_QUERY } from "../../../test/graphql";
import nock from "nock";
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
    beforeEach(() => {
        nock("http://127.0.0.1:3000")
            .get("/data/data_usage?0=MONTH")
            .reply(200, nockResponse);
    });
    it("Get Data Usage", async () => {
        const response = await gCall({
            source: GET_DATA_USAGE_QUERY,
            variableValues: {
                data: TIME_FILTER.MONTH,
            },
        });

        expect(response).toMatchObject({
            data: {
                getDataUsage: {
                    id: expect.any(String),
                    dataConsumed: expect.any(String),
                    dataPackage: expect.any(String),
                },
            },
        });
    });
});

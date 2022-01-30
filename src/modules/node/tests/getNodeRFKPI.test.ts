import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_NODE_RF_KPI_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";
import { PaginationDto } from "../../../common/types";

const nockResponse = {
    status: "success",
    data: [
        {
            qam: 8,
            rfOutput: 11,
            rssi: 3,
            timestamp: 1643632371925,
        },
        {
            qam: 7,
            rfOutput: 2,
            rssi: 7,
            timestamp: 1643632371925,
        },
        {
            qam: 11,
            rfOutput: 19,
            rssi: 7,
            timestamp: 1643632371925,
        },
    ],
};

const meta: PaginationDto = {
    pageNo: 1,
    pageSize: 3,
};

describe("Get Node RF KPIs ", () => {
    beforeEachGetCall("/node/rf_kpis", nockResponse, 200);
    it("get node RF KPIs", async () => {
        const response = await gCall({
            source: GET_NODE_RF_KPI_QUERY,
            variableValues: {
                input: meta,
            },
            contextValue: {
                req: HEADER,
            },
        });

        expect(response).toMatchObject({
            data: {
                getNodeRFKPI: nockResponse.data,
            },
        });
    });
});

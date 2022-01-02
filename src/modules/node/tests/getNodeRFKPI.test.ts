import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_NODE_RF_KPI_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: { qam: 8, rfOutput: 12, rssi: 9 },
};

describe("Get Node RF KPIs ", () => {
    beforeEachGetCall("/node/rf_kpis", nockResponse, 200);
    it("get node RF KPIs", async () => {
        const response = await gCall({
            source: GET_NODE_RF_KPI_QUERY,
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

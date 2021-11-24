import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_NETWORK_QUERY } from "../../../common/graphql";
import { NETWORK_STATUS, NETWORK_TYPE } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "22aa60e1-0726-488f-81e4-20b47cadcf65",
        status: NETWORK_STATUS.ONLINE,
        description: "21 days 5 hours 1 minute",
    },
};

describe("Get Network", () => {
    beforeEachGetCall("/network/get_network?0=PUBLIC", nockResponse, 200);
    it("Get network", async () => {
        const response = await gCall({
            source: GET_NETWORK_QUERY,
            variableValues: {
                data: NETWORK_TYPE.PUBLIC,
            },
            contextValue: {
                req: {
                    headers: {
                        csrf_token: "test",
                        kratos_session: "test",
                    },
                },
            },
        });

        expect(response).toMatchObject({
            data: {
                getNetwork: nockResponse.data,
            },
        });
    });
});

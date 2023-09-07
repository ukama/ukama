import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_ESIM_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: [
        { esim: "# 36501-14-11-1971-2807124", active: false },
        { esim: "# 93799-26-01-1975-6275029", active: false },
        { esim: "# 72040-20-03-2009-9716913", active: false },
        { esim: "# 37284-05-02-2009-9628947", active: false },
    ],
};

describe("Get Users", () => {
    beforeEachGetCall("/esims/get_esims", nockResponse, 200);
    it("get users", async () => {
        const response = await gCall({
            source: GET_ESIM_QUERY,
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                getEsims: nockResponse.data,
            },
        });
    });
});

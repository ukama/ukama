import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_USER_BY_ID_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "awesrdf1234",
        status: "ACTIVE",
        name: "Ms. Herta Balistreri",
        node: "Default Data Plan",
        dataPlan: "NA",
        dataUsage: 161,
        dlActivity: "Table cell",
        ulActivity: "Table cell",
    },
};

const reqParam = {
    id: "awesrdf1234",
};

describe("Get User By Id", () => {
    beforeEachGetCall("/user/get_user?id=awesrdf1234", nockResponse, 200);
    it("get user by id", async () => {
        const response = await gCall({
            source: GET_USER_BY_ID_QUERY,
            variableValues: {
                input: reqParam.id,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                getUser: nockResponse.data,
            },
        });
    });
});

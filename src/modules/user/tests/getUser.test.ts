import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_USER_BY_ID_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "srdtfghj",
        status: "ACTIVE",
        name: "Mrs. Amelie Klein",
        eSimNumber: "# 80577-31-01-1977-7357651",
        iccid: "539682432387695",
        email: "Jarod_Kiehn@hotmail.com",
        phone: "622-058-1593",
        roaming: false,
        dataPlan: 6,
        dataUsage: 4,
    },
};

const reqParam = {
    id: "srdtfghj",
};

describe("Get User By Id", () => {
    beforeEachGetCall("/user/get_user?id=srdtfghj", nockResponse, 200);
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

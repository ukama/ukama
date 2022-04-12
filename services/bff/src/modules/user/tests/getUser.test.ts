import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_USER_BY_ID_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "acasc",
        status: "ACTIVE",
        name: "Mrs. Amelie Klein",
        eSimNumber: "acasc",
        iccid: "539682432387695",
        email: "Jarod_Kiehn@hotmail.com",
        phone: "622-058-1593",
        roaming: "ACTIVE",
        dataPlan: 6,
        dataUsage: 4,
    },
};

const userInput = {
    orgId: "sadasdas",
    userId: "zxczc",
};

describe("Get User By Id", () => {
    beforeEachGetCall("/user/get_user", nockResponse, 200);
    it("get user by id", async () => {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const response = await gCall({
            source: GET_USER_BY_ID_QUERY,
            variableValues: {
                userInput: userInput,
            },
            contextValue: {
                req: HEADER,
            },
        });
        // console.log("Response: ", response);
        // expect(response).toMatchObject({
        //     data: {
        //         getUser: nockResponse.data,
        //     },
        // });
    });
});

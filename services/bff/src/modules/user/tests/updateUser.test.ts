import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { POST_UPDATE_USER_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "1234567890vb",
        name: "Carmel Yundt",
        sim: "# 53648-05-06-2023-8772112",
        email: "Sandy_Bins@Kris.net",
        phone: "740-635-9371",
    },
};
const reqBody = {
    id: "1234567890vb",
    eSimNumber: "# 53648-05-06-2023-8772112",
};

describe("Post Update Users", () => {
    beforeEachPostCall("/user/update_user", reqBody, nockResponse, 200);
    it("post update users", async () => {
        const response = await gCall({
            source: POST_UPDATE_USER_MUTATION,
            variableValues: {
                input: reqBody,
            },
            contextValue: {
                req: HEADER,
            },
        });
        // expect(response).toMatchObject({
        //     data: {
        //         updateUser: nockResponse.data,
        //     },
        // });
    });
});

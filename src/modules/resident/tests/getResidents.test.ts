import "reflect-metadata";
import { gCall, beforeEachCall } from "../../../common/utils";
import { GET_RESIDENTS_QUERY } from "../../../common/graphql";
import { PaginationDto } from "../../../common/types";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "653a9188-7101-43e7-8e7e-39cdc7bbde3d",
            name: "Miss Freddy Barton",
            usage: "255GB",
        },
        {
            id: "6ef62e64-0903-4e93-a308-4930cf08715f",
            name: "Mr. Faustino Murphy",
            usage: "650GB",
        },
    ],
    length: 4,
};

const meta: PaginationDto = {
    pageNo: 1,
    pageSize: 2,
};

describe("Get Residents", () => {
    beforeEachCall(
        "/resident/get_residents?pageNo=1&pageSize=2",
        nockResponse,
        200
    );
    it("get residents", async () => {
        const response = await gCall({
            source: GET_RESIDENTS_QUERY,
            variableValues: {
                input: meta,
            },
        });

        expect(response).toMatchObject({
            data: {
                getResidents: {
                    residents: nockResponse.data,
                    meta: {
                        page: 1,
                        count: 4,
                        pages: 2,
                        size: 2,
                    },
                },
            },
        });
    });
});

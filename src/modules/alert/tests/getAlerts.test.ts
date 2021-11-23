import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_ALERTS_QUERY } from "../../../common/graphql";
import { PaginationDto } from "../../../common/types";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "41ce94c1-1c60-4874-9447-dbe70b1ac689",
            type: "WARNING",
            title: "Magnam hic veniam",
            description:
                "Voluptatibus perspiciatis odit at autem est quam. Culpa omnis iste aliquam voluptatum id fuga et.",
            alertDate: "2021-11-15T09:45:50.334Z",
        },
    ],
    length: 4,
};

const meta: PaginationDto = {
    pageNo: 2,
    pageSize: 3,
};

describe("Get Alerts", () => {
    beforeEachGetCall(
        "/alert/get_alerts?pageNo=2&pageSize=3",
        nockResponse,
        200
    );
    it("get alerts", async () => {
        const response = await gCall({
            source: GET_ALERTS_QUERY,
            variableValues: {
                input: meta,
            },
            contextValue: {
                req: {
                    headers: {
                        authorization: "test",
                    },
                },
            },
        });

        expect(response).toMatchObject({
            data: {
                getAlerts: {
                    alerts: nockResponse.data,
                    meta: {
                        page: 2,
                        count: 4,
                        pages: 2,
                        size: 3,
                    },
                },
            },
        });
    });
});

describe("Get Alerts (Fail)", () => {
    beforeEachGetCall("/alert/get?pageNo=2&pageSize=3", nockResponse, 404);
    it("get alerts(fail)", async () => {
        const response = await gCall({
            source: GET_ALERTS_QUERY,
            variableValues: {
                input: meta,
            },
        });
        expect(response).toEqual(expect.objectContaining({ data: null }));
    });
});

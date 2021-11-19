import "reflect-metadata";
import { gCall, beforeEachCall } from "../../../common/utils";
import { GET_NODES_QUERY } from "../../../common/graphql";
import { PaginationDto } from "../../../common/types";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "7d9ae8e7-4372-4fdc-8645-40447a96ccaa",
            title: "A est",
            description:
                "Nostrum qui hic aut rerum aperiam maiores laborum ut. Architecto dolorem quaerat ut. Qui omnis est.",
            totalUser: 33,
        },
        {
            id: "2654c9ef-0b58-4ad1-aeae-d7f9eb2788e2",
            title: "Non quo",
            description:
                "Laborum omnis aut quia neque est. Dicta non nemo aspernatur deserunt vero quia omnis veniam. Veniam ipsum excepturi sint debitis qui non dolorem consequuntur.",
            totalUser: 92,
        },
    ],
    length: 2,
};

const meta: PaginationDto = {
    pageNo: 1,
    pageSize: 3,
};

describe("Get Nodes", () => {
    beforeEachCall("/node/get_nodes?pageNo=1&pageSize=3", nockResponse, 200);
    it("get nodes", async () => {
        const response = await gCall({
            source: GET_NODES_QUERY,
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
                getNodes: {
                    nodes: nockResponse.data,
                    meta: {
                        page: 1,
                        count: 2,
                        pages: 1,
                        size: 3,
                    },
                },
            },
        });
    });
});

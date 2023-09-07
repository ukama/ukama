import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_NODES_QUERY } from "../../../common/graphql";
import { PaginationDto } from "../../../common/types";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "e3bceb73-2fea-436b-a1b5-b74094d5c98d",
            title: "Voluptas rerum animi",
            description: "Home node",
            status: "PENDING",
            totalUser: 46,
        },
        {
            id: "caf63f96-ae16-411f-9ae4-8a4a201f53d2",
            title: "Et dolores est",
            description: "Work node",
            status: "PENDING",
            totalUser: 30,
        },
        {
            id: "19724dc1-f2f0-4d18-bd6a-6dfd29fd8ebf",
            title: "Adipisci aut magnam",
            description: "Work node",
            status: "ONBOARDED",
            totalUser: 12,
        },
    ],
    length: 6,
};

const meta: PaginationDto = {
    pageNo: 1,
    pageSize: 3,
};

describe("Get Nodes", () => {
    beforeEachGetCall("/node/get_nodes?pageNo=1&pageSize=3", nockResponse, 200);
    it("get nodes", async () => {
        const response = await gCall({
            source: GET_NODES_QUERY,
            variableValues: {
                input: meta,
            },
            contextValue: {
                req: HEADER,
            },
        });

        // expect(response).toMatchObject({
        //     data: {
        //         getNodes: {
        //             nodes: {
        //                 nodes: nockResponse.data,
        //                 activeNodes: 1,
        //                 totalNodes: 6,
        //             },
        //             meta: {
        //                 page: 1,
        //                 count: 6,
        //                 pages: 2,
        //                 size: 3,
        //             },
        //         },
        //     },
        // });
    });
});

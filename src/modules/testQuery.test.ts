import { gCall } from "../test/utils";
import { TEST_QUERY } from "../test/graphql";

describe("Test Query", () => {
    it("Test Query", async () => {
        const response = await gCall({
            source: TEST_QUERY,
        });
        expect(response).toMatchObject({
            data: {
                testQuery: "Hello World",
            },
        });
    });
});

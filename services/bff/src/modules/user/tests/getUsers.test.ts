import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_USER_QUERY } from "../../../common/graphql";
import { GET_USER_TYPE, HEADER } from "../../../constants";
import { GetUserPaginationDto } from "../types";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "4759be09-9cb9-4c78-934d-a63e25299c92",
            status: "INACTIVE",
            name: "Mrs. Tyrel Lakin",
            eSimNumber: "# 32308-06-07-2004-7400525",
            iccid: "162815153316738",
            email: "Rhett.Bahringer@Amelia.io",
            phone: "284-556-9183",
            roaming: false,
            dataPlan: 4,
            dataUsage: 3,
        },
        {
            id: "20f342bb-86b1-4d66-99b0-d67e49ba87b8",
            status: "INACTIVE",
            name: "Ms. Lauryn Watsica",
            eSimNumber: "# 91415-07-12-1989-6505666",
            iccid: "259998240647646",
            email: "Friedrich_Erdman@gmail.com",
            phone: "374-742-8879",
            roaming: true,
            dataPlan: 8,
            dataUsage: 4,
        },
        {
            id: "7d8bf660-94c5-4994-80c3-4de0663b174c",
            status: "ACTIVE",
            name: "Miss Fae Pagac",
            eSimNumber: "# 32731-31-12-1983-1203385",
            iccid: "641806096971438",
            email: "Lisette.Pfannerstill@Laura.com",
            phone: "707-653-0441",
            roaming: true,
            dataPlan: 6,
            dataUsage: 3,
        },
        {
            id: "3534db79-ff75-4ae6-849e-c9a1dbc13573",
            status: "INACTIVE",
            name: "Miss Norene Moen",
            eSimNumber: "# 94618-19-05-2003-3387247",
            iccid: "275134607033563",
            email: "Vincenza.Leuschke@hotmail.com",
            phone: "068-828-1950",
            roaming: false,
            dataPlan: 6,
            dataUsage: 4,
        },
        {
            id: "69363f50-77b6-47a3-a93b-257a63c5d516",
            status: "ACTIVE",
            name: "Miss Cristina Gislason",
            eSimNumber: "# 38202-04-09-1994-4723481",
            iccid: "779986801124261",
            email: "Meagan_Rosenbaum@Clare.co.uk",
            phone: "575-301-8570",
            roaming: false,
            dataPlan: 3,
            dataUsage: 5,
        },
    ],
    length: 13,
};
const meta: GetUserPaginationDto = {
    type: GET_USER_TYPE.ALL,
    pageNo: 1,
    pageSize: 5,
};

describe("Get Users", () => {
    beforeEachGetCall(
        "/user/get_users?type=ALL&pageNo=1&pageSize=5",
        nockResponse,
        200
    );
    it("get users", async () => {
        const response = await gCall({
            source: GET_USER_QUERY,
            variableValues: {
                input: meta,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                getUsers: {
                    users: nockResponse.data,
                    meta: {
                        page: 1,
                        count: 13,
                        pages: 3,
                        size: 5,
                    },
                },
            },
        });
    });
});

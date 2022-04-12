import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_RESIDENTS_QUERY } from "../../../common/graphql";
import { PaginationDto } from "../../../common/types";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: [
        {
            id: "a569ca1f-7300-4b48-8c15-0624d3674990",
            status: "ACTIVE",
            name: "Ms. Devante Jakubowski",
            eSimNumber: "# 87097-07-04-1971-2800559",
            iccid: "472921491040510",
            email: "Blanda.Ole@gmail.com",
            phone: "458-726-0233",
            roaming: false,
            dataPlan: 4,
            dataUsage: 2,
        },
        {
            id: "45d5ba7d-e9b5-4dac-91c7-11388682a374",
            status: "ACTIVE",
            name: "Mr. Al McGlynn",
            eSimNumber: "# 93370-22-07-1985-6962679",
            iccid: "884557740152457",
            email: "Malinda_McGlynn@yahoo.com",
            phone: "082-544-4850",
            roaming: false,
            dataPlan: 6,
            dataUsage: 2,
        },
        {
            id: "42186d07-a760-4f03-b3c9-aaeea7fc66c7",
            status: "ACTIVE",
            name: "Mrs. Pat Kling",
            eSimNumber: "# 40880-22-05-1992-7790086",
            iccid: "746427112466782",
            email: "Bailey.Pierre@Dominique.me",
            phone: "957-748-2082",
            roaming: true,
            dataPlan: 5,
            dataUsage: 3,
        },
        {
            id: "16b36205-0141-4f1f-9178-802ed39e8278",
            status: "ACTIVE",
            name: "Mr. Doug Hoeger",
            eSimNumber: "# 26765-07-08-1970-6650271",
            iccid: "476563532176996",
            email: "Boyer.Rahul@Blanda.biz",
            phone: "153-935-8156",
            roaming: true,
            dataPlan: 7,
            dataUsage: 3,
        },
        {
            id: "26a7b115-befd-43b7-afe9-ea5ef0222bb6",
            status: "ACTIVE",
            name: "Mr. Colin Zboncak",
            eSimNumber: "# 76263-03-12-1989-2143464",
            iccid: "972025378875045",
            email: "Rosenbaum_Winona@Dixie.biz",
            phone: "533-771-2342",
            roaming: true,
            dataPlan: 5,
            dataUsage: 1,
        },
    ],
    length: 21,
};

const meta: PaginationDto = {
    pageNo: 1,
    pageSize: 5,
};

describe("Get Residents", () => {
    beforeEachGetCall("/user/get_users?pageNo=1&pageSize=5", nockResponse, 200);
    it("get residents", async () => {
        const response = await gCall({
            source: GET_RESIDENTS_QUERY,
            variableValues: {
                input: meta,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                getResidents: {
                    residents: {
                        residents: nockResponse.data,
                        activeResidents: 5,
                        totalResidents: 21,
                    },
                    meta: {
                        page: 1,
                        count: 21,
                        pages: 5,
                        size: 5,
                    },
                },
            },
        });
    });
});

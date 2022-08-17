import * as defaultCasual from "casual";
import { ALERT_TYPE, ORG_NODE_STATE, NODE_TYPE } from "../../constants";
import { LoremIpsum } from "lorem-ipsum";
import { AlertDto } from "../../modules/alert/types";
import { BillHistoryDto, CurrentBillDto } from "../../modules/billing/types";
import { DataBillDto, DataUsageDto } from "../../modules/data/types";
import { EsimDto } from "../../modules/esim/types";
import { NetworkDto } from "../../modules/network/types";
import AWS from "aws-sdk";
import fs from "fs";
const s3 = new AWS.S3({
    accessKeyId: process.env.AWS_S3_ACCESS_KEY_ID,
    secretAccessKey: process.env.AWS_S3_SECRET_ACCESS_KEY,
});
import {
    NodeAppResponse,
    NodeDto,
    NodeAppsVersionLogsResponse,
} from "../../modules/node/types";
import { GetUserDto } from "../../modules/user/types";
import { createInvoice } from "../../config/createInvoice";

function randomArray<T>(
    minLength: number,
    maxLength: number,
    elementGenerator: (index?: number, length?: number) => T
): T[] {
    const length = casual.integer(minLength, maxLength);
    const result = [];
    for (let i = 0; i < length; i++) {
        result.push(elementGenerator(i + 1, maxLength));
    }
    return result;
}
const dataUsage = (): DataUsageDto => {
    return {
        id: defaultCasual._uuid(),
        dataConsumed: defaultCasual.integer(1, 2048),
        dataPackage: ``,
    };
};

const dataBill = (): DataBillDto => {
    return {
        id: defaultCasual._uuid(),
        dataBill: defaultCasual.integer(1, 39),
        billDue: defaultCasual.integer(1, 29),
    };
};

const alert = (): AlertDto => {
    return {
        id: defaultCasual._uuid(),
        type: defaultCasual.random_value(ALERT_TYPE),
        title: defaultCasual._title(),
        description: defaultCasual._short_description(),
        alertDate: new Date(),
    };
};

const node = (): NodeDto => {
    return {
        type: "TOWER",
        id: defaultCasual._uuid(),
        name: defaultCasual._title(),
        description: `${defaultCasual.random_value(NODE_TYPE)} node`,
        status: defaultCasual.random_value(ORG_NODE_STATE),
        totalUser: defaultCasual.integer(1, 99),
        isUpdateAvailable: Math.random() < 0.7,
        updateShortNote:
            "Software update available. Estimated 10 minutes, and will (be/not be) disruptive. ",
        updateDescription:
            "Short introduction.\n\n TL;DR\n\n*** NEW ***\nPoint 1\nPoint 2\nPoint 3\n\n*** IMPROVEMENTS ***\nPoint 1\nPoint 2\nPoint 3\n\n*** FIXES ***\nPoint 1\nPoint 2\nPoint 3\n\nWe would love to hear your feedback -- if you have anything to share, please xyz.",
        updateVersion: "12.4",
    };
};
const esim = (): EsimDto => {
    const boolean = {
        true: true,
        false: false,
    };
    return {
        esim: `# ${defaultCasual.integer(11111, 99999)}-${defaultCasual.date(
            "DD-MM-YYYY"
        )}-${defaultCasual.integer(1111111, 9999999)}`,
        active: defaultCasual.random_value(boolean),
    };
};

const currentBill = (): CurrentBillDto => {
    const data = defaultCasual.integer(1, 10);
    const rate = defaultCasual.integer(3, 6);
    const subtotal = data * rate;
    return {
        id: defaultCasual._uuid(),
        name: defaultCasual._name(),
        dataUsed: data,
        rate: rate,
        subtotal: subtotal,
    };
};
const uploadInvoiceToS3 = (pdfFile: any) => {
    const fileContent = fs.readFileSync(pdfFile);
    const params = {
        Bucket: `${process.env.AWS_S3_BUCKET_NAME}`,
        Key: "invoice.pdf",
        Body: fileContent,
    };

    return s3.upload(params, function (err: any, data: any) {
        if (err) {
            throw err;
        }
        return `${data?.Location}`;
    });
};
const billHistory = (): BillHistoryDto => {
    const totalUsage = defaultCasual.integer(1, 10);
    const subtotal = totalUsage * 3;
    const invoiceUrl = uploadInvoiceToS3(
        createInvoice(totalUsage, subtotal, "invoice.pdf")
    );
    return {
        id: defaultCasual._uuid(),
        date: defaultCasual.date("MM-DD-2021"),
        description: `Bill for month`,
        invoice: invoiceUrl,
        totalUsage,
        subtotal: subtotal,
    };
};

const softwareLogs = (): NodeAppsVersionLogsResponse[] => {
    const lorem = new LoremIpsum();
    const logs: NodeAppsVersionLogsResponse[] = [];
    for (let i = 1; i < 6; i++) {
        logs.push({
            version: `0.${i}`,
            notes: lorem.generateParagraphs(1),
            date: Math.floor(Date.now() - 1000000000 * (6 - i)),
        });
    }
    return logs;
};
const nodeApps = (): NodeAppResponse[] => {
    const lorem = new LoremIpsum();
    const logs: NodeAppResponse[] = [];
    for (let i = 5; i > 0; i--) {
        logs.push({
            id: defaultCasual._uuid(),
            title: lorem.generateWords(1),
            version: `0.${i}`,
            cpu: `${defaultCasual.double(0.1, 100)}`,
            memory: `${defaultCasual.double(0.1, 1024)}`,
        });
    }
    return logs;
};

interface Generators extends Casual.Generators {
    _randomArray: <T>(
        minLength: number,
        maxLength: number,
        elementGenerator: (index?: number) => T
    ) => Array<T>;

    _dataUsage: () => DataUsageDto;
    _dataBill: () => DataBillDto;
    _alert: () => AlertDto;
    _node: () => NodeDto;
    _esim: () => EsimDto;
    _getUser: () => GetUserDto;
    _currentBill: () => CurrentBillDto;
    _billHistory: () => BillHistoryDto;
    _network: () => NetworkDto;
    _nodeNetwork: () => NetworkDto;
    _softwareLogs: () => [NodeAppsVersionLogsResponse];
    _nodeApps: () => [NodeAppResponse];
    functions(): Functions;
}
interface Functions extends Casual.functions {
    randomArray: <T>(
        minLength: number,
        maxLength: number,
        elementGenerator: (index?: number) => T
    ) => Array<T>;
    dataUsage: () => DataUsageDto;
    dataBill: () => DataBillDto;
    alert: () => AlertDto;
    node: () => NodeDto;
    getUser: () => GetUserDto;
    esim: () => EsimDto;
    currentBill: () => CurrentBillDto;
    billHistory: () => BillHistoryDto;
    network: () => NetworkDto;
    nodeNetwork: () => NetworkDto;
    softwareLogs: () => [NodeAppsVersionLogsResponse];
    nodeApps: () => [NodeAppResponse];
}

defaultCasual.define("randomArray", randomArray);
defaultCasual.define("dataUsage", dataUsage);
defaultCasual.define("dataBill", dataBill);
defaultCasual.define("alert", alert);
defaultCasual.define("node", node);
defaultCasual.define("esim", esim);
defaultCasual.define("currentBill", currentBill);
defaultCasual.define("billHistory", billHistory);
defaultCasual.define("softwareLogs", softwareLogs);
defaultCasual.define("nodeApps", nodeApps);
const casual = defaultCasual as Generators & Functions & Casual.Casual;

export default casual;

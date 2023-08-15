import { UserResDto, WhoamiDto } from "../../user/resolver/types";
import { HTTP401Error, Messages } from "../errors";
import { Meta, THeaders } from "../types";

const getTimestampCount = (count: string) =>
  parseInt((Date.now() / 1000).toString()) + "-" + count;

const parseHeaders = (reqHeader: any): THeaders => {
  const headers: THeaders = {
    auth: {
      Authorization: "",
      Cookie: "",
    },
    orgId: "",
    userId: "",
    orgName: "",
  };
  if (reqHeader.get("introspection") === "true") return headers;
  if (reqHeader.get("org-id")) {
    headers.orgId = reqHeader.get("org-id") as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_ORG);
  }
  if (reqHeader.get("user-id")) {
    headers.userId = reqHeader.get("user-id") as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_USER);
  }
  if (reqHeader.get("org-name")) {
    headers.orgName = reqHeader.get("org-name") as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_ORG_NAME);
  }

  if (reqHeader.get("x-session-token") || reqHeader.get("cookie")) {
    if (reqHeader.get("x-session-token")) {
      headers.auth.Authorization = reqHeader["x-session-token"] as string;
    } else {
      const cookie: string = reqHeader.get("cookie");
      const cookies = cookie.split(";");
      const session: string =
        cookies.find(item => (item.includes("ukama_session") ? item : "")) ||
        "";
      headers.auth.Cookie = session;
    }
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_AUTH);
  }
  return headers;
};

const parseGatewayHeaders = (reqHeader: any): THeaders => {
  return {
    auth: {
      Authorization: reqHeader["x-session-token"] || "",
      Cookie: reqHeader["cookie"] || "",
    },
    orgId: reqHeader["orgid"] || "",
    userId: reqHeader["userid"] || "",
    orgName: reqHeader["orgname"] || "",
  };
};

const getStripeIdByUserId = (uid: string): string => {
  switch (uid) {
    case "d0a36c51-6a66-4187-b786-72a9e09bf7a4":
      return "cus_MFTZKUVOGtI2fU";
    default:
      return "cus_MFTZKUVOGtI2fU";
  }
};

const getPaginatedOutput = (
  page: number,
  pageSize: number,
  count: number
): Meta => {
  return {
    count,
    page: page ? page : 1,
    size: pageSize ? pageSize : count,
    pages: pageSize ? Math.ceil(count / pageSize) : 1,
  };
};

export {
  getPaginatedOutput,
  getStripeIdByUserId,
  getTimestampCount,
  parseGatewayHeaders,
  parseHeaders,
};

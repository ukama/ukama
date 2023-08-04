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

  if (reqHeader["org-id"]) {
    headers.orgId = reqHeader["org-id"] as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_ORG);
  }
  if (reqHeader["user-id"]) {
    headers.userId = reqHeader["user-id"] as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_USER);
  }
  if (reqHeader["org-name"]) {
    headers.orgName = reqHeader["org-name"] as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_ORG_NAME);
  }

  if (
    reqHeader["x-session-token"] ||
    (reqHeader["cookie"] && reqHeader["cookie"]["ukama_session"])
  ) {
    if (reqHeader["x-session-token"]) {
      headers.auth.Authorization = reqHeader["x-session-token"] as string;
    } else {
      headers.auth.Cookie = `ukama_session=${reqHeader["cookie"]["ukama_session"]}`;
    }
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_AUTH);
  }
  return headers;
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
  parseHeaders,
};

import { Meta, THeaders } from "../types";

const getTimestampCount = (count: string) =>
  parseInt((Date.now() / 1000).toString()) + "-" + count;

const parseHeaders = (): THeaders => {
  const headers: THeaders = {
    auth: {
      Cookie:
        "ukama_session=MTY4NTk1NDYxOXxkRWpsRnZtbVFBZ1lHV0Jnc1ppeDFJYnhuMWtsbTJGeERRNVZFRWJrakxsdk5MN2ptYjA3UEdEWXkzSklzZHlRVmxEdTBjcElMQmswUFNIMzNHTTNYYkJCZV81R2tSRG5UTUFxem9qQzRzUXpwdkNaQUJTSmlpSkRrbGFRMGNKSHIxd3VrXzdFTlFkWEhITXpueFFaekctdW5paDRXZDJ0aGtuLWpobU1LTmFSQlRvbEk3WE5YWFRSc1k0OV9JaUd2TFBCVkFySHZSVDlrR2lRWkJFQThPUURtdjlCYnBYeU9XQkNHLTFrcEdwWEZQamdfWGFpcEd1MnM5dGxFdk1xT1FPd0ZsWXFJNW1taVptYW44Q3JUQT09fGG25BGdPznBtGqyNf65yGWgB_tfIyOHMHhbnP_ITXEo",
      Authorization: "",
    },
    orgId: "aac1ed88-2546-4f9c-a808-fb9c4d0ef24b",
    userId: "851e0abc-7e91-4206-8c85-498e16f91e67",
    orgName: "ukama",
  };
  // const orgId = ctx.req.headers["org-id"];
  // const userId = ctx.req.headers["user-id"];
  // const orgName = ctx.req.headers["org-name"];
  // if (!orgId || !userId || !orgName)
  //   throw new HTTP401Error(Messages.ERR_REQUIRED_HEADER_NOT_FOUND);
  // else {
  //   headers.orgId = orgId as string;
  //   headers.userId = userId as string;
  //   headers.orgName = orgName as string;
  // }
  // if (ctx.authType === "token") {
  //   headers.auth.Authorization = ctx.req.headers["x-session-token"] as string;
  // } else {
  //   const cookieObj: any = convertCookieToObj(ctx.req.headers["cookie"] || "");
  //   headers.auth.Cookie = `ukama_session=${cookieObj["ukama_session"]}`;
  // }
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

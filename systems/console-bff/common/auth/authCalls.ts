import { API_METHOD_TYPE } from "../enums";
import { asyncRestCall } from "./../axiosClient";
import { AUTH_URL } from "./../configs/index";

const updateAttributes = async (
  userId: string,
  email: string,
  name: string,
  cookie: string,
  isFirstVisit: boolean
) => {
  return await asyncRestCall({
    method: API_METHOD_TYPE.PUT,
    url: `${AUTH_URL}/admin/identities/${userId}`,
    body: {
      schema_id: "default",
      state: "active",
      traits: {
        email: email,
        name: name,
        firstVisit: isFirstVisit,
      },
    },
    headers: {
      cookie: cookie,
    },
  });
};

const getIdentity = async (userId: string, cookie: string) => {
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${AUTH_URL}/admin/identities/${userId}`,
    headers: {
      cookie: cookie,
    },
  });
};

export { getIdentity, updateAttributes };

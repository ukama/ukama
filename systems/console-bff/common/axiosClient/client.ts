import axios from "axios";

import { ApiMethodDataDto } from "../types";

class ApiMethods {
  constructor() {
    axios.create({
      timeout: 10000,
    });
  }
  fetch = async (req: ApiMethodDataDto) => {
    return axios(req as any).catch(err => {
      throw err;
    });
  };
}

export default new ApiMethods();

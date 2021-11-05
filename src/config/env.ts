import { config } from "dotenv";
import { cleanEnv, str, num } from "envalid";

const dotenvResult = config(); //to load env into process.env

if (dotenvResult.error) {
    throw dotenvResult.error;
}

const env = cleanEnv(process.env, {
    NODE_ENV: str({
        choices: ["development", "production", "test"],
        desc: `Node environement - choices are ['development', 'production']`,
        default: "development",
    }),
    PORT: num({
        default: 5000,
        desc: "Port of the express server",
        example: "5000",
    }),
    BASE_URL: str({
        desc: `base url for REST Api calls`,
        default: "localhost:5000",
    }),
});

export default env;

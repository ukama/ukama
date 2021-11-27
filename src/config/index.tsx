const dev = {
    ENVIROMENT: "development",
    REACT_APP_AUTH_URL: "http://localhost:4455",
    REACT_APP_API: "https://bff.dev.ukama.com/graphql",
    REACT_APP_KRATOS_BASE_URL: "http://localhost:4433",
};

const prod = {
    ENVIROMENT: "production",
    REACT_APP_AUTH_URL: "https://auth.dev.ukama.com/",
    REACT_APP_API: "https://bff.dev.ukama.com/graphql",
    REACT_APP_KRATOS_BASE_URL: "https://auth.dev.ukama.com/.api/",
};

const config = process.env.NODE_ENV === "production" ? prod : dev;

export default {
    ...config,
};

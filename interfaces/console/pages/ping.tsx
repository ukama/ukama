import Router from "next/router";
import { useEffect } from "react";

const Ping = () => {
    useEffect(() => {
        const { pathname } = Router;
        if (pathname == "/ping") {
            Router.push("/");
        }
    }, []);
    return null;
};

export default Ping;

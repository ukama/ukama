import React from "react";
import { IRoute } from "./config";
import { Switch } from "react-router-dom";
import RouteWithSubRoutes from "./RouteWithSubRoutes";

interface IProps {
    routes: IRoute[];
}

const Router: React.FC<IProps> = ({ routes }) => {
    return (
        <Switch>
            {routes &&
                routes.map((route: IRoute) => (
                    <RouteWithSubRoutes key={route.path} {...route} />
                ))}
        </Switch>
    );
};

export default Router;

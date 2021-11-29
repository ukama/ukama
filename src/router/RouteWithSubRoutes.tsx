import { IRoute } from "./config";
import { Suspense } from "react";
import { Redirect, Route } from "react-router-dom";

const RouteWithSubRoutes = (route: IRoute) => {
    return (
        <Suspense fallback={route.fallback}>
            <Route
                path={route.path}
                render={props =>
                    route.redirect ? (
                        <Redirect to={route.redirect} />
                    ) : (
                        route.private &&
                        route.component && (
                            <route.component {...props} routes={route.routes} />
                        )
                    )
                }
            />
        </Suspense>
    );
};

export default RouteWithSubRoutes;

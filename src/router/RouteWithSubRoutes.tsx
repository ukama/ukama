import { IRoute } from "./config";
import { Suspense } from "react";
import { Redirect, Route } from "react-router-dom";

const RouteWithSubRoutes = (route: IRoute) => {
    const authenticated: boolean = false;
    return (
        <Suspense fallback={route.fallback}>
            <Route
                path={route.path}
                render={props =>
                    route.redirect ? (
                        <Redirect to={route.redirect} />
                    ) : route.private ? (
                        authenticated ? (
                            route.component && (
                                <route.component
                                    {...props}
                                    routes={route.routes}
                                />
                            )
                        ) : (
                            <Redirect to="/login" />
                        )
                    ) : authenticated ? (
                        <Redirect to="/dashboard" />
                    ) : (
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

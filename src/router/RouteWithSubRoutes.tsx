import { IRoute } from "./config";
import { Suspense } from "react";
import { Redirect, Route } from "react-router-dom";
import { useRecoilValue } from "recoil";
import { isLoginAtom } from "../recoil";

const RouteWithSubRoutes = (route: IRoute) => {
    const isLogin = useRecoilValue(isLoginAtom);
    const authRoute = (props: any) =>
        isLogin ? (
            route.component && (
                <route.component {...props} routes={route.routes} />
            )
        ) : (
            <Redirect to="/login" />
        );
    return (
        <Suspense fallback={route.fallback}>
            <Route
                path={route.path}
                render={props =>
                    route.redirect ? (
                        <Redirect to={route.redirect} />
                    ) : route.private ? (
                        authRoute(props)
                    ) : isLogin ? (
                        <Redirect to="/home" />
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

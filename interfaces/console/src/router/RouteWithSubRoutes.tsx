import Layout from "../layout";
import { IRoute } from "./config";
import { Suspense } from "react";
import { FullscreenContainer } from "../styles";
import { Redirect, Route } from "react-router-dom";
import { useGetAccountDetailsQuery } from "../generated";

const RouteWithSubRoutes = (route: IRoute) => {
    const { data: user } = useGetAccountDetailsQuery();

    const fullScreenRoute = (props: any) =>
        route.private &&
        route.component && (
            <FullscreenContainer>
                <route.component {...props} routes={route.routes} />
            </FullscreenContainer>
        );

    const routesWithLayout = (props: any) =>
        route.private &&
        route.component &&
        (user?.getAccountDetails?.isFirstVisit ? (
            <route.component {...props} routes={route.routes} />
        ) : (
            <Layout>
                <route.component {...props} routes={route.routes} />
            </Layout>
        ));

    const getRouteByType = (props: any) =>
        route.isFullScreen ? fullScreenRoute(props) : routesWithLayout(props);

    return (
        <Suspense fallback={route.fallback}>
            <Route
                path={route.path}
                render={props =>
                    route.redirect ? (
                        <Redirect to={route.redirect} />
                    ) : (
                        getRouteByType(props)
                    )
                }
            />
        </Suspense>
    );
};

export default RouteWithSubRoutes;

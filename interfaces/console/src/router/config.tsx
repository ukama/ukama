import { CenterContainer } from "../styles";
import { CircularProgress } from "@mui/material";
import { lazy, ReactNode, ComponentType, LazyExoticComponent } from "react";

export interface IRoute {
    // Path, like in basic prop
    path: string;
    // Exact, like in basic prop
    exact: boolean;
    // Preloader for lazy loading
    fallback: NonNullable<ReactNode> | null;
    // Lazy Loaded component
    component?: LazyExoticComponent<ComponentType<any>>;
    // Sub routes
    routes?: IRoute[];
    // Redirect path
    redirect?: string;
    // If router is private, this is going to be true
    private?: boolean;

    isFullScreen?: boolean;
}

const Loader = (
    <CenterContainer>
        <CircularProgress />
    </CenterContainer>
);

const getRouteObject = (
    path = "/",
    component = "Home",
    isPrivate = true,
    isFullScreen = false
) => {
    return {
        path: path,
        exact: true,
        fallback: Loader,
        private: isPrivate,
        isFullScreen: isFullScreen,
        component: lazy(() => import(`../pages/${component}`)),
    };
};

export const routes = {
    Root: getRouteObject("/", "Home", true),
    Nodes: getRouteObject("/nodes", "Nodes", true),
    Users: getRouteObject("/users", "Users", true),
    Settings: getRouteObject("/settings", "Settings", true, true),
    Billing: getRouteObject("/billing", "Billing", true),

    //Public Routes//
    Error: getRouteObject("/*", "ErrorPage", true, true),
    //
};

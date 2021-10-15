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
}

const Loader = <div> Loading... </div>;

export const routes: IRoute[] = [
    //Default routes//
    {
        path: "/",
        exact: true,
        private: false,
        redirect: "/login",
        fallback: Loader,
    },
    //

    //Privatte Routes//
    //

    //Public Routes//
    {
        path: "/forgotPasswordConfirmation",
        component: lazy(() => import("../pages/ForgotPasswordCofirmation")),
        exact: false,
        private: false,
        fallback: Loader,
    },
    {
        path: "/login",
        component: lazy(() => import("../pages/Login")),
        exact: false,
        private: false,
        fallback: Loader,
    },
    {
        path: "/signUp",
        component: lazy(() => import("../pages/SignUp")),
        exact: false,
        private: false,
        fallback: Loader,
    },
    {
        path: "/forgot-password",
        component: lazy(() => import("../pages/ForgotPassword")),
        exact: false,
        private: false,
        fallback: Loader,
    },
    {
        path: "/*",
        component: lazy(() => import("../pages/ErrorPage")),
        exact: false,
        private: false,
        fallback: Loader,
    },
    //
];

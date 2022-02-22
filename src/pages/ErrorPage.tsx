import React from "react";
import { IRoute } from "../router/config";
import { CenterContainer } from "../styles";
const PageNotFound = React.lazy(() =>
    import("../assets/svg").then(module => ({
        default: module.PageNotFound,
    }))
);

interface IProps {
    routes: IRoute[];
}

const ErrorPage: React.FC<IProps> = () => {
    return (
        <CenterContainer>
            <PageNotFound />
        </CenterContainer>
    );
};

export default ErrorPage;

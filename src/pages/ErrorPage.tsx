import React from "react";
import { IRoute } from "../router/config";
import { CenterContainer } from "../styles";
import { PageNotFound } from "../assets/svg";

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

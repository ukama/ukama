import React from "react";
import { IRoute } from "../router/config";
import { Typography } from "@mui/material";
import { CenterContainer } from "../styles";

interface IProps {
    routes: IRoute[];
}

const ErrorPage: React.FC<IProps> = () => {
    return (
        <CenterContainer>
            <Typography variant="h5" color="error">
                <b>404:</b> Page Not found!
            </Typography>
        </CenterContainer>
    );
};

export default ErrorPage;

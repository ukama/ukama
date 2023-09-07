import React from "react";
import { IRoute } from "../router/config";
import { CenterContainer } from "../styles";
import { Stack, Typography } from "@mui/material";

interface IProps {
    routes: IRoute[];
}

const ErrorPage: React.FC<IProps> = () => {
    return (
        <CenterContainer>
            <Stack spacing={2} alignItems="center">
                <Typography variant="h2">404</Typography>
                <Typography variant="h5">Not Found</Typography>
                <Typography variant="subtitle2">
                    The resource requested could not be found on this server!
                </Typography>
            </Stack>
        </CenterContainer>
    );
};

export default ErrorPage;

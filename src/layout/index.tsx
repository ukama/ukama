import Header from "./Header";
import Sidebar from "./Sidebar";
import { colors } from "../theme";
import { useState } from "react";
import { useRecoilValue } from "recoil";
import { useHistory } from "react-router";
import { getTitleFromPath } from "../utils";
import { isSkeltonLoading } from "../recoil";
import { Box, CssBaseline } from "@mui/material";
const Layout = (props: any) => {
    const { children } = props;
    const history = useHistory();
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [page, setPage] = useState(
        getTitleFromPath(history?.location?.pathname || "/")
    );
    const [isOpen, setIsOpen] = useState(false);

    const handlePageChange = (page: string) => {
        setPage(page);
    };

    const handleDrawerToggle = () => setIsOpen(() => !isOpen);

    return (
        <Box
            sx={{
                display: "flex",
                height: "100%",
                backgroundColor: colors.solitude,
            }}
        >
            <CssBaseline />
            <Sidebar
                page={page}
                isOpen={isOpen}
                isLoading={isSkeltonLoad}
                handlePageChange={handlePageChange}
                handleDrawerToggle={handleDrawerToggle}
            />
            <Box
                component="main"
                sx={{
                    pl: { xs: 2, md: 3 },
                    pr: { xs: 2, md: 3 },
                    width: "100%",
                }}
            >
                <Header
                    pageName={page}
                    isLoading={isSkeltonLoad}
                    handlePageChange={handlePageChange}
                    handleDrawerToggle={handleDrawerToggle}
                />
                {children}
            </Box>
        </Box>
    );
};
export default Layout;

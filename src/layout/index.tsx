import Header from "./Header";
import Sidebar from "./Sidebar";
import { useState } from "react";
import { colors } from "../theme";
import { Box } from "@mui/material";
import { isSkeltonLoading, pageName } from "../recoil";
import { useRecoilState, useRecoilValue } from "recoil";

const Layout = (props: any) => {
    const { children } = props;
    const [page, setPage] = useRecoilState(pageName);
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);

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
                    width: "100%",
                    pl: { xs: 2, md: 3 },
                    pr: { xs: 2, md: 3 },
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

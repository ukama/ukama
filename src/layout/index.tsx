import Header from "./Header";
import Sidebar from "./Sidebar";
import { useState } from "react";
import { Box, Stack } from "@mui/material";
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
            component="div"
            sx={{
                display: "flex",
                height: "100%",
            }}
        >
            <Sidebar
                page={page}
                isOpen={isOpen}
                isLoading={isSkeltonLoad}
                handlePageChange={handlePageChange}
                handleDrawerToggle={handleDrawerToggle}
            />

            <Stack
                spacing={{ xs: 2, md: 4 }}
                sx={{
                    width: "100%",
                    pl: { xs: 2, md: 3, xl: 5 },
                    pr: { xs: 2, md: 3, xl: 5 },
                }}
            >
                <Header
                    pageName={page}
                    isLoading={isSkeltonLoad}
                    handlePageChange={handlePageChange}
                    handleDrawerToggle={handleDrawerToggle}
                />

                {children}
            </Stack>
        </Box>
    );
};
export default Layout;

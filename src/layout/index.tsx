import Header from "./Header";
import Sidebar from "./Sidebar";
import { useState } from "react";
import { colors } from "../theme";
import { useRecoilValue } from "recoil";
import { useHistory } from "react-router";
import { getTitleFromPath } from "../utils";
import { isSkeltonLoading } from "../recoil";
import { Box, CssBaseline } from "@mui/material";
const Layout = (props: any) => {
    const { children } = props;
    const history = useHistory();
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [path, setPath] = useState(
        getTitleFromPath(history?.location?.pathname || "/")
    );
    const [isOpen, setIsOpen] = useState(false);
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
                path={path}
                isOpen={isOpen}
                setPath={setPath}
                isLoading={isSkeltonLoad}
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
                    pageName={path}
                    isLoading={isSkeltonLoad}
                    handleDrawerToggle={handleDrawerToggle}
                />
                {children}
            </Box>
        </Box>
    );
};
export default Layout;

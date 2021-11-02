import Header from "./Header";
import Sidebar from "./Sidebar";
import { useState } from "react";
import { useHistory } from "react-router";
import { getTitleFromPath } from "../utils";
import { Box, CssBaseline, Toolbar } from "@mui/material";
import { colors } from "../theme";
const Layout = (props: any) => {
    const { children } = props;
    const history = useHistory();
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
            <Header pageName={path} handleDrawerToggle={handleDrawerToggle} />
            <Sidebar
                path={path}
                isOpen={isOpen}
                setPath={setPath}
                handleDrawerToggle={handleDrawerToggle}
            />
            <Box component="main" sx={{ p: 3, flexGrow: 1 }}>
                <Toolbar />
                {children}
            </Box>
        </Box>
    );
};
export default Layout;

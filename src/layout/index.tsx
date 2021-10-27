import Sidebar from "./Sidebar";
import { Box, CssBaseline, Toolbar } from "@mui/material";
import { useState } from "react";
import Header from "./Header";
const Layout = (props: any) => {
    const { children } = props;
    const [isOpen, setIsOpen] = useState(false);
    const [path, setPath] = useState("Home");

    const handleDrawerToggle = () => setIsOpen(() => !isOpen);
    return (
        <Box sx={{ display: "flex" }} style={{ backgroundColor: "#F5F6F8" }}>
            <CssBaseline />
            <Sidebar
                path={path}
                isOpen={isOpen}
                setPath={setPath}
                handleDrawerToggle={handleDrawerToggle}
            />
            <Header pageName={path} />
            <Box
                component="main"
                sx={{ width: "100%", margin: "5% 2% 0px 4%" }}
            >
                <Toolbar />
                {children}
            </Box>
        </Box>
    );
};
export default Layout;

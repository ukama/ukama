import { Box, CssBaseline, Toolbar } from "@mui/material";
import Sidebar from "./Sidebar";
import { useState } from "react";

const Layout = (props: any) => {
    const { children } = props;
    const [isOpen, setIsOpen] = useState(false);

    const handleDrawerToggle = () => setIsOpen(() => !isOpen);
    return (
        <Box sx={{ display: "flex" }}>
            <CssBaseline />
            <Sidebar isOpen={isOpen} handleDrawerToggle={handleDrawerToggle} />
            <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
                <Toolbar />
                {children}
            </Box>
        </Box>
    );
};
export default Layout;

import Header from "./Header";
import Sidebar from "./Sidebar";
import { useState } from "react";
import { colors } from "../theme";
import { useHistory } from "react-router";
import { getTitleFromPath } from "../utils";
import { Box, CssBaseline } from "@mui/material";
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
            <Sidebar
                path={path}
                isOpen={isOpen}
                setPath={setPath}
                handleDrawerToggle={handleDrawerToggle}
            />
            <Box component="main" sx={{ pl: 3, pr: 3 }}>
                <Header
                    pageName={path}
                    handleDrawerToggle={handleDrawerToggle}
                />
                {children}
            </Box>
        </Box>
    );
};
export default Layout;

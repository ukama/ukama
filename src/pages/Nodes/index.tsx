import { Box } from "@mui/material";
import { PagePlaceholder } from "../../components";

const Nodes = () => {
    return (
        <Box sx={{ height: "calc(100vh - 8vh)", p: "28px 0px" }}>
            <PagePlaceholder
                showActionButton={false}
                description={
                    "Your nodes have not arrived yet. View their status"
                }
                hyperlink={"#"}
                linkText={"here"}
            />
        </Box>
    );
};

export default Nodes;

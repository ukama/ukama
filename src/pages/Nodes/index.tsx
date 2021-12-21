import { Box } from "@mui/material";
import { useRecoilValue } from "recoil";
import { LoadingWrapper, PagePlaceholder } from "../../components";
import { isSkeltonLoading } from "../../recoil";

const Nodes = () => {
    const _isSkeltonLoading = useRecoilValue(isSkeltonLoading);
    return (
        <Box sx={{ height: "calc(100vh - 8vh)", p: "28px 0px" }}>
            <LoadingWrapper isLoading={_isSkeltonLoading}>
                <PagePlaceholder
                    showActionButton={false}
                    description={
                        "Your nodes have not arrived yet. View their status"
                    }
                    hyperlink={"#"}
                    linkText={"here"}
                />
            </LoadingWrapper>
        </Box>
    );
};

export default Nodes;

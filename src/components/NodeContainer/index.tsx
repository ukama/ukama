import { Box } from "@mui/material";
import RouterIcon from "@mui/icons-material/Router";
import { GetNodesByOrgQuery } from "../../generated";
import { EmptyView, MultiSlideCarousel, NodeCard } from "..";
type NodeContainerProps = {
    items?: GetNodesByOrgQuery["getNodesByOrg"]["nodes"];
    slidesToShow: number;
    count: number | undefined;
    handleItemAction: Function;
};

const NodeContainer = ({
    count = 0,
    items,
    slidesToShow,
    handleItemAction,
}: NodeContainerProps) => {
    return (
        <Box
            component="div"
            sx={{ minHeight: "208px", display: "flex", alignItems: "center" }}
        >
            {count > 1 ? (
                <MultiSlideCarousel
                    numberOfSlides={slidesToShow}
                    disableArrows={count < 3}
                >
                    {items?.map(({ id, title, totalUser, description }) => (
                        <NodeCard
                            key={id}
                            title={title}
                            loading={false}
                            users={totalUser}
                            subTitle={description}
                            handleOptionItemClick={(type: string) =>
                                handleItemAction(id, type)
                            }
                        />
                    ))}
                </MultiSlideCarousel>
            ) : (
                <EmptyView
                    size="large"
                    title="No nodes yet!"
                    icon={RouterIcon}
                />
            )}
        </Box>
    );
};

export default NodeContainer;

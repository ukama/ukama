import { MultiSlideCarousel, NodeCard } from "..";
import { GetNodesByOrgQuery } from "../../generated";

type NodeContainerProps = {
    items?: GetNodesByOrgQuery["getNodesByOrg"]["nodes"];
    slidesToShow: number;
    count: number | undefined;
    handleItemAction: Function;
};

const NodeContainer = ({
    count = 0,
    items = [],
    slidesToShow,
    handleItemAction,
}: NodeContainerProps) => {
    return (
        <>
            {count > 1 ? (
                <MultiSlideCarousel
                    numberOfSlides={slidesToShow}
                    disableArrows={count < 3}
                >
                    {items.map(({ id, title, totalUser, description }) => (
                        <NodeCard
                            key={id}
                            title={title}
                            users={totalUser}
                            loading={false}
                            subTitle={description}
                            handleOptionItemClick={(type: string) =>
                                handleItemAction(id, type)
                            }
                        />
                    ))}
                </MultiSlideCarousel>
            ) : (
                <NodeCard isConfigure={true} />
            )}
        </>
    );
};

export default NodeContainer;

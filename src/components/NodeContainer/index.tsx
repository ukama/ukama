import { MultiSlideCarousel, NodeCard } from "..";

type NodeContainerProps = {
    items: any;
    slidesToShow: number;
    count: number | undefined;
};

const NodeContainer = ({
    count = 0,
    items = [],
    slidesToShow,
}: NodeContainerProps) => {
    return (
        <>
            {count > 1 ? (
                <MultiSlideCarousel
                    numberOfSlides={slidesToShow}
                    disableArrows={count < 3 ? true : false}
                >
                    {items.map(({ id, title, totalUser, description }: any) => (
                        <NodeCard
                            key={id}
                            title={title}
                            users={totalUser}
                            loading={false}
                            subTitle={description}
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

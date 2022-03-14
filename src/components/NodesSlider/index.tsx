import { Stack, Theme, useMediaQuery } from "@mui/material";
import { useEffect, useState } from "react";
import Carousel from "react-material-ui-carousel";
import { NodeCard } from "..";
import { NodeDto } from "../../generated";

interface INodeSlider {
    items: NodeDto[];
    handleItemAction: Function;
}

const NodeSlider = ({ items, handleItemAction }: INodeSlider) => {
    const small = useMediaQuery((theme: Theme) => theme.breakpoints.up("sm"));
    const medium = useMediaQuery((theme: Theme) => theme.breakpoints.up("md"));
    const [list, setList] = useState<any>([]);

    useEffect(() => {
        const slides = [];
        const chunk = medium ? 3 : small ? 2 : 1;
        chunk;
        for (let i = 0; i < items.length; i += chunk) {
            slides.push({
                cid: `chunk-${i}`,
                items: items.slice(i, i + chunk),
            });
        }
        setList(slides);
    }, [small, medium]);

    return (
        <Carousel
            swipe={true}
            animation="slide"
            autoPlay={false}
            indicators={false}
            cycleNavigation={false}
            sx={{ width: "100%", minHeight: "246px", py: "20px" }}
            navButtonsAlwaysVisible
        >
            {list.map(({ cid, items }: any) => (
                <Stack
                    key={cid}
                    direction={"row"}
                    spacing={2}
                    justifyContent={items.length > 2 ? "center" : "flex-start"}
                >
                    {items.map(({ id, title, description }: any, i: number) => (
                        <NodeCard
                            key={i}
                            users={3}
                            title={title}
                            loading={false}
                            subTitle={description}
                            handleOptionItemClick={(type: string) =>
                                handleItemAction(id, type)
                            }
                        />
                    ))}
                </Stack>
            ))}
        </Carousel>
    );
};

export default NodeSlider;

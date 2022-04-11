import { Stack, Theme, useMediaQuery } from "@mui/material";
import { useEffect, useState } from "react";
import Carousel from "react-material-ui-carousel";
import { NodeCard } from "..";
import { NodeDto } from "../../generated";

interface INodeSlider {
    items: NodeDto[];
    handleItemAction: Function;
    handleNodeUpdate: Function;
}

const NodeSlider = ({
    items,
    handleItemAction,
    handleNodeUpdate,
}: INodeSlider) => {
    const small = useMediaQuery((theme: Theme) => theme.breakpoints.up("sm"));
    const medium = useMediaQuery((theme: Theme) => theme.breakpoints.up("md"));
    const [list, setList] = useState<any>([]);

    useEffect(() => {
        const slides = [];
        const isSmall = small ? 2 : 1;
        const chunk = medium ? 3 : isSmall;
        for (let i = 0; i < items.length; i += chunk) {
            slides.push({
                cid: `chunk-${i}`,
                item: items.slice(i, i + chunk),
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
            navButtonsAlwaysVisible
            sx={{ width: "100%", minHeight: "240px", pt: 3, pb: 0 }}
            navButtonsProps={{
                style: {
                    margin: 0,
                },
            }}
        >
            {list.map(({ cid, item }: any) => (
                <Stack
                    key={cid}
                    spacing={2}
                    direction={"row"}
                    sx={{
                        justifyContent: {
                            xs: "center",
                            md: items.length > 1 ? "center" : "flex-start",
                        },
                    }}
                >
                    {item.map(
                        (
                            {
                                id,
                                type,
                                name,
                                description,
                                updateShortNote,
                                isUpdateAvailable,
                            }: any,
                            i: number
                        ) => (
                            <NodeCard
                                key={i}
                                id={id}
                                users={3}
                                type={type}
                                title={name}
                                loading={false}
                                subTitle={description}
                                updateShortNote={updateShortNote}
                                isUpdateAvailable={isUpdateAvailable}
                                handleOptionItemClick={(type: string) =>
                                    handleItemAction(id, type)
                                }
                                handleNodeUpdate={handleNodeUpdate}
                            />
                        )
                    )}
                </Stack>
            ))}
        </Carousel>
    );
};

export default NodeSlider;

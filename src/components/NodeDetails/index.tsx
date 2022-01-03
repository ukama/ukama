import { useState } from "react";
import { Grid, Stack, Tab, Tabs } from "@mui/material";
import { NodeDetailsPanel, NodeDetailsCard } from "../../components";

interface NodeDetailsProps {
    detailsList: Array<{
        title: string;
        renderPropertyStats?: boolean;
        image?: { src: string; alt: string };
        button?: { label: string; onClick: Function };
        properties: Array<{
            name: string;
            value: string | number;
        }>;
    }>;
}

const NodeDetails = ({ detailsList }: NodeDetailsProps) => {
    const [selectedDetail, setSelectedDetail] = useState(0);

    return (
        <>
            <Grid item xs={12} md={4}>
                <Tabs
                    sx={{
                        display: {
                            xs: "unset",
                            md: "none",
                        },
                    }}
                    variant="fullWidth"
                    value={selectedDetail}
                    aria-label="wrapped label tabs example"
                    onChange={(_, value) => setSelectedDetail(value)}
                >
                    {detailsList.map(({ title }, index) => (
                        <Tab key={index} value={index} label={title} />
                    ))}
                </Tabs>
                <Stack spacing={2}>
                    {detailsList.map(({ title, properties }, index) => (
                        <NodeDetailsCard
                            key={index}
                            index={index}
                            title={title}
                            value={selectedDetail}
                            properties={properties}
                            onClick={() => setSelectedDetail(index)}
                        />
                    ))}
                </Stack>
            </Grid>
            <Grid item xs={12} md={8}>
                {detailsList.map(
                    (
                        {
                            title,
                            image,
                            button,
                            properties,
                            renderPropertyStats,
                        },
                        index
                    ) => (
                        <NodeDetailsPanel
                            key={index}
                            index={index}
                            title={title}
                            image={image}
                            button={button}
                            value={selectedDetail}
                            properties={properties}
                            renderPropertyStats={renderPropertyStats}
                        />
                    )
                )}
            </Grid>
        </>
    );
};

export default NodeDetails;

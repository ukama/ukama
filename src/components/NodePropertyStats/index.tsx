import React, { useState } from "react";
import {
    Grid,
    Typography,
    ToggleButton,
    ToggleButtonGroup,
} from "@mui/material";
import { colors } from "../../theme";
import { RechartsData } from "../../constants/stubData";
import { statsPeriodItemType } from "../../types";
import {
    Bar,
    Line,
    XAxis,
    YAxis,
    Tooltip,
    ComposedChart,
    CartesianGrid,
    ResponsiveContainer,
} from "recharts";

interface NodePropertyStatsPorps {
    propery: { name: string; value: string | number };
    periodOptions: Array<{
        id: number;
        label: string;
    }>;
}

const NodePropertyStats = ({
    propery,
    periodOptions,
}: NodePropertyStatsPorps) => {
    const [selectedBtn, setSelectedBtn] = useState("DAY");

    const handleSelectedButtonChange = (
        event: React.MouseEvent<HTMLElement>,
        newSelected: string
    ) => {
        setSelectedBtn(newSelected);
    };

    return (
        <Grid container spacing={1} marginTop={2}>
            <Grid item container xs={12}>
                <Grid item xs={12} sm={6}>
                    <Typography variant="body1">{propery.name}</Typography>
                </Grid>
                <Grid
                    item
                    xs={12}
                    sm={6}
                    display="flex"
                    justifyContent={"flex-end"}
                >
                    <ToggleButtonGroup
                        exclusive
                        size="small"
                        color="primary"
                        value={selectedBtn}
                        onChange={handleSelectedButtonChange}
                    >
                        {periodOptions.map(
                            ({ id, label }: statsPeriodItemType) => (
                                <ToggleButton
                                    fullWidth
                                    key={id}
                                    value={label}
                                    style={{
                                        height: "32px",
                                        color: colors.primaryMain,
                                        border: `1px solid ${colors.primaryMain}`,
                                    }}
                                >
                                    <Typography
                                        variant="body2"
                                        sx={{
                                            p: "0px 2px",
                                            fontWeight: 600,
                                        }}
                                    >
                                        {label}
                                    </Typography>
                                </ToggleButton>
                            )
                        )}
                    </ToggleButtonGroup>
                </Grid>
            </Grid>
            <Grid item xs={12}>
                <ResponsiveContainer width="100%" height={300}>
                    <ComposedChart data={RechartsData}>
                        <CartesianGrid stroke="#f5f5f5" />
                        <XAxis dataKey="name" scale="band" />
                        <YAxis />
                        <Tooltip />
                        <Bar dataKey="uv" barSize={20} fill="#413ea0" />
                        <Line type="monotone" dataKey="uv" stroke="#ff7300" />
                    </ComposedChart>
                </ResponsiveContainer>
            </Grid>
        </Grid>
    );
};

export default NodePropertyStats;

import { useState } from "react";
import { styled } from "@mui/styles";
import { LoadingWrapper } from "../../components";
import { Button, Paper, Typography, Grid } from "@mui/material";

const StyledPaper = styled(Paper)({
    width: "100%",
    borderRadius: 4,
    cursor: "pointer",
    textAlign: "start",
    padding: "19px 24px",
    boxShadow: "2px 2px 6px rgba(0, 0, 0, 0.05)",
});

interface NodeDetailsCardProps {
    title: string;
    loading?: any;
    index: number;
    value: number;
    onClick: Function;
    properties: Array<{ name: string; value: string | number }>;
}

const NodeDetailsCard = ({
    title,
    index,
    value,
    onClick,
    loading,
    properties,
}: NodeDetailsCardProps) => {
    const selected = index === value;
    const [viewMore, setViewMore] = useState(false);
    return (
        <StyledPaper
            onClick={() => onClick()}
            sx={{
                display: {
                    xs: selected ? "unset" : "none",
                    md: "unset",
                },
                borderLeft: {
                    xs: "none",
                    md: `8px solid ${selected ? "#2190F6" : "#FFFFFF"}`,
                },
            }}
        >
            <LoadingWrapper width="100%" height="100px" isLoading={loading}>
                <Typography variant="h6">{title}</Typography>
                {(viewMore ? properties : properties.slice(0, 4)).map(
                    ({ name, value: properyValue }, keyIndex) => (
                        <Grid container spacing={5} key={keyIndex}>
                            <Grid item xs={5}>
                                <Typography
                                    variant="subtitle1"
                                    fontWeight={500}
                                >
                                    {name}
                                </Typography>
                            </Grid>
                            <Grid item xs={7}>
                                <Typography variant="subtitle1">
                                    {properyValue}
                                </Typography>
                            </Grid>
                        </Grid>
                    )
                )}
                {properties.length > 4 && (
                    <Button
                        variant="text"
                        style={{ textTransform: "none" }}
                        onClick={e => {
                            e.stopPropagation();
                            e.preventDefault();
                            setViewMore(val => !val);
                        }}
                    >
                        {viewMore ? "View less" : "View more"}
                    </Button>
                )}
            </LoadingWrapper>
        </StyledPaper>
    );
};

export default NodeDetailsCard;

import { useState } from "react";
import { styled } from "@mui/styles";
import { Button, Paper, Stack, Typography } from "@mui/material";

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
            <Typography variant="h6">{title}</Typography>
            {(viewMore ? properties : properties.slice(0, 4)).map(
                ({ name, value: properyValue }, keyIndex) => (
                    <Stack
                        key={keyIndex}
                        direction="row"
                        justifyContent="space-between"
                    >
                        <Typography variant="subtitle1" fontWeight={500}>
                            {name}
                        </Typography>
                        <Typography variant="subtitle1">
                            {properyValue}
                        </Typography>
                    </Stack>
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
        </StyledPaper>
    );
};

export default NodeDetailsCard;

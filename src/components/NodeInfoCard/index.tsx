import { useState } from "react";
import { styled } from "@mui/styles";
import { TObject } from "../../types";
import { LoadingWrapper } from "..";
import { Button, Paper, Typography, Grid, Container } from "@mui/material";
import { colors } from "../../theme";

const StyledPaper = styled(Paper)({
    width: "100%",
    display: "inherit",
    cursor: "pointer",
    textAlign: "start",
    boxShadow: "2px 2px 6px rgba(0, 0, 0, 0.05)",
});

interface INodeInfoCard {
    title: string;
    index: number;
    loading: boolean;
    isSelected: boolean;
    onSelected: Function;
    properties: Array<TObject>;
}

const NodeInfoCard = ({
    index,
    onSelected,
    loading = true,
    properties = [],
    isSelected = false,
    title = "Node Detail",
}: INodeInfoCard) => {
    const [viewMore, setViewMore] = useState(false);
    const handleClick = () => onSelected(index);
    return (
        <StyledPaper onClick={handleClick}>
            <LoadingWrapper height="100px" radius="small" isLoading={loading}>
                <Container
                    sx={{
                        margin: 0,
                        padding: "19px 24px",
                        borderRadius: "4px",
                        borderLeft: {
                            xs: "none",
                            md: `8px solid ${
                                isSelected ? colors.primary : colors.white
                            }`,
                        },
                    }}
                >
                    <Typography variant="h6" pb={"8px"}>
                        {title}
                    </Typography>
                    {(viewMore ? properties : properties.slice(0, 4)).map(
                        ({ name, value: properyValue }, keyIndex) => (
                            <Grid container spacing={3} key={keyIndex}>
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
                </Container>
            </LoadingWrapper>
        </StyledPaper>
    );
};

export default NodeInfoCard;

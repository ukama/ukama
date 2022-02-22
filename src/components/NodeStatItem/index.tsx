import Zoom from "@mui/material/Zoom";
import { TVariant } from "../../types";
import { HorizontalContainer } from "../../styles";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import { Grid, Typography, Tooltip, IconButton } from "@mui/material";

interface INodeStatItem {
    name: string;
    value: string;
    variant?: TVariant;
    nameInfo?: string; //Tooltip info about stat
    valueInfo?: string; //Tooltip info about stat value
    showAlertInfo?: boolean; //Pass true if its an alert value
}

interface ITextWithToolTip {
    title: string;
    isAlert?: boolean;
    tooltipText: string;
    showToottip?: boolean;
}

const TextWithToolTip = ({
    title,
    tooltipText,
    isAlert = false,
    showToottip = false,
}: ITextWithToolTip) => (
    <HorizontalContainer>
        <Typography
            fontWeight={500}
            variant="subtitle1"
            color={isAlert ? "error" : "default"}
            sx={{ textTransform: "capitalize" }}
        >
            {title}
        </Typography>
        {showToottip && (
            <Tooltip
                arrow
                placement="right"
                title={tooltipText}
                TransitionComponent={Zoom}
            >
                <IconButton
                    color={isAlert ? "error" : "default"}
                    sx={{
                        "&:hover svg path": {
                            fill: "inherit",
                        },
                    }}
                >
                    <InfoOutlinedIcon
                        sx={{
                            width: "16px",
                            height: "16px",
                        }}
                    />
                </IconButton>
            </Tooltip>
        )}
    </HorizontalContainer>
);

const variants = (variant: TVariant, key: string) => {
    switch (key) {
        case "NG":
            return variant === "small" ? 4 : variant === "large" ? 8 : 6;
        case "VG":
            return variant === "small" ? 8 : variant === "large" ? 4 : 6;
    }
};

const NodeStatItem = ({
    name,
    value,
    nameInfo = "",
    valueInfo = "",
    variant = "medium",
    showAlertInfo = false,
}: INodeStatItem) => {
    return (
        <Grid container spacing={3}>
            <Grid item xs={variants(variant, "NG")}>
                <TextWithToolTip
                    title={name}
                    showToottip={!!nameInfo}
                    tooltipText={nameInfo}
                />
            </Grid>
            <Grid item xs={variants(variant, "VG")}>
                <TextWithToolTip
                    title={value}
                    isAlert={showAlertInfo}
                    tooltipText={valueInfo}
                    showToottip={!!valueInfo}
                />
            </Grid>
        </Grid>
    );
};

export default NodeStatItem;

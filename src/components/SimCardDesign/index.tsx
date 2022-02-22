import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";
import { LinkStyle } from "../../styles";
import { Typography, Card, Grid, Theme } from "@mui/material";
import SimCardOutlinedIcon from "@mui/icons-material/SimCardOutlined";

type StyleProps = {
    isSelected?: boolean;
};

const useStyles = makeStyles<Theme, StyleProps>(() => ({
    cardStyle: {
        marginBottom: 16,
        cursor: "pointer",
        padding: "13px 18px",
        border: ({ isSelected }) =>
            isSelected ? `2px solid ${colors.primaryMain}` : "none",
    },
}));

type SimCardDesignProps = {
    id: number;
    title: string;
    serial: string;
    isActivate?: boolean;
    isSelected: boolean;
    handleItemClick: Function;
};

const SimCardDesign = ({
    id,
    title,
    serial,
    isSelected,
    isActivate,
    handleItemClick,
}: SimCardDesignProps) => {
    const classes = useStyles({ isSelected });
    return (
        <Card className={classes.cardStyle} onClick={() => handleItemClick(id)}>
            <Grid spacing={1} item container>
                <Grid item xs={1}>
                    <SimCardOutlinedIcon />
                </Grid>
                <Grid item xs={2}>
                    <LinkStyle underline="hover" sx={{ fontSize: "14px" }}>
                        {title}
                    </LinkStyle>
                </Grid>
                <Grid item xs={7}>
                    <Typography variant="body1">{serial}</Typography>
                </Grid>
                <Grid item xs={2}>
                    {isActivate && (
                        <Typography
                            variant="body1"
                            sx={{ color: colors.primaryMain }}
                        >
                            Activated
                        </Typography>
                    )}
                </Grid>
            </Grid>
        </Card>
    );
};

export default SimCardDesign;

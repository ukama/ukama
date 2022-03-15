import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";
import ShareIcon from "@mui/icons-material/Share";
import {
    Typography,
    Button,
    IconButton,
    Divider,
    Stack,
    Card,
    Theme,
} from "@mui/material";
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
            <Stack
                direction="row"
                spacing={2}
                sx={{ alignItems: "center", width: "100%" }}
            >
                <Stack
                    direction="row"
                    spacing={1}
                    sx={{ alignItems: "center" }}
                    divider={<Divider orientation="vertical" flexItem />}
                >
                    <Typography
                        variant="body1"
                        sx={{
                            fontSize: "14px",
                            color: colors.black,
                            fontWeight: "bold",
                        }}
                    >
                        {title}
                    </Typography>
                    <Typography variant="body1">{serial}</Typography>
                </Stack>
                <Stack
                    direction="row"
                    spacing={2}
                    justifyContent="flex-end"
                    alignItems="center"
                >
                    {isActivate && (
                        <Button sx={{ color: colors.black70 }}>
                            AWAITING ACTIVATION
                        </Button>
                    )}
                    <IconButton>
                        <ShareIcon sx={{ color: colors.black70 }} />
                    </IconButton>
                </Stack>
            </Stack>
        </Card>
    );
};

export default SimCardDesign;

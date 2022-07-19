import { colors } from "../../theme";
import { Alert, AlertColor, Button, Stack, Typography } from "@mui/material";

interface IBillingAlerts {
    title: string;
    btnText: string;
    type: AlertColor;
    onActionClick: Function;
}

const BillingAlerts = ({
    type,
    title,
    btnText,
    onActionClick,
}: IBillingAlerts) => {
    return (
        <Alert icon={false} severity={type}>
            <Stack direction={"row"} p={0} spacing={1}>
                <Typography variant="body1">{title}</Typography>
                <Button
                    variant="text"
                    sx={{
                        textTransform: "none",
                        color: colors.primaryMain,
                        ":hover": {
                            color: theme => theme.palette.text.primary,
                            backgroundColor: "none",
                        },
                    }}
                    onClick={() => onActionClick()}
                >
                    {btnText}
                </Button>
            </Stack>
        </Alert>
    );
};

export default BillingAlerts;

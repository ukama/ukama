import React, { useEffect } from "react";
import {
    Typography,
    OutlinedInput,
    Box,
    FormControl,
    SelectChangeEvent,
    Button,
    MenuItem,
    Select,
    Stack,
    InputLabel,
} from "@mui/material";
import { makeStyles } from "@mui/styles";
import { useRecoilValue } from "recoil";
import { isDarkmode } from "../../recoil";
import colors from "../../theme/colors";
const useStyles = makeStyles(() => ({
    "&.MuiFormHelperText-root.Mui-error": {
        color: "red",
    },
    selectStyle: () => ({
        width: "100%",
        height: "48px",
    }),
    formControl: {
        width: "100%",
        height: "48px",
        paddingBottom: "55px",
    },
}));

const getSelected = (list: any) => {
    if (list.length > 1) return list[0].value;
    return "no_payment_method_Set";
};

interface IPaymentProps {
    title: string;
    paymentMethodData: any;
    onAddPaymentMethod: any;
}
const PaymentCard = ({
    title,
    paymentMethodData,
    onAddPaymentMethod,
}: IPaymentProps) => {
    const classes = useStyles();
    const [paymentMethod, setPaymentMethod] = React.useState("");

    useEffect(() => {
        setPaymentMethod(getSelected(paymentMethodData));
    }, [paymentMethodData]);

    const _isDarkMod = useRecoilValue(isDarkmode);

    const handleChange = (event: SelectChangeEvent) => {
        setPaymentMethod(event.target.value as string);
    };
    return (
        <Box>
            <Typography variant="h6" sx={{ pb: 3 }}>
                {title}
            </Typography>

            <FormControl variant="outlined" className={classes.formControl}>
                <InputLabel
                    shrink
                    variant="outlined"
                    htmlFor="outlined-age-always-notched"
                >
                    PAYMENT METHOD
                </InputLabel>
                <Select
                    value={paymentMethod}
                    variant="outlined"
                    onChange={handleChange}
                    IconComponent={() => null}
                    sx={{
                        "& legend": { width: "135px" },
                    }}
                    input={
                        <OutlinedInput
                            notched
                            label="NODE TYPE"
                            name="node_type"
                            id="outlined-age-always-notched"
                        />
                    }
                    MenuProps={{
                        disablePortal: false,
                        PaperProps: {
                            sx: {
                                boxShadow:
                                    "0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)",
                                borderRadius: "4px",
                            },
                        },
                    }}
                    className={classes.selectStyle}
                >
                    {paymentMethodData.map(({ id, value, label }: any) => (
                        <MenuItem
                            key={id}
                            value={value}
                            sx={{
                                m: 0,
                                p: "6px 16px",
                            }}
                        >
                            <Stack direction="row" spacing={1}>
                                <Typography variant="body1">{label}</Typography>
                                <Button
                                    variant="text"
                                    onClick={onAddPaymentMethod}
                                    sx={{
                                        display:
                                            value === "no_payment_method_Set"
                                                ? "block"
                                                : "none",
                                        textTransform: "none",
                                        color: colors.primaryMain,
                                        ":hover": {
                                            color: theme =>
                                                theme.palette.text.primary,
                                        },
                                    }}
                                >
                                    <Typography variant="body1">
                                        {"Enter now"}
                                    </Typography>
                                </Button>
                            </Stack>
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>
            <Typography
                variant="caption"
                sx={{ color: _isDarkMod ? colors.white : colors.black54 }}
            >
                *Automatically charged to card EOD on the last day of the
                billing cycle
            </Typography>
        </Box>
    );
};

export default PaymentCard;

import React from "react";
import {
    Typography,
    OutlinedInput,
    Box,
    FormControl,
    SelectChangeEvent,
    Divider,
    Button,
    MenuItem,
    Select,
    InputLabel,
} from "@mui/material";
import { makeStyles } from "@mui/styles";

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

interface IPaymentProps {
    title: string;
    handlePaymentMethod: Function;
    paymentMethodData: any;
    onAddPaymentMethod: Function;
}
const PaymentCard = ({
    title,
    handlePaymentMethod,
    paymentMethodData,
    onAddPaymentMethod,
}: IPaymentProps) => {
    const classes = useStyles();
    const [paymentMethod, setPaymentMethod] = React.useState("None set up");

    const handleChange = (event: SelectChangeEvent) => {
        setPaymentMethod(event.target.value as string);
        handlePaymentMethod(event.target.value as string);
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
                    {paymentMethodData.map(({ id, value }: any) => (
                        <MenuItem
                            key={id}
                            value={value}
                            sx={{
                                m: 0,
                                p: "6px 16px",
                            }}
                        >
                            <Typography variant="body1">{value}</Typography>
                        </MenuItem>
                    ))}
                    <Divider />
                    <MenuItem
                        onClick={e => {
                            onAddPaymentMethod();
                            e.stopPropagation();
                        }}
                    >
                        <Button
                            variant="text"
                            sx={{
                                typography: "body1",
                                textTransform: "none",
                            }}
                        >
                            + Add payment method
                        </Button>
                    </MenuItem>
                </Select>
            </FormControl>
            <Typography variant="caption" sx={{ color: colors.black54 }}>
                *Automatically charged to card EOD on the last day of the
                billing cycle
            </Typography>
        </Box>
    );
};

export default PaymentCard;

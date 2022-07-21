import {
    Button,
    Dialog,
    IconButton,
    Typography,
    DialogActions,
    DialogContent,
    Stack,
    Divider,
    Box,
    DialogTitle,
    Grid,
} from "@mui/material";
import colors from "../../../theme/colors";
import CloseIcon from "@mui/icons-material/Close";
import React from "react";
import { INVOICE_LABEL_LIST } from "../../../constants";

type BasicDialogProps = {
    isOpen: boolean;
    handlePrint: any;
    isClosable?: boolean;
    handleCloseAction: any;
    ref: any;
};
const Logo = React.lazy(() =>
    import("../../../assets/svg").then(module => ({
        default: module.Logo,
    }))
);

const InvoiceViewerDialog = ({
    isOpen,
    handleCloseAction,
    ref,
    isClosable = true,
    handlePrint,
}: BasicDialogProps) => {
    return (
        <Dialog
            fullWidth
            open={isOpen}
            maxWidth="sm"
            onClose={handleCloseAction}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
            onBackdropClick={() => isClosable && handleCloseAction()}
            ref={ref}
        >
            <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
            >
                <DialogTitle>
                    <Button variant="text" onClick={handlePrint}>
                        {"Print invoice"}
                    </Button>
                </DialogTitle>
                <IconButton
                    onClick={handleCloseAction}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton>
            </Stack>

            <DialogContent>
                <Grid container spacing={1} alignItems="center">
                    <Grid item xs={5}>
                        <Logo
                            style={{ position: "relative", bottom: 100 }}
                            width={"100%"}
                            height={"36px"}
                            color={colors.primaryMain}
                        />
                    </Grid>
                    <Grid
                        container
                        item
                        xs={7}
                        justifyContent="space-between"
                        spacing={6}
                    >
                        <Grid item xs={8}>
                            <Stack direction="column">
                                {INVOICE_LABEL_LIST.map((i: any) => (
                                    <Box key={i.id}>
                                        <Typography
                                            variant="body2"
                                            sx={{ color: colors.black54 }}
                                        >
                                            {i.label}
                                        </Typography>
                                    </Box>
                                ))}
                            </Stack>
                        </Grid>
                        <Grid item xs={4}>
                            <Stack direction="column">
                                <Typography variant="body2">
                                    023-3932
                                </Typography>
                                <Typography variant="body2">
                                    023-3932
                                </Typography>
                                <Typography variant="body2">
                                    July 11 22
                                </Typography>
                                <Typography variant="body2">
                                    PamojaNet
                                </Typography>
                                <Typography variant="body2">Active</Typography>
                            </Stack>
                        </Grid>
                        <Grid item xs={6}>
                            <Typography
                                variant="body2"
                                sx={{ color: colors.black54 }}
                            >
                                {"Bill from "}
                            </Typography>
                            <Typography
                                variant="body2"
                                sx={{ fontWeight: "600" }}
                            >
                                {"Ukama"}
                            </Typography>
                            <Typography variant="body2">
                                {"Address line 0022"}
                            </Typography>
                        </Grid>

                        <Grid item xs={6}>
                            <Typography
                                variant="body2"
                                sx={{ color: colors.black54 }}
                            >
                                {"Paid by "}
                            </Typography>
                            <Typography
                                variant="body2"
                                sx={{ fontWeight: "600" }}
                            >
                                {"Customer"}
                            </Typography>
                            <Typography variant="body2">
                                {"Visa ending soon"}
                            </Typography>
                        </Grid>
                    </Grid>

                    <Box sx={{ width: "100%", mt: 2 }}>
                        <Divider />
                        <Grid container spacing={3}>
                            <Grid item xs={12}>
                                <Typography
                                    variant="body2"
                                    sx={{
                                        color: colors.black54,
                                        fontWeight: "600",
                                    }}
                                >
                                    {"DESCRIPTION"}
                                </Typography>
                            </Grid>
                            <Grid item xs={8}>
                                <Stack direction="column">
                                    <Typography
                                        variant="body2"
                                        sx={{ fontWeight: "500" }}
                                    >
                                        {"Standard roaming"}
                                    </Typography>
                                    <Typography
                                        variant="body2"
                                        sx={{
                                            color: colors.black54,
                                        }}
                                    >
                                        {
                                            "5 GB; $5/GB; from June 11 2022 to July 11 2022 "
                                        }
                                    </Typography>
                                </Stack>
                            </Grid>
                            <Grid
                                item
                                xs={4}
                                container
                                justifyContent="flex-end"
                            >
                                USD $20
                            </Grid>
                            <Grid item xs={6}>
                                <Stack direction="column">
                                    <Typography
                                        variant="body2"
                                        sx={{ fontWeight: "500" }}
                                    >
                                        {"Taxes & regulatory feeds"}
                                    </Typography>
                                    <Typography
                                        variant="body2"
                                        sx={{
                                            color: colors.black54,
                                        }}
                                    >
                                        {"See more "}
                                    </Typography>
                                </Stack>
                            </Grid>
                            <Grid
                                item
                                xs={6}
                                container
                                justifyContent="flex-end"
                            >
                                USD $2
                            </Grid>
                            <Grid item xs={6} />
                            <Grid item xs={6}>
                                <Stack
                                    direction="row"
                                    justifyContent="space-between"
                                >
                                    <Typography variant="body1" color="initial">
                                        {"Total"}
                                    </Typography>{" "}
                                    <Typography variant="body1" color="initial">
                                        $20
                                    </Typography>{" "}
                                </Stack>
                                <Divider />
                                <Stack
                                    direction="row"
                                    justifyContent="space-between"
                                >
                                    <Typography
                                        variant="body1"
                                        sx={{ fontWeight: "600" }}
                                    >
                                        {"Amount Due"}
                                    </Typography>{" "}
                                    <Typography
                                        variant="body1"
                                        sx={{ fontWeight: "600" }}
                                    >
                                        $20
                                    </Typography>{" "}
                                </Stack>
                            </Grid>
                        </Grid>
                    </Box>
                </Grid>
            </DialogContent>

            <DialogActions>
                <Stack direction={"row"} alignItems="center" spacing={2}>
                    <Button
                        variant="text"
                        color={"primary"}
                        onClick={handleCloseAction}
                    >
                        {"Close"}
                    </Button>
                </Stack>
            </DialogActions>
        </Dialog>
    );
};

export default InvoiceViewerDialog;

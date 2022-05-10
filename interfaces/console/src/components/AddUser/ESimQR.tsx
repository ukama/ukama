import {
    Grid,
    Accordion,
    Typography,
    AccordionSummary,
    AccordionDetails,
} from "@mui/material";
import { colors } from "../../theme";
import QrCodeIcon from "@mui/icons-material/QrCode2";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";

interface IESimQR {
    description: String;
}

const ESimQR = ({ description }: IESimQR) => {
    return (
        <Grid container mb={2}>
            <Grid item xs={12}>
                <Typography variant="body1">{description}</Typography>
            </Grid>
            <Grid item xs={12}>
                <Accordion sx={{ boxShadow: "none" }}>
                    <AccordionSummary
                        expandIcon={<ExpandMoreIcon color="primary" />}
                        sx={{
                            p: 0,
                            m: 0,
                            justifyContent: "flex-start",
                            "& .MuiAccordionSummary-content": {
                                flexGrow: 0.02,
                            },
                        }}
                    >
                        <Typography
                            fontWeight={500}
                            variant="caption"
                            color={colors.primaryMain}
                        >
                            SHOW ESIM QR CODE
                        </Typography>
                    </AccordionSummary>
                    <AccordionDetails
                        sx={{ p: 0, display: "flex", justifyContent: "center" }}
                    >
                        <QrCodeIcon sx={{ height: 164, width: 164 }} />
                    </AccordionDetails>
                </Accordion>
            </Grid>
        </Grid>
    );
};

export default ESimQR;

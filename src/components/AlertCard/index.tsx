import * as React from "react";
import { RoundedCard } from "../../styles";
import {
    Box,
    Typography,
    CardActions,
    IconButton,
    Grid,
    Button,
    CardHeader,
    Card,
} from "@mui/material";
import FavoriteIcon from "@mui/icons-material/Favorite";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ShareIcon from "@mui/icons-material/Share";
import { useTranslation } from "react-i18next";
import { CloudOffIcon } from "../../assets/svg";
import { colors } from "../../theme";

import "../../i18n/i18n";
type AlertCardProps = {
    Icon: React.ReactElement;
    action: any;
    title: string;
    subheader: string;
};

const AlertCard = ({ Icon, action, title, subheader }: AlertCardProps) => {
    const { t } = useTranslation();
    return (
        <>
            <RoundedCard>
                <Box mb={2}>
                    <Typography variant="h6">{t("ALERT.Title")}</Typography>
                </Box>

                <Card elevation={2}>
                    <CardActions>
                        <CloudOffIcon />
                        <div style={{ padding: 7 }}>{title}</div>

                        <Typography variant="caption" color={colors.empress}>
                            08/16/21 1PM
                        </Typography>
                    </CardActions>
                    <Box pl={5.5}>
                        <Typography variant="subtitle2" color="initial">
                            Short description of alert.
                        </Typography>
                    </Box>
                </Card>
            </RoundedCard>
        </>
    );
};
export default AlertCard;

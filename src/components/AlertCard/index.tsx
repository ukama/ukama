import * as React from "react";
import { RoundedCard } from "../../styles";
import { Box, Typography, Grid, CardHeader, Card } from "@mui/material";
import { useTranslation } from "react-i18next";
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
                    <Grid item xs={12} container justifyContent="flex-start">
                        <Typography variant="h6">{t("ALERT.Title")}</Typography>
                    </Grid>
                </Box>

                <Card elevation={2}>
                    <CardHeader
                        avatar={Icon}
                        action={action}
                        title={title}
                        subheader={subheader}
                    />
                </Card>
            </RoundedCard>
        </>
    );
};
export default AlertCard;

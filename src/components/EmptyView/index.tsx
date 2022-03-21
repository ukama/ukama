import { ElementType } from "react";
import { Stack, Typography } from "@mui/material";

interface IEmptyView {
    icon: ElementType;
    title: string;
    size?: "small" | "medium" | "large";
}
const EmptyView = ({ title, icon: Icon, size = "medium" }: IEmptyView) => {
    return (
        <Stack
            spacing={1}
            sx={{
                height: "100%",
                width: "100%",
                alignItems: "center",
            }}
        >
            <Typography variant="body1">{title}</Typography>
            <Icon
                fontSize={size}
                color="textPrimary"
                style={{ opacity: 0.6 }}
            />
        </Stack>
    );
};

export default EmptyView;

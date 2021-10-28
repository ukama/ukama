import { colors } from "../../theme";
import { RoundedCard } from "../../styles";
import { SelectItemType } from "../../types";
import { Grid, MenuItem, Select, Typography } from "@mui/material";

type StatusCardProps = {
    Icon: any;
    title: string;
    option: string;
    subtitle1: string;
    subtitle2: string;
    handleSelect: Function;
    options: SelectItemType[];
};

const StatusCard = ({
    Icon,
    title,
    option,
    options,
    subtitle1,
    subtitle2,
    handleSelect,
}: StatusCardProps) => (
    <RoundedCard>
        <Grid spacing={2} container direction="row" justifyContent="center">
            <Grid item display="flex" alignItems="center">
                <Icon />
            </Grid>
            <Grid xs={12} item sm container spacing={1} direction="column">
                <Grid
                    sm
                    item
                    container
                    spacing={2}
                    display="flex"
                    direction="row"
                    alignItems="center"
                >
                    <Grid item xs={7}>
                        <Typography variant="body1">{title}</Typography>
                    </Grid>
                    <Grid item xs={5} display="flex" alignItems="flex-end">
                        <Select
                            value={option}
                            disableUnderline
                            variant="standard"
                            sx={{ width: "128px", color: colors.empress }}
                            onChange={e => handleSelect(e.target.value)}
                        >
                            {options.map(
                                ({ id, label, value }: SelectItemType) => (
                                    <MenuItem key={id} value={value}>
                                        {label}
                                    </MenuItem>
                                )
                            )}
                        </Select>
                    </Grid>
                </Grid>
                <Grid item sm container>
                    <Grid item xs={12} container>
                        <Typography variant="h3" paddingRight="6px">
                            {subtitle1}
                        </Typography>
                        <Typography
                            variant="body1"
                            color={colors.empress}
                            sx={{ display: "flex", alignItems: "center" }}
                        >
                            {subtitle2}
                        </Typography>
                    </Grid>
                </Grid>
            </Grid>
        </Grid>
    </RoundedCard>
);

export default StatusCard;

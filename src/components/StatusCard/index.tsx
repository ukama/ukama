import { colors } from "../../theme";
import { SelectItemType } from "../../types";
import { RoundedCard, SkeletonRoundedCard } from "../../styles";
import { Grid, MenuItem, Select, Typography } from "@mui/material";

type StatusCardProps = {
    Icon: any;
    title: string;
    option: string;
    loading: boolean;
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
    loading,
    subtitle1,
    subtitle2,
    handleSelect,
}: StatusCardProps) => (
    <>
        {loading ? (
            <SkeletonRoundedCard variant="rectangular" height={104} />
        ) : (
            <RoundedCard>
                <Grid
                    spacing={2}
                    container
                    direction="row"
                    justifyContent="center"
                >
                    <Grid item display="flex" alignItems="center">
                        <Icon />
                    </Grid>
                    <Grid xs={12} item sm container direction="column">
                        <Grid
                            sm
                            item
                            container
                            spacing={2}
                            display="flex"
                            direction="row"
                            alignItems="center"
                        >
                            <Grid item xs={8}>
                                <Typography variant="subtitle2">
                                    {title}
                                </Typography>
                            </Grid>
                            <Grid
                                item
                                xs={4}
                                display="flex"
                                justifyContent="flex-end"
                            >
                                <Select
                                    value={option}
                                    disableUnderline
                                    variant="standard"
                                    sx={{
                                        width: "100%",
                                        color: colors.empress,
                                    }}
                                    onChange={e => handleSelect(e.target.value)}
                                >
                                    {options.map(
                                        ({
                                            id,
                                            label,
                                            value,
                                        }: SelectItemType) => (
                                            <MenuItem key={id} value={value}>
                                                <Typography variant="subtitle2">
                                                    {label}
                                                </Typography>
                                            </MenuItem>
                                        )
                                    )}
                                </Select>
                            </Grid>
                        </Grid>
                        <Grid item sm container>
                            <Grid item xs={12} container>
                                <Typography variant="h5" paddingRight="6px">
                                    {subtitle1}
                                </Typography>
                                <Typography
                                    variant="body1"
                                    color={colors.empress}
                                    sx={{
                                        display: "flex",
                                        alignItems: "center",
                                    }}
                                >
                                    {subtitle2}
                                </Typography>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </RoundedCard>
        )}
    </>
);

export default StatusCard;

import { makeStyles } from "@mui/styles";
import { RoundedCard } from "../../styles";
import { SelectItemType } from "../../types";
import LoadingWrapper from "../LoadingWrapper";
import { Grid, MenuItem, Select, Theme, Typography } from "@mui/material";

const useStyles = makeStyles<Theme>(theme => ({
    selectStyle: {
        width: "108px",
        textAlign: "end",
        "& p": {
            color: theme?.palette?.text?.secondary,
            fontWeight: 500,
            fontSize: "14px",
            lineHeight: "157%",
        },
        "& .MuiSelect-iconStandard": {
            paddingBottom: "4px",
        },
        "& .MuiSelect-iconOpen": {
            paddingBottom: "0px",
        },
    },
}));

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
    subtitle1 = "0",
    subtitle2 = "",
    handleSelect,
}: StatusCardProps) => {
    const classes = useStyles();
    return (
        <LoadingWrapper height={100} isLoading={loading}>
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
                            <Grid item xs={7}>
                                <Typography variant="subtitle2">
                                    {title}
                                </Typography>
                            </Grid>
                            <Grid
                                item
                                xs={5}
                                display="flex"
                                justifyContent="flex-end"
                            >
                                <Select
                                    value={option}
                                    disableUnderline
                                    variant="standard"
                                    className={classes.selectStyle}
                                    MenuProps={{
                                        sx: {
                                            maxHeight: "194px",
                                        },
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
                                                <Typography variant="body1">
                                                    {label}
                                                </Typography>
                                            </MenuItem>
                                        )
                                    )}
                                </Select>
                            </Grid>
                        </Grid>
                        <Grid item sm container>
                            <Grid item xs={12} container alignItems="flex-end">
                                <Typography
                                    variant="h5"
                                    paddingRight="6px"
                                    sx={{
                                        position: "relative",
                                        bottom: "2px",
                                    }}
                                >
                                    {subtitle1}
                                </Typography>
                                {title === "Data Usage" && (
                                    <Typography
                                        variant="body1"
                                        paddingRight="4px"
                                        sx={{
                                            position: "relative",
                                            bottom: "5px",
                                        }}
                                    >
                                        GB
                                    </Typography>
                                )}

                                <Typography
                                    variant="body1"
                                    color="textSecondary"
                                    sx={{
                                        display: "flex",
                                        alignItems: "center",
                                        position: "relative",
                                        bottom: "6px",
                                    }}
                                >
                                    {subtitle2}
                                </Typography>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </RoundedCard>
        </LoadingWrapper>
    );
};
export default StatusCard;

import { Theme } from "@mui/material";
import { makeStyles } from "@mui/styles";

const useStyles = makeStyles<Theme>(() => ({
    selectStyle: () => ({
        width: "fit-content",
    }),
}));

const SelectDisplayProps = {
    style: {
        fontWeight: 600,
        display: "flex",
        fontSize: "20px",
        marginLeft: "4px",
        alignItems: "center",
        minWidth: "fit-content",
    },
};

const PaperProps = {
    boxShadow:
        "0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)",
    borderRadius: "4px",
};

export { useStyles, SelectDisplayProps, PaperProps };

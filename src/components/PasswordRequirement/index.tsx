import { Typography } from "@mui/material";

type RequirementProps = {
    isvalid: boolean;
    validMessage: any;
    invalidMessage: any;
    label: string;
};
const PasswordRequirement = ({
    isvalid,
    validMessage,
    invalidMessage,
    label,
}: RequirementProps) => {
    return (
        <>
            <Typography variant="caption">
                {!isvalid ? invalidMessage : validMessage}

                {label}
            </Typography>
        </>
    );
};

export default PasswordRequirement;

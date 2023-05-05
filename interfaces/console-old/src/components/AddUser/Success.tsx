import { Typography } from "@mui/material";

interface ISuccess {
    description: string;
}

const Success = ({ description }: ISuccess) => {
    return (
        <Typography variant="body1" sx={{ mb: 2 }}>
            {description}
        </Typography>
    );
};

export default Success;

import { HorizontalContainerJustify } from '@/styles/global';
import { Button, Typography } from '@mui/material';

type TableHeaderProps = {
  title: string;
  buttonTitle?: string;
  showSecondaryButton: boolean;
  handleButtonAction?: any;
};

const TableHeader = ({
  title,
  buttonTitle,
  handleButtonAction,
  showSecondaryButton,
}: TableHeaderProps) => {
  return (
    <HorizontalContainerJustify>
      <Typography variant="body2" fontWeight={600}>
        {title}
      </Typography>
      {showSecondaryButton && (
        <Button
          variant="outlined"
          sx={{ width: '144px' }}
          onClick={() => handleButtonAction()}
        >
          {buttonTitle}
        </Button>
      )}
    </HorizontalContainerJustify>
  );
};

export default TableHeader;

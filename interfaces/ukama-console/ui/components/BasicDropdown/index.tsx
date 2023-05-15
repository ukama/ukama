import colors from '@/styles/theme/colors';
import { SelectItemType } from '@/types';
import { FormControl, MenuItem, Select } from '@mui/material';

interface IBasicDropdown {
  value: string;
  isLoading?: boolean;
  handleOnChange: Function;
  networkList: SelectItemType[];
}
const BasicDropdown = ({
  value,
  networkList,
  handleOnChange,
}: IBasicDropdown) => (
  <FormControl sx={{ width: '100%' }} size="small" >
    <Select
      value={value}
      disableUnderline
      variant="standard"
      onChange={(e) => handleOnChange(e.target.value)}
      sx={{
        p: 0,

        color: colors.primaryMain,
      }}
      SelectDisplayProps={{
        style: {
          fontWeight: 400,
          fontSize: '16px',
        },
      }}
    >
      {networkList?.map((item: SelectItemType) => (
        <MenuItem key={item.id} value={item.value}>
          {item.label}
        </MenuItem>
      ))}
    </Select>
  </FormControl>
);

export default BasicDropdown;

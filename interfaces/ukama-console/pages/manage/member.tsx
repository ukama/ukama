import { MANAGE_TABLE_COLUMN } from '@/constants';
import { OrgMembersResDto } from '@/generated';
import { EmptyView, SimpleDataTable } from '@/ui/components';
import PageContainerHeader from '@/ui/components/PageContainerHeader';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import { Paper } from '@mui/material';

interface IMember {
  data: any;
  search: string;
  setSearch: (value: string) => void;
  handleButtonAction: () => void;
}

const structureData = (data: OrgMembersResDto) =>
  data.members?.map((member) => ({
    name: member.user.name,
    email: member.user.email,
    role: 'member',
    uuid: member.uuid,
  }));

const Member = ({ data, search, setSearch, handleButtonAction }: IMember) => {
  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        overflow: 'hidden',
        borderRadius: '5px',
        height: 'calc(100vh - 200px)',
      }}
    >
      <PageContainerHeader
        search={search}
        title={'My members'}
        buttonTitle={'Invite member'}
        onSearchChange={(e: string) => setSearch(e)}
        handleButtonAction={handleButtonAction}
      />
      <br />
      {data && data.length > 0 ? (
        <SimpleDataTable
          dataKey="uuid"
          dataset={structureData(data)}
          columns={MANAGE_TABLE_COLUMN}
        />
      ) : (
        <EmptyView icon={PeopleAltIcon} title="No members yet!" />
      )}
    </Paper>
  );
};

export default Member;

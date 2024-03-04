import React from 'react';
import {
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
} from '@mui/material';
import { InvoiceDto, PackageDto, Subscription } from '@/generated';
import colors from '@/styles/theme/colors';
import LoadingWrapper from '../LoadingWrapper';

interface suBbillingProps {
  billingCycle: InvoiceDto[];
  dataPlans: PackageDto[];
  billingCycleLoading: boolean;
  dataPlanLoading: boolean;
}

const BillingCycle: React.FC<suBbillingProps> = ({
  billingCycle = [],
  dataPlans = [],
  billingCycleLoading,
  dataPlanLoading,
}) => (
  <TableContainer>
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>
            <strong style={{ fontWeight: 'bold' }}> Billing cycle</strong>
          </TableCell>
          <TableCell>
            <strong style={{ fontWeight: 'bold' }}> Data usage</strong>
          </TableCell>
          <TableCell>
            <strong style={{ fontWeight: 'bold' }}> Data plan</strong>
          </TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        <LoadingWrapper
          radius="small"
          width={'100%'}
          isLoading={billingCycleLoading || dataPlanLoading}
          cstyle={{
            overflow: 'auto',
            backgroundColor: false ? colors.white : 'transparent',
          }}
        >
          {billingCycle &&
            billingCycle.map((cycleData, index) => (
              <TableRow key={index}>
                <TableCell>{cycleData.period}</TableCell>
                <TableCell>
                  {cycleData?.rawInvoice?.metadata?.find(
                    (item: any) => item.key === 'bytes_used',
                  )?.value || 'N/A'}
                </TableCell>
                <TableCell>
                  {cycleData &&
                    cycleData.rawInvoice.subscriptions?.map(
                      (subscription: Subscription, subIndex) => {
                        const matchingPlan = dataPlans.find(
                          (plan) => plan.uuid === subscription.planCode,
                        );
                        return (
                          <div key={subIndex}>
                            {matchingPlan?.name || 'Unknown Plan'}
                          </div>
                        );
                      },
                    )}
                </TableCell>
              </TableRow>
            ))}
        </LoadingWrapper>
      </TableBody>
    </Table>
  </TableContainer>
);

export default BillingCycle;

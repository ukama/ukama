// interface IBillingDialog {
//   isOpen: boolean;
//   initPaymentFlow: boolean;
//   handleCloseAction: Function;
//   handleSuccessAction: Function;
// }

const BillingDialog = () => {
  // const setSnackbarMessage = useSetRecoilState(snackbarMessage);

  // const [
  //   attachPaymentWithCustomer,
  //   { loading: attachPaymentWithCustomerLoading },
  // ] = useAttachPaymentWithCustomerMutation({
  //   onCompleted: () => {
  //     handleFlowChange(flow + 1);
  //   },
  //   onError: () => {
  //     setSnackbarMessage({
  //       id: 'pm-link-failed',
  //       message: 'Failed to link payment method',
  //       type: 'error',
  //       show: true,
  //     });
  //   },
  // });

  // const [flow, setFlow] = useState(initPaymentFlow ? 2 : 0);
  // const handleFlowChange = (i: number) => {
  //   if (flow === 2) handleSuccessAction();
  //   setFlow(i);
  // };

  // const handleClose = () => {
  //   setFlow(0);
  //   handleCloseAction();
  // };

  // const handleIsPaymentSuccess = (id: string) => {
  //   if (id) {
  //     attachPaymentWithCustomer({
  //       variables: { paymentId: id },
  //     });
  //   }
  // };

  return (
    <></>
    // <Dialog
    //   fullWidth
    //   open={isOpen}
    //   maxWidth="sm"
    //   aria-labelledby="alert-dialog-title"
    //   aria-describedby="alert-dialog-description"
    // >
    //   <Stack direction="row" alignItems="center" justifyContent="space-between">
    //     <DialogTitle>{BillingDialogList[flow].title}</DialogTitle>
    //     <IconButton
    //       onClick={() => handleClose()}
    //       sx={{ position: 'relative', right: 8 }}
    //     >
    //       <CloseIcon />
    //     </IconButton>
    //   </Stack>

    //   <DialogContent>
    //     <Typography variant="body1">
    //       {BillingDialogList[flow].description}
    //     </Typography>
    //     {flow === 0 && <ChoosePlan />}
    //     {flow === 1 && <CustomizePref />}
    //     {flow === 2 && (
    //       <PaymentForm
    //         isPaymentOnly={initPaymentFlow}
    //         loading={attachPaymentWithCustomerLoading}
    //         handleIsPaymentSuccess={handleIsPaymentSuccess}
    //         handleBackAction={() => handleFlowChange(flow - 1)}
    //       />
    //     )}
    //     {flow === 3 && <></>}
    //   </DialogContent>

    //   {flow !== 2 && (
    //     <DialogActions>
    //       <HorizontalContainerJustify>
    //         <Button
    //           variant="text"
    //           color={'primary'}
    //           sx={{
    //             visibility:
    //               flow !== 0 && flow !== 3 && !initPaymentFlow
    //                 ? 'visible'
    //                 : 'hidden',
    //           }}
    //           onClick={() => handleFlowChange(flow - 1)}
    //         >
    //           Back
    //         </Button>

    //         {flow !== 3 && (
    //           <Stack spacing={2} direction={'row'} alignItems="center">
    //             <Button
    //               variant="contained"
    //               onClick={() => handleFlowChange(flow + 1)}
    //             >
    //               Next
    //             </Button>
    //           </Stack>
    //         )}
    //       </HorizontalContainerJustify>
    //     </DialogActions>
    //   )}
    // </Dialog>
  );
};

export default BillingDialog;

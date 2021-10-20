const AddEmail = (message: string, email: string) =>
    message.replace("/EmailTag", email);

export { AddEmail };

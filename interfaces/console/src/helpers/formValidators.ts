import * as Yup from "yup";
const nameValidationRule = Yup.string().required("Name is required.");
const emailValidatonRule = Yup.string()
    .required("Email is required.")
    .email("Please enter a valid email.");
// const iccidValidator = Yup.string().required("Iccid is required.");
const iccidValidator = Yup.string()
    .required()
    .matches(/^[0-9]+$/, "Must be only digits")
    .min(18, "Iccid must be 18 digits")
    .max(18, "Iccid must be 18 digits");
const securitycodeValidator = Yup.string().required(
    "Security code is required."
);
const ESIM_FORM_SCHEMA = {
    email: emailValidatonRule,
    name: nameValidationRule,
};
const PHYSICAL_SIM_FORM_SCHEMA = {
    iccid: iccidValidator,
    securityCode: securitycodeValidator,
};
export { ESIM_FORM_SCHEMA, PHYSICAL_SIM_FORM_SCHEMA };

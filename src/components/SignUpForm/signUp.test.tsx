import { mount, ReactWrapper } from "enzyme";
import { SignUpForm } from "../index";
import { act } from "react-dom/test-utils";
import Adapter from "@wojtekmaj/enzyme-adapter-react-17";
// eslint-disable-next-line
import Enzyme from "enzyme";
import { TextField } from "@mui/material";
Enzyme.configure({ adapter: new Adapter() });
describe("SignupForm", () => {
    let component: ReactWrapper;
    beforeEach(() => {
        component = mount(<SignUpForm />);
    });
    it("should update email field on change", async () => {
        const emailInput = component.find("input[name='email']");
        await act(async () => {
            emailInput.simulate("change", {
                persist: () => {},
                target: {
                    name: "email",
                    value: "email@gmail.com",
                },
            });
        });
        expect(emailInput.html()).toMatch("email@gmail.com");
    });

    it("should update display error onclick without filling the form ", async () => {
        const inputField = component.find(TextField).at(1);
        const submitButton = component
            .find(".MuiButton-root#signUpButton")
            .at(1);
        submitButton.simulate("click");
        expect(inputField.props().helperText).toBeUndefined();
    });

    it("should update password field on change and should display the password indicator compoenent", async () => {
        const tree = mount(<SignUpForm />);

        const passwordInput = tree.find("input[name='password']");
        const passwordRequirement = component.find(".MuiSvgIcon-root");
        await act(async () => {
            passwordInput.simulate("change", {
                persist: () => {},
                target: {
                    name: "password",
                    value: "Pass12",
                },
            });
        });
        expect(passwordInput.html()).toMatch("Pass12");
        expect(passwordRequirement).toBeTruthy();
    });
});

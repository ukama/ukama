import ReactDOM from "react-dom";
import { RecoilRoot } from "recoil";
import App from "./App";
it("Dummy test ", () => {
    const div = document.createElement("div");
    ReactDOM.render(
        <RecoilRoot>
            <App />
        </RecoilRoot>,
        div
    );
    ReactDOM.unmountComponentAtNode(div);
});

import './main.scss';

import './components/app/app';
import './components/board/board';
import './components/column/column';
import './components/card/card';

import type {HtmxBeforeSwapDetail} from "./types/htmx";

function enforceComponentSwap(evt: CustomEvent<HtmxBeforeSwapDetail>) {
    const detail = evt.detail;
    let elt = detail.elt;
    let root = elt.getRootNode();

    if (root instanceof ShadowRoot) {
        detail.target = root.host as HTMLElement;
        detail.swapOverride = "outerHTML";
    }
}

document.body.addEventListener("htmx:beforeSwap", enforceComponentSwap as EventListener);
document.body.addEventListener("htmx:oobBeforeSwap", enforceComponentSwap as EventListener);



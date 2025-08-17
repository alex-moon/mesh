import './main.scss';

import './components/app/app';
import './components/board/board';
import type {HtmxBeforeSwapDetail} from "./types/htmx";

document.body.addEventListener("htmx:beforeSwap", function(evt: CustomEvent<HtmxBeforeSwapDetail>) {
    const detail = evt.detail;
    let elt = detail.elt;
    let root = elt.getRootNode();
    if (root instanceof ShadowRoot) {
        detail.target = root.host as HTMLElement;
        detail.swapOverride = 'outerHTML';
    }
} as EventListener);

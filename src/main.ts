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

function findInShadow(root: any, id: string): any {
    const element = root.getElementById?.(id);
    if (element) {
        return element;
    }
    const allElements = root.querySelectorAll('*');
    for (let element of allElements) {
        if (element.shadowRoot) {
            const found = findInShadow(element.shadowRoot, id);
            if (found) {
                return found;
            }
        }
    }
    return null;
}

function enableOobSwap(evt: CustomEvent<any>) {
    const id = evt.detail.content.id;
    const found = findInShadow(document, id);
    if (found) {
        found.outerHTML = evt.detail.content.outerHTML;
        evt.preventDefault();
    }
}

document.body.addEventListener("htmx:beforeSwap", enforceComponentSwap as EventListener);
document.body.addEventListener("htmx:oobErrorNoTarget", enableOobSwap as EventListener);


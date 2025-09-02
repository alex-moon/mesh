import './main.scss';

import 'htmx.org';
import 'htmx-ext-sse';

import './components/app/app';
import './components/board/board';
import './components/column/column';
import './components/card/card';

import {findInShadow} from "./ts/shadow";

function enforceComponentSwap(evt: CustomEvent<any>) {
    const detail = evt.detail;
    let elt = detail.elt;
    let root = elt.getRootNode();

    if (root instanceof ShadowRoot) {
        detail.target = root.host as HTMLElement;
        detail.swapOverride = "outerHTML";
    }
}

function handleOobSwap(evt: CustomEvent<any>) {
    const template = document.createElement('template');
    template.innerHTML = evt.detail.data.trim();
    const content = template.content.firstElementChild;

    if (content && content.hasAttribute('hx-swap-oob')) {
        const id = content.id;
        const found = findInShadow(document, id);

        if (found) {
            found.outerHTML = content.outerHTML;
            evt.preventDefault();
        }
    }
}

// function handleSse(evt: CustomEvent<any>) {
//     console.log('SSE message', evt.detail);
//     evt.preventDefault();
// }

function enableOobSwap(evt: CustomEvent<any>) {
    console.log('OOB swap', evt);
    const id = evt.detail.content.id;
    const found = findInShadow(document, id);
    if (found) {
        setTimeout(() => {
            found.outerHTML = evt.detail.content.outerHTML;
        });
    } else {
        console.error('Target not found in shadow DOM', evt);
    }
}

document.body.addEventListener("htmx:beforeSwap", enforceComponentSwap as EventListener);
document.body.addEventListener("htmx:sseBeforeMessage", handleOobSwap as EventListener);
document.body.addEventListener("htmx:oobErrorNoTarget", enableOobSwap as EventListener);

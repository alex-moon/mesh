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

// SSE client for real-time collaboration
function initializeSSE() {
    const eventSource = new EventSource('/sse');
    
    eventSource.onopen = function() {
        console.log('SSE connection established for real-time collaboration');
    };
    
    eventSource.onmessage = function(event) {
        try {
            // Check if it's a connection message
            const data = JSON.parse(event.data);
            if (data.type === 'connected') {
                console.log('SSE client connected:', data.clientID);
                return;
            }
        } catch (e) {
            // Not JSON, likely HTML content - process as OOB update
            const html = event.data;
            
            // Create temporary container to parse HTML
            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = html;
            
            // Find elements with hx-swap-oob attribute (OOB updates)
            const oobElements = tempDiv.querySelectorAll('[hx-swap-oob]');
            
            oobElements.forEach((oobElement) => {
                const id = oobElement.id;
                if (id) {
                    // Find target element in shadow DOM
                    const target = findInShadow(document, id);
                    if (target) {
                        target.outerHTML = oobElement.outerHTML;
                        console.log('Applied SSE OOB update to element:', id);
                    }
                }
            });
        }
    };
    
    eventSource.onerror = function(error) {
        console.error('SSE connection error:', error);
    };
    
    eventSource.onclose = function() {
        console.log('SSE connection closed, attempting to reconnect...');
        setTimeout(initializeSSE, 5000); // Reconnect after 5 seconds
    };
}

// Initialize SSE when page loads
document.addEventListener('DOMContentLoaded', initializeSSE);


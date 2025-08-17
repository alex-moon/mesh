export interface HtmxBeforeSwapDetail {
    elt: HTMLElement;
    target: HTMLElement;
    swapStyle: string;
    swapOverride?: string;
    shouldSwap?: boolean;
    serverResponse?: string;
    xhr?: XMLHttpRequest;
    requestConfig?: Record<string, any>;
}

export interface HtmxApi {
    process(element: Element | Document | ShadowRoot): void;
    // Add other HTMX methods you might use
    ajax(verb: string, path: string, element?: Element): void;
    find(selector: string): Element | null;
    findAll(selector: string): NodeList;
    remove(element: Element): void;
    addClass(element: Element, className: string): void;
    removeClass(element: Element, className: string): void;
    toggleClass(element: Element, className: string): void;
    takeClass(element: Element, className: string): void;
}

declare global {
    interface Window {
        htmx: HtmxApi;
    }
}

// This export makes it a module (required for declare global to work)
export {};

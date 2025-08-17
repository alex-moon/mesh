import {createIcons} from "lucide";

export class MeshElement extends HTMLElement {
    protected icons: any;

    connectedCallback() {
        if (!this.shadowRoot) {
            const root = this.attachShadow({ mode: 'open' });
            const template = this.querySelector('template[shadowrootmode="open"]');
            if (template) {
                root.appendChild((template as any).content.cloneNode(true));
            }
        }
        if (window.htmx) {
            window.htmx.process(this);
            if (this.shadowRoot) {
                window.htmx.process(this.shadowRoot);
            }
        }
        this.bindListeners();
        this.createIcons();
    }

    protected createIcons() {
        if (this.icons) {
            createIcons({
                icons: this.icons,
                attrs: {
                    width: 16,
                    height: 16,
                },
                root: this.shadowRoot as any,
            });
        }
    }

    protected bindListeners() {
        const supportedEvents = [
            'click', 'change', 'input', 'submit', 'focus', 'blur',
            'mouseenter', 'mouseleave', 'keydown', 'keyup', 'keypress',
        ];

        supportedEvents.forEach(eventName => {
            const attribute = "mesh-" + eventName;
            this.all('[' + attribute + ']', el => {
                const methodName = el.getAttribute(attribute);
                if (!methodName) {
                    return;
                }
                const method = (this as any)[methodName];
                if (!method || typeof method !== 'function') {
                    console.error(`Method ${methodName} is not a function`);
                    return;
                }

                el.addEventListener(eventName, method.bind(this));
            });
        });
    }

    all(selector: string, cb: (el: HTMLElement) => void) {
        return this.shadowRoot!.querySelectorAll(selector).forEach(e => cb(e as HTMLElement));
    }

    show(selector: string) {
        this.all(selector, el => {
            el.classList.remove('hide');
        });
    }

    hide(selector: string) {
        this.all(selector, el => {
            el.classList.add('hide');
        });
    }
}

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
        this.bindListeners();
        this.bindFormHandlers();
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
        const supported = [
            'click', 'change', 'input', 'submit', 'focus', 'blur',
            'mouseenter', 'mouseleave', 'keydown', 'keyup', 'keypress',
        ];

        supported.forEach(eventName => {
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

    protected bindFormHandlers() {
        const supported = [
            'get', 'post', 'put', 'patch', 'delete',
        ];

        supported.forEach(verb => {
            const attribute = "mesh-" + verb;
            this.all('[' + attribute + ']', el => {
                const form = el as HTMLFormElement;
                form.addEventListener('submit', (event: Event) => {
                    event.preventDefault();
                    const method = verb.toUpperCase();
                    const url = form.getAttribute(attribute);

                    if (!url) {
                        console.error('No URL specified for form submission');
                        return;
                    }

                    const formData = new FormData(form);
                    this.makeRequest(method, url, formData)
                        .then(response => {
                            if (response.ok) {
                                return response.text();
                            } else {
                                throw new Error('Form submission failed: ' + response.statusText);
                            }
                        })
                        .then(html => this.outerHTML = html)
                        .catch(error => console.error('Form submission failed:', error));
                });
            });
        });
    }

    protected async makeRequest(method: string, url: string, formData: FormData): Promise<Response> {
        const options: RequestInit = {
            method,
            headers: {
                'X-Requested-With': 'XMLHttpRequest',
            },
        };

        if (method === 'GET') {
            const params = new URLSearchParams(formData as any);
            url += (url.includes('?') ? '&' : '?') + params.toString();
        } else {
            options.body = formData;
        }

        return fetch(url, options);
    }

    one(selector: string, cb: (el: HTMLElement) => void) {
        const el = this.shadowRoot!.querySelector(selector);
        if (el) {
            cb(el as HTMLElement);
        }
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

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

    protected bindFormHandlers() {
        // Handle forms with mesh-* HTTP attributes
        this.all('form[mesh-get], form[mesh-post], form[mesh-put], form[mesh-patch], form[mesh-delete]', form => {
            form.addEventListener('submit', this.handleFormSubmit.bind(this));
        });
    }

    protected async handleFormSubmit(event: Event) {
        event.preventDefault();
        const form = event.target as HTMLFormElement;
        
        // Determine HTTP method from mesh-* attributes
        let method = 'GET';
        let url = '';
        
        if (form.hasAttribute('mesh-get')) {
            method = 'GET';
            url = form.getAttribute('mesh-get') || '';
        } else if (form.hasAttribute('mesh-post')) {
            method = 'POST';
            url = form.getAttribute('mesh-post') || '';
        } else if (form.hasAttribute('mesh-put')) {
            method = 'PUT';
            url = form.getAttribute('mesh-put') || '';
        } else if (form.hasAttribute('mesh-patch')) {
            method = 'PATCH';
            url = form.getAttribute('mesh-patch') || '';
        } else if (form.hasAttribute('mesh-delete')) {
            method = 'DELETE';
            url = form.getAttribute('mesh-delete') || '';
        }

        if (!url) {
            console.error('No URL specified for form submission');
            return;
        }

        try {
            const formData = new FormData(form);
            const response = await this.makeRequest(method, url, formData);

            if (response.ok) {
                const html = await response.text();
                this.handleResponse(html, form);
            } else {
                console.error('Form submission failed:', response.status, response.statusText);
            }
        } catch (error) {
            console.error('Form submission error:', error);
        }
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

    protected handleResponse(html: string, form: HTMLFormElement) {
        // Check for mesh-target attribute to determine where to place response
        const target = form.getAttribute('mesh-target');
        
        if (target) {
            const targetElement = this.shadowRoot?.querySelector(target) || document.querySelector(target);
            if (targetElement) {
                targetElement.innerHTML = html;
                return;
            }
        }

        // Check for mesh-swap attribute to determine swap behavior
        const swapMode = form.getAttribute('mesh-swap') || 'outerHTML';
        
        switch (swapMode) {
            case 'innerHTML':
                this.innerHTML = html;
                break;
            case 'outerHTML':
            default:
                this.outerHTML = html;
                break;
        }
    }

    one(selector: string, cb: (el: HTMLElement) => void) {
        const el = this.shadowRoot!.querySelector(selector)
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

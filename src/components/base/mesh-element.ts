export class MeshElement extends HTMLElement {
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
    }
}
import {LitElement} from 'lit';
import {customElement} from 'lit/decorators.js';

@customElement('mesh-board')
export class Board extends LitElement {
    createRenderRoot() {
        // Use existing shadow root if it exists (from DSD)
        return this.shadowRoot || this.attachShadow({mode: 'open'});
    }

    connectedCallback() {
        super.connectedCallback();
        console.log('Board is connected');
    }
}

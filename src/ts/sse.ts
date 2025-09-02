export class SSEManager {
    private eventSource: EventSource | null = null;
    private isConnected = false;

    constructor(private url: string = '/sse') {
        this.connect();
    }

    private connect() {
        if (this.eventSource) {
            this.eventSource.close();
        }

        this.eventSource = new EventSource(this.url);

        this.eventSource.onopen = () => {
            console.log('SSE connection opened');
            this.isConnected = true;
        };

        this.eventSource.addEventListener('oob-update', (event) => {
            this.handleOOBUpdate(event as MessageEvent);
        });

        this.eventSource.onerror = (error) => {
            console.error('SSE connection error:', error);
            this.isConnected = false;
            // Reconnect after 5 seconds
            setTimeout(() => this.connect(), 5000);
        };
    }

    private handleOOBUpdate(event: MessageEvent) {
        const template = document.createElement('template');
        template.innerHTML = event.data.trim();

        for (const content of template.content.querySelectorAll('[mesh-swap-oob]')) {
            const id = content.id;
            const target = this.findInShadow(document, id);

            if (target) {
                target.outerHTML = content.outerHTML;
            } else {
                console.warn('OOB target not found:', id);
            }
        }
    }

    private findInShadow(root: Document | ShadowRoot | Element, id: string): Element | null {
        // First try normal document search
        let element = root.querySelector(`#${id}`);
        if (element) return element;

        // Search in shadow DOMs
        const allElements = root.querySelectorAll('*');
        for (const el of allElements) {
            if (el.shadowRoot) {
                element = this.findInShadow(el.shadowRoot, id);
                if (element) return element;
            }
        }

        return null;
    }

    public close() {
        if (this.eventSource) {
            this.eventSource.close();
            this.eventSource = null;
        }
        this.isConnected = false;
    }

    public get connected(): boolean {
        return this.isConnected;
    }
}

// Global SSE manager instance
export const sseManager = new SSEManager();
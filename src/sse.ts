interface BatchedUpdate {
    id: string;
    html: string;
}

interface UpdateBatch {
    batchId: string;
    updates: BatchedUpdate[];
}

export class SSEManager {
    private eventSource: EventSource | null = null;

    constructor(private url: string = '/sse?stream=oob-updates') {
        this.connect();
    }

    private connect() {
        if (this.eventSource) {
            this.eventSource.close();
        }

        this.eventSource = new EventSource(this.url);

        this.eventSource.addEventListener('oob-batch', (event) => {
            this.handleOOBBatch(event as MessageEvent);
        });

        this.eventSource.onerror = (error) => {
            console.error('SSE connection error:', error);
            setTimeout(() => this.connect(), 5000);
        };
    }

    private handleOOBBatch(event: MessageEvent) {
        const batch: UpdateBatch = JSON.parse(event.data);

        for (const update of batch.updates) {
            this.processOOBUpdate(update.html);
        }
    }

    private processOOBUpdate(html: string) {
        const template = document.createElement('template');
        template.innerHTML = html.trim();

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
        let element = root.querySelector(`#${id}`);
        if (element) {
            return element;
        }

        const allElements = root.querySelectorAll('*');
        for (const el of allElements) {
            if (el.shadowRoot) {
                element = this.findInShadow(el.shadowRoot, id);
                if (element) {
                    return element;
                }
            }
        }

        return null;
    }
}

new SSEManager();
// SSE Module for Mesh Elements
// Handles Server-Sent Events and provides event bus functionality

export interface SSEMessage {
    type: string;
    data: any;
    target?: string;
}

export class SSEManager {
    private eventSource: EventSource | null = null;
    private listeners: Map<string, Set<(data: any) => void>> = new Map();
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

        this.eventSource.onmessage = (event) => {
            this.handleMessage(event);
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

    private handleMessage(event: MessageEvent) {
        try {
            const message: SSEMessage = JSON.parse(event.data);
            this.emit(message.type, message.data);
        } catch (e) {
            console.error('Failed to parse SSE message:', e);
        }
    }

    private handleOOBUpdate(event: MessageEvent) {
        const template = document.createElement('template');
        template.innerHTML = event.data.trim();
        const content = template.content.firstElementChild;

        if (content && content.hasAttribute('mesh-swap-oob')) {
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

    public on(event: string, callback: (data: any) => void) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, new Set());
        }
        this.listeners.get(event)!.add(callback);
    }

    public off(event: string, callback: (data: any) => void) {
        const eventListeners = this.listeners.get(event);
        if (eventListeners) {
            eventListeners.delete(callback);
        }
    }

    public emit(event: string, data: any) {
        const eventListeners = this.listeners.get(event);
        if (eventListeners) {
            eventListeners.forEach(callback => callback(data));
        }
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
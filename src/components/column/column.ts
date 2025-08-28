import {MeshElement} from "../base/mesh-element.ts";

export class Column extends MeshElement {
    private dropIndicator: HTMLElement | null = null;

    connectedCallback() {
        super.connectedCallback();
        this.setupDropTarget();
    }

    setupDropTarget() {
        this.addEventListener('dragover', this.handleDragOver.bind(this));
        this.addEventListener('drop', this.handleDrop.bind(this));
        this.addEventListener('dragenter', this.handleDragEnter.bind(this));
        this.addEventListener('dragleave', this.handleDragLeave.bind(this));
    }

    handleDragOver(e: any) {
        e.preventDefault();
        e.dataTransfer.dropEffect = 'move';

        this.updateDropIndicator(e);
    }

    handleDragEnter() {
        this.classList.add('drag-over');

        this.createDropIndicator();
    }

    handleDragLeave(e: any) {
        // Only remove if we're actually leaving the column
        if (!this.contains(e.relatedTarget)) {
            this.classList.remove('drag-over');

            this.removeDropIndicator();
        }
    }

    handleDrop(e: any) {
        e.preventDefault();
        this.classList.remove('drag-over');
        this.removeDropIndicator();

        const cardId = e.dataTransfer.getData('text/plain');
        const columnId = this.dataset.id;
        if (!cardId || !columnId) {
            throw new Error('Missing card or column ID');
        }

        // Calculate position within column
        const position = this.calculateDropPosition(e);

        this.moveCard(cardId, +columnId, position);
    }

    createDropIndicator() {
        if (this.dropIndicator) return;

        this.dropIndicator = document.createElement('div');
        this.dropIndicator.className = 'drop-indicator';
        this.dropIndicator.style.cssText = `
            height: 4px;
            background: #007acc;
            border: 2px dashed #0056b3;
            border-radius: 4px;
            margin: 8px 0;
            opacity: 0.8;
            position: relative;
            transition: all 0.15s ease;
        `;
    }

    updateDropIndicator(e: any) {
        if (!this.dropIndicator) return;

        const cardsContainer = this.shadowRoot!.querySelector('.cards'); // Adjust selector as needed
        if (!cardsContainer) return;

        // Find where to insert the indicator
        const afterElement = this.getAfterElement(e);

        if (afterElement) {
            // Insert before this card
            cardsContainer.insertBefore(this.dropIndicator, afterElement);
        } else {
            // Insert at the end
            cardsContainer.appendChild(this.dropIndicator);
        }
    }

    removeDropIndicator() {
        if (this.dropIndicator) {
            this.dropIndicator.remove();
            this.dropIndicator = null;
        }
    }

    calculateDropPosition(e: any) {
        const afterElement = this.getAfterElement(e);

        if (!afterElement) {
            return 0;
        }

        return this.getCards().indexOf(afterElement);
    }

    private getCards() {
        return Array.from(this.shadowRoot!.querySelectorAll('mesh-card'));
    }

    private getAfterElement(e: any) {
        return this.getCards().find(card => {
            const rect = card.getBoundingClientRect();
            return e.clientY < rect.top + rect.height / 2;
        });
    }

    async moveCard(cardId: number, columnId: number, position: number) {
        window.htmx.ajax('put', '/card', {
            swap: 'none',
            values: {
                action: 'move',
                cardID: cardId,
                columnID: columnId,
                position: position,
            }
        } as any);
    }
}
window.customElements.define('mesh-column', Column);
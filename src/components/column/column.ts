import {MeshElement} from "../base/mesh-element.ts";

export class Column extends MeshElement {

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

    handleDragOver(e) {
        e.preventDefault();
        e.dataTransfer.dropEffect = 'move';
    }

    handleDragEnter(e) {
        this.classList.add('drag-over');
    }

    handleDragLeave(e) {
        // Only remove if we're actually leaving the column
        if (!this.contains(e.relatedTarget)) {
            this.classList.remove('drag-over');
        }
    }

    handleDrop(e) {
        e.preventDefault();
        this.classList.remove('drag-over');

        const cardId = e.dataTransfer.getData('text/plain');
        const columnId = this.id;

        // Calculate position within column
        const position = this.calculateDropPosition(e);

        this.moveCard(cardId, columnId, position);
    }

    calculateDropPosition(e: any) {
        const cards = Array.from(this.querySelectorAll('mesh-card:not(.dragging)'));
        const afterElement = cards.find(card => {
            const rect = card.getBoundingClientRect();
            return e.clientY < rect.top + rect.height / 2;
        });

        return afterElement ? afterElement.id : 'end';
    }

    async moveCard(cardId: number, columnId: number, position: number) {
        window.htmx.ajax('put', '/card', {cardId, columnId, position});
    }
}
window.customElements.define('mesh-column', Column);
import {MeshElement} from "../base/mesh-element.ts";

import {ArrowLeft, ArrowRight, CircleX, Pencil} from 'lucide';

export class Card extends MeshElement {
    protected icons = {
        ArrowLeft,
        ArrowRight,
        CircleX,
        Pencil,
    };

    edit() {
        this.show('[data-form]');
        this.hide('[data-view]');
    }

    cancel() {
        this.hide('[data-form]');
        this.show('[data-view]');
    }

    connectedCallback() {
        super.connectedCallback();
        this.setupDragAndDrop();
    }

    setupDragAndDrop() {
        this.draggable = true;

        this.addEventListener('dragstart', this.handleDragStart.bind(this));
        this.addEventListener('dragend', this.handleDragEnd.bind(this));
    }

    handleDragStart(e: any) {
        e.dataTransfer.setData('text/plain', this.id);
        this.classList.add('dragging');
        e.dataTransfer.effectAllowed = 'move';
    }

    handleDragEnd() {
        this.classList.remove('dragging');
    }

}
window.customElements.define('mesh-card', Card);
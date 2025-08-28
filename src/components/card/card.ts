import {MeshElement} from "../base/mesh-element.ts";

import {ArrowLeft, ArrowRight, CircleX, Pencil, Grip} from 'lucide';

export class Card extends MeshElement {
    protected icons = {
        ArrowLeft,
        ArrowRight,
        CircleX,
        Pencil,
        Grip,
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
        this.one('.grip', grip => {
            grip.draggable = true;
            this.addEventListener('dragstart', this.handleDragStart.bind(this));
            this.addEventListener('dragend', this.handleDragEnd.bind(this));
        });
    }

    handleDragStart(e: any) {
        const dragImage = this.createDragImage();
        e.dataTransfer.setDragImage(dragImage, 0, 0);

        e.dataTransfer.setData('text/plain', this.dataset.id);
        this.classList.add('dragging');
        e.dataTransfer.effectAllowed = 'move';
    }

    handleDragEnd() {
        this.classList.remove('dragging');
    }

    createDragImage() {
        console.log('piss');

        // Clone the card element
        const clone = this.cloneNode(true) as HTMLElement;

        // Style it for dragging
        clone.style.position = 'absolute';
        clone.style.top = '-1000px'; // Off-screen
        clone.style.left = '-1000px';
        clone.style.width = this.offsetWidth + 'px';
        clone.style.height = this.offsetHeight + 'px';
        clone.style.transform = 'rotate(5deg) scale(1.05) translate(-50%, -50%)'; // Jaunty angle
        clone.style.opacity = '1';
        clone.style.boxShadow = 'none';
        clone.style.zIndex = '9999';
        clone.style.pointerEvents = 'none';

        // Add it to the body temporarily
        document.body.appendChild(clone);

        // Remove it after drag starts (browser captures it)
        setTimeout(() => {
            if (document.body.contains(clone)) {
                document.body.removeChild(clone);
            }
        }, 0);

        return clone;
    }
}
window.customElements.define('mesh-card', Card);
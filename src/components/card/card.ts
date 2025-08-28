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

    promote() {
        console.log('promote');
    }

    demote() {
        console.log('demote');
    }
}
window.customElements.define('mesh-card', Card);
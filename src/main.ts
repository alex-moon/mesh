import './main.scss';

// Import components to register custom elements
import './components/app/app';
import './components/board/board';
import './components/column/column';
import './components/card/card';

// Import SSE manager
import {sseManager} from './sse.ts';

console.log('MESH initialized - HTMX has been replaced');
console.log('SSE connection status:', sseManager.connected);

export function findInShadow(root: any, id: string): any {
    const element = root.getElementById?.(id);
    if (element) {
        return element;
    }
    const allElements = root.querySelectorAll('*');
    for (let element of allElements) {
        if (element.shadowRoot) {
            const found = findInShadow(element.shadowRoot, id);
            if (found) {
                return found;
            }
        }
    }
    return null;
}

/**
 * The one dynamic import of the code editor. Monaco is in its chunk: no
 * module that the pages load first may import `./DAGEditor` in a static way.
 */
export const loadDAGEditor = () => import('./DAGEditor');

/** Starts the load of the code editor, for example when a pointer nears the button that opens it. */
export function preloadDAGEditor(): void {
  void loadDAGEditor();
}

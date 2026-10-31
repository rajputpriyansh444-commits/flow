export function normalizeWikiPagePathFromURL(path: string): string {
  return path.replace(/\.md$/i, '');
}

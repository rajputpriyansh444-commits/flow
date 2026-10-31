export function encodeWikiPagePathForURL(wikiPagePath: string): string {
  const unambiguousPath = wikiPagePath.toLowerCase().endsWith('.md')
    ? `${wikiPagePath}.md`
    : wikiPagePath;
  return unambiguousPath.split('/').map(encodeURIComponent).join('/');
}

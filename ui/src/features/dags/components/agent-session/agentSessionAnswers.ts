export function combineQuestionAnswer(
  selected: string[],
  custom: string,
  multiple: boolean
): string[] {
  const value = custom.trim();
  if (!value) return selected;
  return multiple ? [...selected, value] : [value];
}

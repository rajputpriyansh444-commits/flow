import { useI18n } from './I18nProvider';
import {
  translateStatic,
  type StaticTranslationValues,
} from './staticMessages';

export function I18nText({
  text,
  values,
}: {
  text: string;
  values?: StaticTranslationValues;
}): string {
  const { locale } = useI18n();
  return translateStatic(locale, text, values);
}

import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';

import { Markdown } from '../markdown';

describe('Markdown', () => {
  it('renders inline and fenced code with theme-aware contrast styles', () => {
    const { container } = render(
      <Markdown content={'Use `flow start`.\n\n```yaml\nsteps: []\n```'} />
    );

    const pre = container.querySelector('pre');
    const blockCode = pre?.querySelector('code');
    const inlineCode = Array.from(container.querySelectorAll('code')).find(
      (code) => code.textContent === 'flow start'
    );

    expect(pre?.getAttribute('style')).toContain(
      'background-color: var(--muted)'
    );
    expect(pre?.getAttribute('style')).toContain('color: var(--foreground)');
    expect(blockCode?.getAttribute('style')).toContain(
      'background-color: transparent'
    );
    expect(blockCode?.getAttribute('style')).toContain(
      'color: var(--foreground)'
    );
    expect(inlineCode?.getAttribute('style')).toContain(
      'background-color: var(--muted)'
    );
    expect(inlineCode?.getAttribute('style')).toContain(
      'color: var(--foreground)'
    );
  });
});

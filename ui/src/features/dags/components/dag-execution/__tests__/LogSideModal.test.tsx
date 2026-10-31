import { fireEvent, render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { dispatchHandledEscape } from '@/test/keyboard';
import LogSideModal from '../LogSideModal';

describe('LogSideModal', () => {
  it('closes on Escape unless a dialog inside it handled the key', () => {
    const onClose = vi.fn();
    render(
      <LogSideModal isOpen onClose={onClose} title="stdout">
        log output
      </LogSideModal>
    );

    dispatchHandledEscape();
    expect(onClose).not.toHaveBeenCalled();

    fireEvent.keyDown(window, { key: 'Escape' });
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});

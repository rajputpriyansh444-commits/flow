import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it } from 'vitest';
import { UpdateBanner } from '@/components/UpdateBanner';
import { ConfigContext, type Config } from '@/contexts/ConfigContext';

const updateConfig = {
  version: '2.0.0',
  updateAvailable: true,
  latestVersion: 'v2.1.0',
  upgradeCommand: 'flow upgrade',
  restartRequired: false,
  restartCommand: '',
} satisfies Partial<Config>;

// The banner reads only the update fields, so the rest of Config is omitted.
function renderBanner(overrides: Partial<Config> = {}) {
  const config = { ...updateConfig, ...overrides } as Config;
  return render(
    <ConfigContext.Provider value={config}>
      <UpdateBanner />
    </ConfigContext.Provider>
  );
}

describe('UpdateBanner', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('renders nothing when Flow is up to date', () => {
    const { container } = renderBanner({ updateAvailable: false });

    expect(container).toBeEmptyDOMElement();
  });

  it('shows the update command for this installation', () => {
    renderBanner({ upgradeCommand: 'brew upgrade flow' });

    expect(screen.getByText(/v2\.0\.0 → v2\.1\.0/)).toBeInTheDocument();
    expect(screen.getByText('brew upgrade flow')).toBeInTheDocument();
    expect(screen.queryByText('flow upgrade')).not.toBeInTheDocument();
  });

  it('omits the run instruction when no single command updates Flow', () => {
    renderBanner({ upgradeCommand: '' });

    expect(screen.getByText(/v2\.0\.0 → v2\.1\.0/)).toBeInTheDocument();
    expect(screen.queryByText(/to update/)).not.toBeInTheDocument();
  });

  it('asks for a restart once a new binary is installed', () => {
    renderBanner({
      restartRequired: true,
      restartCommand: 'sudo systemctl restart flow',
    });

    expect(
      screen.getByText(/A new Flow binary is installed/)
    ).toBeInTheDocument();
    expect(screen.getByText('sudo systemctl restart flow')).toBeInTheDocument();
    expect(screen.queryByText('flow upgrade')).not.toBeInTheDocument();
  });

  it('asks for a restart without a command when no service is known', () => {
    renderBanner({ updateAvailable: false, restartRequired: true });

    expect(
      screen.getByText(
        'A new Flow binary is installed. Restart Flow to use it.'
      )
    ).toBeInTheDocument();
  });

  it('stays dismissed for the same version only', async () => {
    const user = userEvent.setup();
    const { unmount } = renderBanner();

    await user.click(
      screen.getByRole('button', { name: 'Dismiss update notification' })
    );
    expect(screen.queryByText(/v2\.1\.0/)).not.toBeInTheDocument();

    unmount();
    renderBanner();
    expect(screen.queryByText(/v2\.1\.0/)).not.toBeInTheDocument();

    renderBanner({ latestVersion: 'v2.2.0' });
    expect(screen.getByText(/v2\.0\.0 → v2\.2\.0/)).toBeInTheDocument();
  });

  it('shows the restart notice again after dismissing the update notice', async () => {
    const user = userEvent.setup();
    const { unmount } = renderBanner();

    await user.click(
      screen.getByRole('button', { name: 'Dismiss update notification' })
    );
    unmount();
    renderBanner({ restartRequired: true });

    expect(
      screen.getByText(/A new Flow binary is installed/)
    ).toBeInTheDocument();
  });
});

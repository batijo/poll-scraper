import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { expect, test, vi, beforeEach } from 'vitest';
import ScraperActions from '../ScraperActions.svelte';

const mockStartScraper = vi.fn().mockResolvedValue(undefined);
const mockStopScraper = vi.fn().mockResolvedValue(undefined);
const mockPreviewScrape = vi.fn().mockResolvedValue({ rawData: [], data: [], statuses: [] });

vi.mock('../../../../wailsjs/go/main/App', () => ({
  StartScraper: (...args: unknown[]) => mockStartScraper(...args),
  StopScraper: (...args: unknown[]) => mockStopScraper(...args),
  PreviewScrape: (...args: unknown[]) => mockPreviewScrape(...args),
}));

beforeEach(() => {
  vi.clearAllMocks();
  mockStartScraper.mockResolvedValue(undefined);
  mockStopScraper.mockResolvedValue(undefined);
  mockPreviewScrape.mockResolvedValue({ rawData: [], data: [], statuses: [] });
});

// ConfirmDialog renders a "Stop" button inside <dialog> that's always in the DOM
// (jsdom doesn't hide non-open dialogs), so we query the action buttons area only.

function getActionButtons(container: HTMLElement) {
  const actionsDiv = container.querySelector('.flex.gap-2.justify-end');
  return actionsDiv!;
}

test('shows Start button when stopped', () => {
  const { container } = render(ScraperActions, { scraperState: 'stopped', onPreviewComplete: () => {} });
  const actions = getActionButtons(container);
  expect(actions.textContent).toContain('Start');
  expect(actions.textContent).not.toContain('Stop');
});

test('shows Stop button when scraping', () => {
  const { container } = render(ScraperActions, { scraperState: 'scraping', onPreviewComplete: () => {} });
  const actions = getActionButtons(container);
  expect(actions.textContent).toContain('Stop');
  expect(actions.textContent).not.toContain('Start');
});

test('calls StartScraper on Start click', async () => {
  const user = userEvent.setup();
  render(ScraperActions, { scraperState: 'stopped', onPreviewComplete: () => {} });
  await user.click(screen.getByText('Start'));
  expect(mockStartScraper).toHaveBeenCalledOnce();
});

test('shows error when StartScraper fails', async () => {
  mockStartScraper.mockRejectedValueOnce('network error');
  const user = userEvent.setup();
  render(ScraperActions, { scraperState: 'stopped', onPreviewComplete: () => {} });
  await user.click(screen.getByText('Start'));

  await waitFor(() => {
    expect(screen.getByText(/Failed to start/)).toBeInTheDocument();
  });
});

test('calls onPreviewComplete after successful preview', async () => {
  const previewResult = { rawData: [], data: [{ name: 'A', value: '1' }], statuses: [] };
  mockPreviewScrape.mockResolvedValueOnce(previewResult);
  const onPreviewComplete = vi.fn();
  const user = userEvent.setup();
  render(ScraperActions, { scraperState: 'stopped', onPreviewComplete });
  await user.click(screen.getByText('Preview'));

  await waitFor(() => {
    expect(onPreviewComplete).toHaveBeenCalledWith(previewResult);
  });
});

test('shows error when preview fails', async () => {
  mockPreviewScrape.mockRejectedValueOnce('timeout');
  const user = userEvent.setup();
  render(ScraperActions, { scraperState: 'stopped', onPreviewComplete: () => {} });
  await user.click(screen.getByText('Preview'));

  await waitFor(() => {
    expect(screen.getByText(/Preview failed/)).toBeInTheDocument();
  });
});

test('Preview button is always visible', () => {
  render(ScraperActions, { scraperState: 'scraping', onPreviewComplete: () => {} });
  expect(screen.getByText('Preview')).toBeInTheDocument();
});

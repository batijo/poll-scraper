import { render, screen, waitFor } from '@testing-library/svelte';
import { expect, test, vi, beforeEach } from 'vitest';
import { tick } from 'svelte';
import DataDisplay from '../DataDisplay.svelte';

const eventCallbacks: Record<string, Function> = {};

vi.mock('../../../../wailsjs/runtime/runtime', () => ({
  EventsOn: vi.fn((eventName: string, callback: Function) => {
    eventCallbacks[eventName] = callback;
    return () => {};
  }),
  EventsOnMultiple: vi.fn(),
  EventsOff: vi.fn(),
  EventsEmit: vi.fn(),
}));

beforeEach(() => {
  for (const key of Object.keys(eventCallbacks)) {
    delete eventCallbacks[key];
  }
  vi.clearAllMocks();
});

test('shows "Scraper is stopped" message when stopped with no data', () => {
  render(DataDisplay, { currentState: 'stopped' });
  expect(screen.getByText(/Scraper is stopped/)).toBeInTheDocument();
});

test('shows "Waiting for data" spinner after polled:state sets scraping', async () => {
  render(DataDisplay, {});

  // Internal scraperState starts as 'stopped', must set via event
  const stateCallback = eventCallbacks['polled:state'];
  stateCallback('scraping');
  await tick();

  await waitFor(() => {
    expect(screen.getByText('Waiting for data...')).toBeInTheDocument();
  });
});

test('shows ErrorCard after polled:state sets error', async () => {
  render(DataDisplay, {});

  const stateCallback = eventCallbacks['polled:state'];
  stateCallback('error');
  await tick();

  await waitFor(() => {
    expect(screen.getByText('Error Loading Data')).toBeInTheDocument();
  });
});

test('shows StatusIndicator with "Stopped" by default', () => {
  render(DataDisplay, {});
  expect(screen.getByText('Stopped')).toBeInTheDocument();
});

test('registers all event listeners on mount', async () => {
  render(DataDisplay, {});
  const { EventsOn } = await import('../../../../wailsjs/runtime/runtime');
  expect(EventsOn).toHaveBeenCalledWith('polled:data', expect.any(Function));
  expect(EventsOn).toHaveBeenCalledWith('polled:state', expect.any(Function));
  expect(EventsOn).toHaveBeenCalledWith('polled:url-status', expect.any(Function));
  expect(EventsOn).toHaveBeenCalledWith('polled:error', expect.any(Function));
  expect(EventsOn).toHaveBeenCalledWith('polled:log', expect.any(Function));
});

test('line count change: polled:error event shows ErrorCard with message', async () => {
  render(DataDisplay, { currentState: 'scraping' });

  // Simulate backend sending error event (line count changed)
  const errorCallback = eventCallbacks['polled:error'];
  expect(errorCallback).toBeDefined();

  errorCallback({
    message: 'URL line count changed for http://example.com: expected 5, got 3',
    timestamp: String(Date.now()),
  });

  await tick();

  await waitFor(() => {
    expect(screen.getByText('Error Loading Data')).toBeInTheDocument();
    expect(screen.getByText('URL line count changed for http://example.com: expected 5, got 3')).toBeInTheDocument();
  });
});

test('line count change: polled:error event sets status to error (red)', async () => {
  const { container } = render(DataDisplay, { currentState: 'scraping' });

  const errorCallback = eventCallbacks['polled:error'];
  errorCallback({
    message: 'URL line count changed',
    timestamp: String(Date.now()),
  });

  await tick();

  await waitFor(() => {
    // StatusIndicator should show "Error" label
    expect(screen.getByText('Error')).toBeInTheDocument();
    // Red dot should be present
    expect(container.querySelector('.bg-red-500')).toBeInTheDocument();
  });
});

test('polled:state event updates scraper state', async () => {
  render(DataDisplay, { currentState: 'scraping' });

  const stateCallback = eventCallbacks['polled:state'];
  expect(stateCallback).toBeDefined();

  stateCallback('error');
  await tick();

  await waitFor(() => {
    expect(screen.getByText('Error')).toBeInTheDocument();
  });
});

test('polled:data event displays data in grid', async () => {
  render(DataDisplay, { currentState: 'scraping' });

  const dataCallback = eventCallbacks['polled:data'];
  expect(dataCallback).toBeDefined();

  dataCallback({
    data: [
      { name: 'Score', value: '42' },
      { name: 'Count', value: '7' },
    ],
    rawData: [
      { name: 'Score', value: '42' },
      { name: 'Count', value: '7' },
    ],
    timestamp: new Date().toISOString(),
  });

  await tick();

  await waitFor(() => {
    expect(screen.getByText('Score:')).toBeInTheDocument();
    expect(screen.getByText('42')).toBeInTheDocument();
    expect(screen.getByText('Count:')).toBeInTheDocument();
    expect(screen.getByText('7')).toBeInTheDocument();
  });
});

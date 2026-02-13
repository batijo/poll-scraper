import { render, screen } from '@testing-library/svelte';
import { expect, test } from 'vitest';
import StatusSection from '../StatusSection.svelte';
import { createDefaultConfig } from '$lib/types/config';

function renderStatus(propOverrides: Record<string, unknown> = {}) {
  const config = {
    ...createDefaultConfig(),
    links: ['http://example.com'],
  };
  return render(StatusSection, { config, ...propOverrides });
}

test('shows last error message when lastError is set', () => {
  renderStatus({
    scraperState: 'error',
    lastError: 'URL line count changed for http://example.com: expected 5, got 3',
  });
  expect(screen.getByText('Last Error')).toBeInTheDocument();
  expect(screen.getByText('URL line count changed for http://example.com: expected 5, got 3')).toBeInTheDocument();
});

test('does not show last error when lastError is null', () => {
  renderStatus({ scraperState: 'scraping' });
  expect(screen.queryByText('Last Error')).not.toBeInTheDocument();
});

test('shows URL ok/fail breakdown', () => {
  renderStatus({
    scraperState: 'scraping',
    urlStatusList: [
      { url: 'http://ok.com', hasData: true, lineCount: 3, error: false },
      { url: 'http://bad.com', hasData: false, lineCount: 0, error: true },
    ],
  });
  expect(screen.getByText('ok')).toBeInTheDocument();
  expect(screen.getByText('fail')).toBeInTheDocument();
});

test('does not show fail count when all URLs are healthy', () => {
  renderStatus({
    scraperState: 'scraping',
    urlStatusList: [
      { url: 'http://ok.com', hasData: true, lineCount: 3, error: false },
    ],
  });
  expect(screen.getByText('ok')).toBeInTheDocument();
  expect(screen.queryByText('fail')).not.toBeInTheDocument();
});

test('shows configured URL count', () => {
  renderStatus({});
  expect(screen.getByText(/1 configured/)).toBeInTheDocument();
});

test('shows update interval from config', () => {
  renderStatus({});
  expect(screen.getByText('1000ms')).toBeInTheDocument();
});

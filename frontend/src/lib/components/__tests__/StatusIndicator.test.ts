import { render, screen } from '@testing-library/svelte';
import { expect, test } from 'vitest';
import StatusIndicator from '../display/StatusIndicator.svelte';

test('shows red dot and "Error" label for error status', () => {
  const { container } = render(StatusIndicator, { effectiveStatus: 'error' });
  expect(screen.getByText('Error')).toBeInTheDocument();
  expect(container.querySelector('.bg-red-500')).toBeInTheDocument();
});

test('shows green dot and "Scraping" label for ok status', () => {
  const { container } = render(StatusIndicator, { effectiveStatus: 'ok' });
  expect(screen.getByText('Scraping')).toBeInTheDocument();
  expect(container.querySelector('.bg-green-500')).toBeInTheDocument();
});

test('shows gray dot and "Stopped" label for stopped status', () => {
  const { container } = render(StatusIndicator, { effectiveStatus: 'stopped' });
  expect(screen.getByText('Stopped')).toBeInTheDocument();
  expect(container.querySelector('.bg-gray-500')).toBeInTheDocument();
});

test('defaults to stopped when no prop given', () => {
  render(StatusIndicator);
  expect(screen.getByText('Stopped')).toBeInTheDocument();
});

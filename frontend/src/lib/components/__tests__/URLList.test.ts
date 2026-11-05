import { render, screen } from '@testing-library/svelte';
import { expect, test, vi } from 'vitest';
import URLList from '../forms/URLList.svelte';

vi.mock('../../../../wailsjs/go/main/App', () => ({
  PreviewURL: vi.fn().mockResolvedValue([]),
}));

test('shows red dot with "Line count changed" title when URL has error', () => {
  const { container } = render(URLList, {
    links: [{ url: 'http://example.com' }],
    initialLinks: [{ url: 'http://example.com' }],
    urlStatusList: [{ url: 'http://example.com', hasData: true, lineCount: 5, error: true }],
  });

  const dot = container.querySelector('.bg-red-500');
  expect(dot).toBeInTheDocument();
  expect(dot?.getAttribute('title')).toBe('Line count changed');
});

test('shows green dot when URL has data and no error', () => {
  const { container } = render(URLList, {
    links: [{ url: 'http://example.com' }],
    initialLinks: [{ url: 'http://example.com' }],
    urlStatusList: [{ url: 'http://example.com', hasData: true, lineCount: 5, error: false }],
  });

  const dot = container.querySelector('.bg-green-500');
  expect(dot).toBeInTheDocument();
  expect(dot?.getAttribute('title')).toBe('Producing data');
});

test('shows gray dot when no status available', () => {
  const { container } = render(URLList, {
    links: [{ url: 'http://example.com' }],
    initialLinks: [{ url: 'http://example.com' }],
    urlStatusList: [],
  });

  const dot = container.querySelector('.bg-gray-500');
  expect(dot).toBeInTheDocument();
  expect(dot?.getAttribute('title')).toBe('No data');
});

test('shows red dot only for the URL with error', () => {
  const { container } = render(URLList, {
    links: [{ url: 'http://ok.com' }, { url: 'http://bad.com' }],
    initialLinks: [{ url: 'http://ok.com' }, { url: 'http://bad.com' }],
    urlStatusList: [
      { url: 'http://ok.com', hasData: true, lineCount: 3, error: false },
      { url: 'http://bad.com', hasData: true, lineCount: 5, error: true },
    ],
  });

  const redDots = container.querySelectorAll('.bg-red-500');
  const greenDots = container.querySelectorAll('.bg-green-500');
  expect(redDots).toHaveLength(1);
  expect(greenDots).toHaveLength(1);
});

test('shows the configured username without exposing the password', () => {
  render(URLList, {
    links: [{ url: 'http://example.com', name: 'alice', password: 'secret' }],
    initialLinks: [{ url: 'http://example.com', name: 'alice', password: 'secret' }],
    urlStatusList: [],
  });

  expect(screen.getByText('Username: alice')).toBeInTheDocument();
  expect(screen.queryByText('secret')).not.toBeInTheDocument();
});

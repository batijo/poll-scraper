import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { expect, test, vi } from 'vitest';
import ErrorCard from '../display/ErrorCard.svelte';

test('renders error message and heading', () => {
  render(ErrorCard, { message: 'Connection failed', reset: () => {} });
  expect(screen.getByText('Error Loading Data')).toBeInTheDocument();
  expect(screen.getByText('Connection failed')).toBeInTheDocument();
});

test('calls reset when Try Again is clicked', async () => {
  const reset = vi.fn();
  const user = userEvent.setup();
  render(ErrorCard, { message: 'Error', reset });
  await user.click(screen.getByText('Try Again'));
  expect(reset).toHaveBeenCalledOnce();
});

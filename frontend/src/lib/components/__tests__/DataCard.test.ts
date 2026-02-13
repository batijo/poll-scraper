import { render, screen } from '@testing-library/svelte';
import { expect, test } from 'vitest';
import DataCard from '../display/DataCard.svelte';

test('renders name, value, and 1-based index', () => {
  render(DataCard, { name: 'Score', value: '42', index: 0 });
  expect(screen.getByText('1.')).toBeInTheDocument();
  expect(screen.getByText('Score:')).toBeInTheDocument();
  expect(screen.getByText('42')).toBeInTheDocument();
});

test('displays correct index for non-zero position', () => {
  render(DataCard, { name: 'Count', value: '10', index: 4 });
  expect(screen.getByText('5.')).toBeInTheDocument();
});

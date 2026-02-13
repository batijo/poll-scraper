import { render, screen } from '@testing-library/svelte';
import { expect, test } from 'vitest';
import DataGrid from '../display/DataGrid.svelte';

test('renders a card for each data item', () => {
  const data = [
    { name: 'Score', value: '42' },
    { name: 'Count', value: '7' },
  ];
  render(DataGrid, { data });
  expect(screen.getByText('Score:')).toBeInTheDocument();
  expect(screen.getByText('42')).toBeInTheDocument();
  expect(screen.getByText('Count:')).toBeInTheDocument();
  expect(screen.getByText('7')).toBeInTheDocument();
});

test('renders nothing when data is empty', () => {
  const { container } = render(DataGrid, { data: [] });
  expect(container.querySelector('.space-y-2')?.children).toHaveLength(0);
});

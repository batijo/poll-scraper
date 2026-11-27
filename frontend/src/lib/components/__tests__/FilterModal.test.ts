import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { beforeAll, expect, test, vi } from 'vitest';
import FilterModal from '../FilterModal.svelte';

beforeAll(() => {
  HTMLDialogElement.prototype.showModal = function () {
    this.setAttribute('open', '');
  };
  HTMLDialogElement.prototype.close = function () {
    this.removeAttribute('open');
  };
});

const rawScrapedData = [
  { name: 'a', value: '10' },
  { name: 'b', value: '20' },
  { name: 'c', value: '30' },
];

test('symbols follow their line when filters change and save as output positions', async () => {
  const onConfirm = vi.fn();
  render(FilterModal, {
    showModal: true,
    rawScrapedData,
    urlStatusList: [{ url: 'https://example.com', hasData: true, lineCount: 3, error: false }],
    selectedLines: [1, 3],
    customLines: [
      { name: 'hidden', value: '1', filtered: true },
      { name: 'custom', value: '5', filtered: false },
    ],
    symbolLines: [{ number: 2, symbol: ' eur' }],
    onConfirm,
  });
  const user = userEvent.setup();

  expect(screen.getByTitle('Symbol for output line 2')).toHaveValue(' eur');

  await user.click(screen.getByLabelText(/\[2\] b: 20/));
  await user.type(screen.getByTitle('Symbol for output line 4'), '%');
  await user.click(screen.getByRole('button', { name: 'Confirm' }));

  expect(onConfirm).toHaveBeenCalledWith(
    [1, 3, 2],
    [
      { number: 3, symbol: ' eur' },
      { number: 4, symbol: '%' },
    ]
  );
});

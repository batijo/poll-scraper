import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { beforeAll, expect, test, vi } from 'vitest';
import URLModal from '../URLModal.svelte';

beforeAll(() => {
  HTMLDialogElement.prototype.showModal = function () {
    this.setAttribute('open', '');
  };
  HTMLDialogElement.prototype.close = function () {
    this.removeAttribute('open');
  };
});

test('adds a URL without credentials', async () => {
  const onAdd = vi.fn();
  const { component } = render(URLModal, { onAdd });
  const user = userEvent.setup();

  component.open();
  await user.type(screen.getByLabelText('URL'), 'https://example.com/data');
  await user.click(screen.getByRole('button', { name: 'Add' }));

  expect(onAdd).toHaveBeenCalledWith({
    url: 'https://example.com/data',
    name: undefined,
    password: undefined,
  });
});

test('adds a URL with credentials', async () => {
  const onAdd = vi.fn();
  const { component } = render(URLModal, { onAdd });
  const user = userEvent.setup();

  component.open();
  await user.type(screen.getByLabelText('URL'), 'https://example.com/data');
  await user.type(screen.getByLabelText(/Username/), 'alice');
  await user.type(screen.getByLabelText(/Password/), 'secret');
  await user.click(screen.getByRole('button', { name: 'Add' }));

  expect(onAdd).toHaveBeenCalledWith({
    url: 'https://example.com/data',
    name: 'alice',
    password: 'secret',
  });
});

test('prefills credentials when editing', async () => {
  const onEdit = vi.fn();
  const { component } = render(URLModal, { onAdd: vi.fn(), onEdit });
  const user = userEvent.setup();

  component.open({ url: 'https://example.com/data', name: 'alice', password: 'secret' }, 2);
  await waitFor(() => {
    expect(screen.getByLabelText(/Username/)).toHaveValue('alice');
    expect(screen.getByLabelText(/Password/)).toHaveValue('secret');
  });

  await user.click(screen.getByRole('button', { name: 'Save' }));

  expect(onEdit).toHaveBeenCalledWith(2, {
    url: 'https://example.com/data',
    name: 'alice',
    password: 'secret',
  });
});

test('requires both credentials when only one is filled', async () => {
  const { component } = render(URLModal, { onAdd: vi.fn() });
  const user = userEvent.setup();

  component.open();
  await user.type(screen.getByLabelText('URL'), 'https://example.com/data');
  await user.type(screen.getByLabelText(/Username/), 'alice');

  expect(screen.getByText('Enter both username and password or leave both empty')).toBeInTheDocument();
  expect(screen.getByRole('button', { name: 'Add' })).toBeDisabled();
});

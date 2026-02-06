export interface Config {
  links: string[];
  port: number;
  ip: string;
  domains: string[];
  with_eq: boolean;
  filter_lines: number[];
  add_lines: string[];
  add_sum: boolean;
  sum_symbols: string;
  update_interval: number;
  write_to_csv: boolean;
  csv_path: string;
  write_to_txt: boolean;
  txt_path: string;
  dataset_name: string;
  debug: boolean;
}

export function createDefaultConfig(): Config {
  return {
    links: [],
    port: 3000,
    ip: 'localhost',
    domains: [],
    with_eq: false,
    filter_lines: [],
    add_lines: [],
    add_sum: false,
    sum_symbols: '',
    update_interval: 1000,
    write_to_csv: false,
    csv_path: '',
    write_to_txt: false,
    txt_path: '',
    dataset_name: '',
    debug: false,
  };
}

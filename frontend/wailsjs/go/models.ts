export namespace config {
	
	export class Config {
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
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.links = source["links"];
	        this.port = source["port"];
	        this.ip = source["ip"];
	        this.domains = source["domains"];
	        this.with_eq = source["with_eq"];
	        this.filter_lines = source["filter_lines"];
	        this.add_lines = source["add_lines"];
	        this.add_sum = source["add_sum"];
	        this.sum_symbols = source["sum_symbols"];
	        this.update_interval = source["update_interval"];
	        this.write_to_csv = source["write_to_csv"];
	        this.csv_path = source["csv_path"];
	        this.write_to_txt = source["write_to_txt"];
	        this.txt_path = source["txt_path"];
	        this.dataset_name = source["dataset_name"];
	        this.debug = source["debug"];
	    }
	}

}


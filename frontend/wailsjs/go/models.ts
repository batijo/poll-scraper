export namespace config {
	
	export class AddLine {
	    name: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new AddLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class Config {
	    links: string[];
	    port: number;
	    ip: string;
	    domains: string[];
	    enable_server: boolean;
	    with_eq: boolean;
	    filter_lines: number[];
	    add_lines: AddLine[];
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
	        this.enable_server = source["enable_server"];
	        this.with_eq = source["with_eq"];
	        this.filter_lines = source["filter_lines"];
	        this.add_lines = this.convertValues(source["add_lines"], AddLine);
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
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace models {
	
	export class Data {
	    name: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Data(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class URLStatus {
	    url: string;
	    hasData: boolean;
	
	    static createFrom(source: any = {}) {
	        return new URLStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.hasData = source["hasData"];
	    }
	}

}

export namespace utils {
	
	export class LogEntry {
	    level: string;
	    message: string;
	    time: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.message = source["message"];
	        this.time = source["time"];
	    }
	}

}


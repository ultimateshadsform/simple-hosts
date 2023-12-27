export namespace main {
	
	export class Host {
	    host: string;
	    ip: string;
	    comment: string;
	
	    static createFrom(source: any = {}) {
	        return new Host(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.ip = source["ip"];
	        this.comment = source["comment"];
	    }
	}

}


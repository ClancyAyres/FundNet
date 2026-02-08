export namespace models {
	
	export class Fund {
	    id: number;
	    code: string;
	    name: string;
	    current_price: number;
	    change_rate: number;
	    // Go type: time
	    last_updated: any;
	
	    static createFrom(source: any = {}) {
	        return new Fund(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.code = source["code"];
	        this.name = source["name"];
	        this.current_price = source["current_price"];
	        this.change_rate = source["change_rate"];
	        this.last_updated = this.convertValues(source["last_updated"], null);
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
	export class FundPosition {
	    id: number;
	    fund_code: string;
	    shares: number;
	    cost_price: number;
	    group_name: string;
	    // Go type: time
	    last_updated: any;
	
	    static createFrom(source: any = {}) {
	        return new FundPosition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.fund_code = source["fund_code"];
	        this.shares = source["shares"];
	        this.cost_price = source["cost_price"];
	        this.group_name = source["group_name"];
	        this.last_updated = this.convertValues(source["last_updated"], null);
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
	export class FundValue {
	    fund_code: string;
	    current_price: number;
	    change_rate: number;
	    estimate_value: number;
	    estimate_change: number;
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new FundValue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fund_code = source["fund_code"];
	        this.current_price = source["current_price"];
	        this.change_rate = source["change_rate"];
	        this.estimate_value = source["estimate_value"];
	        this.estimate_change = source["estimate_change"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
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
	export class GroupSummary {
	    group_name: string;
	    value: number;
	    cost: number;
	    gain: number;
	    gain_rate: number;
	    daily_gain: number;
	    daily_gain_rate: number;
	
	    static createFrom(source: any = {}) {
	        return new GroupSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.group_name = source["group_name"];
	        this.value = source["value"];
	        this.cost = source["cost"];
	        this.gain = source["gain"];
	        this.gain_rate = source["gain_rate"];
	        this.daily_gain = source["daily_gain"];
	        this.daily_gain_rate = source["daily_gain_rate"];
	    }
	}
	export class PortfolioSummary {
	    total_value: number;
	    total_cost: number;
	    total_gain: number;
	    total_gain_rate: number;
	    daily_gain: number;
	    daily_gain_rate: number;
	
	    static createFrom(source: any = {}) {
	        return new PortfolioSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_value = source["total_value"];
	        this.total_cost = source["total_cost"];
	        this.total_gain = source["total_gain"];
	        this.total_gain_rate = source["total_gain_rate"];
	        this.daily_gain = source["daily_gain"];
	        this.daily_gain_rate = source["daily_gain_rate"];
	    }
	}

}


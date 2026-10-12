export namespace models {
	
	export class FilterSettings {
	    sensitivityMode: string;
	    threshold: number;
	    filterShorts: boolean;
	    enabledCategories: Record<string, boolean>;
	    customBlockedKeywords: string[];
	    customBlockedChannels: string[];
	    whitelistedChannels: string[];
	
	    static createFrom(source: any = {}) {
	        return new FilterSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sensitivityMode = source["sensitivityMode"];
	        this.threshold = source["threshold"];
	        this.filterShorts = source["filterShorts"];
	        this.enabledCategories = source["enabledCategories"];
	        this.customBlockedKeywords = source["customBlockedKeywords"];
	        this.customBlockedChannels = source["customBlockedChannels"];
	        this.whitelistedChannels = source["whitelistedChannels"];
	    }
	}
	export class FilterSignal {
	    ruleName: string;
	    category: string;
	    weight: number;
	    confidence: number;
	    description: string;
	    matchedText?: string;
	
	    static createFrom(source: any = {}) {
	        return new FilterSignal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ruleName = source["ruleName"];
	        this.category = source["category"];
	        this.weight = source["weight"];
	        this.confidence = source["confidence"];
	        this.description = source["description"];
	        this.matchedText = source["matchedText"];
	    }
	}
	export class GlobalFilterStats {
	    totalVideosScanned: number;
	    totalVideosQuarantined: number;
	    totalVideosPassed: number;
	    serverUptimeSeconds: number;
	    // Go type: time
	    lastUpdated: any;
	
	    static createFrom(source: any = {}) {
	        return new GlobalFilterStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalVideosScanned = source["totalVideosScanned"];
	        this.totalVideosQuarantined = source["totalVideosQuarantined"];
	        this.totalVideosPassed = source["totalVideosPassed"];
	        this.serverUptimeSeconds = source["serverUptimeSeconds"];
	        this.lastUpdated = this.convertValues(source["lastUpdated"], null);
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
	export class Video {
	    id: string;
	    title: string;
	    channelName: string;
	    channelId: string;
	    channelHandle?: string;
	    views: string;
	    published: string;
	    duration: string;
	    description: string;
	    thumbnail: string;
	    embedUrl: string;
	    isShort: boolean;
	    aiScore: number;
	    humanScore: number;
	    verdict: string;
	    signals: FilterSignal[];
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.channelName = source["channelName"];
	        this.channelId = source["channelId"];
	        this.channelHandle = source["channelHandle"];
	        this.views = source["views"];
	        this.published = source["published"];
	        this.duration = source["duration"];
	        this.description = source["description"];
	        this.thumbnail = source["thumbnail"];
	        this.embedUrl = source["embedUrl"];
	        this.isShort = source["isShort"];
	        this.aiScore = source["aiScore"];
	        this.humanScore = source["humanScore"];
	        this.verdict = source["verdict"];
	        this.signals = this.convertValues(source["signals"], FilterSignal);
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
	export class SearchResponse {
	    query: string;
	    totalScanned: number;
	    totalClean: number;
	    totalQuarantined: number;
	    humanFeedRatio: number;
	    videos: Video[];
	    quarantined: Video[];
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.totalScanned = source["totalScanned"];
	        this.totalClean = source["totalClean"];
	        this.totalQuarantined = source["totalQuarantined"];
	        this.humanFeedRatio = source["humanFeedRatio"];
	        this.videos = this.convertValues(source["videos"], Video);
	        this.quarantined = this.convertValues(source["quarantined"], Video);
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
	
	export class VideoDetailResponse {
	    video: Video;
	    relatedVideos: Video[];
	    quarantinedRelated: Video[];
	
	    static createFrom(source: any = {}) {
	        return new VideoDetailResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.video = this.convertValues(source["video"], Video);
	        this.relatedVideos = this.convertValues(source["relatedVideos"], Video);
	        this.quarantinedRelated = this.convertValues(source["quarantinedRelated"], Video);
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


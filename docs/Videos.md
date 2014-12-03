#### Videos

---

*Get All Videos*

	GET - http://goapi.curtmfg.com/videos?key=[public api key]

*Get All CDNs*

	GET - http://goapi.curtmfg.com/videos/cdn?key=[public api key]

*Get All CDN Types*

	GET - http://goapi.curtmfg.com/videos/cdn/type?key=[public api key]

*Get All Channels*

	GET - http://goapi.curtmfg.com/videos/channel?key=[public api key]

*Get All Channel Types*

	GET - http://goapi.curtmfg.com/videos/channel/type?key=[public api key]

*Get All Video Types*

	GET - http://goapi.curtmfg.com/videos/type?key=[public api key]

*Get Distinct Videos*

	GET - http://goapi.curtmfg.com/videos/distinct?key=[public api key]

*Get Part Videos*

	GET - http://goapi.curtmfg.com/videos/part/<part id>?key=[public api key]

*Get Video Details*

	GET - http://goapi.curtmfg.com/videos/details/<video id>?key=[public api key]

*Get Video*

	GET - http://goapi.curtmfg.com/videos/<video id>?key=[public api key]

*Get CDN*

	GET - http://goapi.curtmfg.com/videos/cdn/<cdn id>?key=[public api key]

*Get CDN Type*

	GET - http://goapi.curtmfg.com/videos/cdn/type/<cdn type id>?key=[public api key]

*Get Channel*

	GET - http://goapi.curtmfg.com/videos/channel/<channel id>?key=[public api key]

*Get Channel Type*

	GET - http://goapi.curtmfg.com/videos/channel/type/<channel type id>?key=[public api key]

*Get Video Type*

	GET - http://goapi.curtmfg.com/videos/type/<video type id>?key=[public api key]

*Add Video*

	POST - http://goapi.curtmfg.com/videos?key=[public api key]

	JSON Payload:

	{
		"title"       : <video title (string)>,
		"description" : <video description (string)>,
		"isPrimary"   : <video is primary? (ex. "true" or "false") (string)>,
		"thumbnail"   : <video path to thumbnail (string)>,
		"videoType": {
			"id": <video type id (ex. 1) (int)>
		}
	}

*Add CDN*

	POST - http://goapi.curtmfg.com/videos/cdn?key=[public api key]

*Add CDN Type*

	POST - http://goapi.curtmfg.com/videos/cdn/type?key=[public api key]

*Add Channel*

	POST - http://goapi.curtmfg.com/videos/channel?key=[public api key]

*Add Channel Type*

	POST - http://goapi.curtmfg.com/videos/channel/type?key=[public api key]

*Add Video Type*

	POST - http://goapi.curtmfg.com/videos/type?key=[public api key]

*Update Video*

	POST - http://goapi.curtmfg.com/videos/<video id>?key=[public api key]

	JSON Payload:

	{
		"title"       : <video title (string)>,
		"description" : <video description (string)>,
		"isPrimary"   : <video is primary? (ex. "true" or "false") (string)>,
		"thumbnail"   : <video path to thumbnail (string)>,
		"videoType": {
			"id": <video type id (ex. 1) (int)>
		},
	}

*Update CDN*
	
	POST - http://goapi.curtmfg.com/videos/cdn/<cdn id>?key=[public api key]

*Update CDN Type*

	POST - http://goapi.curtmfg.com/videos/cdn/type/<cdn type id>?key=[public api key]

*Update Channel*

	POST - http://goapi.curtmfg.com/videos/channel/<channel id>?key=[public api key]

*Update Channel Type*

	POST - http://goapi.curtmfg.com/videos/channel/type/<channel type id>?key=[public api key]

*Update Video Type*
	
	POST - http://goapi.curtmfg.com/videos/type/<video type id>?key=[public api key]

*Delete Video*

	DELETE - http://goapi.curtmfg.com/videos/<video id>?key=[public api key]

*Delete CDN*

	DELETE - http://goapi.curtmfg.com/videos/cdn/<cdn id>?key=[public api key]

*Delete CDN Type*

	DELETE - http://goapi.curtmfg.com/videos/cdn/type/<cdn type id>?key=[public api key]

*Delete Channel*

	DELETE - http://goapi.curtmfg.com/videos/channel/<channel id>?key=[public api key]

*Delete Channel Type*

	DELETE - http://goapi.curtmfg.com/videos/channel/type/<channel type id>?key=[public api key]

*Delete Video Type*

	DELETE - http://goapi.curtmfg.com/videos/type/<video type id>?key=[public api key]


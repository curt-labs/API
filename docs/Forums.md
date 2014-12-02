#### Forums

---
*Get All Groups*

	GET - http://goapi.curtmfg.com/forum/groups?key=[public api key]

*Get All Topics*

	GET - http://goapi.curtmfg.com/forum/topics?key=[public api key]

*Get All Threads*

	GET - http://goapi.curtmfg.com/forum/threads?key=[public api key]

*Get All Posts*

	GET - http://goapi.curtmfg.com/forum/posts?key=[public api key]

*Get Group*

	GET - http://goapi.curtmfg.com/forum/groups/<group id>?key=[public api key]

*Get Topic*

	GET - http://goapi.curtmfg.com/forum/topics/<topic id>?key=[public api key]

*Get Thread*
	
	GET - http://goapi.curtmfg.com/forum/threads/<thread id>?key=[public api key]

*Get Post*

	GET - http://goapi.curtmfg.com/forum/posts/<post id>?key=[public api key]

*Add Group*

	POST - http://goapi.curtmfg.com/forum/groups?key=[public api key]

	Form Payload:

		"name"        : <group name (string)>,
		"description" : <group description (string)>

*Add Topic*

	POST - http://goapi.curtmfg.com/forum/topics?key=[public api key]

	Form Payload:

		"groupID"     : <topic groupID (ex. "1") (string)>,
		"closed"      : <topic is closed? "true" or "false" (string)>,
		"name"        : <topic name (string)>,
		"description" : <topic description (string)>,
		"image"       : <topic image path (string)>

*Add Post*

	POST - http://goapi.curtmfg.com/forum/posts?key=[public api key]

	Form Payload:

		"topicID" : <post topic id (ex. "1") (string)>,
		"parentID": <post parent id (ex. "1") (string)>,
		"notify"  : <notifies poster of replies (ex. "true" or "false") (string)>,
		"sticky"  : <post is sticky? (ex. "true" or "false") (string)>,
		"title"   : <post title (string)>,
		"post"    : <post text (string)>,
		"name"    : <name of the person that posted this post (string)>,
		"email"   : <email of the person that posted this post (string)>,
		"company" : <company name for the person that posted this post (string)>

>> Note: When adding a new post, topicID is used to create a new post, while parentID is used to reply to an existing post.

*Update Group*

	PUT - http://goapi.curtmfg.com/forum/groups/<group id>?key=[public api key]

	Form Payload:

		"name"        : <group name (string)>,
		"description" : <group description (string)>

*Update Topic*

	PUT - http://goapi.curtmfg.com/forum/topics/<topic id>?key=[public api key]

	Form Payload:

		"groupID"     : <topic groupID (ex. "1") (string),
		"closed"      : <topic is closed? (ex. "true" or "false") (string)>,
		"active"      : <topic is active? (ex. "true" or "false") (string)>,
		"name"        : <topic name (string)>,
		"description" : <topic description (string)>,
		"image"       : <topic image path (string)>

*Update Post*

	PUT - http://goapi.curtmfg.com/forum/posts/<post id>?key=[public api key]

	Form Payload:

		"parentID" : <post parent id (ex. "1") (string)>,
		"threadID" : <post thread id (ex. "1") (string)>,
		"approved" : <post is approved? (ex. "true" or "false") (string)>,
		"active"   : <post is active? (ex. "true" or "false") (string)>,
		"notify"   : <notifies poster of replies (ex. "true" or "false") (string)>,
		"sticky"   : <post is sticky? (ex. "true" or "false") (string)>,
		"flag"     : <post is flagged as spam or inappropriate? (ex. "true" or "false") (string)>,
		"title"    : <post title (string)>,
		"post"     : <post text (string)>,
		"name"     : <name of the person that posted this post (string)>,
		"email"    : <email of the person that posted this post (string)>,
		"company"  : <company name for the person that posted this post (string)>

*Delete Group*

	DELETE - http://goapi.curtmfg.com/forum/groups/<group id>?key=[public api key]

>> Note: Deleting a group will delete the group itself and all tied topics, threads, and posts.

*Delete Topic*

	DELETE - http://goapi.curtmfg.com/forum/topics/<topic id>?key=[public api key]

>> Note: Deleting a topic will delete the topic itself and all tied threads and posts.

*Delete Thread*

	DELETE - http://goapi.curtmfg.com/forum/threads/<thread id>?key=[public api key]

>> Note: Deleting a thread will delete the thread itself and all posts tied to it.

*Delete Post*

	DELETE - http://goapi.curtmfg.com/forum/posts/<post id>?key=[public api key]



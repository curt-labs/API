package forum

var (
	MapToGroupID  = 1
	MapToTopicID  = 2
	MapToThreadID = 3
)

func (t Topics) ToMap(selector ...int) map[interface{}]interface{} {
	topicsmap := make(map[interface{}]interface{})
	for _, topic := range t {
		switch {
		case len(selector) > 0 && selector[0] == MapToGroupID:
			if _, found := topicsmap[topic.GroupID]; !found {
				topicsmap[topic.GroupID] = Topics{}
			}
			topicsmap[topic.GroupID] = append(topicsmap[topic.GroupID].(Topics), topic)
		default:
			topicsmap[topic.ID] = topic
		}
	}
	return topicsmap
}

func (t Threads) ToMap(selector ...int) map[interface{}]interface{} {
	threadsmap := make(map[interface{}]interface{})
	for _, thread := range t {
		switch {
		case len(selector) > 0 && selector[0] == MapToTopicID:
			if _, found := threadsmap[thread.TopicID]; !found {
				threadsmap[thread.TopicID] = Threads{}
			}
			threadsmap[thread.TopicID] = append(threadsmap[thread.TopicID].(Threads), thread)
		default:
			threadsmap[thread.ID] = thread
		}
	}
	return threadsmap
}

func (p Posts) ToMap(selector ...int) map[interface{}]interface{} {
	postsmap := make(map[interface{}]interface{})
	for _, post := range p {
		switch {
		case len(selector) > 0 && selector[0] == MapToThreadID:
			if _, found := postsmap[post.ThreadID]; !found {
				postsmap[post.ThreadID] = Posts{}
			}
			postsmap[post.ThreadID] = append(postsmap[post.ThreadID].(Posts), post)
		default:
			postsmap[post.ID] = post
		}
	}
	return postsmap
}

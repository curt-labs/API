package cache

import (
	"encoding/json"
	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/models/customer"
	"net/http"
)

//TODO check for super user
func GetKeys(w http.ResponseWriter, r *http.Request) {
	if !approveuser(r) {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	namespaces, err := redis.GetNamespaces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(namespaces)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(j))
}

func GetByKey(w http.ResponseWriter, r *http.Request) {
	if !approveuser(r) {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	redis_key := r.URL.Query().Get("redis_key")
	redis_namespace := r.URL.Query().Get("redis_namespace")

	key := redis_namespace + ":" + redis_key
	res, err := redis.GetFullPath(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(j))
}

func DeleteKey(w http.ResponseWriter, r *http.Request) {
	if !approveuser(r) {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	redis_key := r.URL.Query().Get("redis")
	err := redis.DeleteFullPath(redis_key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func approveuser(r *http.Request) bool {
	api := r.URL.Query().Get("key")
	if api == "" {
		return false
	}
	c := customer.Customer{}
	err := c.GetCustomerIdFromKey(api)
	if err != nil || c.Id == 0 {
		return false
	}
	return true
}

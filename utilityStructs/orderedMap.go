package utilityStructs

import (
	"go-api/repository"
	"sync"
)

type OrderedMap struct {
	Data map[string]*repository.User
	Keys []string
	Mu   *sync.RWMutex
}

func (o *OrderedMap) shift() {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	_ = o.Data[o.Keys[0]]
	delete(o.Data, o.Keys[0])
	o.Keys = o.Keys[1:]
}

func (o *OrderedMap) InsertUser(user *repository.User) {
	if len(o.Keys) == 100 {
		o.shift()
	}
	o.Mu.Lock()
	defer o.Mu.Unlock()
	userID := user.ID
	o.Data[userID] = user
	o.Keys = append(o.Keys, userID)
}

func (o *OrderedMap) GetUser(userID string) (*repository.User, bool) {
	user, ok := o.Data[userID]
	if ok {
		return user, true
	} else {
		return nil, false
	}
}

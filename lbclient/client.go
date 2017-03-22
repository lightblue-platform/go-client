package lbclient

import ()

type DataServiceClient interface {
	Find(request *FindRequest, returndata interface{}) (*Response, error)
	Insert(request *InsertRequest, returndata interface{}) (*Response, error)
	Update(request *UpdateRequest, returndata interface{}) (*Response, error)
	Save(request *SaveRequest, returndata interface{}) (*Response, error)
	Delete(request *DeleteRequest) (*Response, error)
}

type LockingClient interface {
	Acquire(domain, callerId, resourceId string, ttl int) (bool, error)
	Release(domain, callerId, resourceId string) (bool, error)
	GetLockCount(domain, callerId, resourceId string) (int, error)
	Ping(domain, callerId, resourceId string) (bool, error)
}

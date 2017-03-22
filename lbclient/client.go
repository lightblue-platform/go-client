package lbclient

import ()

type LBClient interface {
	Find(request *FindRequest, data interface{}) (*Response, error)
	Insert(request *InsertRequest, data interface{}) (*Response, error)
	Update(request *UpdateRequest, data interface{}) (*Response, error)
	Save(request *SaveRequest, data interface{}) (*Response, error)
	Delete(request *DeleteRequest) (*Response, error)
}

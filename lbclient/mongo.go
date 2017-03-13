package lbclient

type ReadPreference string

const (
	NEAREST             ReadPreference = "nearest"
	PRIMARY_PREFERRED   ReadPreference = "primaryPreferred"
	PRIMARY             ReadPreference = "primary"
	SECONDARY           ReadPreference = "secondary"
	SECONDARY_PREFERRED ReadPreference = "secondaryPreferred"
)

type MongoExecutionOptions struct {
	ReadPref     ReadPreference `json:"readPreference"`
	WriteConcern string         `json:"writeConcern"`
	MaxQueryTime int            `json:"maxQueryTimeMS:`
}

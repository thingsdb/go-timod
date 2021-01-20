package timod

// Proto is used as protocol type used by ThingsDB.
type Proto int8

const (
	// ProtoModuleConf - when configuration data for the module is received
	ProtoModuleConf Proto = 64

	// ProtoModuleConfOk - respond after successfully configuring the module
	ProtoModuleConfOk Proto = 65

	// ProtoModuleConfErr - respond with a configuration error
	ProtoModuleConfErr Proto = 66

	// ProtoModuleReq - when a request is received
	ProtoModuleReq Proto = 80

	// ProtoModuleRes is used to respond to a ProtoModuleReq package
	ProtoModuleRes Proto = 81

	// ProtoModuleErr is used to respond to a ProtoModuleReq with an error
	ProtoModuleErr Proto = 82
)

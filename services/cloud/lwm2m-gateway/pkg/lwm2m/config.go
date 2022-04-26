package lwm2m

type Lwm2mop string
type ID int32

type Lwm2mContentType uint32

const (
	READ    = "R"
	WRITE   = "W"
	EXECUTE = "E"
)

const (
	SCHEMA_PREFIX  = "specs/schema/"
	SCHEMA_POSTFIX = ".toml"
)

const (
	Read       = "read"
	Write      = "write"
	Discover   = "disc"
	Observe    = "observe"
	Cancel     = "cancel"
	ClientList = "list"
	Execute    = "exec"
)

const (
	COAP_NO_ERROR                  = 0x00
	COAP_IGNORE                    = 0x01
	COAP_201_CREATED               = 0x41
	COAP_202_DELETED               = 0x42
	COAP_204_CHANGED               = 0x44
	COAP_205_CONTENT               = 0x45
	COAP_231_CONTINUE              = 0x5F
	COAP_400_BAD_REQUEST           = 0x80
	COAP_401_UNAUTHORIZED          = 0x81
	COAP_402_BAD_OPTION            = 0x82
	COAP_404_NOT_FOUND             = 0x84
	COAP_405_METHOD_NOT_ALLOWED    = 0x85
	COAP_406_NOT_ACCEPTABLE        = 0x86
	COAP_408_REQ_ENTITY_INCOMPLETE = 0x88
	COAP_412_PRECONDITION_FAILED   = 0x8C
	COAP_413_ENTITY_TOO_LARGE      = 0x8D
	COAP_500_INTERNAL_SERVER_ERROR = 0xA0
	COAP_501_NOT_IMPLEMENTED       = 0xA1
	COAP_503_SERVICE_UNAVAILABLE   = 0xA3
)

var ObjectIDList = [...]ID{0, 1, 2, 3, 4, 5, 6, 7, 3200, 3201, 3203, 3303, 3311, 3316, 3317, 3328, 34567, 34568, 34569, 34570}

const (
	Lwm2mserver   = "0.0.0.0:3000"
	Gatewayserver = "0.0.0.0:3100"
)

const (
	LWM2M_CONTENT_TEXT       = 0 // Also used as undefined
	LWM2M_CONTENT_LINK       = 40
	LWM2M_CONTENT_OPAQUE     = 42
	LWM2M_CONTENT_TLV_OLD    = 1542 // Keep old value for backward-compatibility
	LWM2M_CONTENT_TLV        = 11542
	LWM2M_CONTENT_JSON_OLD   = 1543 // Keep old value for backward-compatibility
	LWM2M_CONTENT_JSON       = 11543
	LWM2M_CONTENT_SENML_JSON = 110
)

package utils

var JitoEndpoints = map[string]JitoEndpointInfo{
	"AMS": {
		BlockEngineUrl:    "amsterdam.mainnet.block-engine.jito.wtf:443",
		RelayerUrl:        "http://amsterdam.mainnet.relayer.jito.wtf:8100",
		ShredReceiverAddr: "74.118.140.240:1002",
		Ntp:               "ntp.amsterdam.jito.wtf",
	},
	"FFM": {
		BlockEngineUrl:    "frankfurt.mainnet.block-engine.jito.wtf:443",
		RelayerUrl:        "http://frankfurt.mainnet.relayer.jito.wtf:8100",
		ShredReceiverAddr: "145.40.93.84:1002",
		Ntp:               "ntp.frankfurt.jito.wtf",
	},
	"NY": {
		BlockEngineUrl:    "ny.mainnet.block-engine.jito.wtf:443",
		RelayerUrl:        "http://ny.mainnet.relayer.jito.wtf:8100",
		ShredReceiverAddr: "141.98.216.96:1002",
		Ntp:               "ntp.dallas.jito.wtf",
	},
	"TKY": {
		BlockEngineUrl:    "tokyo.mainnet.block-engine.jito.wtf:443",
		RelayerUrl:        "http://tokyo.mainnet.relayer.jito.wtf:8100",
		ShredReceiverAddr: "202.8.9.160:1002",
		Ntp:               "ntp.tokyo.jito.wtf",
	},
}
var (
	LIQUIDITY_FEES_NUMERATOR   = 25
	LIQUIDITY_FEES_DENOMINATOR = 10000
	RAY_V4                     = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
)

var Amsterdam = JitoEndpoints["AMS"]
var Frankfurt = JitoEndpoints["FFM"]
var NewYork = JitoEndpoints["NY"]
var Tokyo = JitoEndpoints["TKY"]

package types

const (
	ModuleName   = "msmp"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = "mem_msmp"
)

var (
	PoolKey             = []byte{0x01}
	StakerKey           = []byte{0x02}
	DistributionKey     = []byte{0x03}
	ActivityPointsKey   = []byte{0x04}
	FoundationAccountKey = []byte{0x05}
)

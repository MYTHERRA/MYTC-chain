package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/baseapp"
"github.com/cosmos/cosmos-sdk/client"
"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
"github.com/cosmos/cosmos-sdk/client/rpc"
"github.com/cosmos/cosmos-sdk/codec"
"github.com/cosmos/cosmos-sdk/codec/types"
"github.com/cosmos/cosmos-sdk/server/api"
"github.com/cosmos/cosmos-sdk/server/config"
servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"

storetypes "github.com/cosmos/cosmos-sdk/store/types"
sdk "github.com/cosmos/cosmos-sdk/types"
"github.com/cosmos/cosmos-sdk/types/module"
"github.com/cosmos/cosmos-sdk/version"
"github.com/cosmos/cosmos-sdk/x/auth"
authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
"github.com/cosmos/cosmos-sdk/x/auth/vesting"
vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
"github.com/cosmos/cosmos-sdk/x/bank"
bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
"github.com/cosmos/cosmos-sdk/x/capability"
capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
"github.com/cosmos/cosmos-sdk/x/crisis"
crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
distr "github.com/cosmos/cosmos-sdk/x/distribution"
distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
"github.com/cosmos/cosmos-sdk/x/evidence"
evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
"github.com/cosmos/cosmos-sdk/x/feegrant"
feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
"github.com/cosmos/cosmos-sdk/x/params"
paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
"github.com/cosmos/cosmos-sdk/x/slashing"
slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
"github.com/cosmos/cosmos-sdk/x/staking"
stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
"github.com/cosmos/cosmos-sdk/x/upgrade"
upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/mytherra/mytc/x/lockup"
	lockupkeeper "github.com/mytherra/mytc/x/lockup/keeper"
	lockuptypes "github.com/mytherra/mytc/x/lockup/types"

	"github.com/mytherra/mytc/x/msmp"
	msmpkeeper "github.com/mytherra/mytc/x/msmp/keeper"
	msmptypes "github.com/mytherra/mytc/x/msmp/types"

"github.com/spf13/cast"
abci "github.com/tendermint/tendermint/abci/types"
tmjson "github.com/tendermint/tendermint/libs/json"
"github.com/tendermint/tendermint/libs/log"
tmos "github.com/tendermint/tendermint/libs/os"
dbm "github.com/tendermint/tm-db"
)

const (
AppName = "mytc"
)

var (
// DefaultNodeHome default home directories for the application daemon
DefaultNodeHome string

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
ModuleBasics = module.NewBasicManager(
auth.AppModuleBasic{},
genutil.AppModuleBasic{},
bank.AppModuleBasic{},
capability.AppModuleBasic{},
staking.AppModuleBasic{},
mint.AppModuleBasic{},
distr.AppModuleBasic{},
gov.NewAppModuleBasic(
paramsclient.ProposalHandler,
distrclient.ProposalHandler,
		upgradeclient.ProposalHandler,
		upgradeclient.CancelProposalHandler,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	vesting.AppModuleBasic{},
	lockup.AppModuleBasic{},
	msmp.AppModuleBasic{},
)
// module account permissions
maccPerms = map[string][]string{
authtypes.FeeCollectorName:     nil,
distrtypes.ModuleName:          nil,
minttypes.ModuleName:           {authtypes.Minter},
stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	lockuptypes.ModuleName:         nil,
}
)

var (
_ servertypes.Application = (*MYTCApp)(nil)
_ simapp.App              = (*MYTCApp)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultNodeHome = filepath.Join(userHomeDir, "."+AppName)

	// Set Mytherra bech32 address prefixes
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("mytc", "mytcpub")
	cfg.SetBech32PrefixForValidator("mytcvaloper", "mytcvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("mytcvalcons", "mytcvalconspub")
	cfg.Seal()
}

// MYTCApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type MYTCApp struct {
*baseapp.BaseApp

cdc               *codec.LegacyAmino
appCodec          codec.Codec
interfaceRegistry types.InterfaceRegistry

invCheckPeriod uint

// keys to access the substores
keys    map[string]*storetypes.KVStoreKey
tkeys   map[string]*storetypes.TransientStoreKey
memKeys map[string]*storetypes.MemoryStoreKey

// keepers
AccountKeeper    authkeeper.AccountKeeper
BankKeeper       bankkeeper.Keeper
CapabilityKeeper *capabilitykeeper.Keeper
StakingKeeper    stakingkeeper.Keeper
SlashingKeeper   slashingkeeper.Keeper
MintKeeper       mintkeeper.Keeper
DistrKeeper      distrkeeper.Keeper
GovKeeper        govkeeper.Keeper
CrisisKeeper     crisiskeeper.Keeper
UpgradeKeeper    upgradekeeper.Keeper
ParamsKeeper     paramskeeper.Keeper
EvidenceKeeper   evidencekeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	LockupKeeper     lockupkeeper.Keeper
	MSMPKeeper       msmpkeeper.Keeper

	// Module Manager
mm *module.Manager

// Configurator
configurator module.Configurator
}

// NewMYTCApp returns a reference to an initialized MYTCApp.
func NewMYTCApp(
logger log.Logger,
db dbm.DB,
traceStore io.Writer,
loadLatest bool,
skipUpgradeHeights map[int64]bool,
homePath string,
invCheckPeriod uint,
encodingConfig simappparams.EncodingConfig,
appOpts servertypes.AppOptions,
baseAppOptions ...func(*baseapp.BaseApp),
) *MYTCApp {
appCodec := encodingConfig.Marshaler
cdc := encodingConfig.Amino
interfaceRegistry := encodingConfig.InterfaceRegistry

bApp := baseapp.NewBaseApp(AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
bApp.SetCommitMultiStoreTracer(traceStore)
bApp.SetVersion(version.Version)
bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey,
		feegrant.StoreKey, evidencetypes.StoreKey, capabilitytypes.StoreKey,
		lockuptypes.StoreKey,
		msmptypes.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, lockuptypes.MemStoreKey, msmptypes.MemStoreKey)

	app := &MYTCApp{

BaseApp:           bApp,
cdc:               cdc,
appCodec:          appCodec,
interfaceRegistry: interfaceRegistry,
invCheckPeriod:    invCheckPeriod,
keys:              keys,
tkeys:             tkeys,
memKeys:           memKeys,
}

app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// Applications that wish to enforce statically created ScopedKeepers should call Seal after creating
	// their scoped keepers in NewApp with ScopeToModule
	app.CapabilityKeeper.Seal()

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.AccountKeeper = authkeeper.NewAccountKeeper(
appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.ModuleAccountAddrs(),
	)
	
	stakingKeeper := stakingkeeper.NewKeeper(
appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
)
	
	app.MintKeeper = mintkeeper.NewKeeper(
appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName), &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
	)
	
	app.DistrKeeper = distrkeeper.NewKeeper(
appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	
	app.SlashingKeeper = slashingkeeper.NewKeeper(
appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	
	app.CrisisKeeper = crisiskeeper.NewKeeper(
app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
appCodec, keys[feegrant.StoreKey], app.AccountKeeper,
)
	
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp,
	)

	app.LockupKeeper = *lockupkeeper.NewKeeper(
		appCodec,
		keys[lockuptypes.StoreKey],
		memKeys[lockuptypes.MemStoreKey],
		app.BankKeeper,
	)

	app.MSMPKeeper = *msmpkeeper.NewKeeper(
		appCodec,
		keys[msmptypes.StoreKey],
		memKeys[msmptypes.MemStoreKey],
		app.GetSubspace(msmptypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
	)

	// Register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it can be modified like below:
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute("params", params.NewParamChangeProposalHandler(app.ParamsKeeper)).

		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter,
	)

	// Create static IBC router, add transfer route, then set and seal it
	// ... IBC setup ...

	/****  Module Options ****/

	// NOTE: we may consider skipping genesis at some point in the future
	// incase we decide to create a custom genesis file.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
auth.NewAppModule(appCodec, app.AccountKeeper, nil),
vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
capability.NewAppModule(appCodec, *app.CapabilityKeeper),
crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		params.NewAppModule(app.ParamsKeeper),
		lockup.NewAppModule(appCodec, app.LockupKeeper),
		msmp.NewAppModule(appCodec, app.MSMPKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName, capabilitytypes.ModuleName, minttypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName,
		evidencetypes.ModuleName, stakingtypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName, govtypes.ModuleName,
		crisistypes.ModuleName, genutiltypes.ModuleName, feegrant.ModuleName, paramstypes.ModuleName, vestingtypes.ModuleName,
		lockuptypes.ModuleName,
		msmptypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName, govtypes.ModuleName, stakingtypes.ModuleName, capabilitytypes.ModuleName, authtypes.ModuleName,
		banktypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName, minttypes.ModuleName,
		genutiltypes.ModuleName, evidencetypes.ModuleName, feegrant.ModuleName, upgradetypes.ModuleName, paramstypes.ModuleName,
		vestingtypes.ModuleName,
		lockuptypes.ModuleName,
		msmptypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		feegrant.ModuleName,
		lockuptypes.ModuleName,
		msmptypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// Name returns the name of the App
func (app *MYTCApp) Name() string { return app.BaseApp.Name() }

// GetBaseApp returns the reference to the BaseApp
func (app *MYTCApp) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }

// BeginBlocker application updates every begin block
func (app *MYTCApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *MYTCApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *MYTCApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *MYTCApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *MYTCApp) ModuleAccountAddrs() map[string]bool {
modAccAddrs := make(map[string]bool)
for acc := range maccPerms {
modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
}

return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *MYTCApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns an app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *MYTCApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns an InterfaceRegistry
func (app *MYTCApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *MYTCApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *MYTCApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *MYTCApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *MYTCApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *MYTCApp) SimulationManager() *module.SimulationManager {
	return nil
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *MYTCApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *MYTCApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *MYTCApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(msmptypes.ModuleName)

	return paramsKeeper
}

// ExportAppStateAndValidators exports the state of the application for a genesis file.
func (app *MYTCApp) ExportAppStateAndValidators(
	forZeroHeight bool, jailAllowedAddrs []string,
) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.PrepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}

	genState := app.mm.ExportGenesis(ctx, app.appCodec)
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, app.StakingKeeper)
	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}

// PrepForZeroHeightGenesis prepares for fresh start at zero height
func (app *MYTCApp) PrepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) {
	// TODO: Implementation
}

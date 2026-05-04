package types

// MSMP module event types and attributes.
const (
	EventTypeCollectFee          = "collect_fee"
	EventTypeDistributeRewards   = "distribute_rewards"
	EventTypeClaimActivityPoints = "claim_activity_points"

	AttributeKeySender      = "sender"
	AttributeKeyDistributor = "distributor"
	AttributeKeyAmount      = "amount"
	AttributeKeyPoints      = "points"
)

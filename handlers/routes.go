package handlers

const V1Prefix = "/api/v1"

const (
	IncrementEventPattern    = V1Prefix + "/:item_type/:item_id/:event_type/:user_id"
	GetMetricsPattern        = V1Prefix + "/metrics/:item_type/:item_id"
	GetMetricsListPattern    = V1Prefix + "/metrics/:item_type"
	GetHistogramPattern      = V1Prefix + "/histogram/:item_type/:item_id"
	GetRecentActivityPattern = V1Prefix + "/activity/:item_type/:item_id"
)

const (
	IncrementEventURL    = V1Prefix + "/%s/%s/%s/%s"
	GetMetricsURL        = V1Prefix + "/metrics/%s/%s"
	GetMetricsListURL    = V1Prefix + "/metrics/%s?item_ids=%s"
	GetHistogramURL      = V1Prefix + "/histogram/%s/%s"
	GetRecentActivityURL = V1Prefix + "/activity/%s/%s"
)

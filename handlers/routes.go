package handlers

const V1Prefix = "/api/v1"

// Fiber route patterns (usa :param para Fiber)
const (
	IncrementEventPattern    = V1Prefix + "/:item_type/:item_id/:event_type/:user_id"
	GetMetricsPattern        = V1Prefix + "/:item_type/:item_id"
	GetHistogramPattern      = V1Prefix + "/histogram/:item_type/:item_id"
	GetRecentActivityPattern = V1Prefix + "/activity/:item_type/:item_id"
)

// URL format strings (usa %s para fmt.Sprintf)
const (
	IncrementEventURL    = V1Prefix + "/%s/%s/%s/%s"
	GetMetricsURL        = V1Prefix + "/%s/%s"
	GetHistogramURL      = V1Prefix + "/histogram/%s/%s"
	GetRecentActivityURL = V1Prefix + "/activity/%s/%s"
)

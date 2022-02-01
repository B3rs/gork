package jobs

type IncreaseArgs struct {
	IncreaseThis int `json:"increase_this"`
}

type IncreaseResult struct {
	Increased int `json:"increased"`
}

type LowerizeArgs struct {
	LowerizeThis string `json:"lowerize_this"`
}
type LowerizeResult struct {
	Lowerized string `json:"lowerized"`
}

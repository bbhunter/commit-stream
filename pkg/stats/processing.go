package stats

type ProcessingStats struct {
	IncomingRate  uint32
	ProcessedRate uint32
	FilteredRate  uint32
	Total         uint32
}

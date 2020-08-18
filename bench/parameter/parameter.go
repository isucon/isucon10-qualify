package parameter

import "time"

const (
	NumOfCheckChairSearchPaging    = 3
	NumOfCheckEstateSearchPaging   = 3
	PerPageOfChairSearch           = 30
	PerPageOfEstateSearch          = 30
	MaxLengthOfNazotteResponse     = 50
	NeighborhoodRadiusOfNazotte    = 1e-6
	SleepTimeOnFailScenario        = 1500 * time.Millisecond
	SleepSwingOnFailScenario       = 500 // * time.Millisecond
	SleepTimeOnUserAway            = 500 * time.Millisecond
	SleepSwingOnUserAway           = 100 // * time.Millisecond
	SleepTimeOnBotInterval         = 500 * time.Millisecond
	SleepSwingOnBotInterval        = 100 // * time.Millisecond
	IntervalForCheckWorkers        = 5 * time.Second
	ThresholdTimeOfAbandonmentPage = 1000 * time.Millisecond
	DefaultAPITimeout              = 2000 * time.Millisecond
	InitializeTimeout              = 30 * time.Second
	VerifyTimeout                  = 10 * time.Second
	LoadTimeout                    = 60 * time.Second
)

// IncListOfWorkers 前のレベルとのWorkerの個数の差分を保持するList
// [level][0]: inc of EstateSearchWorker
// [level][1]: inc of ChairSearchWorker
// [level][2]: inc of EstateNazotteSearchWorker
// [level][3]: inc of BotWorker
var IncListOfWorkers = [][4]int{
	{1, 1, 1, 1}, // level 00
	{1, 1, 1, 1}, // level 01
	{1, 1, 1, 1}, // level 02
	{1, 1, 1, 1}, // level 03
	{1, 1, 1, 1}, // level 04
	{1, 1, 1, 1}, // level 05
	{1, 1, 1, 1}, // level 06
	{1, 1, 1, 1}, // level 07
	{1, 1, 1, 1}, // level 08
	{1, 1, 1, 1}, // level 09
	{1, 1, 1, 1}, // level 10
	{1, 1, 1, 1}, // level 11
	{1, 1, 1, 1}, // level 12
	{1, 1, 1, 1}, // level 13
	{1, 1, 1, 1}, // level 14
	{1, 1, 1, 1}, // level 15
	{1, 1, 1, 1}, // level 16
	{1, 1, 1, 1}, // level 17
	{1, 1, 1, 1}, // level 18
	{1, 1, 1, 1}, // level 19
}

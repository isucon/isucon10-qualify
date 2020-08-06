package paramater

import "time"

const (
	NumOfInitialEstateSearchWorker        = 1
	NumOfInitialChairSearchWorker         = 1
	NumOfInitialEstateNazotteSearchWorker = 1
	NumOfInitialBotWorker                 = 1
	NumOfCheckChairSearchPaging           = 3
	NumOfCheckEstateSearchPaging          = 3
	PerPageOfChairSearch                  = 30
	PerPageOfEstateSearch                 = 30
	MaxLengthOfNazotteResponse            = 50
	NeighborhoodRadiusOfNazotte           = 1E-6
	SleepTimeOnFailScenario               = 1500 * time.Millisecond
	SleepSwingOnFailScenario              = 500 // * time.Millisecond
	SleepTimeOnUserAway                   = 500 * time.Millisecond
	SleepSwingOnUserAway                  = 100 // * time.Millisecond
	SleepTimeOnBotInterval                = 500 * time.Millisecond
	SleepSwingOnBotInterval               = 100 // * time.Millisecond
	IntervalForCheckWorkers               = 5 * time.Second
	ThresholdTimeOfAbandonmentPage        = 1000 * time.Millisecond
	DefaultAPITimeout                     = 2000 * time.Millisecond
	InitializeTimeout                     = 30 * time.Second
	VerifyTimeout                         = 10 * time.Second
	LoadTimeout                           = 60 * time.Second
)

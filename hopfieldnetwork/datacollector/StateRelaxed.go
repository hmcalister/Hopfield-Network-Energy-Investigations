package datacollector

// Representation of the result of a relaxation of a state.
//
// State is the state vector that has been relaxed.
//
// UnitEnergies is a vector representing the energies of each unit.
//
// Stable is a bool representing if the state was stable when relaxation finished.
//
// NumSteps is an int representing the number of steps taken when relaxation finished.
//
// DistancesToAllLearned is an array of distances to all learned states.
type StateRelaxedData struct {
	// TrialIndex int `parquet:"name=TrialIndex, type=INT32"`
	StateIndex int `parquet:"name=StateIndex, type=INT32"`
	// State              []float64 `parquet:"name=State, type=DOUBLE, repetitiontype=REPEATED"`
	// Stable             bool      `parquet:"name=Stable, type=BOOL"`
	// NumSteps           int       `parquet:"name=NumSteps, type=DOUBLE"`
	// StateEnergyVector  []float64 `parquet:"name=StateEnergyVector, type=DOUBLE, repetitiontype=REPEATED"`
	// DistancesToLearned []float64 `parquet:"name=DistancesToLearned, type=DOUBLE, repetitiontype=REPEATED"`
}

type onStateRelaxedHandler struct {
	defaultDataHandler
}

func newOnStateRelaxedCollector(dataFile string) *onStateRelaxedHandler {
	return &onStateRelaxedHandler{
		defaultDataHandler: defaultDataHandler{
			eventID:    DataCollectionEvent_OnStateRelax,
			dataWriter: newParquetWriter(dataFile, new(StateRelaxedData)),
		},
	}
}

func (collector *onStateRelaxedHandler) handleEvent(event interface{}) {
	result := event.(StateRelaxedData)
	collector.dataWriter.Write(result)
}
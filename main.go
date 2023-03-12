package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/pkg/profile"

	"hmcalister/hopfield/hopfieldnetwork"
	"hmcalister/hopfield/hopfieldnetwork/datacollector"
	"hmcalister/hopfield/hopfieldnetwork/states"
	"hmcalister/hopfield/hopfieldutils"
)

const DIMENSION = 100
const TARGET_STATES = 100
const LEARNING_RULE = hopfieldnetwork.DeltaLearningRule
const UNITS_UPDATED = 5

var (
	numTrials                *int
	numTestStates            *int
	numThreads               *int
	relaxationResultDataPath *string
	trialEndDataPath         *string
	InfoLogger               *log.Logger
)

func init() {
	numTrials = flag.Int("trials", 1000, "The number of trials to undertake.")
	numTestStates = flag.Int("testStates", 1000, "The number of test states to use for each trial.")
	relaxationResultDataPath = flag.String("relaxationResultDataFile", "data/relaxationResult.pq", "The file to write relaxation result data to. Data is in a parquet format.")
	trialEndDataPath = flag.String("trialEndDataFile", "data/trialEndData.pq", "The file to write trial end data to. Data is in a parquet format.")
	numThreads = flag.Int("threads", 1, "The number of threads to use for relaxation.")
	var logFilePath = flag.String("logFile", "logs/log.txt", "The file to write logs to.")
	flag.Parse()

	logFile, err := os.Create(*logFilePath)
	if err != nil {
		panic("Could not open log file!")
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Main method for entry point
func main() {
	defer profile.Start(profile.ClockProfile, profile.ProfilePath("./profiles"), profile.NoShutdownHook).Stop()
	keyboardInterrupt := make(chan os.Signal, 1)
	signal.Notify(keyboardInterrupt, os.Interrupt)

	collector := datacollector.NewDataCollector().
		AddHandler(datacollector.NewRelaxationResultHandler(*relaxationResultDataPath)).
		AddHandler(datacollector.NewTrialEndHandler(*trialEndDataPath))
	defer collector.WriteStop()
	go collector.CollectData()

TrialLoop:
	for trial := 0; trial < *numTrials; trial++ {
		select {
		case <-keyboardInterrupt:
			InfoLogger.Printf("RECEIVED KEYBOARD INTERRUPT")
			break TrialLoop
		default:
		}
		InfoLogger.Printf("----- TRIAL: %09d -----", trial)
		InfoLogger.Printf("Goroutines: %d\n", runtime.NumGoroutine())

		network := hopfieldnetwork.NewHopfieldNetworkBuilder().
			SetNetworkDimension(DIMENSION).
			SetRandMatrixInit(false).
			SetNetworkLearningRule(LEARNING_RULE).
			SetEpochs(100).
			SetMaximumRelaxationIterations(100).
			SetMaximumRelaxationUnstableUnits(0).
			SetUnitsUpdatedPerStep(UNITS_UPDATED).
			SetDataCollector(collector).
			Build()

		stateGenerator := states.NewStateGeneratorBuilder().
			SetRandMin(-1).
			SetRandMax(1).
			SetGeneratorDimension(DIMENSION).
			Build()

		targetStates := stateGenerator.CreateStateCollection(TARGET_STATES)
		network.LearnStates(targetStates)

		testStates := stateGenerator.CreateStateCollection(*numTestStates)
		relaxationResults := network.ConcurrentRelaxStates(testStates, *numThreads)

		trialNumStable := 0
		trialStableStepsTaken := 0
		for stateIndex, result := range relaxationResults {
			event := datacollector.RelaxationResultData{
				TrialIndex:         trial,
				StateIndex:         stateIndex,
				Stable:             result.Stable,
				NumSteps:           result.NumSteps,
				DistancesToLearned: result.DistancesToLearned,
				EnergyProfile:      result.EnergyProfile,
			}

			collector.EventChannel <- hopfieldutils.IndexedWrapper[interface{}]{
				Index: datacollector.DataCollectionEvent_RelaxationResult,
				Data:  event,
			}

			if result.Stable {
				trialNumStable += 1
				trialStableStepsTaken += result.NumSteps
			}
		}
		trialResult := datacollector.TrialEndData{
			TrialIndex:                 trial,
			NumTestStates:              *numTestStates,
			NumTargetStates:            TARGET_STATES,
			NumStableStates:            trialNumStable,
			StableStatesMeanStepsTaken: float64(trialStableStepsTaken) / float64(trialNumStable),
		}
		collector.EventChannel <- hopfieldutils.IndexedWrapper[interface{}]{
			Index: datacollector.DataCollectionEvent_TrialEnd,
			Data:  trialResult,
		}
		InfoLogger.Printf("Stable States: %05d/%05d\n", trialNumStable, *numTestStates)
	}
	InfoLogger.Println("TRIAL COMPLETE")

	writeStopError := collector.WriteStop()
	if writeStopError != nil {
		InfoLogger.Fatalf("ERROR: DataWriter finished with error %#v!\n", writeStopError)
	}
	InfoLogger.Println("DATA WRITTEN")
}

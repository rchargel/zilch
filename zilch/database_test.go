package zilch

import (
	"runtime"
	"testing"
	"time"
)

func Test_DatabaseLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
		database, err := NewDatabase("../resources")

		if err != nil {
			t.Error(err.Error())
			return
		}

		if database.IsFullyLoaded() {
			t.Error("Database should not be loaded that quickly")
			return
		}

		start := time.Now()
		for {
			if database.IsFullyLoaded() {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		ellapsed := time.Since(start)
		t.Logf("Database loaded in %s\n", ellapsed)

		if !database.IsFullyLoaded() || ellapsed.Seconds() > float64(8.0) {
			t.Error("Database took too long to load")
		}

		// find the total number of records in the distribution
		var totalRecords uint32
		totalRecords = uint32(0)
		for _, entry := range database.DistributionMap {
			totalRecords += entry.ZipCodes
		}

		// find the total number of records in the system
		var totalRecordsInSystem uint32
		totalRecordsInSystem = uint32(0)
		for _, entry := range database.CountryIndexMap {
			totalRecordsInSystem += uint32(len(entry.Entries))
		}

		if totalRecords > totalRecordsInSystem {
			t.Errorf("There are %v records in the distribution map and %v records in the system, this is not correct\n", totalRecords, totalRecordsInSystem)
		} else {
			t.Logf("There are %v records in the distribution map and %v records in the system, this is expected\n", totalRecords, totalRecordsInSystem)
		}

		distributions := database.GetDistributions()

		for i, entry := range distributions {
			if i < (len(distributions) - 1) {
				if entry.ZipCodes < distributions[i+1].ZipCodes {
					t.Error("Preceeding distributions should always have more or the same number of zip codes")
				}
			}
		}
	}
}

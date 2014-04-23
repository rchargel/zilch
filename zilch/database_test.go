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

		query := map[string]string{
			"ZipCode": "22151",
		}

		if queryResult, err := database.ExecQuery(query); err == nil {
			if queryResult.TotalFound != 1 {
				t.Errorf("Found more than one entry: %v", queryResult.TotalFound)
			} else {
				t.Log("Found " + queryResult.ZipCodeEntries[0].City)
			}
		} else {
			t.Error(err.Error())
		}

		delete(query, "ZipCode")
		query["City"] = "Springfield"
		query["Country"] = "US"

		if queryResult, err := database.ExecQuery(query); err == nil {
			if queryResult.TotalFound != 128 {
				t.Errorf("Not 147 entries: %v", queryResult.TotalFound)
			} else if queryResult.StartIndex != 1 {
				t.Errorf("Not first entry %v", queryResult.StartIndex)
			} else if queryResult.EndIndex != 128 {
				t.Errorf("Not last entry %v", queryResult.EndIndex)
			} else if queryResult.ResultsReturned != 128 {
				t.Errorf("Not 147 entries: %v", queryResult.ResultsReturned)
			} else {
				entryMap := make(map[string]ZilchEntry)
				for _, entry := range queryResult.ZipCodeEntries {
					if mpEntry, found := entryMap[entry.ZipCode]; found {
						t.Errorf("Found duplicate\n%s\n%s\n", entry, mpEntry)
					} else {
						entryMap[entry.ZipCode] = entry
					}
				}
			}
		} else {
			t.Error(err.Error())
		}
	}
}

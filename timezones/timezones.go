package timezones

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanoberholster/timezoneLookup"
	log "github.com/sirupsen/logrus"
)

func InitTimezoneService() (timezoneLookup.TimezoneInterface, error) {
	_, err := os.Stat("assets/timezone.snap.json")

	if os.IsNotExist(err) {
		return rebuildDatabase()
	}

	tzService, err := timezoneLookup.LoadTimezones(timezoneLookup.Config{
		DatabaseType: "memory",
		DatabaseName: "assets/timezone",
		Snappy:       true,
	})

	return tzService, err
}

func rebuildDatabase() (timezoneLookup.TimezoneInterface, error) {
	log.Info("Found no database. Rebuild the database.")
	if _, err := os.Stat("assets/timezones-with-oceans.geojson"); err != nil {
		log.Info("Found no shape files to build new database.")
		err = reloadShapeFiles()
	}

	tz := timezoneLookup.MemoryStorage(true, "assets/timezone")
	if err := tz.CreateTimezones("assets/combined-with-oceans.json"); err != nil {
		return nil, err
	}

	return tz, nil
}

func reloadShapeFiles() error {
	if _, err := os.Stat("assets/timezones-with-oceans.geojson.zip"); os.IsNotExist(err) {
		if err := downloadGeoJSON("assets/timezones-with-oceans.geojson.zip"); err != nil {
			log.Error("Failed to load the basic shape files from github.")
			panic(err)
		}
		log.Info("Download the zipped shape files.")
	}

	if err := unzipShapeFiles(); err != nil {
		log.Error("Failed to unzip shape files.")
		panic(err)
	}

	log.Info("Successfull download the shape file.")
	return nil
}

func unzipShapeFiles() error {
	src := "assets/timezones-with-oceans.geojson.zip"
	dest := "assets"

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		log.Info(fpath)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.Name == "dist/combined-with-oceans.json" {
			outPath := "assets/combined-with-oceans.json"

			outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("Found no combined-with-oceans.json to unzip")
}

// DownloadFile will download a url to a local file.
func downloadGeoJSON(filepath string) error {
	log.Info("Start to download basic shape files.")
	url := "https://github.com/evansiroky/timezone-boundary-builder/releases/download/2019a/timezones-with-oceans.geojson.zip"

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

package timezones

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanoberholster/timezoneLookup"
	log "github.com/sirupsen/logrus"
)

var assetPath = "assets"
var buildDatabasePath = assetPath + "/timezone.snap.json"
var rawDatabasePath = assetPath + "/timezones-with-oceans.geojson"
var zipPath = assetPath + "/timezones-with-oceans.geojson.zip"

func InitTimezoneService() (timezoneLookup.TimezoneInterface, error) {
	_, err := os.Stat(buildDatabasePath)

	// TODO: write a json with infos which version of the timezones are stored and
	// to request a update if the version of the timezones is old

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
	if err := os.MkdirAll(assetPath, os.ModePerm); err != nil {
		log.Error("Failed to build assets folder.")
		panic(err)
	}

	if _, err := os.Stat(rawDatabasePath); err != nil {
		log.Info("Found no shape files to build new database.")
		if err := reloadShapeFiles(); err != nil {
			return nil, err
		}
	}

	tz := timezoneLookup.MemoryStorage(true, "assets/timezone")
	if err := tz.CreateTimezones(rawDatabasePath); err != nil {
		return nil, err
	}

	return tz, nil
}

func reloadShapeFiles() error {
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		if err := downloadGeoJSON(zipPath); err != nil {
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
	src := zipPath
	dest := assetPath

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

		if f.Name == "combined-with-oceans.json" {
			outPath := rawDatabasePath

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
	return errors.New("found no combined-with-oceans.json to unzip")
}

// DownloadFile will download a url to a local file.
func downloadGeoJSON(filepath string) error {
	log.Info("Start to download basic shape files.")
	url := "https://github.com/evansiroky/timezone-boundary-builder/releases/download/2020d/timezones-with-oceans.geojson.zip"

	// Set latest URL when we have it, fallback is hardcoded
	latestUrl, err := latestReleaseAssetUrl("evansiroky/timezone-boundary-builder", "github.com/racemap/tz-service", "timezones-with-oceans.geojson.zip")
	if err == nil {
		url = latestUrl
	}

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

func latestReleaseAssetUrl(repo string, userAgent string, assetName string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/"+repo+"/releases/latest", nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", userAgent)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	type Asset struct {
		Name string `json:"name"`
		Url  string `json:"browser_download_url"`
	}

	type Release struct {
		Assets []Asset
	}

	release := Release{}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(release.Assets); i++ {
		if assetName == release.Assets[i].Name {
			return release.Assets[i].Url, nil
		}
	}

	return "", errors.New("no asset found")
}

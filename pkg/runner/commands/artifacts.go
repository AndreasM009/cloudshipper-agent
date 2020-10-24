package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/andreasM009/cloudshipper-agent/pkg/logs"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/proxy"
	"github.com/andreasM009/cloudshipper-agent/pkg/runner/settings"
)

func downloadAndExtractArtifacts(artifactsURL string, proxy proxy.ControllerProxy) error {
	dest := settings.GetArtifactsDirectory()

	proxy.Report(logs.LogMessage{
		LogType: logs.Info,
		Message: "Downloading artifacts...",
	})

	tokens := strings.Split(artifactsURL, "/")
	fileName := tokens[len(tokens)-1]
	file := path.Join(dest, fileName)

	output, err := os.Create(file)

	if err != nil {
		proxy.Report(logs.NewErrorLog(fmt.Sprintf("Downloading artifacts failed: %s", err)))
	}

	defer output.Close()

	response, err := http.Get(artifactsURL)
	if err != nil {
		proxy.Report(logs.NewErrorLog(fmt.Sprintf("Downloading ertifacts failed: %s", err)))
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		proxy.Report(logs.NewErrorLog(fmt.Sprintf("Downloading ertifacts failed: %s", err)))
		return err
	}

	proxy.Report(logs.NewInfoLog("Artifacts downloaded"))

	_, err = unzip(file, dest)
	if err != nil {
		proxy.Report(logs.NewErrorLog(fmt.Sprintf("Unzipping artifacts failed: %s", err)))
		return err
	}
	return nil
}

func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return filenames, err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

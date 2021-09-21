package mmdb

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	downloadUrl = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=%s&suffix=tar.gz"
)

func Download(filePath, licenseKey string) error {

	downloadURL := fmt.Sprintf(downloadUrl, licenseKey)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	uncompressedStream, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if filepath.Ext(header.Name) != ".mmdb" { // Skip all but the database
				continue
			}
			outFile, err := os.Create(filePath)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			outFile.Close()
		default:
			err = fmt.Errorf(
				"ExtractTarGz: uknown type: %v in %s",
				header.Typeflag,
				header.Name)
			return err
		}
	}
	return nil
}

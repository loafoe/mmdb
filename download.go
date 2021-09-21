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

func Download(filePath, licenseKey string, client ...*http.Client) error {
	c := http.DefaultClient
	if len(client) > 0 {
		c = client[0]
	}
	downloadURL := fmt.Sprintf(downloadUrl, licenseKey)
	resp, err := c.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	uncompressedStream, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("gzip.NewReader: %w", err)
	}
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tarReader.Next: %w", err)
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
				return fmt.Errorf("os.Create: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("io.Copy: %w", err)
			}
			_ = outFile.Close()
		default:
			err = fmt.Errorf(
				"ExtractTarGz: uknown type: %v in %s",
				header.Typeflag,
				header.Name)
			return fmt.Errorf("header.Typeflag: %w", err)
		}
	}
	return nil
}

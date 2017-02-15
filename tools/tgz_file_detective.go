package tools

import (
	"os"
	"compress/gzip"
	"archive/tar"
	"io"
	"io/ioutil"
	"fmt"
	"regexp"
)

func ReadTgzFile(tgzFilePath, tarredFile string) ([]byte, error) {
	// Open the tar archive for reading.
	file, err := os.Open(tgzFilePath)
	if err != nil {
		return nil, fmt.Errorf("Could not read archive %s", tgzFilePath)
	}

	gzf, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid gzip archive", tgzFilePath)
	}

	tr := tar.NewReader(gzf)

	tarredFileRegex, err := regexp.Compile(fmt.Sprintf("^(.\\/)?%s$", tarredFile))
	if err != nil {
		return nil, err
	}

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid tar", tgzFilePath)
		}

		if tarredFileRegex.Match([]byte(hdr.Name)) {
			stemcellFileContents, err := ioutil.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("%s does not contain a valid file %s", tgzFilePath, tarredFile)
			}

			return stemcellFileContents, nil
		}
	}

	return nil, fmt.Errorf("%s does not contain file %s", tgzFilePath, tarredFile)
}
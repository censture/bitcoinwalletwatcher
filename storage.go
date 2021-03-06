package bitcoinwalletwatcher

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// InfoFile is the information file
type InfoFile struct {
	CurrentBlock uint64 `json:"current_block"`
}

const (
	// DefaultInfoFilePath represents the information file path
	DefaultInfoFilePath = ".info"
	// DefaultCurrentBlock represents the genesis block of bitcoin network
	DefaultCurrentBlock = 0
)

var filepath = DefaultInfoFilePath

// NewInfoStorage creates new info storage
func NewInfoStorage(path string) (*InfoFile, error) {
	if path != "" {
		filepath = DefaultInfoFilePath
	}

	var info *InfoFile

	b, _ := ioutil.ReadFile(path)
	if len(b) > 0 {
		if err := json.Unmarshal(b, &info); err != nil {
			return nil, err
		}
	}

	if info == nil {
		info = &InfoFile{
			CurrentBlock: DefaultCurrentBlock,
		}
	}

	return info, nil
}

// Update updates details
func (i *InfoFile) Update(block uint64) error {
	i.CurrentBlock = block

	return nil
}

// Save saves the info file to the storage
func (i *InfoFile) Save() error {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(i)
	if err != nil {
		return err
	}

	if _, err = f.Write(b); err != nil {
		return err
	}

	return nil
}

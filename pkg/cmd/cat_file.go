package cmd

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

const CAT_FILE_COMMAND = "cat-file"

type CatFile struct {
	pretty     bool
	showType   bool
	showSize   bool
	objectHash string
}

func NewCatFile(args []string) *CatFile {
	c := CatFile{}
	fs := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	fs.BoolVar(&c.pretty, "p", false, "pretty print the object content")
	fs.BoolVar(&c.showType, "t", false, "show the object type")
	fs.BoolVar(&c.showSize, "s", false, "show the object size")
	fs.Parse(args)
	c.objectHash = fs.Args()[0]
	return &c
}

func (c *CatFile) Name() string {
	return CAT_FILE_COMMAND
}

func (c *CatFile) Exec() error {
	if c.objectHash == "" {
		return fmt.Errorf("missing object hash")
	}

	var hash = c.objectHash
	var blobFolder = hash[:2]
	var blobFile = hash[2:]
	var blobPath = path.Join(".git", "objects", blobFolder, blobFile)

	blobBuffer, err := os.Open(blobPath)
	if err != nil {
		return fmt.Errorf("open blob object: %w", err)
	}
	defer blobBuffer.Close()

	src, err := zlib.NewReader(blobBuffer)
	if err != nil {
		return fmt.Errorf("create zlib reader: %w", err)
	}
	defer src.Close()

	raw, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	indexContent := bytes.IndexByte(raw, 0)
	if indexContent == -1 {
		return fmt.Errorf("corrupted file at %s", blobPath)
	}

	metadata := strings.Split(string(raw[:indexContent]), " ")
	if len(metadata) < 2 {
		return fmt.Errorf("invalid metadata format")
	}

	size, err := strconv.Atoi(metadata[1])
	if err != nil {
		return fmt.Errorf("convert content size: %w", err)
	}

	// Manejo de flags
	if c.pretty {
		fmt.Print(string(raw[indexContent+1:]))
	} else if c.showType {
		fmt.Println(metadata[0])
	} else if c.showSize {
		fmt.Println(size)
	} else {
		return fmt.Errorf("no operation specified")
	}

	return nil
}

func (c *CatFile) Help() string {
	return `Usage: mygit cat-file (-p | -t | -s) <objectHash>
    
    -p  Pretty print the object content.
    -t  Show the object type.
    -s  Show the object size.`
}

package model

type CompressionAlgo string

const (
	Gzip  CompressionAlgo = "gzip"
	Bzip2 CompressionAlgo = "bzip2"
	Zstd  CompressionAlgo = "zstd"
	Xz    CompressionAlgo = "xz"
	None  CompressionAlgo = ""
)

type Compression struct {
	Algo  CompressionAlgo
	Level int
}

func (c *Compression) IsCompressionEnabled() bool {
	return c.Algo != None
}

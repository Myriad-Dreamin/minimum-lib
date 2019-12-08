package mock

import "io"

type NamedReader struct {
	ContentType_ string
	Reader io.Reader
}

func (n NamedReader) Serialize() (io.Reader, error) {
	return n.Reader, nil
}

func (n NamedReader) ContentType() string {
	return n.ContentType_
}

func NewNamedReader(contentType string, reader io.Reader) *NamedReader {
	return &NamedReader{ContentType_: contentType, Reader: reader}
}


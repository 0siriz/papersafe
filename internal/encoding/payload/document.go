package payload

import (
	"bytes"
	"encoding/base64"
	"strconv"
)

// documentVersion is the version of the DOCUMENT payload wire format.
const documentVersion = 1

// Document is an indexed payload carrying one part of a larger data set.
type Document struct {
	// Index is the 1-based position of this payload when the source data was
	// split across multiple payloads.
	Index int
	// Total is the total number of payloads the source data was split into.
	Total int
	Data  []byte
}

// NewDocument creates an indexed document payload carrying one part of a
// larger data set. Index is 1-based and must be between 1 and total, inclusive.
func NewDocument(index, total int, data []byte) Document {
	return Document{
		Index: index,
		Total: total,
		Data:  data,
	}
}

// Domain returns the document domain.
func (d Document) Domain() Domain {
	return DomainDocument
}

// AppendText appends the text encoding of the document payload to b.
func (d Document) AppendText(b []byte) ([]byte, error) {
	if err := d.validate(); err != nil {
		return nil, err
	}

	b = appendHeader(b, DomainDocument)
	b = strconv.AppendUint(b, documentVersion, 10)
	b = append(b, separator...)
	b = strconv.AppendInt(b, int64(d.Index), 10)
	b = append(b, separator...)
	b = strconv.AppendInt(b, int64(d.Total), 10)
	b = append(b, separator...)
	b = base64.RawURLEncoding.AppendEncode(b, d.Data)

	return b, nil
}

// MarshalText returns the text encoding of the document payload.
func (d Document) MarshalText() ([]byte, error) {
	return d.AppendText(nil)
}

// UnmarshalText decodes a document text encoding into the payload.
func (d *Document) UnmarshalText(b []byte) error {
	domain, rest, err := splitHeader(b)
	if err != nil {
		return err
	}
	if domain != DomainDocument {
		return ErrInvalidDomain
	}

	return d.unmarshalBody(rest)
}

func (d *Document) unmarshalBody(rest []byte) error {
	versionPart, rest, found := bytes.Cut(rest, []byte(separator))
	if !found {
		return ErrInvalidVersion
	}

	version, err := strconv.Atoi(string(versionPart))
	if err != nil {
		return ErrInvalidVersion
	}
	if version != documentVersion {
		return ErrUnsupportedVersion
	}

	indexPart, rest, found := bytes.Cut(rest, []byte(separator))
	if !found {
		return ErrInvalidChunk
	}

	totalPart, dataPart, found := bytes.Cut(rest, []byte(separator))
	if !found {
		return ErrInvalidChunk
	}

	index, err := strconv.Atoi(string(indexPart))
	if err != nil {
		return ErrInvalidChunk
	}

	total, err := strconv.Atoi(string(totalPart))
	if err != nil {
		return ErrInvalidChunk
	}

	if index < 1 || total < 1 || index > total {
		return ErrInvalidChunk
	}

	bufSize := base64.RawURLEncoding.DecodedLen(len(dataPart))
	data, err := base64.RawURLEncoding.AppendDecode(make([]byte, 0, bufSize), dataPart)
	if err != nil {
		return err
	}

	d.Index = index
	d.Total = total
	d.Data = data

	return nil
}

func (d Document) validate() error {
	if d.Index < 1 || d.Total < 1 || d.Index > d.Total {
		return ErrInvalidChunk
	}

	return nil
}

package payload

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
)

const (
	Prefix    = "PAPERSAFE"
	separator = ":"
)

type Domain int

const (
	DomainKeyshard Domain = iota
	DomainDocument
)

var domainName = map[Domain]string{
	DomainKeyshard: "KEYSHARD",
	DomainDocument: "DOCUMENT",
}

var domainID = map[string]Domain{
	domainName[DomainKeyshard]: DomainKeyshard,
	domainName[DomainDocument]: DomainDocument,
}

type Payload struct {
	Domain Domain
	// ShardID is the Shamir share index for keyshard payloads.
	ShardID byte
	// Index is the 1-based position of this payload when the source data was
	// split across multiple payloads.
	Index int
	// Total is the total number of payloads the source data was split into.
	Total int
	Data  []byte
}

// NewKeyshard creates a keyshard payload carrying the shard ID.
func NewKeyshard(id byte, data []byte) Payload {
	return Payload{
		Domain:  DomainKeyshard,
		ShardID: id,
		Data:    data,
	}
}

// NewChunk creates an indexed payload carrying one part of a larger data set.
// Index is 1-based and must be between 1 and total, inclusive.
func NewChunk(domain Domain, index, total int, data []byte) Payload {
	return Payload{
		Domain: domain,
		Index:  index,
		Total:  total,
		Data:   data,
	}
}

func (p *Payload) AppendText(b []byte) ([]byte, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	b = append(b, []byte(Prefix)...)
	b = append(b, []byte(separator)...)
	b = append(b, []byte(domainName[p.Domain])...)
	b = append(b, []byte(separator)...)

	switch p.Domain {
	case DomainKeyshard:
		b = fmt.Appendf(b, "%02x", p.ShardID)
		b = append(b, []byte(separator)...)
	case DomainDocument:
		b = strconv.AppendInt(b, int64(p.Index), 10)
		b = append(b, []byte(separator)...)
		b = strconv.AppendInt(b, int64(p.Total), 10)
		b = append(b, []byte(separator)...)
	}

	b = base64.RawURLEncoding.AppendEncode(b, p.Data)

	return b, nil
}

func (p *Payload) MarshalText() ([]byte, error) {
	bufSize := len(Prefix) + len(domainName[p.Domain]) + base64.RawURLEncoding.EncodedLen(len(p.Data)) + 2*len(separator)

	switch p.Domain {
	case DomainKeyshard:
		bufSize += 2 + len(separator)
	case DomainDocument:
		bufSize += len(strconv.Itoa(p.Index)) + len(strconv.Itoa(p.Total)) + 2*len(separator)
	}

	b, err := p.AppendText(make([]byte, 0, bufSize))
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UnmarshalText decodes the text encoding into the payload.
func (p *Payload) UnmarshalText(b []byte) error {
	buf := b

	if !bytes.HasPrefix(buf, []byte(Prefix+separator)) {
		return ErrInvalidPrefix
	}
	buf = buf[len(Prefix+separator):]

	name, rest, found := bytes.Cut(buf, []byte(separator))
	if !found {
		return ErrInvalidDomain
	}

	domain, ok := domainID[string(name)]
	if !ok {
		return ErrInvalidDomain
	}
	p.Domain = domain
	p.ShardID = 0
	p.Index = 0
	p.Total = 0

	var err error
	switch domain {
	case DomainKeyshard:
		rest, err = p.unmarshalKeyshard(rest)
	case DomainDocument:
		rest, err = p.unmarshalDocument(rest)
	}
	if err != nil {
		return err
	}

	bufSize := base64.RawURLEncoding.DecodedLen(len(rest))
	data, err := base64.RawURLEncoding.AppendDecode(make([]byte, 0, bufSize), rest)
	if err != nil {
		return err
	}
	p.Data = data

	return nil
}

// unmarshalKeyshard parses the required shard ID prefix. Base64url never
// contains a colon, so a separator unambiguously marks the ID.
func (p *Payload) unmarshalKeyshard(rest []byte) ([]byte, error) {
	idPart, data, found := bytes.Cut(rest, []byte(separator))
	if !found {
		return nil, ErrInvalidShardID
	}

	id, err := strconv.ParseUint(string(idPart), 16, 8)
	if err != nil || id == 0 {
		return nil, ErrInvalidShardID
	}

	p.ShardID = byte(id)

	return data, nil
}

// unmarshalDocument parses the required index/total prefix. Base64url never
// contains a colon, so separators unambiguously mark chunk metadata.
func (p *Payload) unmarshalDocument(rest []byte) ([]byte, error) {
	indexPart, after, found := bytes.Cut(rest, []byte(separator))
	if !found {
		return nil, ErrInvalidChunk
	}

	totalPart, data, found := bytes.Cut(after, []byte(separator))
	if !found {
		return nil, ErrInvalidChunk
	}

	index, err := strconv.Atoi(string(indexPart))
	if err != nil {
		return nil, ErrInvalidChunk
	}

	total, err := strconv.Atoi(string(totalPart))
	if err != nil {
		return nil, ErrInvalidChunk
	}

	if index < 1 || total < 1 || index > total {
		return nil, ErrInvalidChunk
	}

	p.Index = index
	p.Total = total

	return data, nil
}

func (p *Payload) validate() error {
	switch p.Domain {
	case DomainKeyshard:
		if p.ShardID == 0 {
			return ErrInvalidShardID
		}
		if p.Index != 0 || p.Total != 0 {
			return ErrInvalidChunk
		}
	case DomainDocument:
		if p.ShardID != 0 {
			return ErrInvalidShardID
		}
		if p.Index < 1 || p.Total < 1 || p.Index > p.Total {
			return ErrInvalidChunk
		}
	default:
		return ErrInvalidDomain
	}

	return nil
}

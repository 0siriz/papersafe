package payload

import "bytes"

const (
	Prefix    = "PAPERSAFE"
	separator = ":"
)

// Domain identifies the kind of payload carried in a PAPERSAFE text encoding.
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

// Payload is implemented by every domain-specific payload type.
type Payload interface {
	Domain() Domain
	AppendText(b []byte) ([]byte, error)
	MarshalText() ([]byte, error)
	UnmarshalText(b []byte) error
}

// appendHeader appends the common PAPERSAFE:<domain>: prefix.
func appendHeader(b []byte, domain Domain) []byte {
	b = append(b, Prefix...)
	b = append(b, separator...)
	b = append(b, domainName[domain]...)
	b = append(b, separator...)

	return b
}

// splitHeader validates the common prefix and extracts the domain and the
// domain-specific body.
func splitHeader(b []byte) (Domain, []byte, error) {
	prefix := Prefix + separator
	if !bytes.HasPrefix(b, []byte(prefix)) {
		return 0, nil, ErrInvalidPrefix
	}
	b = b[len(prefix):]

	name, rest, found := bytes.Cut(b, []byte(separator))
	if !found {
		return 0, nil, ErrInvalidDomain
	}

	domain, ok := domainID[string(name)]
	if !ok {
		return 0, nil, ErrInvalidDomain
	}

	return domain, rest, nil
}

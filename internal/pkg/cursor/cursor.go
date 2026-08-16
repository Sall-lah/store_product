package cursor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrInvalidCursor indicates the provided cursor string could not be parsed or decoded.
	ErrInvalidCursor = errors.New("invalid cursor format")
)

// Cursor represents the composite keyset pointer used for stable, O(1) pagination.
type Cursor struct {
	ID        string
	CreatedAt time.Time
}

// Encode serializes a composite (CreatedAt, ID) tuple into an opaque base64 string.
// Keyset cursors prevent pagination drift (items shifting pages when new items are added)
// and eliminate the quadratic database scanning penalty of OFFSET/LIMIT.
func Encode(createdAt time.Time, id string) string {
	raw := fmt.Sprintf("%s|%d", id, createdAt.UnixNano())
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// Decode deserializes a base64 string back into a Cursor struct.
// Returns ErrInvalidCursor if the string is corrupted or invalid.
func Decode(encoded string) (*Cursor, error) {
	if strings.TrimSpace(encoded) == "" {
		return nil, nil
	}

	bytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	parts := strings.SplitN(string(bytes), "|", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidCursor
	}

	id := parts[0]
	nanos, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	return &Cursor{
		ID:        id,
		CreatedAt: time.Unix(0, nanos),
	}, nil
}

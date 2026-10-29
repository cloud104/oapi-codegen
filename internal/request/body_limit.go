package request

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

const (
	defaultMaxBodySize = humanize.MiByte
	maxBodySize        = math.MaxInt64

	maxBodySizeExtension = "x-max-body-size"
)

type BodyLimit struct {
	bytes int64
}

func FindBodyLimit(op codegen.OperationDefinition) *BodyLimit {
	if op.Spec == nil || op.Spec.RequestBody == nil || op.Spec.RequestBody.Value == nil {
		return defaultBodyLimit()
	}

	raw, ok := op.Spec.RequestBody.Value.Extensions[maxBodySizeExtension]
	if !ok {
		return defaultBodyLimit()
	}

	var value string

	switch v := raw.(type) {
	case float64:
		value = strconv.FormatFloat(v, 'f', -1, 64)
	default:
		value = fmt.Sprint(v)
	}

	limit, err := parseBodyLimit(value)
	if err != nil {
		panic(fmt.Errorf("invalid body limit: %w", err))
	}

	return limit
}

func (l *BodyLimit) Bytes() int64 {
	return l.bytes
}

func (l *BodyLimit) String() string {
	return humanize.IBytes(uint64(l.bytes))
}

func defaultBodyLimit() *BodyLimit {
	return &BodyLimit{bytes: defaultMaxBodySize}
}

func newBodyLimit(n int64) (*BodyLimit, error) {
	if n < 1 || n > maxBodySize {
		return nil, fmt.Errorf("body limit must be between 1 and %d bytes", maxBodySize)
	}

	return &BodyLimit{bytes: n}, nil
}

func newBodyLimitFromUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](n T) (*BodyLimit, error) {
	if uint64(n) > math.MaxInt64 {
		return nil, fmt.Errorf("body limit exceeds int64 range")
	}

	return newBodyLimit(int64(n))
}

func parseBodyLimit(s string) (*BodyLimit, error) {
	n, err := humanize.ParseBytes(strings.TrimSpace(s))
	if err != nil {
		return nil, fmt.Errorf("parse body limit: %w", err)
	}

	return newBodyLimitFromUint(n)
}

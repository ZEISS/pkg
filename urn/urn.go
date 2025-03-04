package urn

import (
	"errors"
	"strings"

	pb "github.com/katallaxie/pkg/proto"

	"github.com/go-playground/validator/v10"
)

const (
	// DefaultSeparator is the separator used to separate the segments of a URN.
	DefaultSeparator = ":"
	// DefaultNamespace is the default namespace used for URNs.
	DefaultNamespace = "urn"
)

// Match is a string that can be used to match a URN segment.
type Match string

var (
	// Wildcard is the wildcard used to match any value.
	Wildcard Match = "*"
	// Empty is the empty string
	Empty Match
)

// String returns the string representation of the match.
func (m Match) String() string {
	return string(m)
}

// ErrorInvalid is returned when parsing an URN with an invalid format.
var ErrorInvalid = errors.New("urn: invalid format")

// use a single instance of Validate, it caches struct info.
var validate = validator.New()

// URN represents a unique, uniform identifier for a resource
type URN struct {
	// Namespace is the namespace segment of the URN.
	Namespace Match `validate:"required"`
	// Partition is the partition segment of the URN.
	Partition Match `validate:"max=256"`
	// Service is the service segment of the URN.
	Service Match `validate:"max=256"`
	// Region is the region segment of the URN.
	Region Match `validate:"max=256"`
	// Identifier is the identifier segment of the URN.
	Identifier Match `validate:"max=64"`
	// Resource is the resource segment of the URN.
	Resource Match `validate:"required,max=256"`
}

// String returns the string representation of the URN.
func (u *URN) String() string {
	return strings.Join([]string{u.Namespace.String(), u.Partition.String(), u.Service.String(), u.Region.String(), u.Identifier.String(), u.Resource.String()}, DefaultSeparator)
}

// Match returns true if the right-hand side URN matches the left-hand side URN.
// The left-hand side URN is assumed to be specific and the right-hand side URN is assumed to be generic.
//
//nolint:gocyclo
func (u *URN) Match(urn *URN) bool {
	return (u.Namespace == urn.Namespace || (u.Namespace == Wildcard && urn.Namespace == Wildcard) || (u.Namespace == Empty && urn.Namespace == Empty) || urn.Namespace == Wildcard || urn.Namespace == Empty) &&
		(u.Partition == urn.Partition || (u.Partition == Wildcard && urn.Partition == Wildcard) || (u.Partition == Empty && urn.Partition == Empty) || urn.Partition == Wildcard || urn.Partition == Empty) &&
		(u.Service == urn.Service || (u.Service == Wildcard && urn.Service == Wildcard) || (u.Service == Empty && urn.Service == Empty) || urn.Service == Wildcard || urn.Service == Empty) &&
		(u.Region == urn.Region || (u.Region == Wildcard && urn.Region == Wildcard) || (u.Region == Empty && urn.Region == Empty) || urn.Region == Wildcard || urn.Region == Empty) &&
		(u.Identifier == urn.Identifier || (u.Identifier == Wildcard && urn.Identifier == Wildcard) || (u.Identifier == Empty && urn.Identifier == Empty) || urn.Identifier == Wildcard || urn.Identifier == Empty) &&
		(u.Resource == urn.Resource || (u.Resource == Wildcard && urn.Resource == Wildcard) || (u.Resource == Empty && urn.Resource == Empty) || urn.Resource == Wildcard || urn.Resource == Empty)
}

// ExactMatch returns true if the URN matches the given URN exactly.
func (u *URN) ExactMatch(urn *URN) bool {
	return u.Namespace == urn.Namespace &&
		u.Partition == urn.Partition &&
		u.Service == urn.Service &&
		u.Region == urn.Region &&
		u.Identifier == urn.Identifier &&
		u.Resource == urn.Resource
}

// New takes a namespace, partition, service, region, identifier and resource and returns a URN.
func New(namespace, partition, service, region, identifier, resource Match) (*URN, error) {
	urn := &URN{
		Namespace:  namespace,
		Partition:  partition,
		Service:    service,
		Region:     region,
		Identifier: identifier,
		Resource:   resource,
	}

	validate = validator.New()

	if err := validate.Struct(urn); err != nil {
		return nil, err
	}

	return urn, nil
}

// Parse takes a string and parses it to a URN.
func Parse(s string) (*URN, error) {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	if Match(s) == Wildcard {
		return &URN{
			Namespace:  Wildcard,
			Partition:  Wildcard,
			Service:    Wildcard,
			Region:     Wildcard,
			Identifier: Wildcard,
			Resource:   Wildcard,
		}, nil
	}

	segments := strings.SplitN(s, DefaultSeparator, 6)
	if len(segments) < 5 {
		return nil, ErrorInvalid
	}

	mm := make([]Match, len(segments))
	for i, segment := range segments {
		mm[i] = Match(segment)
		if Match(segment) == Wildcard || Match(segment) == Empty {
			mm[i] = Wildcard
		}
	}

	urn := &URN{
		Namespace:  mm[0],
		Partition:  mm[1],
		Service:    mm[2],
		Region:     mm[3],
		Identifier: mm[4],
		Resource:   mm[5],
	}

	validate = validator.New()

	if err := validate.Struct(urn); err != nil {
		return nil, err
	}

	return urn, nil
}

// ProtoMessage returns the proto.ResourceURN representation of the URN.
func (u *URN) ProtoMessage() *pb.ResourceURN {
	return &pb.ResourceURN{
		Canonical:  u.String(),
		Namespace:  u.Namespace.String(),
		Partition:  u.Partition.String(),
		Service:    u.Service.String(),
		Region:     u.Region.String(),
		Identifier: u.Identifier.String(),
		Resource:   u.Resource.String(),
	}
}

// FromProto returns the URN representation from a proto.ResourceURN representation.
func FromProto(r *pb.ResourceURN) (*URN, error) {
	return Parse(r.GetCanonical())
}

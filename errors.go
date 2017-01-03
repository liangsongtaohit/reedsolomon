package reedsolomon

import (
	"errors"
)

var errInvalidNumShards = errors.New("too less or too more shards for creating encoding matrix, be sure: data/parity > 0 & data+shards < 255")

var errTooFewShards = errors.New("too few shards given for encoding")
// ErrShardNoData will be returned if there are no shards,
// or if the length of all shards is zero.
var errShardNoData = errors.New("no shard data")

// ErrShardSize is returned if shard length isn't the same for all
// shards.
var errShardSize = errors.New("shard sizes does not match")

// errInvalidRowSize will be returned if attempting to create a matrix with negative or zero row number.
var errInvalidRowSize = errors.New("invalid row size")

// errInvalidColSize will be returned if attempting to create a matrix with negative or zero column number.
var errInvalidColSize = errors.New("invalid column size")

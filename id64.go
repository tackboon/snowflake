package snowflake

import (
	"fmt"
	"sync"
	"time"
)

// Generates a custom 64-bit unique ID.
//
// This ID format is inspired by Twitter Snowflake, using millisecond precision,
// and is suitable for high-throughput distributed environments.
//
// Default Bit layout (total: 64 bits):
// - 41 bits: Timestamp in milliseconds (supports ~69 years from a custom epoch)
// - 13 bits: Machine ID (supports up to 8192 unique nodes)
// - 10 bits: Sequence number (supports up to 1024 IDs per millisecond per machine)
//
// Properties:
// - Ordered by time (timestamp in high bits)
// - Globally unique within a distributed system
// - Suitable for high-throughput systems needing time-sortable unique IDs
//
// Limitations:
// - Max 1024 requests per millisecond per machine
// - Max lifespan ~69 years from the defined epoch (StartTimestamp)
// - Not safe to use as a Redis ZSET score (exceeds float64 53-bit precision)
// - Time must be monotonic or tolerate minor clock skew
type Snowflake64Bit struct {
	StartTimestamp int64 // timestamp in milliseconds
	MachineID      int64

	timestampBit uint8
	machineIDBit uint8
	seqBit       uint8

	timestampBitMask int64
	machineIDBitMask int64
	seqBitMask       int64

	currentSeq       int64
	currentTimestamp int64
	maxSkewTolerance int64 // max tolerance in milliseconds
	mu               *sync.Mutex
}

func NewDefaultSnowflake64Bit(machineID int64, startTimestamp int64) Snowflake64Bit {
	var timestampBit uint8 = 41
	var machineIDBit uint8 = 13
	var seqBit uint8 = 10

	return Snowflake64Bit{
		StartTimestamp:   startTimestamp,
		MachineID:        machineID,
		timestampBit:     timestampBit,
		machineIDBit:     machineIDBit,
		seqBit:           seqBit,
		timestampBitMask: (1 << timestampBit) - 1,
		machineIDBitMask: (1 << machineIDBit) - 1,
		seqBitMask:       (1 << seqBit) - 1,
		currentSeq:       0,
		currentTimestamp: 0,
		maxSkewTolerance: 1000,
		mu:               &sync.Mutex{},
	}
}

func NewCustomSnowflake64Bit(machineID int64, startTimestamp int64, timestampBit uint8, machineIDBit uint8, seqBit uint8, maxSkewTolerance int64) Snowflake64Bit {
	return Snowflake64Bit{
		StartTimestamp:   startTimestamp,
		MachineID:        machineID,
		timestampBit:     timestampBit,
		machineIDBit:     machineIDBit,
		seqBit:           seqBit,
		timestampBitMask: (1 << timestampBit) - 1,
		machineIDBitMask: (1 << machineIDBit) - 1,
		seqBitMask:       (1 << seqBit) - 1,
		currentSeq:       0,
		currentTimestamp: 0,
		maxSkewTolerance: 1,
		mu:               &sync.Mutex{},
	}
}

func (s *Snowflake64Bit) GenerateID() (id uint64, err error) {
	if s.MachineID > s.machineIDBitMask {
		return 0, fmt.Errorf("seq or machine_id exceed upper boundary, machine_id: %d", s.MachineID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Cater clock skew happened
	now := time.Now().UnixMilli()
	if now < s.currentTimestamp {
		skewDiff := s.currentTimestamp - now
		if skewDiff > s.maxSkewTolerance {
			return 0, fmt.Errorf("big clock skew detected: %dms", skewDiff)
		}

		now = s.currentTimestamp
	}

	// Handle sequence incrementation
	if now == s.currentTimestamp {
		s.currentSeq++
		if s.currentSeq > s.seqBitMask {
			// Wait for the next millisecond
			nextMilliSeconds := time.UnixMilli(s.currentTimestamp + 1)
			time.Sleep(time.Until(nextMilliSeconds))
			now = nextMilliSeconds.UnixMilli()
			s.currentSeq = 0
		}
	} else {
		s.currentSeq = 0
	}
	s.currentTimestamp = now

	// Generate unique ID
	diffTimestamp := s.currentTimestamp - s.StartTimestamp
	diffTimestamp = diffTimestamp & s.timestampBitMask           // avoid timestamp overflow
	diffTimestamp = diffTimestamp << (s.machineIDBit + s.seqBit) // shift timestamp
	machineID := s.MachineID << s.seqBit                         // shift machine_id

	id = uint64(diffTimestamp) | uint64(machineID) | uint64(s.currentSeq)
	return id, nil
}

func (s *Snowflake64Bit) DecodeID(id uint64) (timestamp int64, machineID int64, sequence int64) {
	sequence = int64(id & uint64(s.seqBitMask))

	machineID = int64((id >> s.seqBit) & uint64(s.machineIDBitMask))

	timestampShift := s.machineIDBit + s.seqBit
	timestampRaw := int64(id >> timestampShift)

	timestamp = timestampRaw + s.StartTimestamp
	return
}

# Snowflake

A customizable and high-performance Snowflake ID generator written in Go. Supports both:

- **53-bit second-based IDs** (within float64 precision)
- **64-bit millisecond-based IDs** (for high-throughput systems)

---
<br>

## ✨ Features

- ✅ Redis ZSET score-safe 53-bit ID generator
- ✅ 64-bit millisecond-resolution ID generator
- ✅ Customizable bit layout (timestamp, machine ID, sequence)
- ✅ Thread-safe with mutex
- ✅ Clock skew detection and tolerance
- ✅ Decode function to extract timestamp, machine ID, and sequence

---
<br>

## 📦 Installation

```bash
go get github.com/tackboon/snowflake
```

---
<br>

## 🚀 Usage

```go
package main

import (
	"fmt"

	"github.com/tackboon/snowflake"
)

func main() {
	var machineID int64 = 1
	var startTimestamp int64 = 1744615749

	generator := snowflake.NewDefaultSnowflake64Bit(machineID, startTimestamp)
	id, err := generator.GenerateID()

	fmt.Println(id, err)
}
```

---
<br>

## ⚠️ Notes

Use 53-bit mode for systems like Redis where ZSET scores require float64 precision

Use 64-bit mode for internal services or high-volume event streams

Ensure all nodes use the same StartTimestamp and are loosely time-synced

<br>

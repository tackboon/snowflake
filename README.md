# snowflake

A customizable and high-performance Snowflake ID generator written in Go. Supports both:

- **53-bit second-based IDs** (within float64 precision)
- **64-bit millisecond-based IDs** (for high-throughput systems)

---

## ✨ Features

- ✅ Redis ZSET score-safe 53-bit ID generator
- ✅ 64-bit millisecond-resolution ID generator
- ✅ Customizable bit layout (timestamp, machine ID, sequence)
- ✅ Thread-safe with mutex
- ✅ Clock skew detection and tolerance
- ✅ Decode function to extract timestamp, machine ID, and sequence

---

## 📦 Installation

```bash
go get github.com/tackboon/snowflake

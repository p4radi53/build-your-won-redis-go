## Redis from scratch

Implements:
- reading of the RESP Protocol
- commands: ping, get, set
- writing to an AppendOnly file as a way for persistence

Out of scope:
- hset, hget commands

### Running

Prerequisites: golang, redis-cli


```
go run .
```


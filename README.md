# discord-bot
Discord Bot related Go code

## Installation

```
go get github.com/tmpest/discord-bot
```

## Sample Server Code

```go
package main

import "github.com/tmpest/discord-bot/pkg/server"

func main() {
	server.New().Start()
}
```
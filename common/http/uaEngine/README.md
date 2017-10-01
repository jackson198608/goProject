# introduction
- this package is design for the situation when you want to get the useragent for http request 

# Installation

	go get github.com/jackson198608/goProject/common/http/uaEngine 

# Quick Start

- Create uaEngine 

```Go
	uaEngine := NewUaEngine("")
```

- Get a random pc useragent 

```Go
uaEngine.GetPcRandomeEngine
```

- Get a random mobile useragent 

```Go
uaEngine.GetMobileRandomeEngine()
```

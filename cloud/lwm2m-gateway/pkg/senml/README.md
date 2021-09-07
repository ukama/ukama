# go-senml

[![CI Status](https://github.com/nkristek/go-senml/workflows/CI/badge.svg)](https://github.com/nkristek/go-senml/actions)
[![GoDoc](https://godoc.org/github.com/nkristek/go-senml?status.svg)](https://godoc.org/github.com/nkristek/go-senml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nkristek/go-senml)](https://goreportcard.com/report/github.com/nkristek/go-senml)

A go library to parse SenML records. It currently supports JSON and XML.

This library implements [RFC 8428](https://tools.ietf.org/rfc/rfc8428.txt) (SenML version 10).

## Install
```sh
go get github.com/nkristek/go-senml
```

## Import
```go
import(
	"github.com/nkristek/go-senml"
)
```

## Usage
```go
// parse using the encoding format
message, err := senml.Decode(payload, senml.JSON)
if err != nil {
	// process error
}

// resolve the message (resolve base attributes, convert relative to absolute time etc.)
resolvedMessage, err := message.Resolve()
if err != nil {
	// process error
}

// encode a new message
encodedMessage, err := message.Encode(senml.JSON)
if err != nil {
	// process error
}
```

## Error handling

If `Resolve()` returns an error it can have one of the following types:

- `InvalidNameError`
- `UnsupportedVersionError`
- `DifferentVersionError`
- `MissingValueError`

Likewise, the `Encode()` and `Decode()` functions return an error of type `UnsupportedFormatError` if it was called with an unsupported format.

The error types provide extra values to parse the exact reason in code. If you need to check on the specific reason on why resolving the message has failed, the following `switch` statement should suffice: 

```go
_, err := message.Resolve()
if err != nil {
	switch err.(type) {
	case *senml.InvalidNameError:
		// do something
		// for example:
		invalidNameError := err.(*senml.InvalidNameError)
		switch invalidNameError.Reason {
		case senml.FirstCharacterInvalid:
			break
		case senml.ContainsInvalidCharacter:
			break
		case senml.Empty:
			break
		}
	case *senml.UnsupportedVersionError:
		// do something
		break
	case *senml.DifferentVersionError:
		// do something
		break
	case *senml.MissingValueError:
		// do something
		break
	default:
		break
	}
}
```

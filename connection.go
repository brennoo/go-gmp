package gmp

import "context"

type Connection interface {
	Execute(ctx context.Context, command interface{}, response interface{}) error
	Close() error
	RawXML(xml string) (string, error)
}

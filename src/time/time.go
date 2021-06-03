//go:generate mockgen -source=time.go -destination=time_mock.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

package time

import "time"

var Time ITime

type ITime interface {
	Now() time.Time
}

type timeProvider struct {}

func (timeProvider) Now() time.Time {
	return time.Now()
}

func init() {
	Time = &timeProvider{}
}

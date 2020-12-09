module primetools

go 1.15

require (
	github.com/bogem/id3v2 v1.2.0
	github.com/dhowden/itl v0.0.0-20170329215456-9fbe21093131
	github.com/dhowden/plist v0.0.0-20141002110153-5db6e0d9931a // indirect
	github.com/draeron/itunes-win v0.0.0-20170927094139-aeb7f600a3d9
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	gopkg.in/djherbis/times.v1 v1.2.0
)

replace github.com/draeron/itunes-win => ../itunes-win

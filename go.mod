module primetools

go 1.15

require (
	github.com/bogem/id3v2 v1.2.0
	github.com/deepakjois/gousbdrivedetector v0.0.0-20161027045320-4d29e4d6f1b7
	github.com/dhowden/itl v0.0.0-20170329215456-9fbe21093131
	github.com/dhowden/plist v0.0.0-20141002110153-5db6e0d9931a // indirect
	github.com/draeron/itunes-win v0.2.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/karrick/godirwalk v1.16.1
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/pelletier/go-toml v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil/v3 v3.20.11
	github.com/sirupsen/logrus v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	go.mongodb.org/mongo-driver v1.4.4
	golang.org/x/sys v0.0.0-20210110051926-789bb1bd4061 // indirect
	golang.org/x/text v0.3.3
	gopkg.in/djherbis/times.v1 v1.2.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

//replace github.com/draeron/itunes-win => ../itunes-win

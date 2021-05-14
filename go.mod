module primetools

go 1.15

require (
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/bogem/id3v2 v1.2.0
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/deepakjois/gousbdrivedetector v0.0.0-20161027045320-4d29e4d6f1b7
	github.com/dhowden/itl v0.0.0-20170329215456-9fbe21093131
	github.com/dhowden/plist v0.0.0-20141002110153-5db6e0d9931a // indirect
	github.com/dhowden/tag v0.0.0-20201120070457-d52dcb253c63
	github.com/draeron/itunes-win v0.2.3
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gobwas/glob v0.2.3
	github.com/jmoiron/sqlx v1.3.3
	github.com/karrick/godirwalk v1.16.1
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/pelletier/go-toml v1.9.1
	github.com/pkg/errors v0.9.1
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil/v3 v3.21.4
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli/v2 v2.3.0
	go.mongodb.org/mongo-driver v1.5.2
	golang.org/x/sys v0.0.0-20210511113859-b0526f3d8744 // indirect
	golang.org/x/text v0.3.6
	gopkg.in/djherbis/times.v1 v1.2.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

//replace github.com/draeron/itunes-win => ../itunes-win

package options

var (
	dryrun = false
)

func SetDryRun() {
	dryrun = true
}

func DryRun() bool {
	return dryrun
}

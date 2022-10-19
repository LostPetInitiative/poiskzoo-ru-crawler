package version

// these are to be set with
// go build -ldflags "-X 'github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version.GitCommit=abcdef' -X 'github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version.AppVersion=1.2.3'

// A git commit that corresponds to the current running code
var GitCommit string = "dev"

// An assigned sem ver to the current running code. "0.0.0" means development build
var AppVersion string = "0.0.0"

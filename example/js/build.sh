#export GOOS=linux

go list -tags javascript  -f {{.Deps}}
gopherjs build --tags javascript -o js.js -m
# ripeproduce

## Installation ##

Install the single dependency for goquery.
```
go get github.com/PuerkitoBio/goquery
```

To run the tests.
```
go test -v
```

To install the command.
```
go install ./...
```

## Running ##

Directly from project.
```
go run ./cmd/ripeproduce.go
```

Using the command if it has been installed.
```
ripeproduce
```

Use the --url flag to change from the default url
```
ripeproduce --url http://xyz.pqr?123
```

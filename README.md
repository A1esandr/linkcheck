### linkcheck

Checker for not working links on site

## Usage

### Prerequisites
* Go 1.16

### Example
```
go run cmd/main.go -url https://golang.org
```
Or can be used as package for checking whole site
```
linkcheck.New(htmlOnlyParam).Start(url)

where:
url - string, site url for check
htmlOnlyParam - bool, load only .html pages or all
```


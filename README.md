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
Also can be used for check separate page
```
results, err := linkcheck.New(htmlOnlyParam).Check(pageURL)
if err != nil {
  fmt.Println(err)
}
for from, state := range results {
  fmt.Printf("%s : %s \n", state, from)
}

where:
pageURL - string, page url for check
htmlOnlyParam - bool, load only .html pages or all
```

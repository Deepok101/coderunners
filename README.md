# coderunners

CodeRunners is an API allowing users to debug and run code via HTTP requests.

Currently, it only supports Python.

## Setup

Follow these steps to start running and debugging Python code
1. Run `go run main.go` to start the CodeRunners API.

You can now run queries through the `localhost:80` base url.

## Endpoint Documentation

All queries will have a `json` body. Make sure to set the `content-type` of all queries as `application/json`

| Endpoint       | Method           | Body                                                                        | Description|
| ----------------|:-------:| :---------------------------------------------------------------------------|------------|
| `/run`          | `POST` | `language` : "\<PROGRAMMING LANGUAGE>" <br> `content` : "\<CODE TO EXECUTE>" | Execute code given a language. |
| `/debug`        | `POST` | `language` : "\<PROGRAMMING LANGUAGE>" <br> `content` : "\<CODE TO DEBUG>"   | Debug code given a language. Start at line 0. |
| `/debug/stepin` | `GET`  |                                                                              | Execute debug step-in functionality |
| `/debug/stepout`| `GET`  |                                                                              | Execute debug step-in functionality |
| `/debug/stepover`| `GET`  |                                                                              | Execute debug step-over functionality |
| `/debug/continue`| `GET`  |                                                                              | Execute debug continue functionality |
| `/debug/setbreakpoint/{lineNo}`| `GET`| `lineNo` (Query Parameter) : <CODE LINE NUMBER>                 | Add a breakpoint to a specific line.  |
  
If you want more information on what step-in, step-over and continue debug functionalities do, click [here](https://winintro.ru/windowspowershellhelp.en/html/62095f16-dd77-4840-bd65-49cebb354c08.htm#:~:text=In%20the%20Command%20Pane%2C%20type%20O%20and%20press%20ENTER%2C%20or,Debug%20menu%2C%20click%20Step%20Out.&text=Continues%20execution%20to%20the%20end,menu%2C%20click%20Run%2FContinue)

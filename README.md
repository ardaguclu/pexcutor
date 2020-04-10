[![GoDoc](https://godoc.org/github.com/ardaguclu/pexcutor?status.png)](https://godoc.org/github.com/ardaguclu/pexcutor)
[![Go Report Card](https://goreportcard.com/badge/github.com/ardaguclu/pexcutor)](https://goreportcard.com/report/github.com/ardaguclu/pexcutor)


Pexcutor is a Go library for managing external processes.

* It has an ability to stop external process, when the caller is cancelled/completed/returned.
* It has an retry mechanism for possible crash and it restarts again according to the given retry count value by caller.
* Signals can be passed directly to the external processes.
* It listens termination signals handled by caller and cancel external process via context cancelling.

## USAGE

    ctx, cancel := context.WithCancel(context.Background())
 	p := pexcutor.New(ctx, 3, "ls", "-alh")
 	err := p.Start()
 	if err != nil {
 		log.Fatal("Start error ", err)
 	}
 
 	sOut, sErr, err := p.GetResult()
 	if err != nil {
 		log.Fatal("GetResult error ", err)
 	}
 
 	log.Println(sOut)
 	log.Println(sErr)
 
## TESTING

`go test -race -cover ./...`

* Unit tests and integrations tests are implemented into the process_test package and currently code coverage is approximately %83.
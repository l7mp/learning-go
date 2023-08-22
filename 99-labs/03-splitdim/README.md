# A silly web app: SplitDim

In the course of this lab we are going to build a Go web app that allows groups of people to keep track of money transfers between themselves and helps clear debs and credits with minimal money transfer. The app is by and large modeled after the excellent [SplitWise](https://www.splitwise.com) app, but it is much dumber so we will call it SplitDim. 

The below tasks guide you in writing a simple web app that implements the barebones SplitDim functionality with a basic local database. Later we will gradually extend the app to implement the 5 cloud native pillars. The tasks are followed by tests; absolve each to check whether you successfully completed the lab.

## Table of Contents

1. [Basics]([#basics])
1. [API](#database-api)
1. [Local data layer](#local-data-layer)
1. [Reset](#reset)
1. [Transfer](#transfer)
1. [Accounts](#accounts)
1. [Clear](#clear)

## Basics

SplitDim helps housemates, trips, friends, and family members maintain their internal money transfers and keep track of who owns who. Imagine you are at a trip with your friends, you invite one of your friends for a coffee, they pay the taxi fee for the entire group, and then someone else from the group pays your train ticket. After a while, it becomes practically impossible to remember all mutual payments and clear the debts. 

Enter SplitDim, a simple web app that allows friends to register their transfers (e.g., "Joe paid Alice's coffee for 5 USD", and then "Alice paid Joe's train ticket for 3 USD") and see (1) the current balance of each registered user (how much debt or credit they have) and (2) the minimal list of mutual money transfers that would allow them the clear all debts ("Alice would need to pay Joe 2 USD to clear the debt").

We are going to build SplitDim as a Go web app. During this lab we will write only the barebones web service that keeps the balances in memory; later we will extend it into a proper cloud-native app. The web service will implement 4 APIs:
- `POST: /api/transfer`: register a transfer between two users of a given amount (this API uses the POST HTTP method to let users post the transfer's details in JSON format),
- `GET: /api/accounts`: return the list of current balances for each registered user,
- `GET: /api/clear`: return the list of transfers that would allow users to clear their debts between themselves, and
- `GET: /api/reset`: reset all balances to zero.

So let's start, shall we?

1. Initialize a new Go project under `99-labs/code/splitdim`. Make sure you actually use this directory: there are some files placed there for you to help your work. 

   ``` sh
   cd 99-labs/code/splitdim
   go mod init splitdim
   go get github.com/stretchr/testify/assert
   go mod tidy
      ```

1. Open a new file called `main.go` and declare that you are going to build an executable.

   ``` go
   package main
   ```

1. Import the packages to be used.

   ``` go
   import (
       "log"
       "net/http"
   )
   ```
1. Implement 4 empty HTTP handlers: these will be the placeholders for the SplitDim API.

   ``` go
   // TransferHandler is a HTTP handler that implements the money transfer API.
   func TransferHandler(w http.ResponseWriter, r *http.Request) {}
   
   // AccountListHandler is a HTTP handler that returns the current balance of each registered user.
   func AccountListHandler(w http.ResponseWriter, r *http.Request) {}
   
   // ClearHandler is a HTTP handler that returns a list of transfers to clear the balance of each user.
   func ClearHandler(w http.ResponseWriter, r *http.Request) {}
   
   // ResetHandler is a HTTP handler that allows to zero out all balances.
   func ResetHandler(w http.ResponseWriter, r *http.Request) {}
   ```

1. Start the main function:

   ``` go
   func main() {
       // Set the default logger to a fancier log format.
       log.SetFlags(log.LstdFlags | log.Lshortfile)
   
       ... 
   }
   ```

1. Install a HTTP handler to serve a static HTML file with inline JavaScript to interact with the SplitDim API. This page implements the client GUI of our service so that you can connect from a browser, send transfers and see the current balances (of course real men use `curl`, but anyway). The HTML file is provided as part of this lab in order to allow you to concentrate on the server side, but feel free to modify it according to your liking.

   Register a static HTTP handler that will serve the prepackaged HTML file for the default path (`/`).

   ``` go
   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       http.ServeFile(w, r, "static/index.html")
   })
   ```

1. Register the 4 empty HTTP handlers for the 4 API endpoints.

   ``` go
   http.HandleFunc("/api/transfer", TransferHandler)
   http.HandleFunc("/api/accounts", AccountListHandler)
   http.HandleFunc("/api/clear", ClearHandler)
   http.HandleFunc("/api/reset", ResetHandler)
   ```

1. And finally start the HTTP server on port 8080. Remember, `http.ListenAndServe` will block until the program exits or an error is raised: `log.Fatal` will write the error message to the standard output in the latter case.

   ``` go
   log.Println("Server listening on http://:8080")
   log.Fatal(http.ListenAndServe(":8080", nil))
   ```

Once ready, you can run the program with `go run main.go`: if all goes well you should see the output:

```
20XX/YY/ZZ 19:03:59 main.go:48: Server listening on http://:8080
```

This means your server is ready to accept HTTP requests.

> ✅ **Check**
>
> Run the below test to check whether you have successfully completed the exercise. If all goes well, you should see the output `PASS`.
> ``` sh
> go test ./... --tags=httphandler -run '^TestAPIEndpoints' -v
> PASS
> ```
> Make sure the web service is running: the test issues requests to the HTTP server and checks whether the response is as expected.

At the moment all HTTP handlers respond to all HTTP methods, whereas our goal is for each API handler to accept only one HTTP method: `/api/transfer` should accept on HTTP POST requests, and all the other APIs should respond to `GET` requests only. Every other type of access should result a HTTP 405 error code ("Not Allowed"). This can be achieved by adding the following test to the beginning of your HTTP handlers. 

``` go
// SomeHandler accepts only POST requests ()
func SomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        // Return HTTP 405
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
}
```
Substitute `http.MethodPost` with `http.MethodGet` for handlers that accept only GET requests.

> ✅ **Check**
>
> Run the below test to check whether you have successfully completed the exercise. If all goes well, you should see the output `PASS`.
> ``` sh
> go test ./... --tags=httphandler -run '^TestAPIMethods' -v
> PASS
> ```
> Make sure the web service is running: the test issues requests to the HTTP server and checks whether the response is as expected.

> **Note**
> 
> Currently `http.HandleFunc("/api/transfer", TransferHandler)` will route *all* HTTP requests with a path that starts with `/api/transfer` to the `TransferHandler`, including `/api/transfer/random/api` or `/api/transfer/some/malicious/attack`. To make sure *only* the required API is served, you can check for the HTTP path inside the handler.

## API

The next task is to design the public SplitDim API, that is, the Go structs (and the corresponding JSON format) that the API accepts in HTTP requests and generates as response.

1. Create a new (sub)package under `pkg/api`. Recall, in Go this amounts to opening a new file `pkg/api/api.go` (the name can be anything) with the following as the first line.

   ``` go
   package api
   ```

1. Design the format that the API endpoint `/api/transfer` accepts. Recall, this API can be called to register money transactions, but this same data format can also be used by the `/api/clear` API to return the list of transactions that would allow users to clear their debts. Add the below to `api.go`:

   ``` go
   package api
   // Transfer represents a transaction.
   type Transfer struct {
       // The debtor.
       Sender string `json:"sender"`
       // The creditor.
       Receiver string `json:"receiver"`
       // The amount transferred in the transaction.
       Amount int `json:"amount"`
   }
   ```

   > **Warning**
   > 
   > Documenting your public APIs is mandatory. Using the [Godoc](https://go.dev/blog/godoc) format will simplify generating easy-to-browse documentation from your code. Below we will sometime omit the docs for brevity, but you should never!

1. We also 

   ``` go
   // Account represents the balance of a user.
   type Account struct {
       // The name of the account holder.
       Holder string `json:"holder"`
       // Current balance.
       Balance int `json:"balance"`
   }
   ```

   > **Warning**
   > 
   > Always add the JSON tags as above!If you expect your data format to ever be marshaled to JSON.

1. Finally, we design the `DataLayer` API: this will be a Go `interface` that simply declares the functions we want to support. Later, we will create multiple implementations.

   ``` go
   // DataLayer is an API for manipulating a balance sheet.
   type DataLayer interface {
       // Transfer will process a transfer.
       Transfer(t Transfer) error
       // AccountList returns the current acount of each user.
       AccountList() ([]Account, error)
       // Clear returns the list of transfers to clear all debts.
       Clear() ([]Transfer, error)
       // Reset sets all balances to zero.
       Reset() error
   }
   ```

   > **Note**
   > 
   > This is an internal API: we don't want people to import and use these internal data structures. Technically, therefore, we should place these definitions, and the implementations we will create later, into a new package under `internal/`, which, recall, cannot be imported from outside the main package. We will spare this now for simplicity, but in practice always be aware of what's public and what's private in your code.

> ✅ **Check**
>
> Run the below test to check whether you have successfully completed the exercise. If all goes well, you should see the output `PASS`.
> ``` sh
> go test ./... --tags=api -v
> PASS
> ```

## Local data layer

The next step is to define our internal `DataLayer`: the internal representation our service uses to represent accounts and balances. There are countless ways to do this: we chose the below simple in-memory representation as it will be straightforward to rewrite over the persistent key-value store `set-put` API later. 

1. Create a new package `pkg/db/local` and import some useful Go packages:

   ``` go
   package local

   import (
       "fmt"
       "sort"
       "sync"
   
       "splitdim/pkg/api"
   )
   ```

   > **Note **
   > 
   > We want to use our own API, hence the (sub)package import `github.com/<my-user>/splitdim/pkg/api`.

   > **Note **
   > 
   > From now on we will not explicitly write import lists: remember, just add imports when Go complains during compilation.

1. Use the below definition for the internal data layer.

   ``` go
   // localDB is a simple implementation of the DataLayer API.
   type localDB struct {
       // accounts maintains the balance for each user name
       accounts map[string]int
       // The read-write mutex makes sure concurrent access is safe.
       mu sync.RWMutex
   }
   ```

1. Create a constructor. Recall, Go does not have explicit constructors, so conventionally we export a `New*` function from the package that can be called to create a data structure.

   ``` go
   // NewDataLayer creates a new database of accounts.
   func NewDataLayer() api.DataLayer {
       return &localDB{accounts: make(map[string]int)}
   }
   ```

   > **Warning**
   > 
   > It is idiomatic Go to make data structure definitions private and let  constructors return an *interface* as a pointer to the struct instead of the actual struct itself (observe that the above returns `api.DataLayer`, not `localDB`). Since `localDB` is private, the caller would not be able to do much with it anyway.

1. Implement the placeholders for the 4 interface methods.

   ``` go
   func (db *localDB) Transfer(t api.Transfer) error { return nil }
   func (db *localDB) AccountList() ([]api.Account, error) { return []api.Account{}, nil }
   func (db *localDB) Clear() ([]api.Transfer, error) { nil }
   func (db *localDB) Reset() error { nil }
   ```

1. Create an actual in-memory database in `main.go`. Declare the global variable `db` that will hold the database (it is global so that the HTTP handlers can access it):

   ``` go
   var db api.DataLayer
   ```

   Then actually construct the database in the `main` function:
   
   ```go
   func main() {
       ...
       db = local.NewDataLayer()
       ...
   }
   ```

> ✅ **Check**
>
> Run the below test to check whether you have successfully completed the exercise. If all goes well, you should see the output `PASS`.
> ``` sh
<!-- > go test ./... --tags=localconstructor -v -->
<!-- > PASS -->
<!-- > 
> ```

## Reset

Next, we set off to write the HTTP handlers and the actual implementations for our API. First we write the `api/reset` HTTP handler and implementation, since (1) this is our simplest API and (2) it helps us initializing the account database in the tests.

Before we set off, below is a list of useful functions that help dealing with encoding/decoding HTTP requests and responses to/from JSON, plus some additional handy utilities you can use to write HTTP handlers:
   - `json.NewDecoder(r.Body).Decode(&var)`: decode the body from the `r` of type `http.Request` into the variable `var` of the requested type (make sure to import `encoding/json`),
   - `json, err := json.Marshal(var)`: marshal the data in `var` into `json` of type `[]byte` or return an error,
   - `w.Header().Set("Content-Type", "application/json")`: set the HTTP response header `Content-Type` to `application/json`,
   - `_, err = w.Write(json)`: write the byte slice `json` into the HTTP response body represented by `w`,
   - `w.WriteHeader(status)`: write the HTTP status `status` (e.g., `http.StatusBadRequest`, `http.StatusNotFound` or `http.StatusOk`) into the HTTP response `w`,
   - `fmt.Fprintf(w, "API request failed: %s", err)`: write an error message into the HTTP response body,
   - `log.Printf("format", args...)`: log a request.

> **Note**
> 
> If unsure about the use of any of these functions, remember that all Go libraries come with excellent documentation at [`pkg.go.dev`](https://pkg.go.dev). For instance, the documentation of the `net/http` package is [here](https://pkg.go.dev/net/http), [here](https://pkg.go.dev/encoding/json) is the docs for the JSON encoding/decoding functions, etc.

1. Implement the HTTP handler for the `/api/reset` API, i.e., the `ResetHandler` function in `main.go`. This function should should implement the following steps: 
   - check that the request uses the HTTP GET method (already done),
   - log the request, 
   - call `db.Reset()` to process the request (this will currently call the placeholder, we will implement this in the next step),
   - return HTTP 500 error status (`http.StatusInternalServerError`) if something went wrong and write the error message into the response body,
   - if all went well, return HTTP 200 (`http.StatusOK`).
   
1. Actually implement the ``db.Reset()` call from the above. This will require extending the below function in the `pkg/db/local` package:

   ``` go
   func (db *localDB) Reset() error { ... }
   ```

   The implementation should make the following steps to reset the account database:
   - lock the database `db` for *writing* (`db.mu.Lock()`),
   - make sure the database will be unlocked at the end of the call (`defer db.mu.Unlock()`), and
   - re-initialize the accounts database: `db.accounts = make(map[string]int)`.
   
> ✅ **Check**
>
> Run the below test to check whether you have successfully completed the exercise. If all goes well, you should see the output `PASS`.
> ``` sh
> go test ./... --tags=reset -v
> PASS
> ```
> Make sure the web service is running: the test issues requests to the HTTP server and checks whether the response is as expected.

## Transfer

1. Implement the HTTP handler for the `/api/transfer` API, i.e., the `TransferHandler` function in `main.go`. This function should should implement the following steps: 
   - check that the request uses the HTTP POST method (already done),
   - check that the HTTP request body is a valid JSON, 
   - unmarshal it into a `api.Transfer` struct, 
   - return HTTP 400 error status (`http.StatusBadRequest`) if any of these checks fail,
   - log the request, 
   - call `db.Transfer` to process the request (this will currently call the placeholder, we will implement this in the next step),
   - return HTTP 500 error status (`http.StatusInternalServerError`) if something went wrong and write the error message into the response body,
   - if all went well, return HTTP 200 (`http.StatusOK`).
   
1. Actually implement the ``db.Transfer(t)` call from the above. This will require extending the below function in the `pkg/db/local` package:

   ``` go
   func (db *localDB) Transfer(t api.Transfer) error { ... }
   ```

   The implementation should make the following steps to process the transaction in `t`:
   - check if the sender and the receiver in the transfer are different and return an appropriate error if not,
   - lock the database `db` for *writing* (`db.mu.Lock()`),
   - make sure the database will be unlocked at the end of the call (`defer db.mu.Unlock()`), and
   - perform the actual transaction: increase the balance of the sender by the amount of the transfer (meaning that the sender is now a *creditor* in the transaction) and decrease the balance of the receiver with the same amount (meaning that the receiver is now in *debt*: the receiver *owns* "amount" of money to the sender).
   
   If any of the users are not actually registered in the database then the transfer API should silently initialize the balance of these users with zero balance and perform the transaction on after that. This simplifies the API a lot (otherwise we would need an additional `api/register` API as well).

> ✅ **Check**
>
> Run the below test to check whether you have successfully completed the exercise. If all goes well, you should see the output `PASS`.
> ``` sh
> go test ./... --tags=transfer -v
> PASS
> ```
> Make sure the web service is running: the test issues requests to the HTTP server and checks whether the response is as expected.

## Accounts

The `api/accounts` API should return the current balance of each registered user.

1. Implement the HTTP handler for the `/api/accounts` API, i.e., the `AccountListHandler` function in `main.go`. This function should should implement the following steps: 
   - check that the request uses the HTTP GET method (already done),
   - log the request, 
   - call `accountList, err := db.AccountList()` to obtain the current list of user accounts and balances from the data layer (an `[]api.Accuunt` slice) (this will currently call the placeholder, we will implement this in the next step),
   - return HTTP 500 error status (`http.StatusInternalServerError`) if something went wrong and write the error message into the response body,
   - marshal the returned `accountList` to JSON and return HTTP 500 error status (`http.StatusInternalServerError`) if this fails,
   - set the HTTP response header `Content-Type` to `application/json`,
   - write the JSON data into the HTTP response body,
   - if all went well, return HTTP 200 (`http.StatusOK`).
   
1. Implement the ``db.AccountList()` call from the above. This will require extending the below function in the `pkg/db/local` package:

   ``` go
   func (db *localDB) AccountList() ([]api.Account, error) { ... }
   ```

   The implementation should make the following steps to return the account list:
   - lock the database `db` for *reading* (`db.mu.RLock()`),
   - make sure the database will be unlocked at the end of the call (`defer db.mu.RUnlock()`), 
   - initialize an empty account list to return later: `ret := []api.Account{}`,
   - iterate through the local database and copy the use accounts out into `ret` (recall: the database is a map from user names to balances while the returned list is a slice of `struct{Holder string, Balance int}` structs, make sure to do the conversion),
   - to make the returned list more accessible to users, sort it by name (don't forget to import the `sort` package):
    ``` go
    sort.Slice(ret, func(i, j int) bool {
        return ret[i].Holder < ret[j].Holder
    })
    ```
    - return `ret` and a `nil` error.

> ✅ **Check**
>
> Run the below test to check whether you have successfully completed the exercise. If all goes well, you should see the output `PASS`.
> ``` sh
> go test ./... --tags=accounts -v
> PASS
> ```
> Make sure the web service is running: the test issues requests to the HTTP server and checks whether the response is as expected.

## Clear

We left the most difficult API to the end: generating a list of suggested transfers between users how to clear all debts in the database. 

Although not entirely trivial, this algorithm is not that difficult: just find a user with positive balance (meaning the user is a creditor) and another one with negative balance (meaning the user is a debtor) and create a transaction with the debtor as the sender and the creditor as the receiver with the amount hat is the minimum of the balances of the debtor and the creditor. Observe that in this step we have cleared the balance of at least one user, so we made progress. Do this as long as there is a user with positive balance. This is the algortithm we will use below, but first we have to write the HTTP handler.

1. Implement the HTTP handler for the `/api/clear` API, i.e., the `ClearHandler` function in `main.go`. This function should should implement the following steps: 
   - check that the request uses the HTTP GET method (already done),
   - log the request, 
   - call `transfers, err := db.Clear()` to obtain a list of transfers that would clear the debts (this will currently call the placeholder, we will implement this in the next step),
   - marshal the returned slice `transfers` to JSON and return HTTP 500 error status (`http.StatusInternalServerError`) if this fails,
   - set the HTTP response header `Content-Type` to `application/json`,
   - write the JSON data into the HTTP response body,
   - if all went well, return HTTP 200 (`http.StatusOK`).
   
1. Actually implement the ``db.Clear()` call from the above. This will require extending the below function in the `pkg/db/local` package:

   ``` go
   func (db *localDB) Clear() ([]api.Transfer, error) { ... }
   ```

   The implementation should make the following steps to clear the debts:
   - lock the database `db` for *reading* (`db.mu.RLock()`),
   - make sure the database will be unlocked at the end of the call (`defer db.mu.RUnlock()`), 
   - check the consistency of the database (this will prevent infinite loops later) by summing up user balances and checking whether the total balance is zero (return an error if this fails),
   - copy the account database into a local variable called `tempAcc` that the algorithm will work on (we don't want to actually perform the clearing operation, we just want to suggest a way to do this),
   - initialize an empty slice of transfers: `transfers := make([]api.Transfer, 0)`,
   - iterate through all user-balance pairs in the accounts database to find the creditor: `for sender, balance := range tempAcc { ... }`,
   - check if the balance of `sender` is negative, otherwise continue,
   - iterate through all user-balance pairs again to find a debtor: `for sender, balance := range tempAcc { ... }`,
   - check if the balance of `receiver` is positive, otherwise continue,
   - compute the minimum of the balances of the `sender` and the `receiver` and store it in `transferAmount`,
   - create an transfer from the `sender` to the `receiver` of size `transferAmount`,
   - increase the balance of the `sender` by `transferAmount` since they are now paying back their debt,
   - decrease the balance of the `receiver` by `transferAmount` since they are now being payed back their credit,
   - check if the balance of the sender reached zero: if yes, break the loop, this user is now ready.
   
> ✅ **Check**
>
> Run the below test to check whether you have successfully completed the exercise. If all goes well, you should see the output `PASS`.
> ``` sh
> go test ./... --tags=clear -v
> PASS
> ```
> Make sure the web service is running: the test issues requests to the HTTP server and checks whether the response is as expected.

> ✅ **Check**
>
> Run the below test to check whether you have successfully completed *all* the exercises above. If all goes well, you should see the output `PASS`.
> ``` sh
> go test ./... --tags=httphandler,api,localconstructor,reset,transfer,accounts,clear -v
> PASS
> ```
> Make sure the web service is running: the test issues requests to the HTTP server and checks whether the response is as expected.

## Deploy

The last step is to package up everything into a Docker container, deploy into Kubernetes, and test via `curl` or the browser GUI. 

1. Create an image build manifest in `deploy/Dockerfile` and build the container image.
1. Create a Kubernetes manifest called `deploy/kubernetes-local-db.yaml` that contains 
   - a Deployment called `splitdim` to run one replica of the `splitdim` container image and
   - a Service of type `LoadBalancer` that exposes the `splitdim` Deployment for external address on port port 80.
1. Test with `curl`.


> ✅ **Check**
> 
> Test your Kubernetes deployment. Some useful commands for testing from the shell:
> - store the external IP assigned by Kubernetes to the `splitdim` service:
>   ``` sh
>   export EXTERNAL_IP=$(kubectl get service splitdim -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
>   ```
> - register a transfer from sender `a` to receiver `b` of amount 1:
>   ``` sh
>   curl -H "Content-Type: application/json" --request POST --data '{"sender":"a","receiver":"b","amount":1}' http://${EXTERNAL_IP}/api/transfer
>   ```
> - list the account database: 
>   ``` sh
>   curl http://${EXTERNAL_IP}/api/accounts
>   ```
> - reset the account database: 
>   ``` sh
>   curl http://${EXTERNAL_IP}/api/reset
>   ```

<!-- Local Variables: -->
<!-- mode: markdown; coding: utf-8 -->
<!-- auto-fill-mode: nil -->
<!-- visual-line-mode: 1 -->
<!-- markdown-enable-math: t -->
<!-- End: -->

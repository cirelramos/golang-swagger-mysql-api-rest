# Simple GO Lang REST API

> Simple RESTful API to create, read, update and delete books. No database implementation yet

## Quick Start




``` bash
PATH Environment Variable: If you encounter "swag: command not found"

errors, it's likely that your Go bin directory is not in your PATH.

You'll need to add it. You can add this line to your ~/.bashrc or ~/.zshrc file.
  
export PATH=$PATH:$(go env GOPATH)/bin or export PATH=$PATH:$(go env HOME)/go/bin
```



``` bash


# Install mux router
  go install github.com/swaggo/swag/cmd/swag@latest && \
go get -u github.com/gorilla/mux && \
go install github.com/swaggo/swag/cmd/swag && \
go get -u github.com/swaggo/http-swagger && \
go get -u github.com/swaggo/swag/example/celler/docs && \
go get -u github.com/go-sql-driver/mysql
```

``` bash
go run main.go
```

## Endpoints

### Get All Books
``` bash
GET api/books
```
### Get Single Book
``` bash
GET api/books/{id}
```

### Delete Book
``` bash
DELETE api/books/{id}
```

### Create Book
``` bash
POST api/books

# Request sample
# {
#   "isbn":"4545454",
#   "title":"Book Three",
#   "author":{"firstname":"Harry",  "lastname":"White"}
# }
```

### Update Book
``` bash
PUT api/books/{id}

# Request sample
# {
#   "isbn":"4545454",
#   "title":"Updated Title",
#   "author":{"firstname":"Harry",  "lastname":"White"}
# }

```


```

## App Info

### Author

Brad Traversy
[Traversy Media](http://www.traversymedia.com)

### Version

1.0.0

### License

This project is licensed under the MIT License

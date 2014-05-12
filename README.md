# GoWiki

A simple wiki app to learn writing web apps with GoLang, from http://golang.org/doc/articles/wiki/. This simple wiki app stores the entries as text files, with internally linked pages using [page] syntax. This application can currently Create, Read and Update Pages, Delete yet to be implemented. 

## To Do
- [x] Add Twitter Bootstrap for styling, served via netdna
- [x] Seperate wiki text files into `/wikis/` and templates into `/templates/`
- [x] Add internal linking to pages using `[page]` syntax 
- [x] Redirect `/` to `/FrontPage`
- [ ] Change storage backend to a MySQL database
- [ ] Add a page to show list of all Wiki Pages
- [ ] Add a WYSWIG editor 
- [ ] Add a delete option to page

## Running application

1. Clone to `$GOPATH`
2. run `go build wiki.go`
3. execute wiki by `./wiki` if in linux ( run `wiki` in Windows )

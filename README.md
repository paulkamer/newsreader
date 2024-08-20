# Feedly clone

WIP - a Feedly clone built with Go and HTMX.

- Backend (Go https://docs.gofiber.io/ | https://github.com/air-verse/air)
- Frontend (HTMX https://htmx.org/)
- Database (SQLite)
- Job queue: TBD

## Running

Start backend with the command:

```sh
air
```

And open http://127.0.0.1:3001

## Testing

```sh
cd src

go test ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## TODOs

- [x] Use template layout
- [x] Adding news sources
- [x] Editing news sources
- [ ] Validate newssource URL (add, edit)
- [ ] [WIP] Tests
- [x] CSRF
- [ ] [WIP] Background jobs: Fetching news
- [x] Feed/newsource page
- [ ] Feed/newsource page
- [ ] List latest (unread) articles
- [ ] Article page
  - [ ] (un)marking an article as Read
- [ ] Background jobs: handling jobs with different priority
  - [ ] Newly added feed should be fetched with highest priority (ASAP)
- [ ] Login
- [ ] Group news sources by category

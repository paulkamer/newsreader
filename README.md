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
```

## TODOs

- [x] Use template layout
- [x] Adding news sources
- [x] Editing news sources
- [ ] Tests
- [ ] CSRF
- [ ] Background jobs: Fetching news
- [ ] Background jobs: job scheduler
- [ ] Feed/newsource page
- [ ] List latest (unread) articles
- [ ] Article page
  - [ ] (un)marking an article as Read
- [ ] Background jobs: handling jobs with different priority
  - [ ] Newly added feed should be fetched with highest priority (ASAP)
- [ ] Login
- [ ] Group news sources by category

# Feedly clone

WIP - a Feedly clone built with Go and HTMX.

- Backend (Go https://docs.gofiber.io/ | https://github.com/air-verse/air)
- Frontend (HTMX https://htmx.org/)
- Database (SQLite)
- Job queue: TBD

## Running

From the `src/` dir, run:

```sh
air
```

And open http://127.0.0.1:3001

## Testing

From the `src/` dir, run:

```sh
./runtests.sh
```

And open the generated `coverage.html` file in your browser

## To-dos

- [ ] DB migrations (golang-migrate )
- [ ] Input sanitization (XSS) (go-playground/validator)
- [ ] Delve debugger
- [ ] [WIP] Unit tests
- [ ] [WIP] e2e tests with Playwright
- [ ] [WIP] profile memory usage
- [ ] (un)marking an individual article as Read
- [ ] Group news sources by category
- [ ] [WIP] Background jobs: job scheduler
- [ ] Background jobs: handling jobs with different priority
  - [ ] Newly added feed should be fetched with highest priority (ASAP)
- [ ] Login

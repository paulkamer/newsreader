# Feedly clone

WIP - a Feedly clone built with Go and HTMX.

- Backend: Go, using https://docs.gofiber.io/ & https://github.com/air-verse/air
- Frontend: [HTMX](https://htmx.org)
- Database: SQLite
- E2E tests: [Playwright](https://playwright.dev)

## Running

From the `src/` dir, run:

```sh
air
```

And open http://127.0.0.1:3001

## Testing

### Unit tests

From the `src/` dir, run:

```sh
./runtests.sh
```

And open the generated `coverage.html` file in your browser

### E2E tests

From the project root dir, run:

```sh
npx playwright test

# Or with a UI, run
npx playwright test --ui

# View the test report with
npx playwright show-report
```

## To-dos

- [x] DB migrations (golang-migrate )
- [x] Delve debugger
- [x] Input sanitization (XSS) (go-playground/validator
  - [x] add news source
  - [x] edit news source
- [ ] [WIP] Unit tests
- [x] e2e tests with Playwright
- [ ] [WIP] profile memory usage
- [x] Support Atom 1.0 feeds
- [ ] (un)marking an individual article as Read
- [ ] Group news sources by category
- [ ] [WIP] Background jobs: job scheduler
- [ ] Background jobs: handling jobs with different priority
  - [ ] Newly added feed should be fetched with highest priority (ASAP)
- [ ] Multiple users; authentication & authorization

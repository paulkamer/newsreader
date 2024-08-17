# Feedly clone

- Backend (Go https://docs.gofiber.io/)
- Frontend (HTMX https://htmx.org/)
- Database (SQLite)
- Job queue: TBD

# TODOs

- [x] Use template layout
- [x] Adding news sources
- [x] Editing news sources
- [ ] CSRF
- [ ] Tests
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

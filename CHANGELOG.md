# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Hash password before storing user to the database ([@bakku](https://github.com/bakku), [#2](https://github.com/bakku/easyalert/pull/2));
- Improve validation of user creation endpoint ([@bakku](https://github.com/bakku), [#6](https://github.com/bakku/easyalert/pull/6));
- Add application/json Content-Type to all API responses ([@bakku](https://github.com/bakku), [#7](https://github.com/bakku/easyalert/pull/7));
- Add auth endpoint to let a user fetch his token ([@bakku](https://github.com/bakku), [#8](https://github.com/bakku/easyalert/pull/8));
- Add auth refresh endpoint to let a user refresh his token ([@bakku](https://github.com/bakku), [#12](https://github.com/bakku/easyalert/pull/12));
- Add endpoint to update a user ([@bakku](https://github.com/bakku), [#27](https://github.com/bakku/easyalert/pull/27));
- Add alerts migrations and CRUD actions ([@bakku](https://github.com/bakku), [#28](https://github.com/bakku/easyalert/pull/28));
- Add create alert endpoint ([@bakku](https://github.com/bakku), [#29](https://github.com/bakku/easyalert/pull/29));
- Use air to hot reload code ([@bakku](https://github.com/bakku), [#33](https://github.com/bakku/easyalert/pull/33));
- Add endpoint to return all alerts ([@bakku](https://github.com/bakku), [#34](https://github.com/bakku/easyalert/pull/34));
- Add endpoint to delete user account ([@bakku](https://github.com/bakku), [#36](https://github.com/bakku/easyalert/pull/36));

### Changed
- Update schema.sql to latest version ([@bakku](https://github.com/bakku), [#13](https://github.com/bakku/easyalert/pull/13));
- Update documentation to better reflect how to handle subject of emails ([@bakku](https://github.com/bakku), [#15](https://github.com/bakku/easyalert/pull/15));
- Remove admin flag from code and database ([@bakku](https://github.com/bakku), [#18](https://github.com/bakku/easyalert/pull/18));
- Migrate project to go modules ([@bakku](https://github.com/bakku), [#30](https://github.com/bakku/easyalert/pull/30));
- Decouple postgres repositories from http server ([@bakku](https://github.com/bakku), [#35](https://github.com/bakku/easyalert/pull/35));

[Unreleased]: https://github.com/bakku/easyalert/compare/b6283ea...HEAD

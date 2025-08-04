# Details

Date : 2025-08-03 00:51:04

Directory /home/y1nhui/work/github_own/luna-dial/luna_dial_server

Total : 48 files,  6648 codes, 569 comments, 1268 blanks, all 8485 lines

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [cmd/init/main.go](/cmd/init/main.go) | Go | 38 | 5 | 11 | 54 |
| [cmd/main.go](/cmd/main.go) | Go | 4 | 0 | 2 | 6 |
| [configs/config.ini](/configs/config.ini) | Ini | 15 | 0 | 4 | 19 |
| [doc/DDD文档.md](/doc/DDD%E6%96%87%E6%A1%A3.md) | Markdown | 135 | 0 | 36 | 171 |
| [doc/TDD测试计划.md](/doc/TDD%E6%B5%8B%E8%AF%95%E8%AE%A1%E5%88%92.md) | Markdown | 287 | 0 | 40 | 327 |
| [docs/SYSTEM\_INIT.md](/docs/SYSTEM_INIT.md) | Markdown | 92 | 0 | 42 | 134 |
| [examples/main\_with\_init.go](/examples/main_with_init.go) | Go | 33 | 7 | 11 | 51 |
| [go.mod](/go.mod) | Go Module File | 38 | 0 | 4 | 42 |
| [go.sum](/go.sum) | Go Checksum File | 191 | 0 | 1 | 192 |
| [internal/biz/errors.go](/internal/biz/errors.go) | Go | 42 | 5 | 7 | 54 |
| [internal/biz/journal.go](/internal/biz/journal.go) | Go | 201 | 21 | 34 | 256 |
| [internal/biz/journal\_repo.go](/internal/biz/journal_repo.go) | Go | 13 | 0 | 3 | 16 |
| [internal/biz/journal\_test.go](/internal/biz/journal_test.go) | Go | 494 | 53 | 107 | 654 |
| [internal/biz/period.go](/internal/biz/period.go) | Go | 170 | 25 | 29 | 224 |
| [internal/biz/period\_test.go](/internal/biz/period_test.go) | Go | 810 | 21 | 31 | 862 |
| [internal/biz/plan.go](/internal/biz/plan.go) | Go | 82 | 7 | 18 | 107 |
| [internal/biz/plan\_test.go](/internal/biz/plan_test.go) | Go | 489 | 51 | 82 | 622 |
| [internal/biz/task.go](/internal/biz/task.go) | Go | 425 | 67 | 81 | 573 |
| [internal/biz/task\_repo.go](/internal/biz/task_repo.go) | Go | 14 | 0 | 3 | 17 |
| [internal/biz/task\_test.go](/internal/biz/task_test.go) | Go | 560 | 56 | 137 | 753 |
| [internal/biz/user.go](/internal/biz/user.go) | Go | 247 | 27 | 38 | 312 |
| [internal/biz/user\_repo.go](/internal/biz/user_repo.go) | Go | 10 | 0 | 3 | 13 |
| [internal/biz/user\_test.go](/internal/biz/user_test.go) | Go | 633 | 37 | 175 | 845 |
| [internal/biz/util\_test.go](/internal/biz/util_test.go) | Go | 17 | 1 | 5 | 23 |
| [internal/data/converter.go](/internal/data/converter.go) | Go | 196 | 18 | 36 | 250 |
| [internal/data/data.go](/internal/data/data.go) | Go | 27 | 5 | 9 | 41 |
| [internal/data/models.go](/internal/data/models.go) | Go | 38 | 4 | 5 | 47 |
| [internal/data/repo.go](/internal/data/repo.go) | Go | 194 | 8 | 51 | 253 |
| [internal/data/session.go](/internal/data/session.go) | Go | 41 | 13 | 14 | 68 |
| [internal/data/session\_memory.go](/internal/data/session_memory.go) | Go | 153 | 20 | 42 | 215 |
| [internal/data/systemConfig.go](/internal/data/systemConfig.go) | Go | 189 | 33 | 49 | 271 |
| [internal/model/errors.go](/internal/model/errors.go) | Go | 10 | 4 | 5 | 19 |
| [internal/model/request.go](/internal/model/request.go) | Go | 5 | 0 | 2 | 7 |
| [internal/server/server.go](/internal/server/server.go) | Go | 1 | 0 | 1 | 2 |
| [internal/service/auth\_handle.go](/internal/service/auth_handle.go) | Go | 134 | 8 | 16 | 158 |
| [internal/service/journal\_handle.go](/internal/service/journal_handle.go) | Go | 21 | 1 | 5 | 27 |
| [internal/service/plan\_handle.go](/internal/service/plan_handle.go) | Go | 21 | 0 | 5 | 26 |
| [internal/service/request.go](/internal/service/request.go) | Go | 41 | 0 | 6 | 47 |
| [internal/service/response.go](/internal/service/response.go) | Go | 83 | 11 | 15 | 109 |
| [internal/service/service.go](/internal/service/service.go) | Go | 60 | 6 | 19 | 85 |
| [internal/service/session.go](/internal/service/session.go) | Go | 131 | 18 | 26 | 175 |
| [internal/service/task\_handle.go](/internal/service/task_handle.go) | Go | 95 | 11 | 20 | 126 |
| [internal/service/task\_handle\_test.go](/internal/service/task_handle_test.go) | Go | 56 | 4 | 8 | 68 |
| [internal/service/user\_handle.go](/internal/service/user_handle.go) | Go | 47 | 5 | 8 | 60 |
| [migrations/0001\_init\_schema.down.sql](/migrations/0001_init_schema.down.sql) | MS SQL | 4 | 4 | 4 | 12 |
| [migrations/0001\_init\_schema.up.sql](/migrations/0001_init_schema.up.sql) | MS SQL | 59 | 5 | 13 | 77 |
| [migrations/0002\_init\_system\_data.down.sql](/migrations/0002_init_system_data.down.sql) | MS SQL | 2 | 4 | 3 | 9 |
| [migrations/0002\_init\_system\_data.up.sql](/migrations/0002_init_system_data.up.sql) | MS SQL | 0 | 4 | 2 | 6 |

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)
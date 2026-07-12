# cmd/shortener

В данной директории содержится код, который скомпилируется в бинарное приложение.

Рекомендуется помещать только код, необходимый для запуска приложения, но не бизнес-логику.

Название директории должно соответствовать названию приложения.

Директория `cmd/shortener` содержит:
- точку входа в приложение (функция `main`)
- инициализацию зависимостей (можно вынести в отдельный пакет `internal/app`)
- настройку и запуск HTTP-сервера (можно вынести в отдельный пакет `internal/router`)
- обработку сигналов завершения работы приложения




# Профилирование
# Результаты сравнения профилей

```
$ go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
```

Результаты оптимизации:

### Общее улучшение: Память уменьшилась на 20% (с 2563.29kB до 512.69kB)

Type: inuse_space
Time: 2026-07-06 23:21:16 MSK
Showing nodes accounting for 512.69kB, 20.00% of 2563.29kB total
Dropped 1 node (cum <= 12.82kB)
flat  flat%   sum%        cum   cum%
512.69kB 20.00% 20.00%   512.69kB 20.00%  github.com/jackc/pgx/v5/pgconn.connectOne
0     0% 20.00%   512.69kB 20.00%  database/sql.(*DB).PingContext
0     0% 20.00%   512.69kB 20.00%  database/sql.(*DB).PingContext.func1
0     0% 20.00%   512.69kB 20.00%  database/sql.(*DB).conn
0     0% 20.00%   512.69kB 20.00%  database/sql.(*DB).retry
0     0% 20.00%   512.69kB 20.00%  github.com/DNA-Z/url-shortener/internal/storage.DBConnect
0     0% 20.00%   512.69kB 20.00%  github.com/jackc/pgx/v5.ConnectConfig
0     0% 20.00%   512.69kB 20.00%  github.com/jackc/pgx/v5.connect
0     0% 20.00%   512.69kB 20.00%  github.com/jackc/pgx/v5/pgconn.ConnectConfig
0     0% 20.00%   512.69kB 20.00%  github.com/jackc/pgx/v5/pgconn.connectPreferred
0     0% 20.00%   512.69kB 20.00%  github.com/jackc/pgx/v5/stdlib.(*driverConnector).Connect
0     0% 20.00%   512.69kB 20.00%  main.getDB
0     0% 20.00%   512.69kB 20.00%  main.main
0     0% 20.00%   512.01kB 19.97%  runtime.(*scavengerState).sleep
0     0% 20.00%   512.01kB 19.97%  runtime.(*timer).maybeAdd
0     0% 20.00%   512.01kB 19.97%  runtime.(*timer).modify
0     0% 20.00%   512.01kB 19.97%  runtime.(*timer).reset (inline)
0     0% 20.00%   512.01kB 19.97%  runtime.(*timers).addHeap
0     0% 20.00%     -513kB 20.01%  runtime.allocm
0     0% 20.00%   512.01kB 19.97%  runtime.bgscavenge
0     0% 20.00%   512.01kB 19.97%  runtime.growslice
0     0% 20.00%   512.69kB 20.00%  runtime.main
0     0% 20.00%     -513kB 20.01%  runtime.mstart
0     0% 20.00%     -513kB 20.01%  runtime.mstart0
0     0% 20.00%     -513kB 20.01%  runtime.mstart1
0     0% 20.00%     -513kB 20.01%  runtime.newm
0     0% 20.00%     -513kB 20.01%  runtime.newobject
0     0% 20.00%     -513kB 20.01%  runtime.schedule
0     0% 20.00%     -513kB 20.01%  runtime.startm
0     0% 20.00%     -513kB 20.01%  runtime.wakep
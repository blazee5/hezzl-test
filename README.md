# Тестовое задание для компании Hezzl
## Используемые технологии
- <img width="25" src="https://user-images.githubusercontent.com/25181517/192149581-88194d20-1a37-4be8-8801-5dc0017ffbbe.png" alt="Go" title="Go"/> Go
- <img width="35" src="https://upload.wikimedia.org/wikipedia/commons/9/9c/NATS-logo.png" /> Nats
- <img width="35" src="https://user-images.githubusercontent.com/25181517/117208740-bfb78400-adf5-11eb-97bb-09072b6bedfc.png" alt="PostgreSQL" title="PostgreSQL"/> Postgres
- <img width="35" src="https://asset.brandfetch.io/idnezyZEJm/id_CPPYVKt.jpeg" /> Clickhouse
- <img width="35" src="https://user-images.githubusercontent.com/25181517/182884894-d3fa6ee0-f2b4-4960-9961-64740f533f2a.png" alt="redis" title="redis"/> Redis

## Команда для запуска в docker
```
docker-compose up -d
```
## Команды для запуска миграций
```
make migrate-postgres
make migrate-clickhouse
```

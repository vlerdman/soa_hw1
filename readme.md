# Первое домашнее задание

Поддерживаемые форматы: `json`, `xml`, `yaml`

## Устройство репозитория

- `proxy` - папка с кодом udp-прокси сервера
- `server` - папка с кодом udp-сервера для сериализации/десериализации

## Запуск

```sh
docker-compose up --build
```

## Пример работы

Сервер сериализации/десериализации выдаёт ответ по запросу `get_result`

```sh
echo -n "get_result" | nc -4u -q1 localhost 2001
```
```sh
json - 176 - 10.387µs - 20.326µs
```

```sh
echo -n "get_result" | nc -4u -q1 localhost 2002
```
```sh
xml - 454 - 34.309µs - 41.194µs
```

```sh
echo -n "get_result" | nc -4u -q1 localhost 2003
```
```sh
yaml - 186 - 50.785µs - 33.584µs
```

Прокси сервер принимает запрос вида `get_result {format_name}`
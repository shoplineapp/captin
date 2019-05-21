## Usage

Start captin with api trigger

```sh
go run main.go hooks.json
```

Test with an api call

```sh
curl localhost:8080/api/events -X POST -d '{"event_key":"custom","source":"service_one","payload":{"_id":"xxxxxx"}}'
```

You will see the webhook callback sent from captin

```sh
[GIN] 2019/05/21 - 12:44:41 | 201 |      77.358µs |             ::1 | POST     /api/events
Webhook callback received
[GIN] 2019/05/21 - 12:44:41 | 200 |      21.414µs |             ::1 | POST     /callback
response Body: Received
```
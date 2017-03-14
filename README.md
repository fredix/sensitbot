# sensitbot

sensitbot est un bot Telegram qui permet de piloter l'api de https://www.sensit.io

## Compilation

GOOS=linux GOARCH=amd64 go build -o sensitbot_linux sensitbot.go

## Configuration

ce bot utilise mongoDB.
il est nécessaire de renommer le fichier de configuration sensitbot.toml-example en sensitbot.toml puis de le modifier avec les valeurs personnelles (token telegram, token sensit.io, ...)

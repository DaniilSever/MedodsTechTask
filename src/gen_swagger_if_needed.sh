#!/bin/sh

# Скрипт проверяет, изменились ли комментарии swagger
# и пересоздаёт доку только при необходимости

set -e

SWAGGER_HASH_FILE=".swagger_hash"
CURRENT_HASH=$(find ./app -name "*.go" | xargs cat | grep -E "^// @|^// @title|^// @version" | md5sum | cut -d' ' -f1)

if [ ! -f "$SWAGGER_HASH_FILE" ] || [ "$(cat "$SWAGGER_HASH_FILE")" != "$CURRENT_HASH" ]; then
  echo "Swagger docs changed, regenerating..."
  swag init -g ./app/main.go
  echo "$CURRENT_HASH" > "$SWAGGER_HASH_FILE"
else
  echo "Swagger docs unchanged, skipping swag init"
fi

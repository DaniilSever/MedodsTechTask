#!/bin/bash
# --- Дожидаемся постгрес, чтобы соединения установились нормально ---
echo "Waiting for postgres..."
while ! nc -z postgres-db 5432; do
  sleep 0.1
done
echo "Postgres started Ok"
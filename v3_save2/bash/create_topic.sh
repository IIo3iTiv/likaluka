#!/bin/bash

# Получаем настройки из переменных окружения
TOPIC_NAME=${SH_KAFKA_TOPIC:-"default-topic"}
BOOTSTRAP_SERVER=${SH_KAFKA_BOOTSTRAP_SERVER:-"localhost:9092"}
MAX_RETRIES=${SH_MAX_RETRIES:-30}
SLEEP_TIME=${SH_SLEEP_TIME:-2}

# Функция для проверки доступности Kafka
wait_for_kafka() {
  local retries=0
  until kafka-topics.sh --list --bootstrap-server "$BOOTSTRAP_SERVER" >/dev/null 2>&1; do
    retries=$((retries+1))
    if [ $retries -ge $MAX_RETRIES ]; then
      echo "Kafka не доступен после $MAX_RETRIES попыток"
      exit 1
    fi
    echo "Ожидание Kafka ($retries/$MAX_RETRIES)..."
    sleep $SLEEP_TIME
  done
}

# Основной скрипт
echo "Проверка доступности Kafka..."
wait_for_kafka

echo "Создание топика $TOPIC_NAME..."
kafka-topics.sh --create \
  --topic "$TOPIC_NAME" \
  --bootstrap-server "$BOOTSTRAP_SERVER" \
  --partitions 1 \
  --replication-factor 1 \
  --if-not-exists

echo "Топик $TOPIC_NAME успешно создан"
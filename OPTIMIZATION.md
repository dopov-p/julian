# Оптимизация метода UpdateContents для высокой нагрузки

## Оптимизации на уровне кода

Метод `UpdateContents` уже оптимизирован следующими способами:

1. **Один запрос вместо двух**: Используется один `UPDATE` запрос с `WHERE` по `name`, вместо `SELECT` + `UPDATE`
2. **Атомарность**: Операция выполняется атомарно на уровне БД
3. **Минимум сетевых round-trips**: Один запрос вместо двух
4. **RETURNING для проверки**: Используется `RETURNING id` для проверки успешности обновления без дополнительного запроса
5. **Автоматическое обновление updated_at**: Время обновления устанавливается автоматически

## Оптимизации на уровне базы данных

### 1. Индексы (критически важно)

Для обеспечения высокой производительности метода `UpdateContents` необходимо создать следующие индексы:

```sql
-- Основной индекс для поиска по name (уникальный, если name должен быть уникальным)
CREATE UNIQUE INDEX IF NOT EXISTS idx_cells_name_unique 
ON cells(name) 
WHERE deleted_at IS NULL;

-- Или если name не уникален, используйте обычный индекс:
CREATE INDEX IF NOT EXISTS idx_cells_name 
ON cells(name) 
WHERE deleted_at IS NULL;

-- Составной индекс для оптимизации WHERE условий (name + deleted_at)
CREATE INDEX IF NOT EXISTS idx_cells_name_deleted_at 
ON cells(name, deleted_at) 
WHERE deleted_at IS NULL;
```

**Рекомендация**: Используйте частичный индекс (`WHERE deleted_at IS NULL`), так как метод обновляет только неудаленные записи. Это:
- Уменьшает размер индекса
- Ускоряет поиск (меньше данных для сканирования)
- Экономит место на диске

### 2. Индексы для JSONB поля Contents (опционально)

Если планируется поиск или фильтрация по содержимому `Contents`, можно создать GIN индекс:

```sql
-- GIN индекс для быстрого поиска по JSONB полю
CREATE INDEX IF NOT EXISTS idx_cells_contents_gin 
ON cells USING GIN (contents);

-- Или более специфичный индекс для поиска по конкретным полям в JSONB
CREATE INDEX IF NOT EXISTS idx_cells_contents_sku_gin 
ON cells USING GIN ((contents->>'sku'));
```

**Примечание**: GIN индексы увеличивают время на INSERT/UPDATE, но значительно ускоряют поиск по JSONB.

### 3. Оптимизация таблицы

```sql
-- Убедитесь, что поле contents имеет тип JSONB (не JSON)
-- JSONB быстрее для операций обновления и поддерживает индексы
ALTER TABLE cells ALTER COLUMN contents TYPE JSONB USING contents::JSONB;

-- Оптимизация autovacuum для таблицы с частыми обновлениями
ALTER TABLE cells SET (
    autovacuum_vacuum_scale_factor = 0.1,
    autovacuum_analyze_scale_factor = 0.05
);
```

### 4. Connection Pool настройки

Для высокой нагрузки рекомендуется настроить connection pool:

```go
// В config/cluster.go при создании pgxpool.Pool
config, err := pgxpool.ParseConfig(dsn)
if err != nil {
    return nil, err
}

// Оптимизация для высокой нагрузки
config.MaxConns = 25                    // Максимум соединений
config.MinConns = 5                     // Минимум соединений в pool
config.MaxConnLifetime = time.Hour     // Время жизни соединения
config.MaxConnIdleTime = time.Minute * 30 // Время простоя соединения
config.HealthCheckPeriod = time.Minute // Период проверки здоровья

dbpool, err := pgxpool.NewWithConfig(ctx, config)
```

### 5. Подготовленные запросы (Prepared Statements)

Для еще большей оптимизации можно использовать prepared statements:

```go
// В repo.go добавить кэш prepared statements
type Repo struct {
    cluster        *config.Cluster
    updateContents *pgxpool.Pool // или использовать pgxpool с prepared statements
}

// Использовать pgxpool с автоматическим кэшированием prepared statements
// pgxpool автоматически кэширует prepared statements для часто используемых запросов
```

### 6. Мониторинг и анализ производительности

```sql
-- Проверка использования индексов
EXPLAIN ANALYZE 
UPDATE cells 
SET contents = $1, updated_at = NOW() 
WHERE name = $2 AND deleted_at IS NULL 
RETURNING id;

-- Проверка статистики таблицы
SELECT schemaname, tablename, n_tup_upd, n_tup_hot_upd 
FROM pg_stat_user_tables 
WHERE tablename = 'cells';

-- Проверка размера индексов
SELECT indexname, pg_size_pretty(pg_relation_size(indexname::regclass)) as size
FROM pg_indexes 
WHERE tablename = 'cells';
```

### 7. Дополнительные рекомендации

1. **Пакетные обновления**: Если нужно обновить много записей, рассмотрите возможность пакетных обновлений через `COPY` или множественные `UPDATE` в одной транзакции.

2. **Партиционирование**: Если таблица очень большая, рассмотрите партиционирование по `name` или другим критериям.

3. **Read Replicas**: Для чтения используйте read replicas, оставляя write операции на primary.

4. **Логирование медленных запросов**: Включите `log_min_duration_statement` для отслеживания медленных запросов.

```sql
-- В postgresql.conf или через ALTER SYSTEM
SET log_min_duration_statement = 100; -- логировать запросы > 100ms
```

## Пример создания всех необходимых индексов

```sql
-- Миграция для создания индексов
BEGIN;

-- Основной индекс для поиска по name (частичный индекс для неудаленных записей)
CREATE UNIQUE INDEX IF NOT EXISTS idx_cells_name_unique 
ON cells(name) 
WHERE deleted_at IS NULL;

-- Индекс для updated_at (если нужна сортировка по времени обновления)
CREATE INDEX IF NOT EXISTS idx_cells_updated_at 
ON cells(updated_at DESC) 
WHERE deleted_at IS NULL;

COMMIT;

-- Анализ таблицы для обновления статистики
ANALYZE cells;
```

## Метрики для мониторинга

Следите за следующими метриками:
- Время выполнения `UpdateContents` (должно быть < 10ms при наличии индексов)
- Количество dead tuples (должно быть низким при правильном autovacuum)
- Размер индексов (не должен расти бесконечно)
- Hit ratio для индексов (должен быть > 95%)


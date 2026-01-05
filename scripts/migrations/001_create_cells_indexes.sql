-- Миграция для создания индексов для оптимизации UpdateContents
-- Выполнить перед использованием метода UpdateContents в production

BEGIN;

-- Основной уникальный индекс для поиска по name (только для неудаленных записей)
-- Этот индекс критически важен для производительности UpdateContents
CREATE UNIQUE INDEX IF NOT EXISTS idx_cells_name_unique 
ON cells(name) 
WHERE deleted_at IS NULL;

-- Индекс для updated_at (опционально, если нужна сортировка по времени обновления)
CREATE INDEX IF NOT EXISTS idx_cells_updated_at 
ON cells(updated_at DESC) 
WHERE deleted_at IS NULL;

-- GIN индекс для JSONB поля contents (опционально, если планируется поиск по содержимому)
-- ВНИМАНИЕ: Этот индекс замедляет INSERT/UPDATE, но ускоряет поиск по JSONB
-- Раскомментируйте только если нужен поиск по содержимому Contents
-- CREATE INDEX IF NOT EXISTS idx_cells_contents_gin 
-- ON cells USING GIN (contents);

COMMIT;

-- Обновление статистики для оптимизатора запросов
ANALYZE cells;


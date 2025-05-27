-- 创建复制用户
CREATE ROLE repl WITH REPLICATION LOGIN PASSWORD '123456';

-- 可选：创建复制槽（不建议自动创建，可手动执行）
-- SELECT * FROM pg_create_physical_replication_slot('replica_slot');
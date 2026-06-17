CREATE TABLE IF NOT EXISTS command_type (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT OR IGNORE INTO command_type (id, name)
VALUES 
(1, 'Rest'),
(2, 'SQL');

CREATE TABLE IF NOT EXISTS command (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    command_query TEXT NOT NULL,
    type_id INTEGER REFERENCES command_type(id) NOT NULL
);

CREATE VIRTUAL TABLE IF NOT EXISTS command_fts_idx USING fts5(
    name,
    description,
    content='command',
    content_rowid='id'
);

-- Triggers to update the FTS index in case the commands change
CREATE TRIGGER IF NOT EXISTS commandInserted AFTER INSERT ON command BEGIN
  INSERT INTO command_fts_idx(rowid, name, description) VALUES (new.id, new.name, new.description);
END;
CREATE TRIGGER IF NOT EXISTS commandDeleted AFTER DELETE ON command BEGIN
  INSERT INTO command_fts_idx(command_fts_idx, rowid, name, description) VALUES('delete', old.id, old.name, old.description);
END;
CREATE TRIGGER IF NOT EXISTS commandUpdated AFTER UPDATE ON command BEGIN
  INSERT INTO command_fts_idx(command_fts_idx, rowid, name, description) VALUES('delete', old.id, old.name, old.description);
  INSERT INTO command_fts_idx(rowid, name, description) VALUES (new.id, new.name, new.description);
END;


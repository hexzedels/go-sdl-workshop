CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    display_name TEXT NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL DEFAULT 'user',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    owner_id INTEGER NOT NULL,
    locale TEXT NOT NULL DEFAULT 'en',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS admin_secrets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT NOT NULL,
    value TEXT NOT NULL
);

-- Seed users (passwords: admin123, password123)
-- bcrypt hashes generated at cost 10
INSERT OR IGNORE INTO users (id, username, password_hash, display_name, bio, role) VALUES
    (1, 'admin', '$2a$10$GkKaGTG0g6jZ.9ja3MqNHur7z2cbc6YFPgqwgNcH2MkEXkxwsHboG', 'Admin User', 'System administrator', 'admin'),
    (2, 'user', '$2a$10$Cbp3qX4boU.DxD6xNR6sHea.xrmGhaETw2x/uG6y8AydPYwP.yh02', 'Regular User', 'Just a normal user', 'user');

INSERT OR IGNORE INTO documents (id, title, content, owner_id, locale) VALUES
    (1, 'Getting Started', 'Welcome to the document management system.', 2, 'en'),
    (2, 'Project Roadmap', 'Q1: Launch MVP, Q2: Scale infrastructure.', 2, 'en'),
    (3, 'Meeting Notes', 'Discussed quarterly goals and team allocation.', 2, 'en'),
    (4, 'CONFIDENTIAL: Salaries', 'CEO: 500000, CTO: 450000, Engineer: 150000', 1, 'en');

INSERT OR IGNORE INTO admin_secrets (id, key, value) VALUES
    (1, 'documents_flag', '{{B64:R09TREx7NWY3NGU3NDFlODk0MTAyYX0=}}');

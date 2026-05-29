-- ============================================================
--  BASE DE DADOS DE TREINO — TABELA DE USUÁRIOS
--  Para usar com PostgreSQL
-- ============================================================


-- 1. CRIAR O BANCO (rode isso no terminal, não dentro do psql)
-- createdb treino_api

-- 2. CRIAR A TABELA
CREATE TABLE usuarios (
    id               SERIAL PRIMARY KEY,
    nome             VARCHAR(150)        NOT NULL,
    data_nascimento  DATE                NOT NULL,
    email            VARCHAR(255)        NOT NULL UNIQUE,
    senha_hash       VARCHAR(255)        NOT NULL,
    status_contrato  VARCHAR(30)         NOT NULL DEFAULT 'ativo'
                       CHECK (status_contrato IN ('ativo', 'suspenso', 'cancelado', 'pendente')),
    id_contrato      VARCHAR(20)         NOT NULL UNIQUE,
    ativo            BOOLEAN             NOT NULL DEFAULT TRUE,
    ultimo_login     TIMESTAMPTZ,
    criado_em        TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    atualizado_em    TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);

-- Índices úteis para treinar queries
CREATE INDEX idx_usuarios_email          ON usuarios (email);
CREATE INDEX idx_usuarios_status         ON usuarios (status_contrato);
CREATE INDEX idx_usuarios_ativo          ON usuarios (ativo);
CREATE INDEX idx_usuarios_ultimo_login   ON usuarios (ultimo_login);


-- ============================================================
-- 3. FUNÇÃO AUTO-UPDATE de atualizado_em
--    (boa prática — treina trigger também!)
-- ============================================================
CREATE OR REPLACE FUNCTION atualizar_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.atualizado_em = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_usuarios_update
BEFORE UPDATE ON usuarios
FOR EACH ROW EXECUTE FUNCTION atualizar_timestamp();


-- ============================================================
-- 4. DADOS FAKE (20 usuários)
--    Senhas são hashes bcrypt simulados (não funcionais)
-- ============================================================
INSERT INTO usuarios
    (nome, data_nascimento, email, senha_hash, status_contrato, id_contrato, ativo, ultimo_login, criado_em)
VALUES
    ('Ana Clara Souza',       '1992-03-14', 'ana.souza@email.com',       '$2b$12$KIX1fakeHashAna0000001', 'ativo',      'CTR-2024-0001', TRUE,  '2025-05-28 09:15:00+00', '2024-01-10 08:00:00+00'),
    ('Bruno Mendes',          '1988-07-22', 'bruno.mendes@email.com',    '$2b$12$KIX1fakeHashBru0000002', 'ativo',      'CTR-2024-0002', TRUE,  '2025-05-27 14:30:00+00', '2024-02-05 10:30:00+00'),
    ('Carla Ferreira',        '1995-11-01', 'carla.ferreira@email.com',  '$2b$12$KIX1fakeHashCar0000003', 'suspenso',   'CTR-2024-0003', FALSE, '2025-03-10 18:00:00+00', '2024-02-20 09:00:00+00'),
    ('Diego Alves',           '1990-06-30', 'diego.alves@email.com',     '$2b$12$KIX1fakeHashDie0000004', 'ativo',      'CTR-2024-0004', TRUE,  '2025-05-29 07:45:00+00', '2024-03-01 11:00:00+00'),
    ('Eduarda Lima',          '1998-02-18', 'eduarda.lima@email.com',    '$2b$12$KIX1fakeHashEdu0000005', 'cancelado',  'CTR-2024-0005', FALSE, '2024-12-01 20:00:00+00', '2024-03-15 14:00:00+00'),
    ('Felipe Costa',          '1985-09-09', 'felipe.costa@email.com',    '$2b$12$KIX1fakeHashFel0000006', 'ativo',      'CTR-2024-0006', TRUE,  '2025-05-26 16:20:00+00', '2024-04-02 08:30:00+00'),
    ('Gabriela Nunes',        '1993-12-25', 'gabriela.nunes@email.com',  '$2b$12$KIX1fakeHashGab0000007', 'ativo',      'CTR-2024-0007', TRUE,  '2025-05-29 11:00:00+00', '2024-04-18 09:45:00+00'),
    ('Henrique Rocha',        '1991-04-04', 'henrique.rocha@email.com',  '$2b$12$KIX1fakeHashHen0000008', 'pendente',   'CTR-2024-0008', FALSE, NULL,                     '2024-05-03 13:00:00+00'),
    ('Isabela Martins',       '1997-08-16', 'isabela.martins@email.com', '$2b$12$KIX1fakeHashIsa0000009', 'ativo',      'CTR-2024-0009', TRUE,  '2025-05-28 22:10:00+00', '2024-05-20 10:00:00+00'),
    ('João Pedro Oliveira',   '1986-01-31', 'joao.oliveira@email.com',   '$2b$12$KIX1fakeHashJoa0000010', 'ativo',      'CTR-2024-0010', TRUE,  '2025-05-25 08:00:00+00', '2024-06-01 07:30:00+00'),
    ('Karen Batista',         '1994-05-07', 'karen.batista@email.com',   '$2b$12$KIX1fakeHashKar0000011', 'suspenso',   'CTR-2024-0011', FALSE, '2025-01-15 12:00:00+00', '2024-06-14 11:15:00+00'),
    ('Lucas Teixeira',        '1989-10-20', 'lucas.teixeira@email.com',  '$2b$12$KIX1fakeHashLuc0000012', 'ativo',      'CTR-2024-0012', TRUE,  '2025-05-29 06:30:00+00', '2024-07-05 09:00:00+00'),
    ('Mariana Ribeiro',       '1996-03-03', 'mariana.ribeiro@email.com', '$2b$12$KIX1fakeHashMar0000013', 'ativo',      'CTR-2024-0013', TRUE,  '2025-05-27 19:45:00+00', '2024-07-22 14:00:00+00'),
    ('Nicolas Carvalho',      '1992-07-11', 'nicolas.carvalho@email.com','$2b$12$KIX1fakeHashNic0000014', 'cancelado',  'CTR-2024-0014', FALSE, '2024-11-20 10:00:00+00', '2024-08-08 10:30:00+00'),
    ('Olivia Santos',         '1999-09-29', 'olivia.santos@email.com',   '$2b$12$KIX1fakeHashOli0000015', 'ativo',      'CTR-2024-0015', TRUE,  '2025-05-29 10:00:00+00', '2024-08-25 08:00:00+00'),
    ('Pedro Henrique Gomes',  '1987-12-05', 'pedro.gomes@email.com',     '$2b$12$KIX1fakeHashPed0000016', 'ativo',      'CTR-2024-0016', TRUE,  '2025-05-28 15:00:00+00', '2024-09-10 11:00:00+00'),
    ('Rafaela Pereira',       '1993-02-14', 'rafaela.pereira@email.com', '$2b$12$KIX1fakeHashRaf0000017', 'pendente',   'CTR-2024-0017', FALSE, NULL,                     '2024-09-28 09:30:00+00'),
    ('Samuel Araújo',         '1990-06-06', 'samuel.araujo@email.com',   '$2b$12$KIX1fakeHashSam0000018', 'ativo',      'CTR-2024-0018', TRUE,  '2025-05-26 21:00:00+00', '2024-10-15 13:00:00+00'),
    ('Tatiane Moreira',       '1995-11-17', 'tatiane.moreira@email.com', '$2b$12$KIX1fakeHashTat0000019', 'ativo',      'CTR-2024-0019', TRUE,  '2025-05-29 08:45:00+00', '2024-11-01 10:00:00+00'),
    ('Victor Hugo Dias',      '1988-04-23', 'victor.dias@email.com',     '$2b$12$KIX1fakeHashVic0000020', 'suspenso',   'CTR-2024-0020', FALSE, '2025-02-10 17:30:00+00', '2024-11-20 08:15:00+00');


-- ============================================================
-- 5. QUERIES DE EXEMPLO PARA TREINAR
-- ============================================================

-- Todos os usuários ativos
-- SELECT * FROM usuarios WHERE ativo = TRUE;

-- Usuários por status de contrato
-- SELECT nome, email, status_contrato FROM usuarios WHERE status_contrato = 'ativo';

-- Último login nos últimos 7 dias
-- SELECT nome, ultimo_login FROM usuarios WHERE ultimo_login >= NOW() - INTERVAL '7 days';

-- Contar usuários por status
-- SELECT status_contrato, COUNT(*) as total FROM usuarios GROUP BY status_contrato;

-- Buscar por email (para login)
-- SELECT id, nome, senha_hash, ativo FROM usuarios WHERE email = 'ana.souza@email.com';

-- Usuários que nunca fizeram login
-- SELECT nome, email, criado_em FROM usuarios WHERE ultimo_login IS NULL;

-- Atualizar último login (simular login)
-- UPDATE usuarios SET ultimo_login = NOW() WHERE email = 'ana.souza@email.com';

-- Soft delete (desativar usuário)
-- UPDATE usuarios SET ativo = FALSE WHERE id = 1;

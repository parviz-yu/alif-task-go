CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    uuid VARCHAR(36) NOT NULL
);

CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    is_identified BOOLEAN NOT NULL,
    balance BIGINT NOT NULL,
    user_id INT NOT NULL REFERENCES users(id)
);

CREATE TABLE replenishments (
    id SERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    date TIMESTAMP NOT NULL,
    wallet_id INT NOT NULL REFERENCES wallets(id)
);

-- Next operations needed to demonstrate API routes
INSERT INTO users (email, uuid) VALUES ('user-1@example.org', 'da06798f-857c-4891-b807-f26b0157f31d');
INSERT INTO users (email, uuid) VALUES ('user-2@example.org', '35d0c993-f1cb-4576-9078-ca0ebe1576c0');
INSERT INTO users (email, uuid) VALUES ('user-3@example.org', 'f7203bd1-1354-4908-8620-6cc7cd715572');

INSERT INTO wallets (is_identified, balance, user_id) VALUES (false, 152000, 1);
INSERT INTO replenishments (amount, date, wallet_id) VALUES (152000, '2023-02-25 18:50:51' , 1);

INSERT INTO wallets (is_identified, balance, user_id) VALUES (true, 5000000, 2);
INSERT INTO replenishments (amount, date, wallet_id) VALUES (5000000, '2023-03-13 17:00:21' , 2);

INSERT INTO wallets (is_identified, balance, user_id) VALUES (false, 781400, 3);
INSERT INTO replenishments (amount, date, wallet_id) VALUES (781400, '2023-03-05 10:20:31' , 3);
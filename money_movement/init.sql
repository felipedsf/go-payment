CREATE USER IF NOT EXISTS 'money_movement_user'@'localhost' IDENTIFIED BY 'Auth123';

CREATE DATABASE IF NOT EXISTS money_movement;

GRANT ALL PRIVILEGES ON money_movement.* to 'money_movement_user'@'localhost';

USE money_movement;

CREATE TABLE IF NOT EXISTS `wallet` (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE ,
    wallet_type VARCHAR(255) NOT NULL,
    INDEX (user_id)
);

CREATE TABLE IF NOT EXISTS `account` (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    cents INT NOT NULL DEFAULT 0,
    account_type VARCHAR(255) NOT NULL,
    wallet_id INT NOT NULL ,
    FOREIGN KEY (wallet_id) REFERENCES wallet(id)
);

CREATE TABLE IF NOT EXISTS `transaction`(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    pid VARCHAR(255) NOT NULL,

    src_user_id INT NOT NULL,
    dst_user_id INT NOT NULL,

    src_wallet_id INT NOT NULL,
    dst_wallet_id INT NOT NULL,

    src_wallet_type VARCHAR(255) NOT NULL,
    dst_wallet_type VARCHAR(255) NOT NULL,

    final_dst_merchant_wallet_id INT NOT NULL,
    amount INT NOT NULL,

    INDEX (pid)
);

-- merchant and customer wallets
INSERT INTO wallet (id, user_id, wallet_type)
VALUES
    (1, 'customer@email.com', 'CUSTOMER'),
    (2, 'merchant@email.com', 'MERCHANT');

-- customer accounts
INSERT INTO account(cents, account_type, wallet_id)
VALUES
    (5000000, 'DEFAULT', 1),
    (0, 'PAYMENT', 1);

-- merchant accounts
INSERT INTO account(cents, account_type, wallet_id)
VALUES
    (0, 'INCOMING', 2);
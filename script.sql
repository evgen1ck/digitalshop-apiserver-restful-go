SET timezone = 'UTC';

CREATE SCHEMA IF NOT EXISTS account;
CREATE SCHEMA IF NOT EXISTS product;
CREATE SCHEMA IF NOT EXISTS payment;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA account;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA product;



DROP TABLE IF EXISTS account.status CASCADE;
CREATE TABLE account.status
(
	status_no	smallserial	,	 
	status_name text 		UNIQUE NOT NULL,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	text		NULL,
	PRIMARY KEY (status_no)
);
INSERT INTO account.status(status_name) VALUES ('not activated'), ('active'), ('blocked'), ('deleted');



DROP TABLE IF EXISTS account.type_registration CASCADE;
CREATE TABLE account.type_registration
(
    type_registration_no	smallserial	,
    type_registration_name  text        UNIQUE NOT NULL,
    created_at              timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    PRIMARY KEY (type_registration_no)
);
INSERT INTO account.type_registration(type_registration_name) VALUES ('user'), ('telegram'), ('employee'), ('vendor');



DROP TABLE IF EXISTS account.registration_temp CASCADE;
CREATE TABLE account.registration_temp
(
    registration_temp_no    bigserial   ,
    email                   text        NOT NULL,
    password                text        NOT NULL,
    confirmation_code       numeric     NOT NULL,
    expiration              timestamp   NOT NULL,
    PRIMARY KEY (registration_temp_no)
);



DROP TABLE IF EXISTS account.account CASCADE;
CREATE TABLE account.account
(
    account_id              uuid        DEFAULT account.UUID_GENERATE_V4(),
    account_status			smallint 	NOT NULL DEFAULT 1,
    last_change_status      timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type_registration       smallint    NOT NULL,
	timestamp_last_activity	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at              timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    PRIMARY KEY (account_id),
    FOREIGN KEY (account_status) REFERENCES account.status(status_no),
    FOREIGN KEY (type_registration) REFERENCES account.type_registration(type_registration_no)
);



DROP TABLE IF EXISTS account.user CASCADE;
CREATE TABLE account.user
(
    account_id              uuid        UNIQUE NOT NULL,
    username                text        UNIQUE NOT NULL,
    password 				bytea		NOT NULL,
	salt_for_password       text        NOT NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_id) REFERENCES account.account(account_id)
);



DROP TABLE IF EXISTS account.telegram_user CASCADE;
CREATE TABLE account.telegram_user
(
    account_id              uuid        UNIQUE NOT NULL,
    telegram_id             text        UNIQUE NOT NULL,
    username                text        NOT NULL,
    photo_url               text        NOT NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_id) REFERENCES account.account(account_id)
);



DROP TABLE IF EXISTS account.employee CASCADE;
CREATE TABLE account.employee
(
    account_id              uuid        UNIQUE NOT NULL,
    surname                 text        NOT NULL,
    name                    text        NOT NULL,
    patronymic              text        NULL,
    password 				bytea		NOT NULL,
    salt_for_password       text        NOT NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_id) REFERENCES account.account(account_id)
);


















DROP TABLE IF EXISTS product.status CASCADE;
CREATE TABLE product.status
(
    status_no	smallserial	,
    status_name text 		UNIQUE NOT NULL,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	text		NULL,
    PRIMARY KEY (status_no)
);

-- type for products
DROP TABLE IF EXISTS product.type CASCADE;
CREATE TABLE product.type
(
    type_no	    smallserial	,
    type_name   text 		UNIQUE NOT NULL,
    created_at              timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    PRIMARY KEY (type_no)
);

-- platform for products
DROP TABLE IF EXISTS product.platform CASCADE;
CREATE TABLE product.platform
(
    platform_no	    smallserial	,
    platform_name   text 		UNIQUE NOT NULL,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at    	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	    text		NULL,
    PRIMARY KEY (platform_no)
);

-- products
DROP TABLE IF EXISTS product.product CASCADE;
CREATE TABLE product.product
(
    product_id      uuid        DEFAULT product.UUID_GENERATE_V4(),
    product_type    smallint    NOT NULL,
    product_name    text        UNIQUE NOT NULL,
    product_desc    text        NOT NULL,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		text		NULL,
    PRIMARY KEY (product_id),
    FOREIGN KEY (product_type) REFERENCES product.type(type_no)
);

-- product variants
DROP TABLE IF EXISTS product.variant CASCADE;
CREATE TABLE product.variant
(
    product_id      uuid        NOT NULL,
    variant_name    text        UNIQUE NOT NULL,
    variant_desc    text        NOT NULL,
    status_no       smallint    NOT NULL,
    quantity        integer     NOT NULL CHECK ( quantity >= -1 ) DEFAULT 0,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		text		NULL,
    PRIMARY KEY (product_id, variant_name),
    FOREIGN KEY (product_id) REFERENCES product.product(product_id),
    FOREIGN KEY (status_no) REFERENCES product.status(status_no)
);







--GRANT USAGE ON SCHEMA xxxx TO user;
--GRANT SELECT, UPDATE, INSERT, DELETE ON ALL TABLES IN SCHEMA xxxx TO user;
--GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA xxxx TO user;


ALTER SYSTEM SET max_wal_size = '1GB';
ALTER SYSTEM SET autovacuum_vacuum_scale_factor = 0.0;
ALTER SYSTEM SET autovacuum_analyze_scale_factor = 0.0;
ALTER SYSTEM SET autovacuum_vacuum_cost_limit = 10000;
ALTER SYSTEM SET autovacuum_vacuum_cost_delay = 0;

CREATE OR REPLACE FUNCTION account.cleanup_registration_temp() RETURNS VOID AS $$
BEGIN
    DELETE FROM account.registration_temp WHERE expiration < NOW();
END;
$$ LANGUAGE plpgsql;
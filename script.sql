SET timezone = 'UTC';

CREATE SCHEMA IF NOT EXISTS optional;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA optional;
CREATE SCHEMA IF NOT EXISTS account;
CREATE SCHEMA IF NOT EXISTS employee;
CREATE SCHEMA IF NOT EXISTS product;
CREATE SCHEMA IF NOT EXISTS vendor;
CREATE SCHEMA IF NOT EXISTS payment;



-- status for accounts
DROP TABLE IF EXISTS account.status CASCADE;
CREATE TABLE account.status
(
	status_no	smallserial	,	 
	status_name text 		UNIQUE NOT NULL,
    commentary	text 	   	NULL,
    modified_at	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (status_no)
);
INSERT INTO account.status(status_name) VALUES ('account not verified');

-- accounts
DROP TABLE IF EXISTS account.account CASCADE;
CREATE TABLE account.account
(
    account_id              uuid        DEFAULT optional.UUID_GENERATE_V4(),
	account_email			text		UNIQUE NOT NULL,
	account_status			smallint 	NOT NULL DEFAULT 1,
	account_nickname		text		UNIQUE NOT NULL,
	password 				bytea		NOT NULL,
	saltForPassword         varchar     NOT NULL,
	timestamp_reg			timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
	timestamp_last_activity	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
	reg_ip					inet		NOT NULL,
    commentary			    text		NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (account_id),
	FOREIGN KEY (account_status) REFERENCES account.status(status_no)
);

-- type for products
DROP TABLE IF EXISTS product.type CASCADE;
CREATE TABLE product.type
(
    type_no	    smallserial	,
    type_name   text 		UNIQUE NOT NULL,
    commentary	text 	   	NULL,
    modified_at	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (type_no)
);

-- platform for products
DROP TABLE IF EXISTS product.platform CASCADE;
CREATE TABLE product.platform
(
    platform_no	    smallserial	,
    platform_name   text 		UNIQUE NOT NULL,
    commentary	    text 	   	NULL,
    modified_at	    timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (platform_no)
);

-- products
DROP TABLE IF EXISTS product.product CASCADE;
CREATE TABLE product.product
(
    product_id      uuid        DEFAULT optional.UUID_GENERATE_V4(),
    product_type    smallint    NOT NULL,
    product_name    text        UNIQUE NOT NULL,
    product_desc    text        NOT NULL,
    commentary      text        NOT NULL,
    modified_at	    timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
    quantity        integer     NOT NULL CHECK ( quantity >= -1 ) DEFAULT 0,
    commentary      text        NOT NULL,
    modified_at	    timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, variant_name),
    FOREIGN KEY (product_id) REFERENCES product.product(product_id)
);







--GRANT USAGE ON SCHEMA xxxx TO user;
--GRANT SELECT, UPDATE, INSERT, DELETE ON ALL TABLES IN SCHEMA xxxx TO user;
--GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA xxxx TO user;
SET timezone = 'UTC';

CREATE SCHEMA IF NOT EXISTS account;
CREATE SCHEMA IF NOT EXISTS product;
CREATE SCHEMA IF NOT EXISTS payment;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA account;

select * from account.registration_temp;
select * from account.user;
select * from account.account;



DROP TABLE IF EXISTS account.state CASCADE;
CREATE TABLE account.state
(
	state_no	smallserial	,
	state_name  text 		NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	text		NULL,
	PRIMARY KEY (state_no)
);
INSERT INTO account.state(state_name) VALUES ('active'), ('blocked'), ('deleted');



DROP TABLE IF EXISTS account.registration_method CASCADE;
CREATE TABLE account.registration_method
(
    registration_method_no	    smallserial	,
    registration_method_name    text        NOT NULL UNIQUE,
    created_at                  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         	    timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			        text		NULL,
    PRIMARY KEY (registration_method_no)
);
INSERT INTO account.registration_method(registration_method_name) VALUES ('web application'), ('telegram account'), ('google account');



CREATE TABLE account.role
(
    role_no      smallserial ,
    role_name    text        NOT NULL UNIQUE,
    created_at   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at  timestamp	 NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary   text		 NULL,
    PRIMARY KEY (role_no)
);
INSERT INTO account.role(role_name) VALUES ('user'), ('admin');



DROP TABLE IF EXISTS account.registration_temp_data CASCADE;
CREATE TABLE account.registration_temp_data
(
    confirmation_token      text        ,
    nickname                text        NOT NULL,
    email                   text        NOT NULL,
    password                text        NOT NULL,
    expiration              timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP + interval '10 minute',
    PRIMARY KEY (confirmation_token)
);



DROP TABLE IF EXISTS account.account CASCADE;
CREATE TABLE account.account
(
    account_id              uuid        DEFAULT account.UUID_GENERATE_V4(),
    account_state   	    smallint 	NOT NULL DEFAULT 1,
    account_role            smallint    NOT NULL DEFAULT 1,
    last_change_state       timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    registration_method     smallint    NOT NULL DEFAULT 1,
	timestamp_last_activity	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at              timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    PRIMARY KEY (account_id),
    FOREIGN KEY (account_state) REFERENCES account.status(status_no),
    FOREIGN KEY (account_role) REFERENCES account.role(role_no),
    FOREIGN KEY (registration_method) REFERENCES account.registration_method(registration_method_no)
);
SELECT EXISTS(select account_id from account.account aa inner join account.role ar on AA.account_role = ar.role_no where account_id = '' and account_role = '2');
DROP TABLE IF EXISTS account.user CASCADE;
CREATE TABLE account.user
(
    account_id              uuid        NOT NULL UNIQUE,
    email                   text        NOT NULL UNIQUE,
    nickname                text        NOT NULL UNIQUE,
    password 				text		NOT NULL,
	salt_for_password       text        NOT NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_id) REFERENCES account.account(account_id)
);

-- DROP TABLE IF EXISTS account.telegram_user CASCADE;
-- CREATE TABLE account.telegram_user
-- (
--     account_id               uuid        NOT NULL UNIQUE,
--     telegram_id              text        NOT NULL UNIQUE,
--     username                 text        NOT NULL,
--     photo_url                text        NOT NULL,
--     modified_at         	    timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     commentary			    text		NULL,
--     FOREIGN KEY (account_id) REFERENCES account.account(account_id)
-- );

DROP TABLE IF EXISTS account.employee CASCADE;
CREATE TABLE account.employee
(
    account_id              uuid        NOT NULL UNIQUE,
    surname                 text        NOT NULL,
    name                    text        NOT NULL,
    patronymic              text        NULL,
    password 				bytea		NOT NULL,
    salt_for_password       text        NOT NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_id) REFERENCES account.account(account_id)
);










DROP TABLE IF EXISTS product.state CASCADE;
CREATE TABLE product.state
(
    state_no	smallserial	,
    state_name  text 	    NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	text		NULL,
    PRIMARY KEY (state_no)
);
INSERT INTO account.state(state_name) VALUES ('active'), ('temporarily unavailable'), ('blocked'), ('deleted');



DROP TABLE IF EXISTS product.type CASCADE;
CREATE TABLE product.type
(
    type_no	    smallserial	,
    type_name   text 		NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary  text		NULL,
    PRIMARY KEY (type_no)
);
INSERT INTO product.type(type_name) VALUES ('game'), ('other');



DROP TABLE IF EXISTS product.category CASCADE;
CREATE TABLE product.category
(
    category_no	    smallserial	,
    category_name   text 		NOT NULL UNIQUE,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary      text		NULL,
    PRIMARY KEY (category_no)
);
INSERT INTO product.category(category_name) VALUES ('key'), ('gift'), ('text'), ('link'), ('code');



DROP TABLE IF EXISTS product.service CASCADE;
CREATE TABLE product.service
(
    service_no	    smallserial	,
    service_name    text 		NOT NULL UNIQUE,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at    	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	    text		NULL,
    PRIMARY KEY (service_no)
);
INSERT INTO product.service(service_name) VALUES ('Steam'), ('Ubisoft'), ('Epic Games'), ('Electronic Arts'), ('Ozon'), ('Wildberries'), ('Ivi'), ('YouTube');



DROP TABLE IF EXISTS product.product CASCADE;
CREATE TABLE product.product
(
    product_id      uuid        DEFAULT account.UUID_GENERATE_V4(),
    product_type    smallint    NOT NULL,
    product_name    text        NOT NULL UNIQUE,
    product_desc    text        NOT NULL,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		text		NULL,
    PRIMARY KEY (product_id),
    FOREIGN KEY (product_type) REFERENCES product.type(type_no)
);



DROP TABLE IF EXISTS product.variant CASCADE;
CREATE TABLE product.variant
(
    product_id          uuid        NOT NULL,
    variant_name        text        NOT NULL UNIQUE,
    variant_desc        text        NOT NULL,
    product_service     smallint    NOT NULL,
    product_category    smallint    NOT NULL,
    product_state       smallint    NOT NULL,
    last_change_state   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    quantity            integer     NOT NULL CHECK ( quantity >= -1 ) DEFAULT 0,
    created_at          timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		    text		NULL,
    PRIMARY KEY (product_id, variant_name),
    FOREIGN KEY (product_id) REFERENCES product.product(product_id),
    FOREIGN KEY (product_state) REFERENCES product.state(state_no),
    FOREIGN KEY (product_service) REFERENCES product.service(service_no),
    FOREIGN KEY (product_category) REFERENCES product.category(category_no)
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
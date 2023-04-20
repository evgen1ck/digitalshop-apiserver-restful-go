SET timezone = 'UTC';

CREATE SCHEMA IF NOT EXISTS account;
CREATE SCHEMA IF NOT EXISTS product;
CREATE SCHEMA IF NOT EXISTS payment;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA account;

select * from account.user;
select * from account.account;



DROP TABLE IF EXISTS account.state CASCADE;
CREATE TABLE account.state
(
	state_no	smallserial	PRIMARY KEY,
	state_name  text 		NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	text		NULL
);
INSERT INTO account.state(state_name) VALUES ('active'), ('blocked'), ('deleted');



DROP TABLE IF EXISTS account.registration_method CASCADE;
CREATE TABLE account.registration_method
(
    registration_method_no	    smallserial	PRIMARY KEY,
    registration_method_name    text        NOT NULL UNIQUE,
    created_at                  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         	    timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			        text		NULL
);
INSERT INTO account.registration_method(registration_method_name) VALUES ('web application'), ('telegram account'), ('google account');



DROP TABLE IF EXISTS account.role CASCADE;
CREATE TABLE account.role
(
    role_no      smallserial PRIMARY KEY,
    role_name    text        NOT NULL UNIQUE,
    created_at   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at  timestamp	 NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary   text		 NULL
);
INSERT INTO account.role(role_name) VALUES ('user'), ('admin');



DROP TABLE IF EXISTS account.account CASCADE;
CREATE TABLE account.account
(
    account_id              uuid        PRIMARY KEY DEFAULT account.UUID_GENERATE_V4(),
    account_state   	    smallint 	NOT NULL DEFAULT 1,
    account_role            smallint    NOT NULL DEFAULT 1,
    last_change_state       timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    registration_method     smallint    NOT NULL DEFAULT 1,
	last_activity	        timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at              timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_state) REFERENCES account.state(state_no),
    FOREIGN KEY (account_role) REFERENCES account.role(role_no),
    FOREIGN KEY (registration_method) REFERENCES account.registration_method(registration_method_no)
);



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



DROP TABLE IF EXISTS account.telegram_user CASCADE;
CREATE TABLE account.telegram_user
(
    account_id               uuid        NOT NULL UNIQUE,
    telegram_id              text        NOT NULL UNIQUE,
    username                 text        NOT NULL,
    photo_url                text        NOT NULL,
    modified_at         	    timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_id) REFERENCES account.account(account_id)
);



DROP TABLE IF EXISTS account.employee CASCADE;
CREATE TABLE account.employee
(
    account_id              uuid        NOT NULL UNIQUE,
    surname                 text        NOT NULL,
    name                    text        NOT NULL,
    patronymic              text        NULL,
    password 				text		NOT NULL,
    salt_for_password       text        NOT NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_id) REFERENCES account.account(account_id)
);










DROP TABLE IF EXISTS product.state CASCADE;
CREATE TABLE product.state
(
    state_no	smallserial	PRIMARY KEY,
    state_name  text 	    NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	text		NULL
);
INSERT INTO product.state(state_name) VALUES ('unavailable'), ('active'), ('blocked'), ('deleted');



DROP TABLE IF EXISTS product.type CASCADE;
CREATE TABLE product.type
(
    type_no	    smallserial	PRIMARY KEY,
    type_name   text 		NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary  text		NULL
);
INSERT INTO product.type(type_name) VALUES ('games'), ('software'), ('media content'), ('e-tickets'), ('virtual gifts');



DROP TABLE IF EXISTS product.subtype CASCADE;
CREATE TABLE product.subtype
(
    type_no         smallint    ,
    subtype_no      smallserial UNIQUE ,
    subtype_name    text        NOT NULL UNIQUE,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary      text		NULL,
    PRIMARY KEY (type_no, subtype_no),
    FOREIGN KEY (type_no) REFERENCES product.type(type_no)
);
INSERT INTO product.subtype(type_no, subtype_name) VALUES
                                                       ('1', 'computer games'),
                                                       ('1', 'mobile games'),
                                                       ('1', 'console games'),
                                                       ('2', 'office applications'),
                                                       ('2', 'antivirus software'),
                                                       ('2', 'design software'),
                                                       ('2', 'video editing software'),
                                                       ('2', 'audio recording software'),
                                                       ('3', 'music'),
                                                       ('3', 'movies'),
                                                       ('3', 'books'),
                                                       ('3', 'audiobooks'),
                                                       ('3', 'magazines');



DROP TABLE IF EXISTS product.service CASCADE;
CREATE TABLE product.service
(
    service_no	    smallserial	PRIMARY KEY,
    service_name    text 		NOT NULL UNIQUE,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at    	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	    text		NULL
);
INSERT INTO product.service(service_name) VALUES ('steam'), ('ubisoft'), ('epic games'), ('electronic arts'), ('ozon'), ('wildberries'), ('ivi'), ('youtube');



DROP TABLE IF EXISTS product.item CASCADE;
CREATE TABLE product.item
(
    item_no         smallserial	PRIMARY KEY,
    item_name       text 		NOT NULL UNIQUE,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at    	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	    text		NULL
);
INSERT INTO product.item(item_name) VALUES ('key'), ('code'), ('usual link'), ('gift as link');



DROP TABLE IF EXISTS product.product CASCADE;
CREATE TABLE product.product
(
    product_id      uuid        PRIMARY KEY DEFAULT account.UUID_GENERATE_V4(),
    product_subtype smallint    NOT NULL,
    product_name    text        NOT NULL UNIQUE,
    description     text        NOT NULL,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		text		NULL,
    FOREIGN KEY (product_subtype) REFERENCES product.subtype(subtype_no)
);



DROP TABLE IF EXISTS product.variant CASCADE;
CREATE TABLE product.variant
(
    product_id          uuid        ,
    variant_id          uuid        UNIQUE DEFAULT account.UUID_GENERATE_V4(),
    variant_name        text        ,
    product_service     smallint    NOT NULL,
    product_state       smallint    NOT NULL DEFAULT 1,
    product_item        smallint    NOT NULL,
    mask                text        NOT NULL,
    last_change_state   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    quantity            integer     NOT NULL CHECK ( quantity >= -1 ) DEFAULT 0,
    price               money       NOT NULL CHECK ( price >= 0::money ),
    discount_money      money       NOT NULL DEFAULT 0,
    discount_percent    smallint    NOT NULL DEFAULT 0,
    account_id          uuid        NOT NULL,
    created_at          timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		    text		NULL,
    PRIMARY KEY (product_id, variant_id),
    FOREIGN KEY (product_id) REFERENCES product.product(product_id),
    FOREIGN KEY (product_state) REFERENCES product.state(state_no),
    FOREIGN KEY (product_service) REFERENCES product.service(service_no),
    FOREIGN KEY (product_item) REFERENCES product.item(item_no),
    FOREIGN KEY (account_id) REFERENCES account.account(account_id),
    CHECK (
            (discount_money = 0::money AND discount_percent > 0::smallint)
            OR
            (discount_money > 0::money AND discount_percent = 0::smallint)
        )
);



--GRANT USAGE ON SCHEMA xxxx TO user;
--GRANT SELECT, UPDATE, INSERT, DELETE ON ALL TABLES IN SCHEMA xxxx TO user;
--GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA xxxx TO user;
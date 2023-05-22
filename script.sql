SET timezone = 'UTC';

CREATE SCHEMA IF NOT EXISTS account;
CREATE SCHEMA IF NOT EXISTS product;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA account;

--select * from account.user;
--select * from account.account;
--select * from product.variant;



DROP TABLE IF EXISTS account.state CASCADE;
CREATE TABLE account.state
(
	state_no	smallserial	PRIMARY KEY,
	state_name  text 		NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	text		NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS account_state_name_idx ON account.state (lower(state_name));
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
CREATE UNIQUE INDEX IF NOT EXISTS account_registration_method_name_idx ON account.registration_method (lower(registration_method_name));
INSERT INTO account.registration_method(registration_method_name) VALUES ('web application'), ('telegram account'), ('google account'), ('from admin panel');



DROP TABLE IF EXISTS account.role CASCADE;
CREATE TABLE account.role
(
    role_no      smallserial PRIMARY KEY,
    role_name    text        NOT NULL UNIQUE,
    created_at   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at  timestamp	 NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary   text		 NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS account_role_name_idx ON account.role (lower(role_name));
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
INSERT INTO account.account(account_id) VALUES ('4ad0f276-b11b-4c17-a160-3671699f0694');
INSERT INTO account.account(account_id, account_role) VALUES ('4ad0f276-b11b-4c17-a160-3671699f0693', '2');



DROP TABLE IF EXISTS account.user CASCADE;
CREATE TABLE account.user
(
    user_account            uuid        NOT NULL UNIQUE,
    email                   text        NOT NULL UNIQUE,
    nickname                text        NOT NULL UNIQUE,
    password 				text		NOT NULL,
	salt_for_password       text        NOT NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (user_account) REFERENCES account.account(account_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS account_email_idx ON account.user (lower(email));
INSERT INTO account.user(user_account, email, nickname, password, salt_for_password) VALUES ('4ad0f276-b11b-4c17-a160-3671699f0694', '77lm@mail.ru', 'Evgenick', 'QDmOn45b1pvrdIeKpGo/QWhoh3Yk4SW6ohlqlmnEeY0', 'Q/04YJ4R9L2n8ZVMszEe+w');



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
    login                   text        NOT NULL UNIQUE,
    password 				text		NOT NULL,
    salt_for_password       text        NOT NULL,
    modified_at         	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary			    text		NULL,
    FOREIGN KEY (account_id) REFERENCES account.account(account_id)
);
INSERT INTO account.employee(account_id, surname, name, patronymic, login, password, salt_for_password) VALUES ('4ad0f276-b11b-4c17-a160-3671699f0693', 'Kovalev', 'Dmitry', NULL, 'administrator', 'QDmOn45b1pvrdIeKpGo/QWhoh3Yk4SW6ohlqlmnEeY0', 'Q/04YJ4R9L2n8ZVMszEe+w');



DROP TABLE IF EXISTS product.state CASCADE;
CREATE TABLE product.state
(
    state_no	smallserial	PRIMARY KEY,
    state_name  text 	    NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	text		NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS product_state_name_idx ON product.state (lower(state_name));
INSERT INTO product.state(state_name) VALUES ('unavailable without price'), ('active'), ('deleted'), ('unavailable with price'), ('invisible');



DROP TABLE IF EXISTS product.type CASCADE;
CREATE TABLE product.type
(
    type_no	    smallserial	PRIMARY KEY,
    type_name   text 		NOT NULL UNIQUE,
    created_at  timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary  text		NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS product_type_name_idx ON product.type (lower(type_name));
INSERT INTO product.type(type_name) VALUES ('games'), ('software'), ('media content'), ('e-tickets'), ('virtual gifts'), ('replenishment of in-game currency');



DROP TABLE IF EXISTS product.subtype CASCADE;
CREATE TABLE product.subtype
(
    type_no         smallint    ,
    subtype_no      serial      UNIQUE,
    subtype_name    text        NOT NULL UNIQUE,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary      text		NULL,
    PRIMARY KEY (type_no, subtype_no),
    FOREIGN KEY (type_no) REFERENCES product.type(type_no)
);
CREATE UNIQUE INDEX IF NOT EXISTS product_subtype_name_idx ON product.subtype (lower(subtype_name));
INSERT INTO product.subtype(type_no, subtype_name) VALUES
                                                       ('1', 'computer version'),
                                                       ('1', 'mobile version'),
                                                       ('1', 'console version'),
                                                       ('2', 'antivirus software'),
                                                       ('2', 'design software'),
                                                       ('3', 'music'),
                                                       ('3', 'books'),
                                                       ('6', 'g-coins'),
                                                       ('1', 'downloadable game content');



DROP TABLE IF EXISTS product.service CASCADE;
CREATE TABLE product.service
(
    service_no	    smallserial	PRIMARY KEY,
    service_name    text 		NOT NULL UNIQUE,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at    	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	    text		NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS product_service_name_idx ON product.service (lower(service_name));
INSERT INTO product.service(service_name) VALUES ('steam'), ('ubisoft'), ('epic games'), ('electronic arts'), ('discord'), ('youtube'), ('playstation'), ('xbox'), ('nintendo'), ('universal');



DROP TABLE IF EXISTS product.item CASCADE;
CREATE TABLE product.item
(
    item_no         smallserial	PRIMARY KEY,
    item_name       text 		NOT NULL UNIQUE,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at    	timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary	    text		NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS product_item_name_idx ON product.item (lower(item_name));
INSERT INTO product.item(item_name) VALUES ('key/code'), ('link');



DROP TABLE IF EXISTS product.product CASCADE;
CREATE TABLE product.product
(
    product_id      uuid        PRIMARY KEY DEFAULT account.UUID_GENERATE_V4(),
    product_name    text        NOT NULL UNIQUE,
    description     text        NOT NULL,
    tags            text        NULL,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		text		NULL
);
INSERT INTO product.product(product_id, product_name, tags, description) VALUES
('9beaf75e-2925-4815-bcb6-1dd364293848', 'Grand Theft Auto 5', 'gta 5, gta5', 'Лос-Сантос – город солнца, старлеток и вышедших в тираж звезд. Некогда предмет зависти всего западного мира, ныне это пристанище дрянных реалити-шоу, задыхающееся в тисках экономических проблем. В центре всей заварухи – троица совершенно разных преступников, отчаянно пытающихся ухватить удачу за хвост в непрекращающейся борьбе за место под солнцем. Бывший член уличной банды Франклин старается завязать с прошлым. Отошедший от дел грабитель банков Майкл обнаруживает, что в честной жизни все не так радужно, как представлялось. Повернутый на насилии псих Тревор перебивается от одного дельца к другому в надежде сорвать крупный куш. Исчерпав варианты, эти трое ставят на кон собственные жизни и учиняют серию дерзких ограблений, в которых – или пан, или пропал.'),
('af153d1a-263c-4e77-a4d0-fb11ce781365', 'Red Dead Redemption 2', 'rdr2, rdr 2, rdo', 'Америка, 1899 год. Артур Морган и другие подручные Датча ван дер Линде вынуждены пуститься в бега. Их банде предстоит участвовать в кражах, грабежах и перестрелках в самом сердце Америки. За ними по пятам идут федеральные агенты и лучшие в стране охотники за головами, а саму банду разрывают внутренние противоречия. Артуру предстоит выбрать, что для него важнее: его собственные идеалы или же верность людям, которые его взрастили.'),
('85f8d115-ca4b-4db5-b416-6828e4c0e90a', 'Warframe', NULL, 'Пробудитесь в роли неудержимого воина и сражайтесь вместе с друзьями в этой сюжетной бесплатной онлайн-игре. Столкнитесь с враждующими фракциями в обширной межпланетной системе, следуя указаниям загадочной Лотос, повышайте уровень своего Варфрейма, создайте арсенал разрушительной огневой мощи, и откройте свой истинный потенциал в огромных открытых мирах этого захватывающего сражения от третьего лица.'),
('573b8cea-bbfa-4415-8f16-1b793a97c85f', 'PUBG: BATTLEGROUNDS', NULL, 'Высаживайтесь в стратегически важных местах, добывайте оружие и припасы и постарайтесь выжить и остаться последней командой на одном из многочисленных полей боя.'),
('7a33fa78-df96-4b7e-ac64-4f152ca2022f', 'Superliminal', NULL, 'Восприятие – это реальность. В этой умопомрачительной головоломке от первого лица вам предстоит сбежать из сюрреалистического мира снов, решая невозможные загадки при помощи перспективы.');



DROP TABLE IF EXISTS product.variant CASCADE;
CREATE TABLE product.variant
(
    product_id          uuid        ,
    variant_id          uuid        UNIQUE DEFAULT account.UUID_GENERATE_V4(),
    variant_name        text        ,
    variant_service     smallint    NOT NULL,
    variant_state       smallint    NOT NULL DEFAULT 1,
    variant_subtype     integer     NOT NULL,
    last_change_state   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    variant_item        smallint    NOT NULL,
    mask                text        NOT NULL,
    quantity_current    integer     NOT NULL CHECK ( quantity_current >= 0 ) DEFAULT 0,
    quantity_sold       integer     NOT NULL CHECK ( quantity_sold >= 0 ) DEFAULT 0,
    price               numeric     NOT NULL CHECK ( price >= 0 ),
    discount_money      numeric     NOT NULL CHECK ( discount_money >= 0 AND discount_money <= price ) DEFAULT 0,
    discount_percent    smallint    NOT NULL CHECK ( discount_percent >= 0 AND discount_percent <= 100 ) DEFAULT 0,
    variant_account     uuid        NOT NULL,
    created_at          timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at         timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		    text		NULL,
    PRIMARY KEY (product_id, variant_id),
    FOREIGN KEY (product_id) REFERENCES product.product(product_id),
    FOREIGN KEY (variant_state) REFERENCES product.state(state_no),
    FOREIGN KEY (variant_subtype) REFERENCES product.subtype(subtype_no),
    FOREIGN KEY (variant_service) REFERENCES product.service(service_no),
    FOREIGN KEY (variant_item) REFERENCES product.item(item_no),
    FOREIGN KEY (variant_account) REFERENCES account.account(account_id),
    UNIQUE (variant_name, variant_service, variant_subtype),
    CHECK (
        (discount_money = 0 AND discount_percent = 0::smallint)
            OR
        (discount_money = 0 AND discount_percent > 0::smallint)
            OR
        (discount_money > 0 AND discount_percent = 0::smallint)
        )
);
INSERT INTO product.variant(product_id, variant_id, variant_name, variant_service, variant_state, variant_subtype, variant_item, mask, price, discount_money, discount_percent, variant_account)VALUES
('9beaf75e-2925-4815-bcb6-1dd364293848', '96c7074d-a205-4489-ba38-e5ebbcb676cd', 'Grand Theft Auto 5: Premium Edition', 1, 2, 1, 1, 'XXXXX-XXXXX-XXXXX', 1199, 0, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('9beaf75e-2925-4815-bcb6-1dd364293848', 'e8fac47e-cff7-4a98-9ff9-1524eb826bc3', 'Grand Theft Auto 5: Premium Edition', 7, 2, 3, 1, 'XXXXX-00000000-YYYYY', 1755, 100, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('af153d1a-263c-4e77-a4d0-fb11ce781365', '7193a698-7d38-41c1-93aa-0ba82942a87a', 'Red Dead Redemption 2', 1, 2, 1, 1, 'XXXXX-XXXXX-XXXXX', 1199, 0, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('af153d1a-263c-4e77-a4d0-fb11ce781365', '7f3ef0ea-d506-4300-890d-78810187db96', 'Red Dead Redemption 2: Ultimate Edition', 9, 2, 3, 1, 'XXXXX-YYYYY-XXXXX', 1999, 0, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('af153d1a-263c-4e77-a4d0-fb11ce781365', '5eac26f5-92e2-48ba-81e8-e03c86a90814', 'Red Dead Online', 3, 2, 1, 1, 'XXXXX-XXXXX-XXXXX', 699, 0, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('85f8d115-ca4b-4db5-b416-6828e4c0e90a', 'c58f405a-8b38-4afd-8ffc-757243498702', 'Warframe', 1, 2, 1, 1, 'YYYYY-XXXXX-XXXXX', 299, 0, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('573b8cea-bbfa-4415-8f16-1b793a97c85f', 'fb8d23d5-cbce-450e-8600-1acc14981fcb', 'PUBG: BATTLEGROUNDS: Deluxe Edition', 1, 2, 1, 1, 'XXXXX-YYYYY-XXXXX', 2299, 0, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('573b8cea-bbfa-4415-8f16-1b793a97c85f', '646d24f6-6434-454e-91b6-daf572d177bc', 'PUBG: BATTLEGROUNDS: Ultimate Edition', 1, 2, 1, 1, 'XXXXX-YYYYY-XXXXX', 2999, 0, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('7a33fa78-df96-4b7e-ac64-4f152ca2022f', 'ea40a18c-c2d2-41aa-9de3-113c483fc22d', 'Superliminal', 3, 2, 1, 1, 'XXXXX-XXXXX-YYYYY', 699, 100, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('573b8cea-bbfa-4415-8f16-1b793a97c85f', '5eb83209-1645-467b-9d0f-6819fc8ca2bb', '200 G-Coins', 1, 2, 8, 1, 'XXXXX-XXXXX-XXXXX-XXXXX', 199, 0, 50, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('573b8cea-bbfa-4415-8f16-1b793a97c85f', 'b0210718-7571-4428-86b5-5bed693ed2d4', '300 G-Coins', 1, 2, 8, 1, 'XXXXX-XXXXX-XXXXX-XXXXX', 259, 0, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('573b8cea-bbfa-4415-8f16-1b793a97c85f', '8c6f1ae5-d33a-42f1-b0d3-6bcb25ce9d94', '400 G-Coins', 1, 2, 8, 1, 'XXXXX-XXXXX-XXXXX-XXXXX', 339, 200, 0, '4ad0f276-b11b-4c17-a160-3671699f0694'),
('573b8cea-bbfa-4415-8f16-1b793a97c85f', '4407a26e-96ac-4880-afad-c80220d3218d', '600 G-Coins', 1, 2, 8, 1, 'XXXXX-XXXXX-XXXXX-XXXXX', 599, 0, 10, '4ad0f276-b11b-4c17-a160-3671699f0694');



DROP TABLE IF EXISTS product.order CASCADE;
CREATE TABLE product.order
(
    order_id        uuid        PRIMARY KEY DEFAULT account.UUID_GENERATE_V4(),
    order_account   uuid        NOT NULL,
    price           numeric     NOT NULL CHECK ( price >= 0 ),
    paid            bool        NOT NULL DEFAULT false,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		text		NULL,
    FOREIGN KEY (order_account) REFERENCES account.account(account_id)
);



DROP TABLE IF EXISTS product.content CASCADE;
CREATE TABLE product.content
(
    content_id      uuid        PRIMARY KEY DEFAULT account.UUID_GENERATE_V4(),
    content_variant uuid        NOT NULL,
    content_order   uuid        UNIQUE NULL DEFAULT NULL,
    data            text        NOT NULL,
    created_at      timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp	NOT NULL DEFAULT CURRENT_TIMESTAMP,
    commentary		text		NULL,
    FOREIGN KEY (content_variant) REFERENCES product.variant(variant_id),
    FOREIGN KEY (content_order) REFERENCES product.order(order_id)
);
INSERT INTO product.content (content_variant, data)
VALUES
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '2T7Q9-4F2S5-WP6R4-DHGK1-BE8K9'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '6U5B7-3G8N2-KV4M1-JCDR3-AF2V9'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '9Y3Z1-5H6M4-XR7G9-PK2L8-BF1N6'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '1C2X3-7W5R9-KF9G2-JV4N6-BL6Y1'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '5P7O2-9S6V3-XL1J4-GZ1H6-BR4T8'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '4A5R8-6T9S3-VN7G1-PK1H9-BZ3L6'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '2K7F4-9H3D1-XT5R6-DL8J2-BP6Y1'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '8U3B5-6F9N4-KZ1X7-PV4G2-BR7T1'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '3S6C9-7G4N1-LB2X5-DH5R8-BZ1V3'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '6T7Q4-2F5S9-WR8P6-JD3G2-BN1K7'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '9Y2Z7-5H4M1-XL6G9-PC2N3-BV8F4'),
    ('96c7074d-a205-4489-ba38-e5ebbcb676cd', '1C5X3-3W6R9-KG2V8-JL4D7-BF9N5');



DROP MATERIALIZED VIEW IF EXISTS product.product_variants_summary_all_data;
CREATE MATERIALIZED VIEW product.product_variants_summary_all_data AS
SELECT
    t.type_name,
    st.subtype_name,
    s.service_name,
    p.product_name,
    pv.variant_name,
    ps.state_name,
    pv.price,
    pv.discount_money,
    pv.discount_percent,
    CASE
        WHEN pv.discount_money > 0::numeric THEN pv.price - pv.discount_money
        WHEN pv.discount_percent > 0 THEN pv.price * (1 - pv.discount_percent / 100.0)
        ELSE pv.price
        END AS final_price,
    i.item_name,
    pv.mask,
    pv.quantity_current,
    pv.quantity_sold,
    CASE
        WHEN pv.quantity_current = 0 THEN 'out of stock'
        WHEN pv.quantity_current = 1 THEN 'last in stock'
        WHEN pv.quantity_current > 1 AND pv.quantity_current < 10 THEN 'limited stock'
        WHEN pv.quantity_current >= 10 AND pv.quantity_current < 30 THEN 'adequate stock'
        WHEN pv.quantity_current >= 30 THEN 'large stock'
        ELSE 'error'
        END AS text_quantity,
    p.description,
    p.tags,
    p.product_id,
    pv.variant_id,
    pv.variant_account
FROM
    product.variant pv
        JOIN product.product p ON pv.product_id = p.product_id
        JOIN product.service s ON pv.variant_service = s.service_no
        JOIN product.state ps ON pv.variant_state = ps.state_no
        JOIN product.item i ON pv.variant_item = i.item_no
        JOIN product.subtype st ON pv.variant_subtype = st.subtype_no
        JOIN product.type t ON st.type_no = t.type_no;



REFRESH MATERIALIZED VIEW product.product_variants_summary_all_data;
SELECT * FROM product.product_variants_summary_all_data;



--GRANT USAGE ON SCHEMA xxxx TO user;
--GRANT SELECT, UPDATE, INSERT, DELETE ON ALL TABLES IN SCHEMA xxxx TO user;
--GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA xxxx TO user;
--
-- Individual Microservice Databases
--

-- common schema
CREATE DATABASE commondb TEMPLATE template0;
\c commondb
CREATE TABLE events
(
    entity_name    text        NOT NULL,
    entity_id      text        NOT NULL,
    correlation_id text        NOT NULL,
    causation_id   text        NOT NULL,
    event_version  int         NOT NULL,
    event_name     text        NOT NULL,
    event_data     bytea       NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entity_name, entity_id, event_version)
);
CREATE TABLE messages
(
    message_id  text        NOT NULL,
    destination text        NOT NULL,
    payload     bytea       NOT NULL,
    headers     bytea       NOT NULL,
    published   boolean     NOT NULL DEFAULT false,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id)
);
CREATE INDEX unpublished_idx ON messages (created_at) WHERE not published;
CREATE INDEX published_idx ON messages (modified_at) WHERE published;
CREATE TABLE saga_instances
(
    saga_name      text        NOT NULL,
    saga_id        text        NOT NULL,
    saga_data_name text        NOT NULL,
    saga_data      bytea       NOT NULL,
    current_step   int         NOT NULL,
    end_state      boolean     NOT NULL,
    compensating   boolean     NOT NULL,
    modified_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (saga_name, saga_id)
);
CREATE TABLE snapshots
(
    entity_name      text        NOT NULL,
    entity_id        text        NOT NULL,
    snapshot_name    text        NOT NULL,
    snapshot_data    bytea       NOT NULL,
    snapshot_version int         NOT NULL,
    modified_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entity_name, entity_id)
);

-- accounting service
CREATE DATABASE accounting TEMPLATE commondb;
\c accounting;
CREATE USER accounting_user WITH ENCRYPTED PASSWORD 'accounting_pass';
GRANT CONNECT ON DATABASE accounting TO accounting_user;
GRANT USAGE ON SCHEMA public TO accounting_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public TO accounting_user;

-- consumer service
CREATE DATABASE consumer TEMPLATE commondb;
\c consumer;
CREATE USER consumer_user WITH ENCRYPTED PASSWORD 'consumer_pass';
GRANT CONNECT ON DATABASE consumer TO consumer_user;
GRANT USAGE ON SCHEMA public TO consumer_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public TO consumer_user;

-- delivery service
CREATE DATABASE delivery TEMPLATE commondb;
\c delivery;
CREATE TABLE couriers
(
    id          text,
    plan        bytea,
    available   boolean,
    modified_at timestamptz,
    PRIMARY KEY (id)
);
CREATE INDEX first_available_courier_idx ON couriers (modified_at DESC) WHERE available;
CREATE TABLE deliveries
(
    id               text,
    restaurant_id    text,
    courier_id       text,
    pickup_address   bytea,
    delivery_address bytea,
    pickup_time      timestamptz,
    ready_by         timestamptz,
    status           text,
    PRIMARY KEY (id)
);
CREATE TABLE restaurants
(
    id      text,
    name    text,
    address bytea,
    menu    bytea,
    PRIMARY KEY (id)
);
CREATE USER delivery_user WITH ENCRYPTED PASSWORD 'delivery_pass';
GRANT CONNECT ON DATABASE delivery TO delivery_user;
GRANT USAGE ON SCHEMA public TO delivery_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public TO delivery_user;

-- kitchen service
CREATE DATABASE kitchen TEMPLATE commondb;
\c kitchen;
CREATE TABLE restaurants
(
    id      text,
    name    text,
    address bytea,
    menu    bytea,
    PRIMARY KEY (id)
);
CREATE USER kitchen_user WITH ENCRYPTED PASSWORD 'kitchen_pass';
GRANT CONNECT ON DATABASE kitchen TO kitchen_user;
GRANT USAGE ON SCHEMA public TO kitchen_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public TO kitchen_user;

-- ordering service
CREATE DATABASE ordering TEMPLATE commondb;
\c ordering;
CREATE TABLE restaurants
(
    id      text,
    name    text,
    address bytea,
    menu    bytea,
    PRIMARY KEY (id)
);
CREATE USER ordering_user WITH ENCRYPTED PASSWORD 'ordering_pass';
GRANT CONNECT ON DATABASE ordering TO ordering_user;
GRANT USAGE ON SCHEMA public TO ordering_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public TO ordering_user;

-- order history service
CREATE DATABASE order_history TEMPLATE commondb;
\c order_history;
CREATE TABLE orders
(
    id              text,
    consumer_id     text,
    restaurant_id   text,
    restaurant_name text,
    line_items      bytea,
    order_total     int,
    status          int,
    keywords        text[],
    created_at      timestamptz,
    PRIMARY KEY (id)
);
CREATE INDEX consumer_orders_seek_idx ON orders (consumer_id, created_at DESC, id DESC);
CREATE USER order_history_user WITH ENCRYPTED PASSWORD 'order_history_pass';
GRANT CONNECT ON DATABASE ordering TO order_history_user;
GRANT USAGE ON SCHEMA public TO order_history_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public TO order_history_user;

-- restaurant service
CREATE DATABASE restaurant TEMPLATE commondb;
\c restaurant;
CREATE TABLE restaurants
(
    id      text,
    name    text,
    address bytea,
    menu    bytea,
    PRIMARY KEY (id)
);
CREATE USER restaurant_user WITH ENCRYPTED PASSWORD 'restaurant_pass';
GRANT CONNECT ON DATABASE restaurant TO restaurant_user;
GRANT USAGE ON SCHEMA public TO restaurant_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public TO restaurant_user;

--
-- Monolith Database with Schemas
--
CREATE DATABASE ftgogo;
\c ftgogo
CREATE TABLE events
(
    entity_name    text        NOT NULL,
    entity_id      text        NOT NULL,
    correlation_id text        NOT NULL,
    causation_id   text        NOT NULL,
    event_version  int         NOT NULL,
    event_name     text        NOT NULL,
    event_data     bytea       NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entity_name, entity_id, event_version)
);
CREATE TABLE messages
(
    message_id  text        NOT NULL,
    destination text        NOT NULL,
    payload     bytea       NOT NULL,
    headers     bytea       NOT NULL,
    published   boolean     NOT NULL DEFAULT false,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id)
);
CREATE INDEX unpublished_idx ON messages (created_at) WHERE not published;
CREATE INDEX published_idx ON messages (modified_at) WHERE published;
CREATE TABLE saga_instances
(
    saga_name      text        NOT NULL,
    saga_id        text        NOT NULL,
    saga_data_name text        NOT NULL,
    saga_data      bytea       NOT NULL,
    current_step   int         NOT NULL,
    end_state      boolean     NOT NULL,
    compensating   boolean     NOT NULL,
    modified_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (saga_name, saga_id)
);
CREATE TABLE snapshots
(
    entity_name      text        NOT NULL,
    entity_id        text        NOT NULL,
    snapshot_name    text        NOT NULL,
    snapshot_data    bytea       NOT NULL,
    snapshot_version int         NOT NULL,
    modified_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entity_name, entity_id)
);

CREATE USER ftgogo_user WITH ENCRYPTED PASSWORD 'ftgogo_pass';
GRANT CONNECT ON DATABASE ftgogo TO ftgogo_user;
REVOKE USAGE ON SCHEMA public FROM ftgogo_user;
REVOKE INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA public FROM ftgogo_user;

-- accounting
CREATE SCHEMA accounting;
CREATE TABLE accounting.events (LIKE public.events INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE accounting.messages (LIKE public.messages INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE accounting.saga_instances (LIKE public.saga_instances INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE accounting.snapshots (LIKE public.snapshots INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
-- accounting access
GRANT USAGE ON SCHEMA accounting TO ftgogo_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA accounting TO ftgogo_user;

-- consumer
CREATE SCHEMA consumer;
CREATE TABLE consumer.events (LIKE public.events INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE consumer.messages (LIKE public.messages INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE consumer.saga_instances (LIKE public.saga_instances INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE consumer.snapshots (LIKE public.snapshots INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
-- consumer access
GRANT USAGE ON SCHEMA consumer TO ftgogo_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA consumer TO ftgogo_user;

-- delivery
CREATE SCHEMA delivery;
CREATE TABLE delivery.events (LIKE public.events INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE delivery.messages (LIKE public.messages INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE delivery.saga_instances (LIKE public.saga_instances INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE delivery.snapshots (LIKE public.snapshots INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE delivery.couriers
(
    id          text,
    plan        bytea,
    available   boolean,
    modified_at timestamptz,
    PRIMARY KEY (id)
);
CREATE INDEX first_available_courier_idx ON delivery.couriers (modified_at DESC) WHERE available;
CREATE TABLE delivery.deliveries
(
    id               text,
    restaurant_id    text,
    courier_id       text,
    pickup_address   bytea,
    delivery_address bytea,
    pickup_time      timestamptz,
    ready_by         timestamptz,
    status           text,
    PRIMARY KEY (id)
);
CREATE TABLE delivery.restaurants
(
    id      text,
    name    text,
    address bytea,
    menu    bytea,
    PRIMARY KEY (id)
);
-- delivery access
GRANT USAGE ON SCHEMA delivery TO ftgogo_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA delivery TO ftgogo_user;

-- kitchen
CREATE SCHEMA kitchen;
CREATE TABLE kitchen.events (LIKE public.events INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE kitchen.messages (LIKE public.messages INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE kitchen.saga_instances (LIKE public.saga_instances INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE kitchen.snapshots (LIKE public.snapshots INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE kitchen.restaurants
(
    id      text,
    name    text,
    address bytea,
    menu    bytea,
    PRIMARY KEY (id)
);
-- kitchen access
GRANT USAGE ON SCHEMA kitchen TO ftgogo_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA kitchen TO ftgogo_user;

-- orders
CREATE SCHEMA orders;
CREATE TABLE orders.events (LIKE public.events INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE orders.messages (LIKE public.messages INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE orders.saga_instances (LIKE public.saga_instances INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE orders.snapshots (LIKE public.snapshots INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE orders.restaurants
(
    id      text,
    name    text,
    address bytea,
    menu    bytea,
    PRIMARY KEY (id)
);
-- orders access
GRANT USAGE ON SCHEMA orders TO ftgogo_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA orders TO ftgogo_user;

-- orderhistory
CREATE SCHEMA orderhistory;
CREATE TABLE orderhistory.events (LIKE public.events INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE orderhistory.messages (LIKE public.messages INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE orderhistory.saga_instances (LIKE public.saga_instances INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE orderhistory.snapshots (LIKE public.snapshots INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE orderhistory.orders
(
    id              text,
    consumer_id     text,
    restaurant_id   text,
    restaurant_name text,
    line_items      bytea,
    order_total     int,
    status          int,
    keywords        text[],
    created_at      timestamptz,
    PRIMARY KEY (id)
);
CREATE INDEX consumer_orders_seek_idx ON orderhistory.orders (consumer_id, created_at DESC, id DESC);
-- orderhistory access
GRANT USAGE ON SCHEMA orderhistory TO ftgogo_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA orderhistory TO ftgogo_user;

-- restaurant
CREATE SCHEMA restaurant;
CREATE TABLE restaurant.events (LIKE public.events INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE restaurant.messages (LIKE public.messages INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE restaurant.saga_instances (LIKE public.saga_instances INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE restaurant.snapshots (LIKE public.snapshots INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES);
CREATE TABLE restaurant.restaurants
(
    id      text,
    name    text,
    address bytea,
    menu    bytea,
    PRIMARY KEY (id)
);
-- restaurant access
GRANT USAGE ON SCHEMA restaurant TO ftgogo_user;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA restaurant TO ftgogo_user;

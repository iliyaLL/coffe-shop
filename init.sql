CREATE TYPE status AS ENUM ('open', 'closed');

CREATE TABLE orders (
    id serial primary key,
    customer_name varchar(255) not null,
    order_status status not null,
    created_at timestamp not null,
    customer_preferences jsonb not null default '{}'::jsonb
);

CREATE TABLE order_status_history (
    id serial primary key,
    order_id int references orders (id) on delete cascade,
    updated_at timestamp not null,
    old_status status not null,
    new_status status not null
);

CREATE TABLE menu_items (
    id serial primary key,
    name varchar(255) not null,
    description varchar(1000) not null,
    price decimal(10, 2) not null
);

CREATE TABLE price_history (
    id serial primary key,
    menu_item_id int references menu_items (id) on delete cascade,
    old_price decimal(10,2) not null,
    new_price decimal(10,2) not null,
    updated_at timestamp not null
);

CREATE TABLE order_item (
    order_id int references orders (id) on delete cascade,
    menu_item_id int references menu_items (id) on delete cascade,
    quantity int not null
);

CREATE TYPE unit AS ENUM ('shots', 'ml', 'g');

CREATE TABLE inventory (
    id serial primary key,
    name varchar(255) not null,
    quantity int not null default 0,
    unit unit not null,
    categories varchar(50)[]
);

CREATE TABLE inventory_transactions (
    id serial primary key,
    inventory_id int references inventory (id) on delete cascade,
    old_quantity int not null,
    new_quantity int not null,
    transaction_date timestamp not null
);

CREATE TABLE menu_item_inventory (
    menu_id int references menu_items (id) on delete cascade,
    inventory_id int references inventory (id) on delete cascade,
    quantity int not null
);
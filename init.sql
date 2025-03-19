CREATE TYPE status AS ENUM ('open', 'in progress', 'closed');

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
    name varchar(255) not null unique,
    description varchar(1000) not null,
    price decimal(10, 2) not null constraint positive_price CHECK (price >= 0)
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
    quantity int not null constraint positive_quantity CHECK (quantity >= 0)
);

CREATE TYPE unit AS ENUM ('shots', 'ml', 'g', 'units');

CREATE TABLE inventory (
    id serial primary key,
    name varchar(255) not null unique,
    quantity int not null default 0 constraint positive_quantity CHECK (quantity >= 0),
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
    quantity int not null constraint positive_quantity CHECK (quantity >= 0)
);

INSERT INTO inventory (name, quantity, unit, categories) VALUES
('Espresso Shot', 500, 'shots', ARRAY['Beverage']),
('Milk', 5000, 'ml', ARRAY['Dairy']),
('Flour', 10000, 'g', ARRAY['Baking']),
('Blueberries', 2000, 'g', ARRAY['Fruit']),
('Raspberry', 2000, 'g', ARRAY['Fruit']),
('Sugar', 5000, 'g', ARRAY['Baking', 'Sweetener']),
('Coffee Beans', 5000, 'g', ARRAY['Beverage', 'Raw Material']),
('Ground Coffee', 3000, 'g', ARRAY['Beverage']),
('Vanilla Syrup', 2000, 'ml', ARRAY['Flavoring']),
('Caramel Syrup', 2000, 'ml', ARRAY['Flavoring']),
('Chocolate Syrup', 2500, 'ml', ARRAY['Flavoring']),
('Whipped Cream', 1000, 'ml', ARRAY['Dairy', 'Topping']),
('Tea Leaves', 1500, 'g', ARRAY['Beverage', 'Raw Material']),
('Honey', 1000, 'ml', ARRAY['Sweetener', 'Flavoring']),
('Pastry Dough', 5000, 'g', ARRAY['Baking']),
('Butter', 2000, 'g', ARRAY['Dairy']),
('Eggs', 300, 'units', ARRAY['Baking', 'Dairy']);

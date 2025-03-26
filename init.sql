CREATE TYPE status AS ENUM ('open', 'in progress', 'closed');

CREATE TABLE orders (
    id serial primary key,
    customer_name varchar(255) not null,
    order_status status not null,
    created_at timestamp not null default now(),
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



INSERT INTO menu_items (id, name, description, price) VALUES
(30, 'Blueberry Muffin', 'Freshly baked muffin with blueberries', 2.00),
(31, 'Raspberry Muffin', 'Muffin with fresh raspberries', 2.00),
(32, 'Strawberry Muffin', 'Freshly baked muffin with strawberries', 2.00),
(33, 'Caffe Latte', 'Espresso with steamed milk', 3.50),
(34, 'Espresso', 'A strong shot of coffee', 2.00),
(35, 'Vanilla Cappuccino', 'Espresso with vanilla syrup and foam', 3.80),
(36, 'Caramel Macchiato', 'Espresso with caramel syrup and steamed milk', 4.20),
(37, 'Chocolate Frappe', 'Blended chocolate drink with whipped cream', 4.50);


-- Blueberry Muffin
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(30, 3, 100),  -- Flour
(30, 4, 50),   -- Blueberries
(30, 6, 10),   -- Sugar
(30, 15, 100), -- Pastry Dough
(30, 16, 20);  -- Butter

-- Raspberry Muffin
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(31, 3, 100),
(31, 5, 50),
(31, 6, 10),
(31, 15, 100),
(31, 16, 20);

-- Strawberry Muffin
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(32, 3, 100),
(32, 1, 30),
(32, 10, 20),
(32, 2, 100),
(32, 16, 20);

-- Caffe Latte
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(33, 1, 1),    -- Espresso Shot
(33, 2, 200);  -- Milk

-- Espresso
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(34, 1, 1);

-- Vanilla Cappuccino
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(35, 1, 1),
(35, 2, 150),
(35, 9, 30); -- Vanilla Syrup

-- Caramel Macchiato
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(36, 1, 1),
(36, 2, 200),
(36, 10, 30); -- Caramel Syrup

-- Chocolate Frappe
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(37, 11, 100), -- Chocolate Syrup
(37, 2, 100),
(37, 12, 50);  -- Whipped Cream



-- Function for inventory quantity tracking
CREATE OR REPLACE FUNCTION log_inventory_transaction()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.quantity <> NEW.quantity THEN
        INSERT INTO inventory_transactions (inventory_id, old_quantity, new_quantity, transaction_date)
        VALUES (NEW.id, OLD.quantity, NEW.quantity, NOW());
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to track quantity updates in inventory
CREATE TRIGGER after_inventory_update
AFTER UPDATE ON inventory
FOR EACH ROW
EXECUTE FUNCTION log_inventory_transaction();


-- Function for price change tracking
CREATE OR REPLACE FUNCTION log_price_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.price <> NEW.price THEN
        INSERT INTO price_history (menu_item_id, old_price, new_price, updated_at)
        VALUES (NEW.id, OLD.price, NEW.price, NOW());
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to track price updates in menu_items
CREATE TRIGGER after_price_update
AFTER UPDATE ON menu_items
FOR EACH ROW
EXECUTE FUNCTION log_price_change();


-- Function for order status tracking
CREATE OR REPLACE FUNCTION log_order_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.order_status <> NEW.order_status THEN
        INSERT INTO order_status_history (order_id, updated_at, old_status, new_status)
        VALUES (NEW.id, NOW(), OLD.order_status, NEW.order_status);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to track status updates in orders
CREATE TRIGGER after_order_status_update
AFTER UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION log_order_status_change();
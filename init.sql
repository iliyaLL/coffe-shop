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
('Eggs', 300, 'units', ARRAY['Baking', 'Dairy']),
('Cinnamon', 1500, 'g', ARRAY['Spice']),
('Nutmeg', 1000, 'g', ARRAY['Spice']),
('Matcha Powder', 800, 'g', ARRAY['Tea', 'Flavoring']);


INSERT INTO menu_items (name, description, price) VALUES
('Blueberry Muffin', 'Freshly baked muffin with blueberries', 2.00),
('Raspberry Muffin', 'Muffin with fresh raspberries', 2.00),
('Strawberry Muffin', 'Freshly baked muffin with strawberries', 2.00),
('Caffe Latte', 'Espresso with steamed milk', 3.50),
('Espresso', 'A strong shot of coffee', 2.00),
('Vanilla Cappuccino', 'Espresso with vanilla syrup and foam', 3.80),
('Caramel Macchiato', 'Espresso with caramel syrup and steamed milk', 4.20),
('Chocolate Frappe', 'Blended chocolate drink with whipped cream', 4.50),
('Matcha Latte', 'Green tea with steamed milk', 3.60),
('Chai Tea Latte', 'Spiced tea with milk', 3.70);


-- Blueberry Muffin
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(1, 3, 100),  -- Flour
(1, 4, 50),   -- Blueberries
(1, 6, 10),   -- Sugar
(1, 15, 100), -- Pastry Dough
(1, 16, 20);  -- Butter

-- Raspberry Muffin
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(2, 3, 100),
(2, 5, 50),
(2, 6, 10),
(2, 15, 100),
(2, 16, 20);

-- Strawberry Muffin
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(3, 3, 100),
(3, 1, 30),
(3, 10, 20),
(3, 2, 100),
(3, 16, 20);

-- Caffe Latte
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(4, 1, 1),    -- Espresso Shot
(4, 2, 200);  -- Milk

-- Espresso
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(5, 1, 1);

-- Vanilla Cappuccino
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(6, 1, 1),
(6, 2, 150),
(6, 9, 30); -- Vanilla Syrup

-- Caramel Macchiato
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(7, 1, 1),
(7, 2, 200),
(7, 10, 30); -- Caramel Syrup

-- Chocolate Frappe
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(8, 11, 100), -- Chocolate Syrup
(8, 2, 100),
(8, 12, 50);  -- Whipped Cream

-- Matcha Latte
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(9, 2, 150),     -- Milk
(9, 18, 10);     -- Matcha Powder

-- Chai Tea Latte
INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES
(10, 2, 150),    -- Milk
(10, 17, 5),     -- Cinnamon
(10, 18, 3);     -- Nutmeg


INSERT INTO orders (customer_name, order_status) VALUES
('Gary Soto', 'closed'),
('Shawn Todd II', 'closed'),
('Cynthia Miller', 'closed'),
('James Johnson', 'open'),
('Charles Henry', 'closed'),
('Andrea Miller', 'open'),
('Patricia Villa', 'closed'),
('Jonathan Hernandez', 'closed'),
('Mindy Reynolds', 'closed'),
('Lori Palmer', 'open'),
('Mitchell Mercer', 'closed'),
('Dr. Sandra Brown DDS', 'closed'),
('Jacqueline Obrien', 'closed'),
('Diana Sanders', 'closed'),
('Stephen Davis', 'open'),
('Carrie Clayton', 'closed'),
('Kyle Randall', 'closed'),
('Ronald Levine', 'open'),
('Rebecca Nixon', 'open'),
('Jaime Robinson', 'closed'),
('Jessica Bell', 'open'),
('David Ramirez', 'closed'),
('Karen Brooks', 'closed'),
('William Bates', 'open'),
('Gerald Benson MD', 'open'),
('Christopher Wolfe', 'closed'),
('Lisa Reynolds', 'open'),
('Michael Sexton', 'closed'),
('Edward Horne', 'open'),
('Bianca Lopez', 'open');



INSERT INTO order_item (order_id, menu_item_id, quantity) VALUES
(1, 2, 1),
(1, 9, 4),
(1, 5, 5),
(2, 4, 5),
(2, 8, 2),
(2, 3, 4),
(3, 6, 3),
(4, 9, 5),
(5, 2, 2),
(6, 10, 2),
(6, 1, 1),
(7, 2, 2),
(7, 8, 5),
(8, 8, 1),
(8, 3, 4),
(9, 7, 2),
(10, 9, 2),
(10, 3, 3),
(10, 2, 4),
(11, 1, 2),
(11, 7, 1),
(12, 9, 1),
(13, 9, 4),
(13, 5, 3),
(13, 7, 4),
(14, 5, 1),
(15, 8, 3),
(15, 1, 1),
(15, 3, 1),
(16, 9, 2),
(17, 6, 4),
(17, 7, 4),
(17, 8, 5),
(18, 2, 1),
(18, 10, 5),
(18, 3, 2),
(19, 5, 3),
(19, 2, 2),
(19, 10, 2),
(20, 6, 2),
(20, 7, 5),
(20, 2, 4),
(21, 4, 1),
(22, 3, 2),
(22, 8, 2),
(22, 2, 4),
(23, 5, 3),
(23, 6, 5),
(23, 1, 3),
(24, 7, 1),
(25, 9, 2),
(25, 1, 5),
(26, 9, 4),
(27, 8, 4),
(28, 5, 4),
(29, 1, 2),
(30, 1, 2);

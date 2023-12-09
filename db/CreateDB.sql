create table product_categories (
    category_name varchar not null,
    description varchar default '',
    id serial primary key
);
insert into product_categories values
('Аткрытка', 'аписание'),
('Стикерпук', 'очень классный');


create table products (
    name varchar not null,
    default_cost numeric(7, 2) not null,
    category integer not null,
    self_cost numeric(6, 2) default 0,
    amount integer default 0,
    id serial primary key,
    foreign key (category) references product_categories (id)
);
insert into products values
('Открытка', 15, 1, 5, 11),
('Стикерпак', 20, 2, 5, 2),
('Ого, еще открытка', 1, 1, 5, 4),
('Новый стикерпак', 11, 2, 5, 3),
('Плоттерная фигнюшка', 4, 1, 5, 7),
('Серьги', 7, 1, 5, 100);


create table statuses (
    description varchar not null,
    id serial primary key
);
insert into statuses values
('OK'),
('Refunding'),
('Successfully refunded');


create table positions (
    product integer not null,
    cost numeric(7, 2) default 0,
    count integer default 1 check ( count > 0 ),
    status integer default 1 check ( status > 0 and status < 4 ),
    id integer primary key,
    foreign key (status) references statuses (id),
    foreign key (product) references products (id)
);


create table receipts (
    date date not null,
    status integer default 1 check ( status > 0 and status < 4 ),
    id serial primary key,
    foreign key (status) references statuses (id)
);


create table positions_in_receipts (
    position integer not null,
    receipt integer not null,
    primary key (position, receipt),
    foreign key (position) references positions (id) on delete cascade,
    foreign key (receipt) references receipts (id)
);


create table markets (
    name varchar not null,
    dates daterange not null,
    fee integer not null default 0,
    id serial primary key
);


create table receipts_on_markets (
    receipt integer not null,
    market integer not null,
    primary key (market, receipt),
    foreign key (market) references markets (id),
    foreign key (receipt) references receipts (id) on delete cascade
);


create table roles (
    name varchar unique not null,
    id serial primary key
);
insert into roles values
('Admin'),
('Salesman');


create table users (
    role_id integer not null,
    username varchar unique not null,
    password varchar not null,
    id serial primary key,
    foreign key (role_id) references roles (id) on delete cascade
);
insert into users values (1, 'Admin', 'Passw0rd');


create table permissions (
    role_id integer,
    table_name varchar,
    permission varchar,
    primary key (role_id, table_name, permission),
    foreign key (role_id) references roles (id) on delete cascade
);
insert into permissions values
(1, 'all', 'all'),
(2, 'product_categories', 'read'),
(2, 'products', 'read'),
(2, 'statuses', 'read'),
(2, 'positions', 'modify'),
(2, 'receipts', 'modify'),
(2, 'positions_in_receipts', 'modify'),
(2, 'markets', 'read'),
(2, 'roles', 'read'),
(2, 'permissions', 'read');
create schema goodbuy;

/*      TABLES      */

create table goodbuy.product_categories (
    category_name varchar not null,
    description varchar default '',
    id serial primary key
);
insert into goodbuy.product_categories values
('Без категории', 'Стандартная категория по умолчанию, либо в случае ошибки', 0);


create table goodbuy.products (
    name varchar not null unique,
    default_cost numeric(7, 2) check ( default_cost >= 0 ),
    category integer not null,
    self_cost numeric(6, 2) default 0 check ( self_cost >= 0 ),
    amount integer default 0 check ( amount >= 0 ),
    id serial primary key,
    foreign key (category) references goodbuy.product_categories (id)
);


create table goodbuy.statuses (
    description varchar not null,
    id serial primary key
);
insert into goodbuy.statuses values
('OK'),
('Refunding'),
('Successfully refunded');


create table goodbuy.positions (
    product integer not null,
    cost numeric(7, 2) default 0 check ( cost >= 0 ),
    count integer default 1 check ( count > 0 ),
    status integer default 1 check ( status > 0 and status < 4 ),
    id serial primary key,
    foreign key (status) references goodbuy.statuses (id),
    foreign key (product) references goodbuy.products (id)
);


create table goodbuy.receipts (
    date date not null,
    status integer default 1 check ( status > 0 and status < 4 ),
    id serial primary key,
    foreign key (status) references goodbuy.statuses (id)
);


create table goodbuy.positions_in_receipts (
    position integer not null,
    receipt integer not null,
    primary key (position, receipt),
    foreign key (position) references goodbuy.positions (id) on delete cascade,
    foreign key (receipt) references goodbuy.receipts (id)
);


create table goodbuy.markets (
    name varchar not null,
    dates daterange not null,
    fee integer not null default 0,
    id serial primary key
);


create table goodbuy.receipts_on_markets (
    receipt integer not null,
    market integer not null,
    primary key (market, receipt),
    foreign key (market) references goodbuy.markets (id),
    foreign key (receipt) references goodbuy.receipts (id) on delete cascade
);

create table goodbuy.roles (
    name varchar unique not null,
    id serial primary key
);
insert into goodbuy.roles values
('Unknown', 0);
insert into goodbuy.roles values
('Admin'),
('Salesman'),
('Analyst');


create table goodbuy.users (
    role_id integer not null,
    username varchar unique not null,
    password varchar not null,
    id serial primary key,
    foreign key (role_id) references goodbuy.roles (id) on delete cascade
);
insert into goodbuy.users values (1, 'Administrator', 'd41e98d1eafa6d6011d3a70f1a5b92f0');



/*      VIEWS       */

create view goodbuy.get_products as
    select
        p.id,
        p.name,
        p.default_cost,
        p.self_cost,
        p.amount,
        pc.category_name
    from goodbuy.products p
    join goodbuy.product_categories pc on p.category = pc.id;


create view goodbuy.get_detailed_receipts as
    select
        r.date,
        r.id as receipt_id,
        r_s.description as receipt_status,
        pir.position as id,
        pr.name as product,
        p.cost as cost,
        p.count as count,
        p_s.description as status
    from goodbuy.receipts r
    join goodbuy.positions_in_receipts pir on r.id = pir.receipt
    join goodbuy.positions p on p.id = pir.position
    join goodbuy.products pr on pr.id = p.product
    join goodbuy.statuses r_s on r.status = r_s.id
    join goodbuy.statuses p_s on p.status = p_s.id;


/*      PROCEDURES      */

create function goodbuy.new_receipt() returns integer as $$
declare
    new_id integer;
begin
    insert into goodbuy.receipts values (current_date) returning id into new_id;
    return new_id;
end;
$$
language plpgsql;

create procedure goodbuy.new_position(
    rec_id integer,
    new_product varchar,
    new_cost numeric,
    new_count integer
) as $$
declare
    pos_id integer;
begin
    insert into goodbuy.positions values
    (
     (select id from goodbuy.products where name = new_product),
     new_cost,
     new_count
    ) returning id into pos_id;

    insert into goodbuy.positions_in_receipts values (pos_id, rec_id);
end;
$$
language plpgsql;


create procedure goodbuy.add_product(
    name varchar,
    default_cost numeric(7, 2),
    category_name_add varchar,
    self_cost numeric(6, 2),
    amount integer
) as $$
declare
    category_id integer;
begin
    select id into category_id
              from goodbuy.product_categories as p_c
              where p_c.category_name = category_name_add;
    if category_id is null then
        category_id = 0;
    end if;
    insert into goodbuy.products
    values (
            add_product.name,
            add_product.default_cost,
            category_id,
            add_product.self_cost,
            add_product.amount
           );
end;
$$
language plpgsql;


create procedure goodbuy.edit_product(
    id_edit integer,
    name varchar,
    default_cost numeric(7, 2),
    category_name_edit varchar,
    self_cost numeric(6, 2),
    amount integer
) as $$
declare
    category_id integer;
begin
    select cats.id into category_id
    from goodbuy.product_categories as cats
    where cats.category_name = category_name_edit;

    update goodbuy.products
    set name = edit_product.name,
        default_cost = edit_product.default_cost,
        category = category_id,
        self_cost = edit_product.self_cost,
        amount = edit_product.amount
    where id = edit_product.id_edit;
end;
$$
language plpgsql;


create function goodbuy.get_filtered_products (
    min_category integer,
    max_category integer,
    min_amount integer,
    max_amount integer,
    min_self_cost numeric,
    max_self_cost numeric,
    min_default_cost numeric,
    max_default_cost numeric
) returns table (
    product varchar,
    default_cost numeric,
    category integer,
    self_cost numeric,
    amount integer,
    id integer,
    category_name varchar
) as $$
begin
    return query
    select
        p.name,
        p.default_cost,
        p.category,
        p.self_cost,
        p.amount,
        p.id,
        pc.category_name
    from goodbuy.products p
    join goodbuy.product_categories pc on pc.id = p.category
    where
        (min_category = 0 or p.category >= min_category)
        and (max_category = 0 or p.category <= max_category)
        and (min_amount = 0 or p.amount >= min_amount)
        and (max_amount = 0 or p.amount <= max_amount)
        and (min_self_cost = 0 or p.self_cost >= min_self_cost)
        and (max_self_cost = 0 or p.self_cost <= max_self_cost)
        and (min_default_cost = 0 or p.default_cost >= min_default_cost)
        and (max_default_cost = 0 or p.default_cost <= max_default_cost);
end;
$$
language plpgsql;


/*      TRIGGERS        */

create function goodbuy.update_products_category()
returns trigger as $$
begin
    update goodbuy.products
    set category = 0
    where category = old.id;
    return old;
end;
$$
language plpgsql;

create trigger update_products_category_trigger_on_delete
    before delete on goodbuy.product_categories
    for each row execute function goodbuy.update_products_category();


create function goodbuy.check_default_cost()
returns trigger as $$
begin
    if new.self_cost > new.default_cost then
        raise exception 'can`t set default cost less than self cost. we`re making money.';
    end if;

    return new;
end;
$$
language plpgsql;

create trigger check_default_cost_trigger_on_insert
    before insert on goodbuy.products
    for each row execute function goodbuy.check_default_cost();


/*      FUNCTIONS       */

create function goodbuy.get_top_N_products_by_sales(past_days integer, n integer)
returns table(
    product varchar,
    sales bigint
    ) as $$
begin
    return query
    select
        pr.name as product_name,
        SUM(p.count) as total_sales
    from goodbuy.positions p
    join goodbuy.products pr on p.product = pr.id
    join goodbuy.positions_in_receipts pir on p.id = pir.position
    join goodbuy.receipts r on pir.receipt = r.id
    where r.date >= (current_date - past_days * interval '1 day')
    group by pr.name
    order by total_sales desc
    limit n;
end;
$$
language plpgsql;

create function goodbuy.get_top_N_products_by_profit(past_days integer, n integer)
returns table(
    product varchar,
    profit numeric
    ) as $$
begin
    return query
    select
        pr.name as product_name,
        SUM((p.cost - pr.self_cost) * p.count) as profit
    from goodbuy.positions p
    join goodbuy.products pr on p.product = pr.id
    join goodbuy.positions_in_receipts pir on p.id = pir.position
    join goodbuy.receipts r on pir.receipt = r.id
    where r.date >= (current_date - past_days * interval '1 day')
    group by pr.name
    order by profit desc
    limit n;
end;
$$
language plpgsql;

create function goodbuy.get_N_popular_products_on_markets(n integer)
returns table (
    market varchar,
    product varchar,
    sales bigint
    ) as $$
begin
    return query
    with RankedProducts as (
        select m.name as this_market,
        pr.name as this_product,
        sum(p.count) as total_sales,
        row_number() over (partition by m.name order by sum(p.count) desc) as rank
        from
            goodbuy.positions p
        join
            goodbuy.products pr on p.product = pr.id
        join
            goodbuy.positions_in_receipts pir on p.id = pir.position
        join
            goodbuy.receipts r on pir.receipt = r.id
        join
            goodbuy.receipts_on_markets rm on r.id = rm.receipt
        join
            goodbuy.markets m on rm.market = m.id
        group by
            m.name, pr.name
    )
    select this_market, this_product, total_sales from RankedProducts where rank <= n;
end;
$$
language plpgsql;

create function goodbuy.get_income_past_N_days(n integer)
returns bigint as $$
declare
    total_income bigint;
begin
    select coalesce(sum(p.count * p.cost), 0) as total
    into total_income
    from goodbuy.positions p
    join goodbuy.positions_in_receipts pr on p.id = pr.position
    join goodbuy.receipts r on pr.receipt = r.id
    where r.status = 1
      and p.status = 1
      and r.date >= (current_date - n * interval '1 day');

    return total_income;
end;
$$
language plpgsql;